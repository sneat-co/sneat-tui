package commands

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-ai-backend/semlint"
	"github.com/sneat-co/sneat-cli/internal/actionstore"
	"github.com/sneat-co/sneat-cli/internal/config"
	"github.com/sneat-co/sneat-cli/internal/session"
	"github.com/sneat-co/sneat-cli/internal/sneatauth"
)

// fakeActionsAPI is a scriptable ActionsAPI for command tests: each method
// records the last request it received and returns a preset response/error.
type fakeActionsAPI struct {
	validateReq  *dto4sneatai.ValidateActionRequest
	validateResp dto4sneatai.ActionResponse
	validateErr  error

	createReq  *dto4sneatai.CreateActionRequest
	createResp dto4sneatai.ActionResponse
	createErr  error

	patchReq  *dto4sneatai.PatchActionRequest
	patchResp dto4sneatai.ActionResponse
	patchErr  error

	getSpaceID, getActionID string
	getResp                 dto4sneatai.ActionResponse
	getErr                  error

	commitReq  *dto4sneatai.CommitActionRequest
	commitResp dto4sneatai.CommitActionResponse
	commitErr  error

	cancelReq *dto4sneatai.ActionRequest
	cancelErr error

	contextSpaceID string
	contextResp    dto4sneatai.ContextResponse
	contextErr     error

	queryReq  *dto4sneatai.QueryRequest
	queryResp dto4sneatai.QueryResponse
	queryErr  error

	schemaResp actionspec.Schema
	schemaErr  error
}

func (f *fakeActionsAPI) ActionCreate(_ context.Context, req dto4sneatai.CreateActionRequest) (dto4sneatai.ActionResponse, error) {
	f.createReq = &req
	return f.createResp, f.createErr
}

func (f *fakeActionsAPI) ActionPatch(_ context.Context, req dto4sneatai.PatchActionRequest) (dto4sneatai.ActionResponse, error) {
	f.patchReq = &req
	return f.patchResp, f.patchErr
}

func (f *fakeActionsAPI) ActionGet(_ context.Context, spaceID, actionID string) (dto4sneatai.ActionResponse, error) {
	f.getSpaceID, f.getActionID = spaceID, actionID
	return f.getResp, f.getErr
}

func (f *fakeActionsAPI) ActionValidate(_ context.Context, req dto4sneatai.ValidateActionRequest) (dto4sneatai.ActionResponse, error) {
	f.validateReq = &req
	return f.validateResp, f.validateErr
}

func (f *fakeActionsAPI) ActionCommit(_ context.Context, req dto4sneatai.CommitActionRequest) (dto4sneatai.CommitActionResponse, error) {
	f.commitReq = &req
	return f.commitResp, f.commitErr
}

func (f *fakeActionsAPI) ActionCancel(_ context.Context, req dto4sneatai.ActionRequest) error {
	f.cancelReq = &req
	return f.cancelErr
}

func (f *fakeActionsAPI) Context(_ context.Context, spaceID string) (dto4sneatai.ContextResponse, error) {
	f.contextSpaceID = spaceID
	return f.contextResp, f.contextErr
}

func (f *fakeActionsAPI) Query(_ context.Context, req dto4sneatai.QueryRequest) (dto4sneatai.QueryResponse, error) {
	f.queryReq = &req
	return f.queryResp, f.queryErr
}

func (f *fakeActionsAPI) Schema(_ context.Context) (actionspec.Schema, error) {
	return f.schemaResp, f.schemaErr
}

// actionsEnv builds an Env wired with a fake ActionsAPI and a client-config
// dir under t's temp dir (via SNEAT_CONFIG_DIR), so drafts never touch the
// real user config directory.
func actionsEnv(api *fakeActionsAPI, configDir string) Env {
	env := testEnv(&fakeStore{load: &session.Session{UID: "u1"}}, sneatauth.Result{})
	env.NewActionsAPI = func(config.Config) (ActionsAPI, error) { return api, nil }
	env.NewSpacesReader = func(config.Config) (SpacesReader, error) {
		return &fakeSpacesReader{spaces: map[string]any{
			"famID":  map[string]any{"type": "family"},
			"privID": map[string]any{"type": "private"},
		}}, nil
	}
	env.Getenv = func(k string) string {
		if k == "SNEAT_CONFIG_DIR" {
			return configDir
		}
		return ""
	}
	return env
}

func semanticJSON(kind actionspec.Kind) string {
	switch kind {
	case actionspec.KindBuy:
		return `{"kind":"buy","operations":[{"objects":[{"text":"shoes"}]}]}`
	default:
		return `{"kind":"` + string(kind) + `"}`
	}
}

// needsInputValidation is a Validation stub for a candidate that still needs
// a question answered (CanCommit false).
func needsInputValidation() semlint.Validation {
	return semlint.Validation{
		Status:    semlint.StatusNeedsInput,
		CanCommit: false,
		Questions: []semlint.Question{{Text: "Who is it for?", Blocking: true}},
	}
}

func canCommitValidation() semlint.Validation {
	return semlint.Validation{Status: semlint.StatusValidated, CanCommit: true}
}

// actionstoreFor opens the on-disk draft store a test's SNEAT_CONFIG_DIR
// points at, matching actionsDir's resolution.
func actionstoreFor(t *testing.T, configDir string) *actionstore.Store {
	t.Helper()
	return actionstore.NewStore(filepath.Join(configDir, "sneat", "actions"))
}
