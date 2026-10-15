---
PLAN: "refactor!: delete the keyring:// .env marker (Scheme, IsReference)"
EXECUTOR: jules
REVIEWER: none
---

> Read AGENTS.md at the repo root before any change.

# Plan — keyring: the `.env` marker is gone

Part of the "Despliegue guiado" wave (`https://github.com/webtyp/docs/blob/main/DEPLOY_WIZARD_MASTER_PLAN.md`,
decision D15). `keyring.Scheme` (`"keyring://"`) and `keyring.IsReference` existed so a project could
write `KEY=keyring://` in `.env`. Secrets are now declared in the project's `config/deploy.go`, and
`webtyp/env` looks a key up in the keyring after the environment and `.env`, with no marker.

**Prerequisite:** the published `webtyp.com/env` no longer imports `keyring.Scheme`/`IsReference`
(its plan of this wave). Check: `go list -m -versions webtyp.com/env` and in the newest version
`grep -rn "IsReference\|keyring.Scheme"` → empty. If it still uses them, **stop** and say so in the PR.

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
```

## Design gate

1. **Prior art.** Not applicable: pure deletion of an exported convention.
2. **Names.** None added.
3. **Ledger.** Concepts −1 (`keyring://`). Ways to declare a local secret: −1.
4. **Where.** Here: the symbols live here.
5. **Deletes.** `reference.go`, `reference_test.go`, the `keyring://` section of `README.md` and of
   any file under `docs/`.

## Stage 1 — delete

- Delete `reference.go` and `reference_test.go`.
- `README.md` and `docs/*`: delete every mention of `keyring://`, `Scheme`, `IsReference` and of
  `KEYRING_DOTENV_MASTER_PLAN.md`. Where a paragraph explained the marker, replace it with one
  line: a project declares its secrets in `config/deploy.go`; `webtyp/env` reads them from the
  module's keyring (`auto.OpenForModule`).

| Stage | Files |
|---|---|
| 1 | delete `reference.go`, `reference_test.go`; `README.md`; `docs/*` that mention the marker |

## Acceptance

- `gotest` green.
- `grep -rn "keyring://\|IsReference\|Scheme" --include=*.go --include=*.md .` → empty.
- `git diff --stat main` deletes only `reference.go` and `reference_test.go`.
