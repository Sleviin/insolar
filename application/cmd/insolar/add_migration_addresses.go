///
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
///

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/insolar/insolar/application/api/sdk"
)

func addMigrationAddresses(adminUrls []string, publicUrls []string, memberKeysDirPath string, addressesPath string, shardsCount int) {
	insSDK, err := sdk.NewSDK(adminUrls, publicUrls, memberKeysDirPath)
	check("SDK is not initialized: ", err)
	var filename string
	// method AddMigrationAddresses in contract use only 10 shards in one call
	for i := 0; i*10 < shardsCount; i++ {
		if shardsCount == 10 {
			filename = addressesPath
		} else {
			addrPath := strings.Split(addressesPath, ".json")
			// filename match with files, generated by generator utility (from insolar/migrationAddressGenerator)
			filename = addrPath[0] + "_" + strconv.Itoa(i) + ".json"
		}
		rawConf, err := ioutil.ReadFile(filepath.Clean(filename))
		check("Error while reading file: ", err)
		var addresses []string
		err = json.Unmarshal(rawConf, &addresses)
		check(fmt.Sprintf("Error while unmarshal content of file %s to list of addresses: ", filename), err)
		_, err = insSDK.AddMigrationAddresses(addresses)
		check(fmt.Sprintf("Error while adding addresses from file %s: ", filename), err)
		fmt.Printf("Addresses to shards with indexes %d-%d were added successfully\n", i*10, i*10+9)
	}
	fmt.Println("All addresses were added successfully")
}
