package commands

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/sneat-co/sneat-ai-backend/dbo4sneatai"
	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-cli/internal/actionstore"
)

func TestActionCommit_ClientHeld_SendsCanonicalSemanticAndStoresResult(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	// A dirty candidate carrying a resolved Title, which Canonical strips.
	dirty := actionspec.Semantic{Kind: actionspec.KindBuy, Operations: []actionspec.Operation{
		{Objects: []actionspec.Text{{Text: "shoes"}}, Contacts: []actionspec.ContactRef{{ContactID: "c1", Title: "Alice"}}},
	}}
	if err := store.Save(actionstore.Draft{ActionID: "act_abc12345", SpaceID: "famID", Semantic: dirty}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	api := &fakeActionsAPI{commitResp: dto4sneatai.CommitActionResponse{
		ActionID: "act_abc12345", Result: dbo4sneatai.CommitResult{HappeningID: "hap1"},
	}}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"action", "commit", "act_abc12345"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.commitReq == nil {
		t.Fatal("ActionCommit not called")
	}
	if api.commitReq.Semantic == nil {
		t.Fatal("commit request missing semantic")
	}
	if api.commitReq.Semantic.Operations[0].Contacts[0].Title != "" {
		t.Fatalf("semantic not canonicalized: %+v", api.commitReq.Semantic)
	}
	got, err := store.Load("act_abc12345")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Committed == nil || got.Committed.HappeningID != "hap1" {
		t.Fatalf("commit result not stored: %+v", got.Committed)
	}
}

func TestActionCommit_ServerHeld_SendsOnlyID(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	if err := store.Save(actionstore.Draft{ActionID: "act_srv000001", SpaceID: "famID", ServerHeld: true}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	api := &fakeActionsAPI{commitResp: dto4sneatai.CommitActionResponse{ActionID: "act_srv000001"}}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "commit", "act_srv000001"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.commitReq.Semantic != nil {
		t.Fatalf("server-held commit should not send a semantic: %+v", api.commitReq.Semantic)
	}
}

func TestActionCommit_Refused_ExitsCode2WithJSONError(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	if err := store.Save(actionstore.Draft{ActionID: "act_abc12345", SpaceID: "famID"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	refusal := errors.New("sneat api POST sneatai/action_commit: http 400: action cannot be committed: validation requires input")
	api := &fakeActionsAPI{commitErr: refusal}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"action", "commit", "act_abc12345"})
	err := root.Execute()
	var exitErr *ExitCodeError
	if !errors.As(err, &exitErr) || exitErr.Code != 2 {
		t.Fatalf("err = %v, want *ExitCodeError{Code: 2}", err)
	}
	if !strings.Contains(buf.String(), `"error"`) || !strings.Contains(buf.String(), "validation requires input") {
		t.Fatalf("stdout = %q", buf.String())
	}
}

func TestActionCommit_OtherError_NotExitCode2(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	_ = store.Save(actionstore.Draft{ActionID: "act_abc12345", SpaceID: "famID"})
	api := &fakeActionsAPI{commitErr: errors.New("sneat api POST sneatai/action_commit: http 500: boom")}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "commit", "act_abc12345"})
	err := root.Execute()
	var exitErr *ExitCodeError
	if errors.As(err, &exitErr) {
		t.Fatalf("unexpected ExitCodeError for a non-refusal error: %v", err)
	}
	if err == nil {
		t.Fatal("expected an error")
	}
}
