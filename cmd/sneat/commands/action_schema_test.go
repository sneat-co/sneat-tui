package commands

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
)

func TestActionSchema_Server_FetchesFromAPI(t *testing.T) {
	api := &fakeActionsAPI{schemaResp: actionspec.Schema{Version: "server-version"}}
	env := actionsEnv(api, t.TempDir())
	root := Root(env)
	root.AddCommand(Action(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"action", "schema"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var got actionspec.Schema
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, buf.String())
	}
	if got.Version != "server-version" {
		t.Fatalf("got = %+v", got)
	}
}

func TestActionSchema_Offline_PrintsCompiledInSchema(t *testing.T) {
	env := actionsEnv(&fakeActionsAPI{}, t.TempDir())
	root := Root(env)
	root.AddCommand(Action(env))
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"action", "schema", "--offline"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var got actionspec.Schema
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, buf.String())
	}
	if got.Version != actionspec.SchemaVersion {
		t.Fatalf("got.Version = %q, want %q", got.Version, actionspec.SchemaVersion)
	}
}
