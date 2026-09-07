package sneatapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func TestActionCreate_Posts(t *testing.T) {
	var gotMethod, gotPath, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		_ = json.NewEncoder(w).Encode(dto4sneatai.ActionResponse{ActionID: "act_srv00001"})
	}))
	defer srv.Close()

	c := New(srv.URL, fakeTS{}, srv.Client())
	out, err := c.ActionCreate(context.Background(), dto4sneatai.CreateActionRequest{
		SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("s1")},
		Semantic:     actionspec.Semantic{Kind: actionspec.KindBuy},
	})
	if err != nil {
		t.Fatalf("ActionCreate: %v", err)
	}
	if gotMethod != http.MethodPost || !strings.HasSuffix(gotPath, "sneatai/action_create") {
		t.Fatalf("%s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `"kind":"buy"`) {
		t.Fatalf("body = %s", gotBody)
	}
	if out.ActionID != "act_srv00001" {
		t.Fatalf("out = %+v", out)
	}
}

func TestActionPatch_Posts(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewEncoder(w).Encode(dto4sneatai.ActionResponse{ActionID: "act_srv00001"})
	}))
	defer srv.Close()

	c := New(srv.URL, fakeTS{}, srv.Client())
	if _, err := c.ActionPatch(context.Background(), dto4sneatai.PatchActionRequest{}); err != nil {
		t.Fatalf("ActionPatch: %v", err)
	}
	if !strings.HasSuffix(gotPath, "sneatai/action_patch") {
		t.Fatalf("path = %s", gotPath)
	}
}

func TestActionGet_SendsQueryParams(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotQuery = r.Method, r.URL.Path, r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(dto4sneatai.ActionResponse{ActionID: "act_srv00001"})
	}))
	defer srv.Close()

	c := New(srv.URL, fakeTS{}, srv.Client())
	out, err := c.ActionGet(context.Background(), "s1", "act_srv00001")
	if err != nil {
		t.Fatalf("ActionGet: %v", err)
	}
	if gotMethod != http.MethodGet || !strings.HasSuffix(gotPath, "sneatai/action_get") {
		t.Fatalf("%s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotQuery, "spaceID=s1") || !strings.Contains(gotQuery, "actionID=act_srv00001") {
		t.Fatalf("query = %s", gotQuery)
	}
	if out.ActionID != "act_srv00001" {
		t.Fatalf("out = %+v", out)
	}
}

func TestActionValidate_Posts(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewEncoder(w).Encode(dto4sneatai.ActionResponse{ActionID: "act_srv00001"})
	}))
	defer srv.Close()

	c := New(srv.URL, fakeTS{}, srv.Client())
	if _, err := c.ActionValidate(context.Background(), dto4sneatai.ValidateActionRequest{}); err != nil {
		t.Fatalf("ActionValidate: %v", err)
	}
	if !strings.HasSuffix(gotPath, "sneatai/action_validate") {
		t.Fatalf("path = %s", gotPath)
	}
}

func TestActionCommit_Posts(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewEncoder(w).Encode(dto4sneatai.CommitActionResponse{ActionID: "act_srv00001"})
	}))
	defer srv.Close()

	c := New(srv.URL, fakeTS{}, srv.Client())
	out, err := c.ActionCommit(context.Background(), dto4sneatai.CommitActionRequest{})
	if err != nil {
		t.Fatalf("ActionCommit: %v", err)
	}
	if !strings.HasSuffix(gotPath, "sneatai/action_commit") {
		t.Fatalf("path = %s", gotPath)
	}
	if out.ActionID != "act_srv00001" {
		t.Fatalf("out = %+v", out)
	}
}

func TestActionCommit_ErrorPropagatesBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"action cannot be committed: validation requires input"}`))
	}))
	defer srv.Close()

	c := New(srv.URL, fakeTS{}, srv.Client())
	_, err := c.ActionCommit(context.Background(), dto4sneatai.CommitActionRequest{})
	if err == nil || !strings.Contains(err.Error(), "validation requires input") {
		t.Fatalf("err = %v", err)
	}
}

func TestActionCancel_Posts(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(srv.URL, fakeTS{}, srv.Client())
	if err := c.ActionCancel(context.Background(), dto4sneatai.ActionRequest{}); err != nil {
		t.Fatalf("ActionCancel: %v", err)
	}
	if !strings.HasSuffix(gotPath, "sneatai/action_cancel") {
		t.Fatalf("path = %s", gotPath)
	}
}

func TestContext_SendsQueryParam(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotQuery = r.Method, r.URL.Path, r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(dto4sneatai.ContextResponse{Today: "2026-09-07"})
	}))
	defer srv.Close()

	c := New(srv.URL, fakeTS{}, srv.Client())
	out, err := c.Context(context.Background(), "s1")
	if err != nil {
		t.Fatalf("Context: %v", err)
	}
	if gotMethod != http.MethodGet || !strings.HasSuffix(gotPath, "sneatai/context") {
		t.Fatalf("%s %s", gotMethod, gotPath)
	}
	if gotQuery != "spaceID=s1" {
		t.Fatalf("query = %s", gotQuery)
	}
	if out.Today != "2026-09-07" {
		t.Fatalf("out = %+v", out)
	}
}

func TestQuery_Posts(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewEncoder(w).Encode(dto4sneatai.QueryResponse{})
	}))
	defer srv.Close()

	c := New(srv.URL, fakeTS{}, srv.Client())
	if _, err := c.Query(context.Background(), dto4sneatai.QueryRequest{}); err != nil {
		t.Fatalf("Query: %v", err)
	}
	if !strings.HasSuffix(gotPath, "sneatai/query") {
		t.Fatalf("path = %s", gotPath)
	}
}

func TestSchema_Gets(t *testing.T) {
	var gotMethod, gotPath, gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotAuth = r.Method, r.URL.Path, r.Header.Get("Authorization")
		_ = json.NewEncoder(w).Encode(actionspec.CurrentSchema())
	}))
	defer srv.Close()

	c := New(srv.URL, fakeTS{}, srv.Client())
	out, err := c.Schema(context.Background())
	if err != nil {
		t.Fatalf("Schema: %v", err)
	}
	if gotMethod != http.MethodGet || !strings.HasSuffix(gotPath, "sneatai/schema") {
		t.Fatalf("%s %s", gotMethod, gotPath)
	}
	if gotAuth != "Bearer idt" {
		t.Fatalf("auth = %q, want token sent even though not required", gotAuth)
	}
	if out.Version != actionspec.SchemaVersion {
		t.Fatalf("out.Version = %q", out.Version)
	}
}
