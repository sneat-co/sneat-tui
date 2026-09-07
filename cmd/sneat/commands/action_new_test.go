package commands

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
)

func TestActionNew_ClientHeld_ValidatesAndSavesDraft(t *testing.T) {
	dir := t.TempDir()
	api := &fakeActionsAPI{validateResp: dto4sneatai.ActionResponse{
		Semantic:   actionspec.Semantic{Kind: actionspec.KindBuy},
		Validation: needsInputValidation(),
	}}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"action", "new", "--json", semanticJSON(actionspec.KindBuy), "--space", "famID"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.validateReq == nil {
		t.Fatal("ActionValidate not called")
	}
	if api.validateReq.Semantic == nil || api.validateReq.Semantic.Kind != actionspec.KindBuy {
		t.Fatalf("validate semantic = %+v", api.validateReq.Semantic)
	}
	var resp dto4sneatai.ActionResponse
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, buf.String())
	}
	if resp.ActionID == "" {
		t.Fatalf("no actionID printed: %+v", resp)
	}

	store := actionstoreFor(t, dir)
	draft, err := store.Load(resp.ActionID)
	if err != nil {
		t.Fatalf("Load draft: %v", err)
	}
	if draft.SpaceID != "famID" || draft.Semantic.Kind != actionspec.KindBuy {
		t.Fatalf("draft = %+v", draft)
	}
	if draft.ServerHeld {
		t.Fatalf("draft should not be server-held")
	}
}

func TestActionNew_MissingJSON_Errors(t *testing.T) {
	env := actionsEnv(&fakeActionsAPI{}, t.TempDir())
	root := Root(env)
	root.AddCommand(Action(env))
	root.SetArgs([]string{"action", "new"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error: --json required")
	}
}

func TestActionNew_Server_CreatesAndIndexesLocally(t *testing.T) {
	dir := t.TempDir()
	api := &fakeActionsAPI{createResp: dto4sneatai.ActionResponse{
		ActionID: "act_srvabc12", Semantic: actionspec.Semantic{Kind: actionspec.KindBuy},
	}}
	env := actionsEnv(api, dir)
	root := Root(env)
	root.AddCommand(Action(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"action", "new", "--json", semanticJSON(actionspec.KindBuy), "--space", "famID", "--server"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.createReq == nil {
		t.Fatal("ActionCreate not called")
	}
	if api.validateReq != nil {
		t.Fatal("ActionValidate should not be called in --server mode")
	}
	store := actionstoreFor(t, dir)
	draft, err := store.Load("act_srvabc12")
	if err != nil {
		t.Fatalf("Load index entry: %v", err)
	}
	if !draft.ServerHeld {
		t.Fatalf("draft should be server-held: %+v", draft)
	}
}
