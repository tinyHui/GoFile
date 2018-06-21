package utils

import (
	"io/ioutil"
	"os"

	"github.com/tinyhui/GoFile/utils/log"
	"gopkg.in/yaml.v2"
)

var logger = log.GetLogger()

type Parameters struct {
	StorageRoot string `yaml:"storageRoot"`
	Port        int    `yaml:"port"`
}

func LoadParameters() Parameters {
	parametersFile := os.Getenv("config")
	if parametersFile == "" {
		logger.Fatalln("Config file path missing")
	}

	yamlFile, err := ioutil.ReadFile(parametersFile)
	if err != nil {
		logger.Fatalf("configFile %s .Get err #%v", parametersFile, err)
	}

	parameters := Parameters{}
	err = yaml.Unmarshal(yamlFile, &parameters)
	if err != nil {
		logger.Fatalf("Unmarshal: %v", err)
	}

	return parameters
}
