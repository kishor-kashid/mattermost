# Environment Variables Setup for AI Features

## Quick Setup

1. **Create a `.env` file in the server directory:**
   ```bash
   cd server
   touch .env
   ```

2. **Add your OpenAI API key to the `.env` file:**
   ```bash
   # Required: Your OpenAI API key
   MM_AISETTINGS_OPENAIAPIKEY=your-openai-api-key-here
   ```

3. **Optional: Configure other AI settings:**
   ```bash
   # OpenAI Model (defaults to gpt-3.5-turbo)
   MM_AISETTINGS_OPENAIMODEL=gpt-3.5-turbo
   
   # Enable/disable AI features
   MM_AISETTINGS_ENABLE=true
   MM_AISETTINGS_ENABLESUMMARIZATION=true
   MM_AISETTINGS_ENABLEANALYTICS=true
   MM_AISETTINGS_ENABLEACTIONITEMS=true
   MM_AISETTINGS_ENABLEFORMATTING=true
   
   # Limits
   MM_AISETTINGS_MAXMESSAGELIMIT=500
   MM_AISETTINGS_APIRATELIMIT=60
   ```

4. **Restart Mattermost:**
   ```bash
   make restart
   ```

## Complete .env Template

```bash
# Mattermost AI Settings
# Environment variables override settings in config.json

# OpenAI API Key for AI features (REQUIRED)
# Get your API key from: https://platform.openai.com/api-keys
MM_AISETTINGS_OPENAIAPIKEY=sk-...your-key-here...

# OpenAI Model (optional - defaults to gpt-3.5-turbo)
# Options: gpt-3.5-turbo, gpt-4, gpt-4-turbo-preview
#MM_AISETTINGS_OPENAIMODEL=gpt-3.5-turbo

# Enable AI Features globally (optional - defaults to true)
#MM_AISETTINGS_ENABLE=true

# Enable specific AI features (optional - all default to true)
#MM_AISETTINGS_ENABLESUMMARIZATION=true
#MM_AISETTINGS_ENABLEANALYTICS=true
#MM_AISETTINGS_ENABLEACTIONITEMS=true
#MM_AISETTINGS_ENABLEFORMATTING=true

# Maximum messages to process for AI operations (optional - defaults to 500)
#MM_AISETTINGS_MAXMESSAGELIMIT=500

# API rate limit per minute (optional - defaults to 60)
#MM_AISETTINGS_APIRATELIMIT=60
```

## How It Works

Mattermost automatically reads environment variables that start with `MM_` and uses them to override configuration values in `config.json`. 

The pattern is:
- `MM_` + `SECTION_NAME` + `_` + `FIELD_NAME`
- Example: `MM_AISETTINGS_OPENAIAPIKEY` overrides `AISettings.OpenAIAPIKey`

## Security Notes

- ✅ The `.env` file is automatically ignored by git (listed in `.gitignore`)
- ✅ Never commit API keys to version control
- ✅ Use environment variables for sensitive data like API keys
- ✅ Keep `config.json` free of secrets

## Verification

After setting up your `.env` file and restarting:

1. **Check if the environment variable is loaded:**
   ```bash
   # On Linux/Mac
   echo $MM_AISETTINGS_OPENAIAPIKEY
   
   # On Windows PowerShell
   $env:MM_AISETTINGS_OPENAIAPIKEY
   ```

2. **Test the AI connection:**
   ```bash
   curl -X POST http://localhost:8065/api/v4/ai/test \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"test_prompt": "Hello"}'
   ```

## Alternative: System Environment Variables

Instead of using a `.env` file, you can also set system environment variables:

**Linux/Mac:**
```bash
export MM_AISETTINGS_OPENAIAPIKEY="sk-...your-key..."
```

**Windows PowerShell:**
```powershell
$env:MM_AISETTINGS_OPENAIAPIKEY = "sk-...your-key..."
```

**Windows Command Prompt:**
```cmd
set MM_AISETTINGS_OPENAIAPIKEY=sk-...your-key...
```

## Troubleshooting

### Environment variable not being read

Make sure you:
1. Created the `.env` file in the `server/` directory
2. Restarted Mattermost after creating/modifying `.env`
3. Used the correct variable name format (`MM_AISETTINGS_OPENAIAPIKEY`)

### Still seeing API key in config.json

That's okay! Environment variables override config.json values. You can:
1. Leave the config.json value empty
2. Or remove it entirely - the environment variable takes precedence

The environment variable will always be used if it's set, regardless of what's in config.json.

