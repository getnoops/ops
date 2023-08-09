package tokenstore

import (
	"fmt"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"os"
	"path/filepath"
	"strings"
)

const (
	noOpsFolder    = ".no_ops"
	configFileName = "no_opsconfig"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
}

func getNoOpsFolderPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, noOpsFolder), nil
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
		"ACCESS_TOKEN=%s\nREFRESH_TOKEN=%s\nTOKEN_TYPE=%s\n",
		tokens.AccessToken, tokens.RefreshToken, tokens.TokenType))

	err = os.WriteFile(configFilePath, tokenData, 0600)
	if err != nil {
		return fmt.Errorf("error writing tokens to 'no_opsconfig' file: %v", err)
	}

	fmt.Println("Tokens written to 'no_opsconfig' file.")
	return nil
}

func Retrieve() (tokens *Tokens, err error) {
	noOpsFolderPath, err := getNoOpsFolderPath()
	if err != nil {
		return nil, fmt.Errorf("error getting user's home directory: %v", err)
	}

	configFilePath := filepath.Join(noOpsFolderPath, configFileName)

	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading 'no_opsconfig' file: %v", err)
	}

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
		}
	}

	return tokens, nil
}
