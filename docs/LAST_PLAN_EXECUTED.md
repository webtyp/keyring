---
PLAN: "feat: keyring:// .env marker convention and OpenForModule scoping"
EXECUTOR: jules
REVIEWER: none
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.

# Plan — the `.env` marker this package owns, plus module-scoped opening

## Part of a multi-repo wave

This is one piece of `KEYRING_DOTENV_MASTER_PLAN.md` (orchestrator:
`webtyp.com/app-releases`, `docs/KEYRING_DOTENV_MASTER_PLAN.md`). This
piece has **no dependency on any other piece of that wave** — dispatch it
immediately, in parallel with the `webtyp/kvdb` and `webtyp/tui` plans.
`webtyp/env` and `webtyp/app` each depend on this publishing first.

## Why

Some `.env` values are credentials (a client secret, a token) that must never
sit in a plaintext file. The ecosystem's answer is: the `.env` line declares
**which variable** needs a secret, never the secret itself —
`IAM_CLIENT_SECRET=keyring://` — and the real value lives in the OS credential
store via this package. Two things are still missing to make that usable by
any consumer:

1. Nothing names the marker string as a constant — every consumer would
   invent its own literal, and per this ecosystem's rule a consumer never
   re-derives a convention a library should own (see
   `webtyp/app-releases/docs/CONSTRUCTION_HARNESS.md`, "Lego pieces, never
   forks").
2. `OpenKeyring(service)` (in `keyring/auto`) requires the caller to already
   know its `service` namespace. `webtyp/env` and `webtyp/app` both need
   to scope by "whatever Go module the current project is" with zero
   configuration, so two different projects using the same variable name
   (`IAM_CLIENT_SECRET` in both `misitio` and some other consumer) never
   collide in the same OS keychain.

## The change

### 1. The marker — new file **`reference.go`** (package `keyring`, repo root)

```go
package keyring

// Scheme is the .env marker meaning "this value lives in the keyring, not on
// disk" — e.g. IAM_CLIENT_SECRET=keyring://. webtyp/env resolves it via
// auto.OpenForModule; webtyp/app's SECRETS section is how a developer fills
// it in. Full convention: KEYRING_DOTENV_MASTER_PLAN.md.
const Scheme = "keyring://"

// IsReference reports whether raw is exactly the Scheme marker rather than a
// literal value already stored inline in the .env file.
func IsReference(raw string) bool { return raw == Scheme }
```

### 2. Module-scoped opening — new file **`auto/module.go`** (package `auto`)

```go
package auto

import (
	"os"
	"path/filepath"
	"strings"

	"webtyp.com/keyring"
)

const (
	// ErrGoModNotFound is returned when no go.mod exists in dir or any parent.
	ErrGoModNotFound = keyring.Error("keyring/auto: no go.mod found walking up from the given directory")
	// ErrNoModuleDirective is returned when go.mod has no `module` line.
	ErrNoModuleDirective = keyring.Error("keyring/auto: go.mod has no module directive")
)

// findModulePath walks up from dir looking for go.mod and returns the value
// of its `module` directive (e.g. "github.com/veltylabs/misitio").
func findModulePath(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for {
		data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if after, ok := strings.CutPrefix(line, "module "); ok {
					return strings.TrimSpace(after), nil
				}
			}
			return "", ErrNoModuleDirective
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrGoModNotFound
		}
		dir = parent
	}
}

// OpenForModule opens (auto-selecting the platform's provider, same as
// OpenKeyring) a Keyring scoped to the Go module found by walking up from
// dir — pass "." for the current working directory. The module's own path is
// the service namespace: it is already globally unique, so two different
// projects' identically named keys never collide, and no project has to
// remember to declare a service name of its own.
func OpenForModule(dir string) (*Keyring, error) {
	service, err := findModulePath(dir)
	if err != nil {
		return nil, err
	}
	return OpenKeyring(service), nil
}
```

`keyring.Error` (the comparable-string error type) already exists in
[`errors.go`](errors.go) — reuse it exactly as shown, do not declare a new
error type for these two sentinels.

## Anti-footguns

- **Do not** make `Scheme` configurable or add a variant syntax
  (`keyring://<key>`, `keyring://<service>/<key>`) — the service is always the
  caller's own module path (from `OpenForModule`) and the key is always the
  `.env` variable name itself; encoding either into the marker string would
  duplicate information every caller already has.
- **`auto/module.go` is `!wasm` implicitly** — `os.ReadFile`/`filepath` already
  behave correctly there; do not add a build tag that does not already exist
  on the other files in `keyring/auto`.
- Do not touch `webtyp/env` or `webtyp/app` in this plan — they are
  separate plans in the same wave that depend on this one.

## Tests

File: **`reference_test.go`** (new, package `keyring`, repo root).

| Test | Asserts |
|---|---|
| `TestIsReference` | `IsReference("keyring://")` is `true`; `IsReference("keyring://foo")` is `false`; `IsReference("")` is `false`; `IsReference("https://example.com")` is `false` |

File: **`auto/module_test.go`** (new, package `auto`).

`findModulePath` is pure filesystem logic — no keychain involved, safe to run
anywhere. `OpenForModule` (and anything that calls `.Get`/`.Set` through it)
opens the **real** platform provider, unlike this repo's own
`tests.RunConformance` harness (see [`tests/conformance.go`](tests/conformance.go)),
which takes an explicit, test-controlled `keyring.Provider` — there is no
injectable fake behind `auto`, by design (its entire point is "whatever this
machine actually has"). So any test that calls `OpenForModule` and then
`Set`/`Get` must:
- Guard with `if _, err := OpenForModule(t.TempDir()); ...` — no, use the
  intended dir; concretely: probe first with `kr, err := OpenForModule(".")`
  and `if err != nil { t.Skip(...) }`, then also skip if the first `Set` call
  fails (`ErrUnavailable`/`ErrUnsupported`) — headless CI often has no D-Bus
  session or keychain.
- `t.Cleanup(func() { kr.Delete(probeKey) })` immediately after a successful
  `Set`, so a failed assertion never leaves a stale entry for the next run.

| Test | Asserts |
|---|---|
| `TestFindModulePathFindsOwnGoMod` | `findModulePath(".")` run from this repo returns `"webtyp.com/keyring"` (this repo's own `go.mod`) — no keychain involved, no skip needed |
| `TestFindModulePathWalksUpFromSubdir` | `findModulePath("./auto")` (a subdirectory with no `go.mod` of its own) still returns `"webtyp.com/keyring"` — no skip needed |
| `TestFindModulePathNotFound` | `findModulePath(t.TempDir())` (an empty temp directory, no `go.mod` anywhere above it up to `/`) returns `ErrGoModNotFound` — no skip needed; if flaky under `/tmp` on a machine where a stray `go.mod` exists above the temp root, assert only that an error is returned, not which one |
| `TestOpenForModuleUsesModulePath` | probe-and-skip as above; then `Set` a probe key through `OpenForModule(".")`'s result, `t.Cleanup` a `Delete` of it, and confirm `OpenKeyring("webtyp.com/keyring")` (this repo's own module path, opened directly) reads back the same value — proving `OpenForModule` scoped to that exact service |

## Acceptance criteria

- [ ] `go build ./...` and `go vet ./...` clean.
- [ ] `go test ./...` green, including all new tests.
- [ ] `grep -n "const Scheme" reference.go` → one match.
- [ ] `grep -n "func OpenForModule" auto/module.go` → one match.
- [ ] `grep -rn "keyring:// \|keyring://<" .` outside `reference.go`/`module.go`/tests → empty: the marker is defined in exactly one place.

## Out of scope

`webtyp/env` resolving the marker, and `webtyp/app`'s interactive section
for filling in a missing secret — both are separate plans in this wave,
dispatched once this one publishes.
