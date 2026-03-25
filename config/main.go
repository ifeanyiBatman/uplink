package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct{
	AuthToken string
	Domain string
	Port string

}

type cfgJSON struct {
	CurrentUser string 			`json: "currentUser"`
	UserCFgs map[string]Config	`json: "userConfigs"`
}



const configFileName = ".uplinkconfig.json"


func GetCurrentUser () (string, error) {
	home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("failed to find Home DIR")
			return "", err
		}
	configFileURL := filepath.Join(home, configFileName)
	content , err := os.ReadFile(configFileURL)
	var configs cfgJSON

	err = json.Unmarshal(content,&configs)
	if err != nil {
		return "" , err
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
	content , err := os.ReadFile(configFileURL)
	var configs cfgJSON

	err = json.Unmarshal(content,&configs)
	if err != nil {
		return Config{} , err
	}
	userCFG , exists := configs.UserCFgs[username]
	if !exists {
		return Config{}, fmt.Errorf("User doesn't exist")
	}
	configs.CurrentUser = username

	return userCFG, nil

}

func (cfg *Config) SetUser (username string) error {
	home, err := os.UserHomeDir()
		if err != nil {
			fmt.Errorf("failed to find Home dir: %w", err)
		}
	configFileURL := filepath.Join(home, configFileName)



	var configs cfgJSON
	content , err := os.ReadFile(configFileURL)
	if err == nil {
		err = json.Unmarshal(content,&configs)
		if err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
	}
	} else if os.IsNotExist(err) {
		configs = cfgJSON{
			UserCFgs: make(map[string]Config),
		}
	} else {
		return fmt.Errorf("failed to rread config: %w", err)
	}

	if configs.UserCFgs == nil {
		configs = cfgJSON{
			UserCFgs: make(map[string]Config),
		}
	}

	configs.UserCFgs[username] = *cfg
	configs.CurrentUser = username

	data , err := json.MarshalIndent(configs,"","  ")
	if err != nil{
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configFileURL, data , 0644)
	if err != nil {
		return err
	}
	return nil
}


func (cfg Config) GetLocalhost() string {
	return fmt.Sprintf("http://localhost:%s", cfg.Port)
} 