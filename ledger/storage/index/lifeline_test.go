/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package index

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/stretchr/testify/assert"
)

func TestCloneObjectLifeline(t *testing.T) {
	t.Parallel()

	currentIdx := objectLifeline()

	clonedIdx := CloneObjectLifeline(currentIdx)

	assert.Equal(t, currentIdx, clonedIdx)
	assert.False(t, &currentIdx == &clonedIdx)
}

func id() (id *core.RecordID) {
	fuzz.New().NilChance(0.5).Fuzz(&id)
	return
}

func delegates() (result map[core.RecordRef]core.RecordRef) {
	fuzz.New().NilChance(0.5).NumElements(1, 10).Fuzz(&result)
	return
}

func state() (state record.State) {
	fuzz.New().NilChance(0).Fuzz(&state)
	return
}

func objectLifeline() ObjectLifeline {
	var index ObjectLifeline
	fuzz.New().NilChance(0).Funcs(
		func(idx *ObjectLifeline, c fuzz.Continue) {
			idx.LatestState = id()
			idx.LatestStateApproved = id()
			idx.ChildPointer = id()
			idx.Delegates = delegates()
			idx.State = state()
			idx.Parent = gen.Reference()
			idx.LatestUpdate = gen.PulseNumber()
			idx.JetID = gen.JetID()
		},
	).Fuzz(&index)

	return index
}
