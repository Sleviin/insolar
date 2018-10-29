/*
 *
 *  *    Copyright 2018 Insolar
 *  *
 *  *    Licensed under the Apache License, Version 2.0 (the "License");
 *  *    you may not use this file except in compliance with the License.
 *  *    You may obtain a copy of the License at
 *  *
 *  *        http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  *    Unless required by applicable law or agreed to in writing, software
 *  *    distributed under the License is distributed on an "AS IS" BASIS,
 *  *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  *    See the License for the specific language governing permissions and
 *  *    limitations under the License.
 *
 */

package logicrunner

import (
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// RegisterExecutor registers an executor for particular `MachineType`
func (lr *LogicRunner) RegisterExecutor(t core.MachineType, e core.MachineLogicExecutor) error {
	lr.Executors[int(t)] = e
	return nil
}

// GetExecutor returns an executor for the `MachineType` if it was registered (`RegisterExecutor`),
// returns error otherwise
func (lr *LogicRunner) GetExecutor(t core.MachineType) (core.MachineLogicExecutor, error) {
	if res := lr.Executors[int(t)]; res != nil {
		return res, nil
	}

	return nil, errors.Errorf("No executor registered for machine %d", int(t))
}

func (lr *LogicRunner) GetExecution(ref Ref) *ExecutionState {
	lr.executionMutex.Lock()
	defer lr.executionMutex.Unlock()
	res, ok := lr.execution[ref]
	if !ok {
		return nil
	}
	return res
}

func (lr *LogicRunner) UpsertExecution(ref Ref) *ExecutionState {
	lr.executionMutex.Lock()
	defer lr.executionMutex.Unlock()
	if _, ok := lr.execution[ref]; !ok {
		lr.execution[ref] = &ExecutionState{}
	}
	return lr.execution[ref]
}

// refreshCaseBind lock CaseBind data, copy it, clean original, unlock original, return copy
func (lr *LogicRunner) refreshCaseBind() map[Ref][]core.CaseRecord {
	lr.caseBindMutex.Lock()
	defer lr.caseBindMutex.Unlock()

	oldObjectsRecords := lr.caseBind.Records

	lr.caseBind = core.CaseBind{Records: make(map[Ref][]core.CaseRecord)}

	return oldObjectsRecords
}

func (lr *LogicRunner) addObjectCaseRecord(ref Ref, cr core.CaseRecord) {
	lr.caseBindMutex.Lock()
	lr.caseBind.Records[ref] = append(lr.caseBind.Records[ref], cr)
	lr.caseBindMutex.Unlock()
}

func (lr *LogicRunner) nextValidationStep(ref Ref) (*core.CaseRecord, int) {
	lr.caseBindReplaysMutex.Lock()
	defer lr.caseBindReplaysMutex.Unlock()
	r, ok := lr.caseBindReplays[ref]
	if !ok {
		return nil, -1
	} else if r.RecordsLen <= r.Step {
		return nil, r.Step
	}
	ret := r.Records[r.Step]
	r.Step++
	lr.caseBindReplays[ref] = r
	return &ret, r.Step
}

func (lr *LogicRunner) pulse() *core.Pulse {
	pulse, err := lr.Ledger.GetPulseManager().Current()
	if err != nil {
		panic(err)
	}
	return pulse
}

func (lr *LogicRunner) GetConsensus(r Ref) (*Consensus, bool) {
	lr.consensusMutex.Lock()
	defer lr.consensusMutex.Unlock()
	c, ok := lr.consensus[r]
	if !ok {
		// arr, err := lr.Ledger.GetJetCoordinator().QueryRole(core.RoleVirtualValidator, r, lr.Pulse.PulseNumber)
		//if err != nil {
		//	panic("cannot QueryRole")
		//}
		c = newConsensus(lr, nil)
		lr.consensus[r] = c
	}
	return c, ok
}
