// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application/testutils/launchnet"
)

func TestGetFreeAddressCount(t *testing.T) {
	migrationShardsMap := getAddressCount(t, 0)

	for _, m := range migrationShardsMap {
		require.True(t, m > 0)
	}
}

func TestGetFreeAddressCount_WithIndex_NotAllRange(t *testing.T) {
	numLeftShards := 2
	numShards, err := launchnet.GetNumShards()
	require.NoError(t, err)
	var migrationShards = getAddressCount(t, numShards-numLeftShards)
	require.Len(t, migrationShards, numLeftShards)
}

func TestGetFreeAddressCount_StartIndexTooBig(t *testing.T) {
	numShards, err := launchnet.GetNumShards()
	require.NoError(t, err)
	_, _, err = makeSignedRequest(launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.getAddressCount",
		map[string]interface{}{"startWithIndex": numShards + 2})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "incorrect start shard index")
}

func TestGetFreeAddressCount_IncorrectIndexType(t *testing.T) {
	_, _, err := makeSignedRequest(launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.getAddressCount",
		map[string]interface{}{"startWithIndex": "0"})
	data := checkConvertRequesterError(t, err).Data
	expectedError(t, data.Trace, "doesn't match the schema")
}

func TestGetFreeAddressCount_FromMember(t *testing.T) {
	member := createMember(t)
	_, _, err := makeSignedRequest(launchnet.TestRPCUrl, member, "migration.getAddressCount",
		map[string]interface{}{"startWithIndex": 0})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "only migration daemon admin can call this method")
}
