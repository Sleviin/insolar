// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/application/testutils/launchnet"
)

func TestGetRequest(t *testing.T) {
	postResp, err := http.Get(launchnet.TestRPCUrl)
	require.NoError(t, err)
	defer postResp.Body.Close()
	require.Equal(t, http.StatusMethodNotAllowed, postResp.StatusCode)
}

func TestWrongUrl(t *testing.T) {
	jsonValue, _ := json.Marshal(postParams{})
	testURL := launchnet.AdminHostPort + "/not_api"
	postResp, err := http.Post(testURL, "application/json", bytes.NewBuffer(jsonValue))
	defer postResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, postResp.StatusCode)
}

func TestWrongJson(t *testing.T) {
	postResp, err := http.Post(launchnet.TestRPCUrl, "application/json", bytes.NewBuffer([]byte("some not json value")))
	defer postResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)

	response := &requester.ContractResponse{}
	unmarshalCallResponse(t, body, response)
	require.NotNil(t, response.Error)

	require.Equal(t, "The JSON received is not a valid request payload.", response.Error.Message)
	require.Nil(t, response.Result)
}
