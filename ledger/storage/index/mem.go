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
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/db"
)

// StorageMem is an in-memory struct for index-storage
type StorageMem struct {
	jetIndex db.JetIndexModifier

	lock   sync.RWMutex
	memory map[core.RecordID]ObjectLifeline
}

// NewStorageMem creates a new instance of Storage.
func NewStorageMem() *StorageMem {
	return &StorageMem{
		memory:   map[core.RecordID]ObjectLifeline{},
		jetIndex: db.NewJetIndex(),
	}
}

func (s *StorageMem) Set(ctx context.Context, id core.RecordID, index ObjectLifeline) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.memory[id]
	if ok {
		return ErrOverride
	}

	s.memory[id] = index
	s.jetIndex.Add(id, index.JetID)

	return nil
}

func (s *StorageMem) ForID(ctx context.Context, id core.RecordID) (idx ObjectLifeline, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	idx, ok := s.memory[id]
	if !ok {
		err = ErrNotFound
		return
	}

	return
}
