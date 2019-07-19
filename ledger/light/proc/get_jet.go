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

package proc

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/pkg/errors"
)

type GetJet struct {
	message  payload.Meta
	objectID insolar.ID
	pulse    insolar.PulseNumber

	Dep struct {
		Jets   jet.Storage
		Sender bus.Sender
	}
}

func NewGetJet(msg payload.Meta, objectID insolar.ID, pulse insolar.PulseNumber) *GetJet {
	return &GetJet{
		message:  msg,
		objectID: objectID,
		pulse:    pulse,
	}
}

func (p *GetJet) Proceed(ctx context.Context) error {
	jetID, actual := p.Dep.Jets.ForID(ctx, p.pulse, p.objectID)
	msg := bus.ReplyAsMessage(ctx, &reply.Jet{ID: insolar.ID(jetID), Actual: actual})

	msg, err := payload.NewMessage(&payload.Jet{
		JetID:  jetID,
		Actual: actual,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	p.Dep.Sender.Reply(ctx, p.message, msg)
	return nil
}
