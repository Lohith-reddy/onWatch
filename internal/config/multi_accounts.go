package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// MultiAccountConfig defines one externally configured account for live quota checks.
// Provider supports "codex" and "anthropic".
type MultiAccountConfig struct {
	Name            string `json:"name"`
	Provider        string `json:"provider"`
	Token           string `json:"token,omitempty"`
	TokenEnv        string `json:"token_env,omitempty"`
	AuthFile        string `json:"auth_file,omitempty"`
	CredentialsFile string `json:"credentials_file,omitempty"`
	AccountID       string `json:"account_id,omitempty"`
}

// ParseMultiAccountConfig parses ONWATCH_MULTI_ACCOUNTS from JSON.
// Invalid entries are skipped so a single bad account does not disable the feature.
func ParseMultiAccountConfig(raw string) []MultiAccountConfig {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	var accounts []MultiAccountConfig
	if err := json.Unmarshal([]byte(raw), &accounts); err != nil {
		return nil
	}

	valid := make([]MultiAccountConfig, 0, len(accounts))
	for _, acct := range accounts {
		acct.Name = strings.TrimSpace(acct.Name)
		acct.Provider = strings.ToLower(strings.TrimSpace(acct.Provider))
		acct.Token = strings.TrimSpace(acct.Token)
		acct.TokenEnv = strings.TrimSpace(acct.TokenEnv)
		acct.AuthFile = normalizePath(acct.AuthFile)
		acct.CredentialsFile = normalizePath(acct.CredentialsFile)
		acct.AccountID = strings.TrimSpace(acct.AccountID)

		if acct.Name == "" {
			continue
		}
		if acct.Provider != "codex" && acct.Provider != "anthropic" {
			continue
		}
		if !hasTokenSource(acct) {
			continue
		}
		valid = append(valid, acct)
	}

	if len(valid) == 0 {
		return nil
	}
	return valid
}

func hasTokenSource(acct MultiAccountConfig) bool {
	if acct.Token != "" || acct.TokenEnv != "" {
		return true
	}
	if acct.Provider == "codex" {
		return acct.AuthFile != ""
	}
	if acct.Provider == "anthropic" {
		return acct.CredentialsFile != ""
	}
	return false
}

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if path == "~" {
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			path = home
		}
	} else if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}
	return filepath.Clean(path)
}
