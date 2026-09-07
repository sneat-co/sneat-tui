package commands

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-cli/internal/actionstore"
)

func TestActionShow_ClientHeld_PrintsDraft(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	_ = store.Save(actionstore.Draft{ActionID: "act_abc12345", SpaceID: "famID", Semantic: actionspec.Semantic{Kind: actionspec.KindBuy}})
	api := &fakeActionsAPI{}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"action", "show", "act_abc12345"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.getActionID != "" {
		t.Fatal("ActionGet should not be called for a client-held draft")
	}
	var got actionstore.Draft
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, buf.String())
	}
	if got.ActionID != "act_abc12345" {
		t.Fatalf("draft = %+v", got)
	}
}

func TestActionShow_ServerHeld_Fetches(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	_ = store.Save(actionstore.Draft{ActionID: "act_srv000001", SpaceID: "famID", ServerHeld: true})
	api := &fakeActionsAPI{getResp: dto4sneatai.ActionResponse{ActionID: "act_srv000001"}}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "show", "act_srv000001"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.getActionID != "act_srv000001" || api.getSpaceID != "famID" {
		t.Fatalf("ActionGet(%q, %q)", api.getSpaceID, api.getActionID)
	}
}

func TestActionShow_UnknownID_Errors(t *testing.T) {
	env := actionsEnv(&fakeActionsAPI{}, t.TempDir())
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "show", "act_missing1"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestActionList_ListsLocalDrafts(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	_ = store.Save(actionstore.Draft{ActionID: "act_aaa00001", SpaceID: "famID", Semantic: actionspec.Semantic{Kind: actionspec.KindBuy}})
	_ = store.Save(actionstore.Draft{ActionID: "act_bbb00002", SpaceID: "famID", ServerHeld: true})
	env := actionsEnv(&fakeActionsAPI{}, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"action", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var got []actionstore.Draft
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, buf.String())
	}
	if len(got) != 2 {
		t.Fatalf("got %d drafts, want 2: %+v", len(got), got)
	}
}

func TestActionList_EmptyByDefault(t *testing.T) {
	env := actionsEnv(&fakeActionsAPI{}, t.TempDir())
	root := Root(env)
	root.AddCommand(Action(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"action", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var got []actionstore.Draft
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, buf.String())
	}
	if len(got) != 0 {
		t.Fatalf("got = %+v, want empty", got)
	}
}

func TestActionCancel_ClientHeld_DeletesDraft(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	_ = store.Save(actionstore.Draft{ActionID: "act_abc12345", SpaceID: "famID"})
	env := actionsEnv(&fakeActionsAPI{}, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "cancel", "act_abc12345"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if _, err := store.Load("act_abc12345"); err == nil {
		t.Fatal("draft should have been deleted")
	}
}

func TestActionCancel_ServerHeld_CallsAPIAndDeletesIndex(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	_ = store.Save(actionstore.Draft{ActionID: "act_srv000001", SpaceID: "famID", ServerHeld: true})
	api := &fakeActionsAPI{}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "cancel", "act_srv000001"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.cancelReq == nil || api.cancelReq.ActionID != "act_srv000001" {
		t.Fatalf("ActionCancel req = %+v", api.cancelReq)
	}
	if _, err := store.Load("act_srv000001"); err == nil {
		t.Fatal("index entry should have been deleted")
	}
}
