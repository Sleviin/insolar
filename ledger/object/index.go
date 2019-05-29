//
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
//

package object

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineIndex -o ./mocks -s _mock.go

// LifelineIndex is a base storage for lifelines.
type LifelineIndex interface {
	// LifelineAccessor provides methods for fetching lifelines.
	LifelineAccessor
	// LifelineModifier provides methods for modifying lifelines.
	LifelineModifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineAccessor -o ./mocks -s _mock.go

// LifelineAccessor provides methods for fetching lifelines.
type LifelineAccessor interface {
	// ForID returns a lifeline from a bucket with provided PN and ObjID
	ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (Lifeline, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineModifier -o ./mocks -s _mock.go

// LifelineModifier provides methods for modifying lifelines.
type LifelineModifier interface {
	// Set set a lifeline to a bucket with provided pulseNumber and ID
	Set(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineStateModifier -o ./mocks -s _mock.go

// LifelineStateModifier provides an interface for changing a state of lifeline.
type LifelineStateModifier interface {
	// SetLifelineUsage updates a last usage fields of a bucket for a provided pulseNumber and an object id
	SetLifelineUsage(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexCleaner -o ./mocks -s _mock.go

// IndexCleaner provides an interface for removing backets from a storage.
type IndexCleaner interface {
	// DeleteForPN method removes indexes from a storage for a provided
	DeleteForPN(ctx context.Context, pn insolar.PulseNumber)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexBucketModifier -o ./mocks -s _mock.go

// IndexBucketModifier provides methods for modifying buckets of index.
// Index contains buckets with pn->objID->Bucket hierarchy.
// With using of IndexBucketModifier there is a possibility to set buckets from outside of an index.
type IndexBucketModifier interface {
	// SetBucket adds a bucket with provided pulseNumber and ID
	SetBucket(ctx context.Context, pn insolar.PulseNumber, bucket IndexBucket) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexBucketAccessor -o ./mocks -s _mock.go

// IndexBucketAccessor provides an interface for fetching buckets from an index.
type IndexBucketAccessor interface {
	// ForPNAndJet returns a collection of buckets for a provided pn and jetID
	ForPNAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) []IndexBucket
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.PendingModifier -o ./mocks -s _mock.go

// PendingModifier provides methods for modifying pending requests.
type PendingModifier interface {
	SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, req record.Request) error
	SetResult(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, req record.Result) error
	SetFilament(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, filPN insolar.PulseNumber, recs []record.Virtual) error
	RefreshState(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error
	SetReadUntil(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, readUntil *insolar.PulseNumber) error
}

type PendingMeta struct {
	PreviousPN        *insolar.PulseNumber
	ReadUntil         *insolar.PulseNumber
	IsStateCalculated bool
}

type PendingAccessor interface {
	MetaForObjID(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) (PendingMeta, error)
	RequestsForObjID(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID, count int) ([]record.Request, error)
	Records(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) ([]record.Virtual, error)
}
