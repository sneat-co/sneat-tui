---
format: https://specscore.md/feature-specification
status: Draft
---

# Feature: Action Protocol CLI

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/sneat-cli/spec/features/action-protocol-cli?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/sneat-cli/spec/features/action-protocol-cli?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/sneat-cli/spec/features/action-protocol-cli?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/sneat-cli/spec/features/action-protocol-cli?op=request-change) |
**Status:** Draft
**Source Ideas:** —

## Summary

sneat action/context/query commands: a thin CLI over the sneat-ai-backend Action Protocol

## Problem

An AI agent (a Claude Code skill from `sneat-co/sneat-ai`, or any other agent) needs to turn a user's
utterance into a real change in Sneat — buy something, schedule something, record a birthday — without
knowing Sneat's storage model, its Firestore schema, or which extension (Listus/Calendarius/Contactus)
owns the result. The agent only understands language; `sneat-ai-backend`'s Action Protocol
(`/v0/sneatai/*`) is what turns a semantic candidate into a committed change, deterministically,
asking for missing input rather than guessing.

Before this Feature, the only way to drive that protocol was raw HTTP. `sneat action`, `sneat context`
and `sneat query` are the CLI surface that makes it a shell command: the actionspec JSON goes in on
`--json`, and every response is JSON printed to stdout, matching how an agent already runs `sneat
contact`/`sneat space` and parses their output.

## Behavior

### Client-held drafts are the default

#### REQ: client-held-draft-default

`sneat action new` MUST default to holding a new action client-side: it generates an id of the form
`act_` + 8 lowercase base36 characters, calls `action_validate` with the semantic candidate inline
(stateless — the server does not persist anything), and saves the id, space, language, utterance, the
*resolved* semantic and validation the server returned, and timestamps as a local draft file at
`<UserConfigDir>/sneat/actions/<actionID>.json`. `--server` instead calls `action_create`, which
persists the action on the server under the id the server returns; the CLI keeps only a minimal local
index entry (id, space, a `serverHeld` marker) so later commands on that id know to route to the
server without an extra round trip to find out.

#### REQ: patch-merges-on-resolved-base

`sneat action add <actionID>` MUST merge the patch onto the *resolved* semantic from the action's last
validate/commit response (`actionspec.Merge`), not onto the raw candidate the agent first sent —
merging onto raw input would discard contact IDs and dates the server already resolved, forcing the
agent to re-supply them. For a client-held draft the merged candidate is re-validated inline
(`action_validate` with the merged `Semantic`) and the draft is updated with the response; passing the
action's stored `committedHappeningID` (if it was already committed) so an amendment does not trip the
duplicate-happening check against its own prior commit. For a server-held action the patch is sent as-is
to `action_patch`; merging and re-validation happen server-side.

### Commit sends the canonical candidate, and a refusal is scriptable

#### REQ: commit-sends-canonical-semantic

`sneat action commit <actionID>` for a client-held draft MUST send `actionspec.Canonical` of the
draft's current semantic — resolved values and display titles stripped — as the commit request's
`semantic`, alongside the client-generated `actionID`. This is the same normalization the server uses
for commit idempotency, so a re-commit of an unchanged candidate is recognized as a repeat rather than
a new mutation. A server-held action instead sends only its space and id; the server already holds the
authoritative candidate.

#### REQ: commit-refusal-exit-code

When the server refuses to commit because validation still needs input (`facade4sneatai.ErrNotCommittable`,
an HTTP 400 whose body contains "validation requires input"), `sneat action commit` MUST print the
server's error body as `{"error": "..."}` JSON to stdout and exit with process code 2, distinct from
the process code 1 used for every other command failure. An agent scripting around this command can
then branch on the exit code alone — "ask the user a question" vs. "something is actually broken" —
without parsing error text.

### Server mode has the same surface

#### REQ: server-mode-parity

Every `sneat action` subcommand (`show`, `validate`, `commit`, `cancel`) MUST behave identically from
the caller's point of view whether the action is client-held or server-held: same flags, same JSON
response shape (the underlying `dto4sneatai.ActionResponse`/`CommitActionResponse`), same exit-code
contract. The CLI decides which mode to use per action from its local index (`serverHeld`), never from
a flag repeated on every subsequent command — an agent that started an action with `new --server` does
not need to remember to pass `--server` again on `add`/`commit`/etc.

### The CLI knows nothing about how Sneat stores anything

#### REQ: no-storage-knowledge

No command in this Feature MUST read or write Firestore, a Listus/Calendarius/Contactus collection, or
any `sneat-core-modules`/`sneat-ext-contracts` domain type directly. Every mutation and read goes
through `internal/sneatapi.Client`'s Action Protocol methods (`ActionCreate`, `ActionPatch`,
`ActionGet`, `ActionValidate`, `ActionCommit`, `ActionCancel`, `Context`, `Query`, `Schema`), which
speak only `dto4sneatai`/`actionspec` types — the same contract boundary `sneat-ai-backend`'s HTTP
handlers enforce server-side. This is what keeps a storage migration on the backend from requiring a
CLI change.

## Acceptance Criteria

### AC: agent-can-complete-an-action-end-to-end

**Requirements:** action-protocol-cli#req:client-held-draft-default, action-protocol-cli#req:patch-merges-on-resolved-base, action-protocol-cli#req:commit-sends-canonical-semantic, action-protocol-cli#req:commit-refusal-exit-code

An agent can go from an utterance to a committed change using only `sneat action new`, zero or more
`sneat action add`, and `sneat action commit`, entirely through JSON on stdin/flags and stdout — never
touching a file format or a Firestore path. If the server needs more information, the agent sees that
as a distinct exit code and the server's own question text, not a stack trace.

### AC: client-and-server-modes-are-interchangeable

**Requirements:** action-protocol-cli#req:server-mode-parity, action-protocol-cli#req:no-storage-knowledge

Switching an integration from client-held drafts (the CLI default, suited to a short-lived agent
session) to server-held actions (suited to a long-running or multi-caller agent) is a single flag on
`action new`, with no other command changing shape — because both modes are thin wrappers over the same
`dto4sneatai` contract the server already defines.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
