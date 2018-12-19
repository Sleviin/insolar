package logicrunner

import (
	"testing"

	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

func TestOnPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	mb := testutils.NewMessageBusMock(t)
	mb.SendMock.Return(&reply.ID{}, nil)

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	lr.MessageBus = mb

	// test empty lr
	pulse := core.Pulse{}

	err := lr.OnPulse(ctx, pulse)
	require.NoError(t, err)

	objectRef := testutils.RandomRef()

	// test empty es
	lr.state[objectRef] = &ObjectState{ExecutionState: &ExecutionState{Behaviour: &ValidationSaver{}}}
	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	require.Nil(t, lr.state[objectRef].ExecutionState)

	// test empty es with query in current
	lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
		},
	}
	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	require.Equal(t, InPending, lr.state[objectRef].ExecutionState.pending)

	// test empty es with query in current and query in queue - es.pending true, message.ExecutorResults.Pending = true, message.ExecutorResults.Queue one element
	result := make(chan ExecutionQueueResult, 1)

	// TODO maybe need do something more stable and easy to debug
	go func() {
		<-result
	}()

	qe := ExecutionQueueElement{
		result: result,
	}

	queue := append(make([]ExecutionQueueElement, 0), qe)

	lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
			Queue:     queue,
		},
	}

	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	require.Equal(t, InPending, lr.state[objectRef].ExecutionState.pending)
}

func TestPendingFinished(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	mb := testutils.NewMessageBusMock(t)
	pulse := core.Pulse{}
	objectRef := testutils.RandomRef()

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	lr.MessageBus = mb

	ps := testutils.NewPulseStorageMock(t)
	ps.CurrentMock.Return(&pulse, nil)
	lr.PulseStorage = ps

	es := &ExecutionState{
		Behaviour: &ValidationSaver{},
		Current:   &CurrentExecution{},
		pending:   NotPending,
	}

	// make sure that if there is no pending finishPendingIfNeeded returns false,
	// doesn't send PendingFinished message and doesn't change ExecutionState.pending
	require.False(t, lr.finishPendingIfNeeded(ctx, es, objectRef))
	require.Zero(t, mb.SendCounter)
	require.Equal(t, NotPending, es.pending)

	// make sure that in pending case finishPendingIfNeeded returns true
	// sends PendingFinished message and sets ExecutionState.pending back to NotPending
	es.pending = InPending
	mb.SendMock.Expect(ctx, &message.PendingFinished{Reference: objectRef}, pulse, nil).Return(&reply.ID{}, nil)
	require.True(t, lr.finishPendingIfNeeded(ctx, es, objectRef))
	require.Equal(t, NotPending, es.pending)
}

func TestStartQueueProcessorIfNeeded_DontStartQueueProcessorWhenPending(
	t *testing.T,
) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	am := testutils.NewArtifactManagerMock(t)
	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	lr.ArtifactManager = am

	objectRef := testutils.RandomRef()

	od := testutils.NewObjectDescriptorMock(t)
	od.HasPendingRequestsMock.Expect().Return(true)

	am.GetObjectMock.Return(od, nil)

	es := &ExecutionState{ArtifactManager: am, Queue: make([]ExecutionQueueElement, 0)}
	es.Queue = append(es.Queue, ExecutionQueueElement{})
	err := lr.StartQueueProcessorIfNeeded(
		ctx,
		es,
		&message.CallMethod{
			ObjectRef: objectRef,
			Method:    "some",
		},
	)
	require.NoError(t, err)
	require.Equal(t, InPending, es.pending)
}

func TestCheckPendingRequests(
	t *testing.T,
) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	objectRef := testutils.RandomRef()

	am := testutils.NewArtifactManagerMock(t)

	od := testutils.NewObjectDescriptorMock(t)
	am.GetObjectMock.Return(od, nil)

	es := &ExecutionState{ArtifactManager: am}
	pending, err := es.CheckPendingRequests(
		ctx, &message.CallConstructor{},
	)
	require.NoError(t, err)
	require.Equal(t, NotPending, pending)

	od.HasPendingRequestsMock.Expect().Return(false)
	am.GetObjectMock.Return(od, nil)
	es = &ExecutionState{ArtifactManager: am}
	pending, err = es.CheckPendingRequests(
		ctx, &message.CallMethod{
			ObjectRef: objectRef,
		},
	)
	require.NoError(t, err)
	require.Equal(t, NotPending, pending)

	od.HasPendingRequestsMock.Expect().Return(true)
	am.GetObjectMock.Return(od, nil)
	es = &ExecutionState{ArtifactManager: am}
	pending, err = es.CheckPendingRequests(
		ctx, &message.CallMethod{
			ObjectRef: objectRef,
		},
	)
	require.NoError(t, err)
	require.Equal(t, InPending, pending)

	am.GetObjectMock.Return(nil, errors.New("some"))
	es = &ExecutionState{ArtifactManager: am}
	pending, err = es.CheckPendingRequests(
		ctx, &message.CallMethod{
			ObjectRef: objectRef,
		},
	)
	require.Error(t, err)
	require.Equal(t, NotPending, pending)
}

func TestPrepareState(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})

	object := testutils.RandomRef()
	msg := &message.ExecutorResults{
		Caller:    testutils.RandomRef(),
		RecordRef: object,
	}

	// not pending
	// it's a first call, it's also initialize lr.state[object].ExecutionState
	// also check for empty Queue
	msg.Pending = false
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, NotPending, lr.state[object].ExecutionState.pending)
	require.Equal(t, 0, len(lr.state[object].ExecutionState.Queue))

	// pending without queue
	lr.state[object].ExecutionState.pending = PendingUnknown
	msg.Pending = true
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, InPending, lr.state[object].ExecutionState.pending)

	// do not change pending status if it isn't unknown
	lr.state[object].ExecutionState.pending = NotPending
	msg.Pending = true
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, NotPending, lr.state[object].ExecutionState.pending)

	// do not change pending status if it isn't unknown
	lr.state[object].ExecutionState.pending = InPending
	msg.Pending = false
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, InPending, lr.state[object].ExecutionState.pending)

	// brand new queue from message
	msg.Queue = []message.ExecutionQueueElement{message.ExecutionQueueElement{}}
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, 1, len(lr.state[object].ExecutionState.Queue))

	// add new element in existing queue
	queueElementRequest := testutils.RandomRef()
	msg.Queue = []message.ExecutionQueueElement{message.ExecutionQueueElement{Request: &queueElementRequest}}
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, 2, len(lr.state[object].ExecutionState.Queue))
	require.Equal(t, &queueElementRequest, lr.state[object].ExecutionState.Queue[0].request)

}
func TestHandlePendingFinishedMessage(
	t *testing.T,
) {
	ctx := inslogger.TestContext(t)

	objectRef := testutils.RandomRef()

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})

	parcel := testutils.NewParcelMock(t).MessageMock.Return(
		&message.PendingFinished{Reference: objectRef},
	)

	re, err := lr.HandlePendingFinishedMessage(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, &reply.OK{}, re)

	st := lr.MustObjectState(objectRef)

	es := st.ExecutionState
	require.NotNil(t, es)
	require.Equal(t, NotPending, es.pending)

	es.Current = &CurrentExecution{}
	re, err = lr.HandlePendingFinishedMessage(ctx, parcel)
	require.Error(t, err)

	es.Current = nil

	re, err = lr.HandlePendingFinishedMessage(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, &reply.OK{}, re)

}
