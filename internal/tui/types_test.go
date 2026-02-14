package tui

import (
	"testing"

	"github.com/mi8bi/ghqx/internal/status"
)

func TestViewState(t *testing.T) {
	states := []ViewState{
		ViewStateList,
		ViewStateLoading,
		ViewStateError,
		ViewStateConfirm,
	}

	for i, state := range states {
		if int(state) != i {
			t.Errorf("ViewState value mismatch: got %d, want %d", state, i)
		}
	}
}

func TestMessageType(t *testing.T) {
	types := []MessageType{
		MessageTypeInfo,
		MessageTypeSuccess,
		MessageTypeWarning,
		MessageTypeError,
	}

	for i, msgType := range types {
		if int(msgType) != i {
			t.Errorf("MessageType value mismatch: got %d, want %d", msgType, i)
		}
	}
}

func TestOperationType(t *testing.T) {
	operations := []OperationType{
		OperationRefresh,
		OperationQuit,
	}

	for i, op := range operations {
		if int(op) != i {
			t.Errorf("OperationType value mismatch: got %d, want %d", op, i)
		}
	}
}

func TestMessage(t *testing.T) {
	msg := Message{
		Text: "test message",
		Type: MessageTypeInfo,
		Hint: "test hint",
	}

	if msg.Text != "test message" {
		t.Error("Message Text mismatch")
	}

	if msg.Type != MessageTypeInfo {
		t.Error("Message Type mismatch")
	}

	if msg.Hint != "test hint" {
		t.Error("Message Hint mismatch")
	}
}

func TestNewProjectRow(t *testing.T) {
	pd := status.ProjectDisplay{
		Repo:       "user/repo",
		Workspace:  "sandbox",
		GitManaged: "Managed",
		Status:     "clean",
		FullPath:   "/path/to/repo",
	}

	row := NewProjectRow(pd)

	if row.Repo != pd.Repo {
		t.Error("ProjectRow Repo mismatch")
	}

	if row.Workspace != pd.Workspace {
		t.Error("ProjectRow Workspace mismatch")
	}

	if row.GitManaged != pd.GitManaged {
		t.Error("ProjectRow GitManaged mismatch")
	}

	if row.Status != pd.Status {
		t.Error("ProjectRow Status mismatch")
	}

	if row.FullPath != pd.FullPath {
		t.Error("ProjectRow FullPath mismatch")
	}
}

func TestProjectRowEmbedding(t *testing.T) {
	pd := status.ProjectDisplay{
		Repo: "test/repo",
	}

	row := NewProjectRow(pd)

	// Test that embedded fields are accessible
	if row.ProjectDisplay.Repo != "test/repo" {
		t.Error("embedded ProjectDisplay not accessible")
	}

	// Test that we can access embedded fields directly
	if row.Repo != "test/repo" {
		t.Error("direct field access failed")
	}
}
