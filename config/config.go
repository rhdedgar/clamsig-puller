/*
Copyright 2020 Doug Edgar.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/rhdedgar/clamsig-puller/models"
)

var (
	// ClamInstallDir is the directory to which we have installed clam,
	// and is the target dir for our downloads.
	ClamInstallDir = "/clam/"
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
