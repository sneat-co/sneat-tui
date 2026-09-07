package commands

import (
	"testing"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-cli/internal/actionstore"
)

func TestActionValidate_ClientHeld_RevalidatesAndSaves(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	_ = store.Save(actionstore.Draft{ActionID: "act_abc12345", SpaceID: "famID", Semantic: actionspec.Semantic{Kind: actionspec.KindBuy}})
	api := &fakeActionsAPI{validateResp: dto4sneatai.ActionResponse{
		Semantic: actionspec.Semantic{Kind: actionspec.KindBuy}, Validation: canCommitValidation(),
	}}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "validate", "act_abc12345"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.validateReq == nil || api.validateReq.Semantic == nil {
		t.Fatal("ActionValidate not called with an inline semantic")
	}
	got, err := store.Load("act_abc12345")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !got.Validation.CanCommit {
		t.Fatalf("draft validation not updated: %+v", got.Validation)
	}
}

func TestActionValidate_ServerHeld_ValidatesByID(t *testing.T) {
	dir := t.TempDir()
	store := actionstoreFor(t, dir)
	_ = store.Save(actionstore.Draft{ActionID: "act_srv000001", SpaceID: "famID", ServerHeld: true})
	api := &fakeActionsAPI{validateResp: dto4sneatai.ActionResponse{ActionID: "act_srv000001"}}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "validate", "act_srv000001"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.validateReq == nil || api.validateReq.Semantic != nil {
		t.Fatalf("server-held validate should be by id only: %+v", api.validateReq)
	}
	if api.validateReq.ActionID != "act_srv000001" {
		t.Fatalf("actionID = %q", api.validateReq.ActionID)
	}
}
