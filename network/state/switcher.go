/*
 *    Copyright 2018 INS Ecosystem
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

package state

import (
	"context"

	"github.com/insolar/insolar/core"
)

// NetworkSwitcher is a network FSM using for bootstrapping
type NetworkSwitcher struct {
	state core.NetworkState
}

// NewNetworkSwitcher creates new NetworkSwitcher
func NewNetworkSwitcher() (*NetworkSwitcher, error) {
	return &NetworkSwitcher{state: core.NoNetworkState}, nil
}

// GetState method returns current network state
func (ns *NetworkSwitcher) GetState() core.NetworkState {
	return core.CompleteNetworkState
}

// OnPulse method checks current state and finds out reasons to update this state
func OnPulse(ctx context.Context, pulse core.Pulse) error {
	// TODO: check discovery nodes is equal to ActiveList
	return nil
}
