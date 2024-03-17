package cli

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Default	config
var defaultConfigYml = map[string]string{
	"model":    "gpt-3.5-turbo",
	"voiceId":  "Joanna",
	"language": "English",
}

var configFolderName = "/jarvis/"
var configFileName = "jarvis.yml"

func getConfigDir() (string, error) {
	var userConfigDir, err = os.UserConfigDir()
	return userConfigDir + configFolderName, err
}

func CreateConfigDirectory() error {
	var err error
	var cligptConfigDir string
	if cligptConfigDir, err = getConfigDir(); err != nil {
		return err
	}
	var cligptConfigFile = cligptConfigDir + configFileName

	os.MkdirAll(cligptConfigDir, 0755)

	if _, err = os.Stat(cligptConfigFile); !os.IsNotExist(err) {
		return nil
	}

	// create default config
	var token string
	if token, err = tokenRequest(); err != nil {
		return err
	}
	var language string
	if language, err = languageRequest(); err != nil {
		return err
	}
	defaultConfigYml["language"] = language
	if defaultConfigYml["language"] == "Spanish" {
		defaultConfigYml["voiceId"] = "Lupe"
	} else {
		defaultConfigYml["voiceId"] = "Joanna"
	}
	defaultConfigYml["auth"] = token

	file, err := os.Create(cligptConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var bffYml = yaml.NewEncoder(file)
	defer bffYml.Close()
	return bffYml.Encode(defaultConfigYml)
}

func ReadYml() (map[string]string, error) {
	var err error
	var cligptConfigDir string
	if cligptConfigDir, err = getConfigDir(); err != nil {
		return nil, err
	}
	var cligptConfigFile = cligptConfigDir + configFileName

	file, err := os.Open(cligptConfigFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config map[string]string
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

var errTokenLenNotValid = errors.New("the token len is not valid")

func tokenRequest() (string, error) {
	var token string
	fmt.Print("Please write your OpenAI token: ")
	fmt.Scanln(&token)
	if len(token) == 0 {
		return token, errTokenLenNotValid
	}
	return token, nil
}

func languageRequest() (string, error) {
	var language string
	fmt.Print("Please write the language (English or Spanish): ")
	fmt.Scanln(&language)
	if language != "English" && language != "Spanish" {
		return language, errors.New("the language is not valid")
	} else {
		return language, nil
	}
}
