///
// Copyright 2020 Insolar Technologies GmbH
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

// +build functest

package functest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
)

func TestActivateDaemonDoubleCall(t *testing.T) {
	activeDaemons := activateDaemons(t, countTwoActiveDaemon)
	for _, daemon := range activeDaemons {
		_, _, err := testrequest.MakeSignedRequest(launchnet.TestRPCUrl, &MigrationAdmin, "migration.activateDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		data := checkConvertRequesterError(t, err).Data
		require.Contains(t, data.Trace, "daemon member already activated")
	}
}

func TestActivateDeactivateDaemon(t *testing.T) {
	activeDaemons := activateDaemons(t, countTwoActiveDaemon)
	for _, daemon := range activeDaemons {
		_, err := testrequest.SignedRequest(t, launchnet.TestRPCUrl, &MigrationAdmin, "migration.deactivateDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		require.NoError(t, err)
	}

	for _, daemon := range activeDaemons {
		res, _, err := testrequest.MakeSignedRequest(launchnet.TestRPCUrl, &MigrationAdmin, "migration.checkDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		require.NoError(t, err)
		status := res.(map[string]interface{})["status"].(string)
		require.Equal(t, status, "inactive")
	}

	for _, daemon := range activeDaemons {
		_, err := testrequest.SignedRequest(t, launchnet.TestRPCUrl, &MigrationAdmin, "migration.activateDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		require.NoError(t, err)
	}

	for _, daemon := range activeDaemons {
		res, _, err := testrequest.MakeSignedRequest(launchnet.TestRPCUrl, &MigrationAdmin, "migration.checkDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		require.NoError(t, err)
		status := res.(map[string]interface{})["status"].(string)
		require.Equal(t, status, "active")
	}
}
func TestDeactivateDaemonDoubleCall(t *testing.T) {
	activeDaemons := activateDaemons(t, countTwoActiveDaemon)
	for _, daemon := range activeDaemons {
		_, _, err := testrequest.MakeSignedRequest(launchnet.TestRPCUrl, &MigrationAdmin, "migration.deactivateDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		require.NoError(t, err)
	}
	for _, daemon := range activeDaemons {
		_, _, err := testrequest.MakeSignedRequest(launchnet.TestRPCUrl, &MigrationAdmin, "migration.deactivateDaemon",
			map[string]interface{}{"reference": daemon.Ref})
		data := checkConvertRequesterError(t, err).Data
		require.Contains(t, data.Trace, "daemon member already deactivated")
	}
}
func TestActivateAccess(t *testing.T) {

	member := createMigrationMemberForMA(t)
	_, _, err := testrequest.MakeSignedRequest(launchnet.TestRPCUrl, member, "migration.activateDaemon",
		map[string]interface{}{"reference": MigrationDaemons[0].Ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "only migration admin can activate migration demons")
}

func TestDeactivateAccess(t *testing.T) {

	member := createMigrationMemberForMA(t)
	_, _, err := testrequest.MakeSignedRequest(launchnet.TestRPCUrl, member, "migration.deactivateDaemon",
		map[string]interface{}{"reference": MigrationDaemons[0].Ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "only migration admin can deactivate migration demons")
}

func TestCheckDaemonAccess(t *testing.T) {

	member := createMigrationMemberForMA(t)
	_, _, err := testrequest.MakeSignedRequest(launchnet.TestRPCUrl, member, "migration.checkDaemon",
		map[string]interface{}{"reference": MigrationDaemons[0].Ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "permission denied to information about migration daemons")
}
