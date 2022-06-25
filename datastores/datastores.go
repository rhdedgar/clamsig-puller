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

package datastores

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rhdedgar/clamsig-puller/models"
)

var (
	AppSecrets = models.AppSecrets{}
	// ClamInstallDir is the directory to which we have installed clam,
	// and is the target dir for our downloads.
	//ClamInstallDir = "/clam/"
	// ConfigPath is the path to the config file containing secrets needed by the application.
	//ConfigPath = "/secrets/clam_update_config.json"
)

func loadConfigFile(filePath string, dest interface{}) error {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error loading secrets json from:  %v %v\n", filePath, err)
	}

	err = json.Unmarshal(fileBytes, dest)
	if err != nil {
		return fmt.Errorf("Error Unmarshaling secrets json: %v\n", err)
	}
	return nil
}

func init() {
	filePath := os.Getenv("CLAM_SECRETS_FILE")

	err := loadConfigFile(filePath, &AppSecrets)
	if err != nil {
		fmt.Println("Error reading file: ", err)
	}

	if AppSecrets.ClamConfigDir == "" {
		AppSecrets.ClamConfigDir = os.Getenv("CLAM_DB_DIRECTORY")

		if AppSecrets.ClamConfigDir == "" {
			AppSecrets.ClamConfigDir = "/var/lib/clamav/"
		}
	}

	// be tolerant of an env var path not already suffixed with a trailing slash
	if !strings.HasSuffix(AppSecrets.ClamConfigDir, "/") {
		AppSecrets.ClamConfigDir = AppSecrets.ClamConfigDir + "/"
	}

	for _, item := range AppSecrets.ClamConfigFiles {
		AppSecrets.ClamConfigFileMap[item] = struct{}{}
	}
}
