///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package smachine

import (
	"context"
	"time"
)

type InitFunc func(ctx InitializationContext) StateUpdate
type StateFunc func(ctx ExecutionContext) StateUpdate
type MigrateFunc func(ctx MigrationContext) StateUpdate
type CreateFunc func(ctx ConstructionContext) StateMachine
type AsyncResultFunc func(ctx AsyncResultContext)
type BroadcastReceiveFunc func(ctx AsyncResultContext, payload interface{}) bool
type ErrorHandlerFunc func(ctx FailureContext)

type BasicContext interface {
	SlotLink() SlotLink
	ParentLink() SlotLink

	GetContext() context.Context
	//GetContainer() SlotMachineState
}

/* During construction SlotLink() will have correct SlotID, but MAY have INVALID status, as slot was not yet created */
type ConstructionContext interface {
	BasicContext
	SetContext(context.Context)
	SetParent(SlotLink)
}

type stepContext interface {
	BasicContext

	SetDefaultMigration(fn MigrateFunc)
	SetDefaultErrorHandler(fn ErrorHandlerFunc)
	SetDefaultFlags(StepFlags)

	JumpExt(SlotStep) StateUpdate
	Jump(StateFunc) StateUpdate

	Share(data interface{}, wakeUpAfterUse bool) SharedDataLink

	Error(error) StateUpdate
	Stop() StateUpdate
}

type InitializationContext interface {
	stepContext

	BargeInWithParam(BargeInApplyFunc) BargeInParamFunc
	BargeIn() BargeInRequester
}

type BargeInApplyFunc func(BargeInContext) StateUpdate
type BargeInParamFunc func(interface{}) bool
type BargeInFunc func() bool

type BargeInRequester interface {
	WithJumpExt(SlotStep) BargeInFunc
	WithJump(StateFunc) BargeInFunc
	WithWakeUp() BargeInFunc
	WithStop() BargeInFunc
}

type MigrationContext interface {
	stepContext

	/* A step the target slot is at */
	AffectedStep() SlotStep

	Replace(CreateFunc) StateUpdate
	ReplaceWith(StateMachine) StateUpdate

	/* Keeps the last state */
	Stay() StateUpdate
	/* Makes active if was waiting or polling */
	WakeUp() StateUpdate
}

type ExecutionContext interface {
	stepContext

	StepLink() StepLink
	GetPendingCallCount() int

	//ListenBroadcast(key string, broadcastFn BroadcastReceiveFunc)
	//SyncOneStep(key string, weight int32, broadcastFn BroadcastReceiveFunc) Syncronizer
	//SyncManySteps(key string)

	// Allocates a new slot and schedules initialization.
	// It is guaranteed that:
	// 1) the child will start at the same migration state as the creator (caller of this function)
	// 1) initialization of the new slot will happen before any migration
	NewChild(context.Context, CreateFunc) SlotLink

	// Same as NewChild, but also grantees that child's initialization will be completed on return.
	// Please avoid using it unless necessary.
	InitChild(context.Context, CreateFunc) SlotLink

	//NewShared(context.Context, SharedState) SharedStateAdapter
	UseShared(SharedDataAccessor) SharedAccessReport // ?

	BargeInWithParam(BargeInApplyFunc) BargeInParamFunc
	BargeInThisStepOnly() BargeInRequester
	BargeIn() BargeInRequester

	Replace(CreateFunc) StateUpdate
	ReplaceWith(StateMachine) StateUpdate
	Repeat(limit int) StateUpdate

	Yield() StateConditionalUpdate
	Poll() StateConditionalUpdate

	WaitActivation(SlotLink) StateConditionalUpdate
	WaitShared(SharedDataLink) StateConditionalUpdate
	WaitAny() StateConditionalUpdate
	WaitAnyUntil(time.Time) StateConditionalUpdate

	Sleep() StateConditionalUpdate
}

type StateConditionalUpdate interface {
	ConditionalUpdate

	/* Returns true when condition is already true */
	IsAvailable() bool
}

type CallConditionalUpdate interface {
	ConditionalUpdate
}

type ConditionalUpdate interface {
	ThenJump(StateFunc) StateUpdate
	ThenJumpExt(SlotStep) StateUpdate
	ThenRepeat() StateUpdate
}

//type Syncronizer interface {
//	IsFirst() bool
//	Broadcast(payload interface{}) (total, accepted int)
//	ReleaseAll()
//
//	Wait() StateUpdate
//	WaitOrDeadline(d time.Time) StateUpdate
//}

type AsyncResultContext interface {
	BasicContext

	WakeUp()
}

type BargeInContext interface {
	BasicContext

	GetBargeInParam() interface{}
	IsAtOriginalStep() bool

	/* A step the target slot is at */
	AffectedStep() SlotStep

	JumpExt(SlotStep) StateUpdate
	Jump(StateFunc) StateUpdate

	/* Keeps the last state */
	Stay() StateUpdate
	/* Makes active if was waiting or polling */
	WakeUp() StateUpdate

	Stop() StateUpdate
}

type FailureContext interface {
	BasicContext

	/* A step the slot is at */
	AffectedStep() SlotStep

	GetError() error
	IsPanic() bool

	IgnoreDefaultHandler()

	// See ExecutionContext.NewChild
	NewChild(context.Context, CreateFunc) SlotLink

	// See ExecutionContext.InitChild
	InitChild(context.Context, CreateFunc) SlotLink
}