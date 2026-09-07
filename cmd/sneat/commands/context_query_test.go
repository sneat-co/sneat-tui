package commands

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
)

func TestContext_PrintsSnapshotAsJSONByDefault(t *testing.T) {
	api := &fakeActionsAPI{contextResp: dto4sneatai.ContextResponse{Today: "2026-09-07"}}
	env := actionsEnv(api, t.TempDir())
	root := Root(env)
	root.AddCommand(Context(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"context", "--space", "famID"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.contextSpaceID != "famID" {
		t.Fatalf("spaceID = %q", api.contextSpaceID)
	}
	var got dto4sneatai.ContextResponse
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, buf.String())
	}
	if got.Today != "2026-09-07" {
		t.Fatalf("got = %+v", got)
	}
}

func TestContext_DefaultsToFamilySpace(t *testing.T) {
	api := &fakeActionsAPI{}
	env := actionsEnv(api, t.TempDir())
	root := Root(env)
	root.AddCommand(Context(env))
	root.SetArgs([]string{"context"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.contextSpaceID != "famID" {
		t.Fatalf("spaceID = %q, want famID default", api.contextSpaceID)
	}
}

func TestQuery_PostsAndPrintsResponse(t *testing.T) {
	api := &fakeActionsAPI{queryResp: dto4sneatai.QueryResponse{Before: "2026-09-30"}}
	env := actionsEnv(api, t.TempDir())
	root := Root(env)
	root.AddCommand(Query(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"query", "--space", "famID", "--json",
		`{"contacts":[{"mention":"Vasilisa"}],"before":{"relative":"end_of_month"}}`})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if api.queryReq == nil {
		t.Fatal("Query not called")
	}
	if string(api.queryReq.SpaceID) != "famID" {
		t.Fatalf("spaceID = %q", api.queryReq.SpaceID)
	}
	if len(api.queryReq.Contacts) != 1 || api.queryReq.Contacts[0].Mention != "Vasilisa" {
		t.Fatalf("contacts = %+v", api.queryReq.Contacts)
	}
	if api.queryReq.Before == nil || api.queryReq.Before.Relative != "end_of_month" {
		t.Fatalf("before = %+v", api.queryReq.Before)
	}
	var got dto4sneatai.QueryResponse
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, buf.String())
	}
	if got.Before != "2026-09-30" {
		t.Fatalf("got = %+v", got)
	}
}

func TestQuery_MissingJSON_Errors(t *testing.T) {
	env := actionsEnv(&fakeActionsAPI{}, t.TempDir())
	root := Root(env)
	root.AddCommand(Query(env))
	root.SetArgs([]string{"query"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error: --json required")
	}
}
