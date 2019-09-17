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

package main

import (
	"github.com/insolar/insolar/api/sdk"
)

type createMembersScenario struct {
	insSDK  *sdk.SDK
	members []*sdk.Member
}

func (s *createMembersScenario) canBeStarted() error {
	return nil
}

func (s *createMembersScenario) prepare() {}

func (s *createMembersScenario) scenario(index int) (string, error) {
	creator := s.members[index]

	_, traceID, err := s.insSDK.CreateMember(creator.Reference)
	return traceID, err
}

func (s *createMembersScenario) checkResult() {}
