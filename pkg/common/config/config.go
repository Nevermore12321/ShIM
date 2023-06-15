package config

import (
	"ShIM/pkg/common/config/confs"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var Root = "/etc/shim"
var Config confs.Config

func unmarshalYamlConfig(config interface{}, configName string) {
	// config path for environment variable
	// default config path is /etc/shim/
	var envName string
	if configName == "shim.yaml" {
		envName = "OS_CONFIG_DIR"
	} else if configName == "user.yaml" {
		envName = "OS_USER_CONFIG_DIR"
	}

	cfgPath := os.Getenv(envName)

	// Configured configuration file path in env
	var bytes []byte = make([]byte, 10)
	if len(cfgPath) != 0 {
		cfgDir := filepath.Join(cfgPath, configName)

		bytes, err := os.ReadFile(cfgDir)
		if err != nil {
			panic("Read config file " + cfgDir + "failed: " + err.Error())
		}
		// marshal yaml file to config object
		if err := yaml.Unmarshal(bytes, config); err != nil {
			panic("Unmarshal yaml file to object failed: " + err.Error())
		}

	} else { // not configured in env，default /etc/shim/
		cfgDir := filepath.Join(Root, configName)
		bytes, err := os.ReadFile(cfgDir)
		if err != nil {
			panic("Read config file " + cfgDir + "Failed: " + err.Error())
		}
		if err := yaml.Unmarshal(bytes, config); err != nil {
			panic("Unmarshal yaml file to object failed: " + err.Error())
		}
	}

}
