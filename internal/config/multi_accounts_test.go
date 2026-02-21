package config

import "testing"

func TestParseMultiAccountConfig(t *testing.T) {
	raw := `[
		{"name":"work-codex","provider":"codex","auth_file":"~/.codex-work/auth.json"},
		{"name":"personal-claude","provider":"anthropic","token_env":"CLAUDE_PERSONAL_TOKEN"},
		{"name":"invalid","provider":"zai","token":"x"},
		{"name":"missing-source","provider":"codex"}
	]`

	accounts := ParseMultiAccountConfig(raw)
	if len(accounts) != 2 {
		t.Fatalf("expected 2 valid accounts, got %d", len(accounts))
	}

	if accounts[0].Provider != "codex" {
		t.Fatalf("first provider = %s, want codex", accounts[0].Provider)
	}
	if accounts[0].AuthFile == "" {
		t.Fatalf("expected normalized auth_file for codex account")
	}
	if accounts[1].Provider != "anthropic" {
		t.Fatalf("second provider = %s, want anthropic", accounts[1].Provider)
	}
	if accounts[1].TokenEnv != "CLAUDE_PERSONAL_TOKEN" {
		t.Fatalf("unexpected token_env: %s", accounts[1].TokenEnv)
	}
}

func TestParseMultiAccountConfig_InvalidJSON(t *testing.T) {
	if got := ParseMultiAccountConfig("{bad}"); got != nil {
		t.Fatalf("expected nil for invalid JSON, got %#v", got)
	}
}
