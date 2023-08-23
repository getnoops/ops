package tokenstore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zitadel/oidc/v2/pkg/oidc"
)

const (
	noOpsFolder    = ".no_ops"
	configFileName = "no_opsconfig"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	IDToken      string
}

func getNoOpsFolderPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, noOpsFolder), nil
}

func validateTokens(tokens Tokens) error {
	if tokens.AccessToken == "" {
		return errors.New("access token not found")
	}
	if tokens.RefreshToken == "" {
		return errors.New("refresh token not found")
	}
	if tokens.TokenType == "" {
		return errors.New("token type not found")
	}
	if tokens.IDToken == "" {
		return errors.New("ID token not found")
	}
	return nil
}

func Store(tokens *oidc.Tokens[*oidc.IDTokenClaims]) error {
	noOpsFolderPath, err := getNoOpsFolderPath()
	if err != nil {
		return fmt.Errorf("error getting user's home directory: %v", err)
	}

	_, err = os.Stat(noOpsFolderPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(noOpsFolderPath, 0755) // Create the folder with permission 0755
		if err != nil {
			return fmt.Errorf("error creating '.no_ops' folder: %v", err)
		}
		fmt.Println("'.no_ops' folder created.")
	}

	configFilePath := filepath.Join(noOpsFolderPath, configFileName)

	tokenData := []byte(fmt.Sprintf(
		"ACCESS_TOKEN=%s\nREFRESH_TOKEN=%s\nTOKEN_TYPE=%s\nID_TOKEN=%s",
		tokens.AccessToken, tokens.RefreshToken, tokens.TokenType, tokens.IDToken))

	err = os.WriteFile(configFilePath, tokenData, 0600)
	if err != nil {
		return fmt.Errorf("error writing tokens to 'no_opsconfig' file: %v", err)
	}

	fmt.Println("Tokens written to 'no_opsconfig' file.")
	return nil
}

func Retrieve() (*Tokens, error) {
	noOpsFolderPath, err := getNoOpsFolderPath()
	if err != nil {
		return nil, fmt.Errorf("error getting user's home directory: %v", err)
	}

	configFilePath := filepath.Join(noOpsFolderPath, configFileName)

	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading 'no_opsconfig' file: %v", err)
	}

	var tokens Tokens
	lines := strings.Split(string(configData), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "ACCESS_TOKEN":
			tokens.AccessToken = value
		case "REFRESH_TOKEN":
			tokens.RefreshToken = value
		case "TOKEN_TYPE":
			tokens.TokenType = value
		case "ID_TOKEN":
			tokens.IDToken = value
		}
	}
	err = validateTokens(tokens)
	if err != nil {
		return nil, err
	}
	return &tokens, nil
}

func UpdateTokens(accessToken string, refreshToken string) error {
	noOpsFolderPath, err := getNoOpsFolderPath()
	if err != nil {
		return fmt.Errorf("error getting user's home directory: %v", err)
	}

	configFilePath := filepath.Join(noOpsFolderPath, configFileName)

	// Read the existing configuration
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading 'no_opsconfig' file: %v", err)
	}

	// Create a map to store the updated key-value pairs
	updatedConfig := make(map[string]string)
	for _, line := range strings.Split(string(configData), "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			updatedConfig[key] = parts[1]
		}
	}

	// Update the key-value pairs with the provided tokens
	updatedConfig["ACCESS_TOKEN"] = accessToken
	updatedConfig["REFRESH_TOKEN"] = refreshToken

	// Generate the new configuration data
	var updatedConfigLines []string
	for key, value := range updatedConfig {
		updatedConfigLines = append(updatedConfigLines, fmt.Sprintf("%s=%s", key, value))
	}
	updatedConfigData := []byte(strings.Join(updatedConfigLines, "\n"))

	// Write the updated configuration back to the file
	err = os.WriteFile(configFilePath, updatedConfigData, 0600)
	if err != nil {
		return fmt.Errorf("error writing updated tokens to 'no_opsconfig' file: %v", err)
	}

	return nil
}
