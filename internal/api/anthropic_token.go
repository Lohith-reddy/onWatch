package api

import (
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"time"
)

// claudeCredentials represents the Claude Code credentials JSON structure.
type claudeCredentials struct {
	ClaudeAiOauth struct {
		AccessToken      string   `json:"accessToken"`
		RefreshToken     string   `json:"refreshToken"`
		ExpiresAt        int64    `json:"expiresAt"` // Unix milliseconds
		Scopes           []string `json:"scopes"`
		SubscriptionType string   `json:"subscriptionType"`
		RateLimitTier    string   `json:"rateLimitTier"`
	} `json:"claudeAiOauth"`
}

// AnthropicCredentials contains the parsed OAuth credentials with computed fields.
type AnthropicCredentials struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	ExpiresIn    time.Duration // time until expiry
	Scopes       []string
}

// IsExpiringSoon returns true if the token expires within the given duration.
func (c *AnthropicCredentials) IsExpiringSoon(threshold time.Duration) bool {
	return c.ExpiresIn < threshold
}

// IsExpired returns true if the token has already expired.
func (c *AnthropicCredentials) IsExpired() bool {
	return c.ExpiresIn <= 0
}

// parseClaudeCredentials extracts the OAuth access token from Claude Code credentials JSON.
func parseClaudeCredentials(data []byte) (string, error) {
	var creds claudeCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return "", err
	}
	return creds.ClaudeAiOauth.AccessToken, nil
}

// parseFullClaudeCredentials extracts all OAuth fields from Claude Code credentials JSON.
func parseFullClaudeCredentials(data []byte) (*AnthropicCredentials, error) {
	var creds claudeCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	oauth := creds.ClaudeAiOauth
	if oauth.AccessToken == "" {
		return nil, nil // no credentials
	}

	// Convert expiresAt from Unix milliseconds to time.Time
	expiresAt := time.UnixMilli(oauth.ExpiresAt)
	expiresIn := time.Until(expiresAt)

	return &AnthropicCredentials{
		AccessToken:  oauth.AccessToken,
		RefreshToken: oauth.RefreshToken,
		ExpiresAt:    expiresAt,
		ExpiresIn:    expiresIn,
		Scopes:       oauth.Scopes,
	}, nil
}

// DetectAnthropicToken attempts to auto-detect the Anthropic OAuth token
// from the Claude Code credentials stored in the system keychain or file.
// Returns empty string if not found.
func DetectAnthropicToken(logger *slog.Logger) string {
	return detectAnthropicTokenPlatform(logger)
}

// DetectAnthropicCredentials attempts to auto-detect the full Anthropic OAuth credentials
// from the Claude Code credentials stored in the system keychain or file.
// Returns nil if not found.
func DetectAnthropicCredentials(logger *slog.Logger) *AnthropicCredentials {
	return detectAnthropicCredentialsPlatform(logger)
}

// DetectAnthropicCredentialsFromFile parses Claude credentials JSON from a specific file path.
func DetectAnthropicCredentialsFromFile(path string, logger *slog.Logger) *AnthropicCredentials {
	if logger == nil {
		logger = slog.Default()
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logger.Debug("Credential file not readable", "path", path, "error", err)
		return nil
	}

	creds, err := parseFullClaudeCredentials(data)
	if err != nil {
		logger.Debug("Failed to parse full credentials", "path", path, "error", err)
		return nil
	}
	return creds
}

// WriteAnthropicAccessToken updates local Claude credentials with a new access token.
// If refresh token or expiry is missing, existing detected values are reused when available.
func WriteAnthropicAccessToken(accessToken, refreshToken string, expiresAt *time.Time, logger *slog.Logger) error {
	accessToken = strings.TrimSpace(accessToken)
	refreshToken = strings.TrimSpace(refreshToken)
	if accessToken == "" {
		return nil
	}

	if refreshToken == "" || expiresAt == nil {
		if existing := DetectAnthropicCredentials(logger); existing != nil {
			if refreshToken == "" {
				refreshToken = strings.TrimSpace(existing.RefreshToken)
			}
			if expiresAt == nil && !existing.ExpiresAt.IsZero() {
				t := existing.ExpiresAt
				expiresAt = &t
			}
		}
	}

	if refreshToken == "" {
		// Keep valid shape even if refresh token is unavailable.
		refreshToken = accessToken
	}

	expiry := 3600
	if expiresAt != nil {
		delta := int(time.Until(*expiresAt).Seconds())
		if delta > 0 {
			expiry = delta
		}
	}
	return WriteAnthropicCredentials(accessToken, refreshToken, expiry)
}
