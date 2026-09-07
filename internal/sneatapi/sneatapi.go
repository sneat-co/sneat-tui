// Package sneatapi is a thin client for the sneat-go HTTP API (mutations),
// authenticated with the user's Firebase ID token.
package sneatapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sneat-co/contactus/backend/dto4contactus"
	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/sneat-co/sneat-ai-backend/api4sneatai"
	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"golang.org/x/oauth2"
)

// TokenSource yields the bearer token for API calls (an oauth2 token whose
// AccessToken is the Firebase ID token).
type TokenSource interface {
	Token() (*oauth2.Token, error)
}

// Client calls the sneat-go API under a base URL like https://api.sneat.cloud/v0/.
type Client struct {
	http    *http.Client
	baseURL string
	ts      TokenSource
}

// New builds a Client. baseURL is normalized to end with "/".
func New(baseURL string, ts TokenSource, hc *http.Client) *Client {
	if hc == nil {
		hc = http.DefaultClient
	}
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return &Client{http: hc, baseURL: baseURL, ts: ts}
}

// CreateContact POSTs contactus/create_contact and returns the raw response.
func (c *Client) CreateContact(ctx context.Context, req dto4contactus.CreateContactRequest) (map[string]any, error) {
	var out map[string]any
	if err := c.do(ctx, http.MethodPost, "contactus/create_contact", req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteContact calls DELETE contactus/delete_contact.
func (c *Client) DeleteContact(ctx context.Context, req dto4contactus.ContactRequest) error {
	return c.do(ctx, http.MethodDelete, "contactus/delete_contact", req, nil)
}

// sneataiPath turns an api4sneatai path constant (e.g. "/v0/sneatai/action_create")
// into the relative path this Client's baseURL (which already ends in "/v0/") expects.
func sneataiPath(p string) string {
	return strings.TrimPrefix(p, "/v0/")
}

// ActionCreate POSTs action_create, starting a server-held action.
func (c *Client) ActionCreate(ctx context.Context, req dto4sneatai.CreateActionRequest) (dto4sneatai.ActionResponse, error) {
	var out dto4sneatai.ActionResponse
	err := c.do(ctx, http.MethodPost, sneataiPath(api4sneatai.PathActionCreate), req, &out)
	return out, err
}

// ActionPatch POSTs action_patch, applying a semantic patch to a server-held action.
func (c *Client) ActionPatch(ctx context.Context, req dto4sneatai.PatchActionRequest) (dto4sneatai.ActionResponse, error) {
	var out dto4sneatai.ActionResponse
	err := c.do(ctx, http.MethodPost, sneataiPath(api4sneatai.PathActionPatch), req, &out)
	return out, err
}

// ActionGet GETs action_get for a server-held action.
func (c *Client) ActionGet(ctx context.Context, spaceID, actionID string) (dto4sneatai.ActionResponse, error) {
	var out dto4sneatai.ActionResponse
	q := url.Values{"spaceID": {spaceID}, "actionID": {actionID}}
	err := c.do(ctx, http.MethodGet, sneataiPath(api4sneatai.PathActionGet)+"?"+q.Encode(), nil, &out)
	return out, err
}

// ActionValidate POSTs action_validate, either against a stored action
// (ActionID set, Semantic nil) or an inline candidate (Semantic set).
func (c *Client) ActionValidate(ctx context.Context, req dto4sneatai.ValidateActionRequest) (dto4sneatai.ActionResponse, error) {
	var out dto4sneatai.ActionResponse
	err := c.do(ctx, http.MethodPost, sneataiPath(api4sneatai.PathActionValidate), req, &out)
	return out, err
}

// ActionCommit POSTs action_commit, either for a stored action (ActionID only)
// or a client-held draft (ActionID + Semantic).
func (c *Client) ActionCommit(ctx context.Context, req dto4sneatai.CommitActionRequest) (dto4sneatai.CommitActionResponse, error) {
	var out dto4sneatai.CommitActionResponse
	err := c.do(ctx, http.MethodPost, sneataiPath(api4sneatai.PathActionCommit), req, &out)
	return out, err
}

// ActionCancel POSTs action_cancel.
func (c *Client) ActionCancel(ctx context.Context, req dto4sneatai.ActionRequest) error {
	return c.do(ctx, http.MethodPost, sneataiPath(api4sneatai.PathActionCancel), req, nil)
}

// Context GETs the Space snapshot an agent may read before acting.
func (c *Client) Context(ctx context.Context, spaceID string) (dto4sneatai.ContextResponse, error) {
	var out dto4sneatai.ContextResponse
	q := url.Values{"spaceID": {spaceID}}
	err := c.do(ctx, http.MethodGet, sneataiPath(api4sneatai.PathContext)+"?"+q.Encode(), nil, &out)
	return out, err
}

// Query POSTs query, asking what to buy/schedule for contacts before a horizon.
func (c *Client) Query(ctx context.Context, req dto4sneatai.QueryRequest) (dto4sneatai.QueryResponse, error) {
	var out dto4sneatai.QueryResponse
	err := c.do(ctx, http.MethodPost, sneataiPath(api4sneatai.PathQuery), req, &out)
	return out, err
}

// Schema GETs the machine-readable Action Protocol schema. No auth is
// required by the server, but the token is sent when available.
func (c *Client) Schema(ctx context.Context) (actionspec.Schema, error) {
	var out actionspec.Schema
	err := c.do(ctx, http.MethodGet, sneataiPath(api4sneatai.PathSchema), nil, &out)
	return out, err
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, r)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.ts != nil {
		tok, err := c.ts.Token()
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sneat api %s %s: http %d: %s", method, path, resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if out != nil && len(data) > 0 {
		return json.Unmarshal(data, out)
	}
	return nil
}
