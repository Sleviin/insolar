// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package goplugintestutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	recursive "github.com/insolar/insolar/application/proxy/recursive_call_one"
)
type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

var INSATTR_Recursive_API = true
func (r *One) Recursive() (error) {
	remoteSelf := recursive.GetObject(r.GetReference())
	err := remoteSelf.Recursive()
	return err
}

`

func TestContractsBuilder_Build(t *testing.T) {
	insgocc, err := BuildPreprocessor()
	assert.NoError(t, err)

	am := artifacts.NewClientMock(t)
	am.RegisterIncomingRequestMock.Set(func(ctx context.Context, request *record.IncomingRequest) (rp1 *payload.RequestInfo, err error) {
		rp1 = &payload.RequestInfo{RequestID: gen.ID()}
		return
	})
	am.DeployCodeMock.Set(func(ctx context.Context, code []byte, machineType insolar.MachineType) (ip1 *insolar.ID, err error) {
		assert.Equal(t, insolar.MachineTypeGoPlugin, machineType)
		id := gen.ID()
		return &id, nil
	})
	am.ActivatePrototypeMock.Set(func(ctx context.Context, request insolar.Reference, parent insolar.Reference, code insolar.Reference, memory []byte) (err error) {
		return nil
	})

	pa := pulse.NewAccessorMock(t)
	pa.LatestMock.Set(func(ctx context.Context) (p1 insolar.Pulse, err error) {
		return *insolar.GenesisPulse, nil
	})

	j := jet.NewCoordinatorMock(t)
	j.MeMock.Set(func() (r1 insolar.Reference) {
		return gen.Reference()
	})

	cb := NewContractBuilder(insgocc, am, pa, j)
	defer cb.Clean()

	contractMap := make(map[string]string)
	contractMap["recursive_call_one"] = contractOneCode

	buildOptions := BuildOptions{PanicIsLogicalError: false}
	err = cb.Build(context.Background(), contractMap, buildOptions)
	assert.NoError(t, err)

	reference := cb.Prototypes["recursive_call_one"]
	PrototypeRef := reference.String()
	assert.NotEmpty(t, PrototypeRef)
}
