package ssh

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type KeyManager struct {
	sshDir string
}

func NewKeyManager() (*KeyManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %v", err)
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	if err := ensureSSHDir(sshDir); err != nil {
		return nil, err
	}

	return &KeyManager{
		sshDir: sshDir,
	}, nil
}

func ensureSSHDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("failed to create SSH directory: %v", err)
		}
	}
	return nil
}

func (km *KeyManager) AddKey(key, comment string) error {
	authorizedKeysPath := filepath.Join(km.sshDir, "authorized_keys")
	
	// Validate the SSH key format
	if !strings.HasPrefix(strings.TrimSpace(key), "ssh-") {
		return fmt.Errorf("invalid SSH key format")
	}

	// Format the key with the comment
	formattedKey := fmt.Sprintf("%s %s\n", strings.TrimSpace(key), comment)

	// Append the key to authorized_keys
	f, err := os.OpenFile(authorizedKeysPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open authorized_keys: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(formattedKey); err != nil {
		return fmt.Errorf("failed to write key: %v", err)
	}

	return nil
}

func (km *KeyManager) ListKeys() ([]string, error) {
	authorizedKeysPath := filepath.Join(km.sshDir, "authorized_keys")
	
	content, err := ioutil.ReadFile(authorizedKeysPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read authorized_keys: %v", err)
	}

	var keys []string
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line != "" && !strings.HasPrefix(line, "#") {
			keys = append(keys, line)
		}
	}

	return keys, scanner.Err()
}

func (km *KeyManager) RemoveKey(keyToRemove string) error {
	authorizedKeysPath := filepath.Join(km.sshDir, "authorized_keys")
	
	content, err := ioutil.ReadFile(authorizedKeysPath)
	if err != nil {
		return fmt.Errorf("failed to read authorized_keys: %v", err)
	}

	var newLines []string
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		// Keep comments, empty lines, and non-matching keys
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") || trimmedLine != keyToRemove {
			newLines = append(newLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %v", err)
	}

	newContent := strings.Join(newLines, "\n")
	if len(newLines) > 0 {
		newContent += "\n"
	}

	if err := ioutil.WriteFile(authorizedKeysPath, []byte(newContent), 0600); err != nil {
		return fmt.Errorf("failed to write updated authorized_keys: %v", err)
	}

	return nil
} 