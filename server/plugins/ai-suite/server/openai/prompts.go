package openai

import "strings"

// Template represents a reusable prompt consisting of system/user text.
type Template struct {
	Name        string
	System      string
	User        string
	Description string
}

// Render builds chat completion messages for the given template and variables.
func (t Template) Render(vars map[string]string) []Message {
	systemContent := replaceVars(t.System, vars)
	userContent := replaceVars(t.User, vars)

	messages := []Message{}
	if strings.TrimSpace(systemContent) != "" {
		messages = append(messages, Message{Role: RoleSystem, Content: systemContent})
	}
	if strings.TrimSpace(userContent) != "" {
		messages = append(messages, Message{Role: RoleUser, Content: userContent})
	}
	return messages
}

func replaceVars(content string, vars map[string]string) string {
	if len(vars) == 0 {
		return content
	}
	out := content
	for key, value := range vars {
		token := "{{" + key + "}}"
		out = strings.ReplaceAll(out, token, value)
	}
	return out
}

var (
	// SummarizeConversationTemplate produces summaries plus action items.
	SummarizeConversationTemplate = Template{
		Name: "summarize_conversation",
		System: "You are GPT-4 assisting users of the Mattermost collaboration platform. " +
			"Produce concise summaries highlighting key updates, decisions, and blockers.",
		User: "Channel: {{channel}}\n" +
			"Timeframe: {{timeframe}}\n\n" +
			"Conversation:\n{{conversation}}\n\n" +
			"Provide:\n1. A short paragraph summary.\n2. Bullet list of action items with owners if mentioned.",
		Description: "General conversation summarization with explicit action items.",
	}

	// ActionItemClassificationTemplate extracts commitments/tasks from content.
	ActionItemClassificationTemplate = Template{
		Name:   "action_item_classification",
		System: "You are GPT-4 identifying commitments and tasks in a Mattermost conversation.",
		User: "Review the following content and return JSON with fields `action_items` (array of {assignee, task, due_date}) " +
			"and `notes` for any important context. Only capture actionable commitments.\n\n{{conversation}}",
		Description: "Detects action items, responsible owners, and due dates.",
	}
)
