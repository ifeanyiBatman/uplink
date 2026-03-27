package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	AuthToken string
	Domain    string
	Port      string
}

type cfgJSON struct {
	CurrentUser string            `json:"currentUser"`
	UserCFgs    map[string]Config `json:"userConfigs"`
}

const configFileName = ".uplinkconfig.json"

func GetCurrentUser() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("failed to find Home DIR")
		return "", err
	}
	configFileURL := filepath.Join(home, configFileName)
	content, err := os.ReadFile(configFileURL)
	var configs cfgJSON

	err = json.Unmarshal(content, &configs)
	if err != nil {
		return "", err
	}
	return configs.CurrentUser, nil

}

func GetUserConfig(username string) (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("failed to find Home DIR")
		return Config{}, err
	}
	configFileURL := filepath.Join(home, configFileName)
	content, err := os.ReadFile(configFileURL)
	var configs cfgJSON

	err = json.Unmarshal(content, &configs)
	if err != nil {
		return Config{}, err
	}
	userCFG, exists := configs.UserCFgs[username]
	if !exists {
		return Config{}, fmt.Errorf("User doesn't exist")
	}
	configs.CurrentUser = username

	return userCFG, nil

}

func (cfg *Config) SetUser(username string) error {
	configs, err := getConfigFile()
	if err != nil {
		return err
	}
	

	if configs.UserCFgs == nil {
		configs = cfgJSON{
			UserCFgs: make(map[string]Config),
		}
	}

	configs.UserCFgs[username] = *cfg
	configs.CurrentUser = username

	setConfigFile(configs)
	return nil
}

func getConfigFile() (cfgJSON, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return cfgJSON{}, fmt.Errorf("failed to find Home dir: %w", err)
	}
	configFileURL := filepath.Join(home, configFileName)

	var configs cfgJSON
	content, err := os.ReadFile(configFileURL)
	if err == nil {
		err = json.Unmarshal(content, &configs)
		if err != nil {
			return configs, fmt.Errorf("failed to parse config: %w", err)
		}
	} else if os.IsNotExist(err) {
		configs = cfgJSON{
			UserCFgs: make(map[string]Config),
		}
	} else {
		return configs, fmt.Errorf("failed to read config: %w", err)
	}

	return configs, nil
}

func setConfigFile(configs cfgJSON) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to find Home dir: %w", err)
	}
	configFileURL := filepath.Join(home, configFileName)
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configFileURL, data, 0644)
	if err != nil {
		return err
	}
	return nil

}

func Login(username string) error {
	configs, err := getConfigFile()
	if err != nil {
		return err
	}
	if configs.CurrentUser == username {
		return nil
	}
	_, ok := configs.UserCFgs[username]
	if ok {
		configs.CurrentUser = username
		setConfigFile(configs)
		return nil
	}
	return fmt.Errorf("can't find user/ or isnt registered")
}

func ListUsers() ([]string, error) {
	var users []string
	configs, err := getConfigFile()
	if err != nil {
		return users, fmt.Errorf("failed to list users : %w", err)
	}
	for user := range configs.UserCFgs {
		users = append(users, user)
		if user == configs.CurrentUser {
			fmt.Printf("%s[x]\n", user)
		} else {
			fmt.Println(user)
		}
	}
	return users, nil
}

func (cfg Config) GetLocalhost() string {
	return fmt.Sprintf("http://localhost:%s", cfg.Port)
}