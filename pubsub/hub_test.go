package pubsub

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/leg100/otf"
	"github.com/leg100/otf/configversion"
	"github.com/leg100/otf/organization"
	"github.com/leg100/otf/sql"
	"github.com/leg100/otf/workspace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPubSub_E2E tests that one pubsub process can publish a message and that
// another pubsub process can receive it.
func TestPubSub_E2E(t *testing.T) {
	db := sql.NewTestDB(t)
	org := organization.CreateTestOrganization(t, db)
	ws := workspace.CreateTestWorkspace(t, db, org.Name())
	cv := configversion.CreateTestConfigurationVersion(t, db, ws, otf.ConfigurationVersionCreateOptions{})
	run := run.CreateTestRun(t, db, ws, cv)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() { cancel() })

	// setup sender
	sender, err := newSpoke(logr.Discard(), db, ChannelName("events_e2e_test"), PID("sender-1"))
	require.NoError(t, err)
	senderGot, err := sender.Subscribe(ctx, "sender-1")
	require.NoError(t, err)
	// setup receiver
	receiver, err := newSpoke(logr.Discard(), db, ChannelName("events_e2e_test"), PID("receiver-1"))
	require.NoError(t, err)
	receiverGot, err := receiver.Subscribe(ctx, "sender-2")
	require.NoError(t, err)

	go func() { sender.Start(ctx) }()
	go func() { receiver.Start(ctx) }()

	// Give Start() time to connect and start listening
	time.Sleep(100 * time.Millisecond)

	// this is the event we're publishing from the sender and expecting to make
	// its way to postgres and then back to the receiver.
	want := otf.Event{
		Type:    otf.EventRunStatusUpdate,
		Payload: run,
	}
	sender.Publish(want)

	// Give time for message to make its way via postgres and back.
	time.Sleep(time.Second)

	// We expect the receiver process to have received a copy
	assert.Equal(t, 1, len(receiverGot))
	assert.Equal(t, want, <-receiverGot)

	// We also expect the sender process to have published a copy locally for local
	// subs
	assert.Equal(t, 1, len(senderGot))
	assert.Equal(t, want, <-senderGot)
}

func TestPubSub_marshal(t *testing.T) {
	ps := &spoke{pid: "process-1"}
	event := otf.Event{
		Type: otf.EventRunStatusUpdate,
		Payload: otf.NewTestRun(t, otf.TestRunCreateOptions{
			ID: otf.String("run-123"),
		}),
	}

	got, err := ps.marshal(event)
	require.NoError(t, err)
	want := "{\"relation\":\"run\",\"action\":\"status_update\",\"id\":\"run-123\",\"pid\":\"process-1\"}"
	assert.Equal(t, want, string(got))
}

func TestPubSub_reassemble(t *testing.T) {
	run := otf.NewTestRun(t, otf.TestRunCreateOptions{
		ID: otf.String("run-123"),
	})
	ps := spoke{
		db: &fakePubSubDB{
			run: run,
		},
	}

	got, err := ps.reassemble(context.Background(), message{
		Table:  "run",
		Action: "status_update",
		ID:     "run-123",
	})
	require.NoError(t, err)
	want := otf.Event{
		Type:    otf.EventRunStatusUpdate,
		Payload: run,
	}
	assert.Equal(t, want, got)
}

type fakePubSubDB struct {
	run       *run.Run
	workspace *workspace.Workspace

	otf.DB
}

func (f *fakePubSubDB) GetRun(context.Context, string) (*run.Run, error) {
	return f.run, nil
}

func (f *fakePubSubDB) GetWorkspace(context.Context, string) (*workspace.Workspace, error) {
	return f.workspace, nil
}