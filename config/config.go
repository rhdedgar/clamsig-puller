package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/rhdedgar/clamd/models"
)

var (
	// ClamInstallDir is the directory to which we have installed clam,
	// and is the target dir for our downloads.
	ClamInstallDir = "/var/lib/clamav/"
	// ConfigPath is the path to the config file containing secrets needed by the application.
	ConfigPath = "/secrets/clam_update_config.json"

	// ConfigFile is a struct to contain the contents of ConfigPath.
	ConfigFile models.ConfigFile
)

func init() {
	fileBytes, err := ioutil.ReadFile(ConfigPath)

	if err != nil {
		fmt.Println("Error loading secrets json: ", err)
	}

	//fmt.Println("Config file contents: ", string(fileBytes))

	err = json.Unmarshal(fileBytes, &ConfigFile)
	if err != nil {
		fmt.Println("Error Unmarshalling secrets json: ", err)
	}
}
