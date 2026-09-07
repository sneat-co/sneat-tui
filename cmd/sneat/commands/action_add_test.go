package commands

import (
	"bytes"
	"testing"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-cli/internal/actionstore"
)

func TestActionAdd_ClientHeld_MergesLocallyAndRevalidates(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	base := actionspec.Semantic{Kind: actionspec.KindBuy, Operations: []actionspec.Operation{
		{Objects: []actionspec.Text{{Text: "shoes"}}},
	}}
	if err := store.Save(actionstore.Draft{ActionID: "act_abc12345", SpaceID: "famID", Semantic: base}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	merged := actionspec.Semantic{Kind: actionspec.KindBuy, Operations: []actionspec.Operation{
		{Objects: []actionspec.Text{{Text: "shoes"}}, Contacts: []actionspec.ContactRef{{ContactID: "c1", Title: "Alice"}}},
	}}
	api := &fakeActionsAPI{validateResp: dto4sneatai.ActionResponse{Semantic: merged, Validation: canCommitValidation()}}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"action", "add", "act_abc12345", "--json", `{"operations":[{"contacts":[{"mention":"Alice"}]}]}`})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.validateReq == nil {
		t.Fatal("ActionValidate not called")
	}
	if api.validateReq.Semantic == nil || len(api.validateReq.Semantic.Operations) != 1 ||
		api.validateReq.Semantic.Operations[0].Objects[0].Text != "shoes" ||
		len(api.validateReq.Semantic.Operations[0].Contacts) != 1 ||
		api.validateReq.Semantic.Operations[0].Contacts[0].Mention != "Alice" {
		t.Fatalf("merged semantic sent to validate = %+v", api.validateReq.Semantic)
	}
	got, err := store.Load("act_abc12345")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got.Semantic.Operations[0].Contacts) != 1 || got.Semantic.Operations[0].Contacts[0].ContactID != "c1" {
		t.Fatalf("draft not updated with resolved semantic: %+v", got.Semantic)
	}
	if !got.Validation.CanCommit {
		t.Fatalf("draft validation not updated: %+v", got.Validation)
	}
}

func TestActionAdd_ServerHeld_Patches(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	if err := store.Save(actionstore.Draft{ActionID: "act_srv000001", SpaceID: "famID", ServerHeld: true}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	api := &fakeActionsAPI{patchResp: dto4sneatai.ActionResponse{ActionID: "act_srv000001"}}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "add", "act_srv000001", "--json", `{"kind":"buy"}`})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.patchReq == nil {
		t.Fatal("ActionPatch not called")
	}
	if api.validateReq != nil {
		t.Fatal("ActionValidate should not be called for a server-held action")
	}
	if api.patchReq.ActionID != "act_srv000001" || string(api.patchReq.SpaceID) != "famID" {
		t.Fatalf("patch req = %+v", api.patchReq)
	}
}

func TestActionAdd_UnknownID_Errors(t *testing.T) {
	env := actionsEnv(&fakeActionsAPI{}, t.TempDir())
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "add", "act_missing1", "--json", `{"kind":"buy"}`})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for unknown action id")
	}
}

func TestActionAdd_MissingJSON_Errors(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	_ = store.Save(actionstore.Draft{ActionID: "act_abc12345", SpaceID: "famID"})
	env := actionsEnv(&fakeActionsAPI{}, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "add", "act_abc12345"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error: --json required")
	}
}
