package summarizer

import "errors"

// Type identifies the summarization context requested by the user.
type Type string

const (
	// TypeThread triggers summarization for a specific thread (root + replies).
	TypeThread Type = "thread"
	// TypeChannel triggers summarization for an entire channel/time range.
	TypeChannel Type = "channel"
)

var (
	// ErrDisabled is returned when summarization is disabled via config.
	ErrDisabled = errors.New("summarizer: feature disabled")
	// ErrUnauthorized indicates the user lacks permission for the target channel.
	ErrUnauthorized = errors.New("summarizer: insufficient permissions")
	// ErrInvalidRequest indicates malformed or missing parameters.
	ErrInvalidRequest = errors.New("summarizer: invalid request")
	// ErrEmptyConversation is returned when no qualifying posts were found.
	ErrEmptyConversation = errors.New("summarizer: no messages to summarize")
	// ErrClientUnavailable signals OpenAI is not configured.
	ErrClientUnavailable = errors.New("summarizer: llm client unavailable")
)

// Request contains the parameters passed to the summarizer service.
type Request struct {
	Type Type

	UserID string

	ChannelID  string
	RootPostID string
	PostID     string

	TimeRange   string
	SinceMillis int64
	UntilMillis int64

	Force bool
}

// Response contains the AI-generated summary plus contextual metadata.
type Response struct {
	ID               string         `json:"id"`
	Type             Type           `json:"type"`
	ChannelID        string         `json:"channel_id"`
	ChannelName      string         `json:"channel_name"`
	RootPostID       string         `json:"root_post_id,omitempty"`
	Title            string         `json:"title"`
	Summary          string         `json:"summary"`
	MessageCount     int            `json:"message_count"`
	ParticipantCount int            `json:"participant_count"`
	Participants     []Participant  `json:"participants"`
	GeneratedAt      int64          `json:"generated_at"`
	Range            SummaryRange   `json:"range"`
	Context          SummaryContext `json:"context"`
	Usage            *Usage         `json:"usage,omitempty"`
	LimitReached     bool           `json:"limit_reached"`
	Cached           bool           `json:"cached"`

	ConversationHash string `json:"-"`
}

// Participant captures metadata surfaced in the summary panel.
type Participant struct {
	ID          string `json:"id"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// SummaryRange represents the time window used for channel summaries.
type SummaryRange struct {
	Since int64  `json:"since"`
	Until int64  `json:"until"`
	Label string `json:"label"`
}

// SummaryContext adds presentation metadata for the RHS panel.
type SummaryContext struct {
	TypeLabel    string `json:"type_label"`
	MessageLimit int    `json:"message_limit"`
	Timeframe    string `json:"timeframe"`
}

// Usage surfaces LLM token accounting metadata.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
