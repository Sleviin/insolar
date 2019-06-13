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

package handle

import (
	"context"
	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
)

type GetCode struct {
	dep     *proc.Dependencies
	message payload.Meta
	passed  bool
}

func NewGetCode(dep *proc.Dependencies, msg payload.Meta, passed bool) *GetCode {
	return &GetCode{
		dep:     dep,
		message: msg,
		passed:  passed,
	}
}

func (s *GetCode) Present(ctx context.Context, f flow.Flow) error {
	msg := payload.GetCode{}
	err := msg.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetCode message")
	}

	ctx = inslogger.WithLoggerLevel(ctx,insolar.ErrorLevel)
	ctx, _ = inslogger.WithField(ctx, "code_id", msg.CodeID.DebugString())

	passIfNotFound := !s.passed
	code := proc.NewGetCode(s.message, msg.CodeID, passIfNotFound)
	s.dep.GetCode(code)
	return f.Procedure(ctx, code, false)
}
