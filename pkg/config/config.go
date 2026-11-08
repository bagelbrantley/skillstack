package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Credentials struct {
	Token     string `json:"token"`
	UID       string `json:"uid"`
	Email     string `json:"email"`
	ServerURL string `json:"server_url"`
	SavedAt   string `json:"saved_at"`
}

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user home directory: %w", err)
	}
	dir := filepath.Join(home, ".config", "skillstack")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}
	return dir, nil
}

func GetCredentialsPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "credentials.json"), nil
}

func SaveCredentials(creds Credentials) error {
	path, err := GetCredentialsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode credentials: %w", err)
	}

	// Write file with strict 0600 permissions (read/write only by owner)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to save credentials file: %w", err)
	}
	return nil
}

func LoadCredentials() (*Credentials, error) {
	path, err := GetCredentialsPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, nil // No credentials file found
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials file: %w", err)
	}

	return &creds, nil
}

func ClearCredentials() error {
	path, err := GetCredentialsPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove credentials file: %w", err)
	}
	return nil
}
