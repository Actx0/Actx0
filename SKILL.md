# actx0 project memory

Durable facts for agents working in this repo. Prefer this over reinventing wrappers or re-learning gotchas.

Always-on rules still apply: `.cursor/rules/strict-coding.mdc`, `.cursor/rules/go-spacing.mdc`. Do not duplicate them here.

## Maintain this file

When any of these happen, **edit root `SKILL.md` in the same turn**:

- User says "remember", "prefer", "we always", or corrects a pattern
- A debt review reveals a reusable convention (delete helper X, use Y)
- A lib gotcha is confirmed (vet false positive, API quirk)
- A one-liner helper is inlined or replaced by stdlib/`lo`

Keep entries short. Remove obsolete notes. Never store secrets.

## Debt review

- Checklist: repo-root `debt.md` (`- [ ]` / `- [x] path`) — Go packages and `web/` frontend.
- Slash command: `/debt` (see `.cursor/commands/debt.md`). Supports a path (e.g. `api/me.go`, `web/src/App.vue`).
- Goal while reviewing: strip AI slop, unused wrappers, defensive code; match happy-path style.

## Prefer / avoid

| Prefer | Avoid |
|--------|--------|
| `lo.FromPtr(p)` | custom `StringValue` / nil-deref helpers |
| pass `Id` / `lo.FromPtr(*Id)` into `Exec`/`Query` (uses `Id.Value()` → SQL NULL when empty) | `ValueOrNil` / `ValueOrNilPtr` helpers |
| `isNotFound(err)` → `nil, nil` on Get miss | repeating `errors.Is(err, sql.ErrNoRows)` at every call site |
| `GetDB()` for reads; `GetDB(true)` / `GetWriteDB()` for writes | treating `GetDB()` as always the primary |
| `new(true)`, `new(false)`, `new(x)` (Go 1.26 expr form) | `lo.ToPtr(literal)` for the same |
| `r.URL.Query().Get("id")` at call site | tiny util wrappers around one query get |
| `for x := range strings.SplitSeq(s, ",")` | `strings.Split` then range when you only iterate once |
| Reflect call to `gotext.Po.Get` for msgid lookup | direct `po.Get(key)` (vet: printf / non-constant format) |
| Stripe metadata camelCase (`workspaceId`, …) | snake_case metadata keys |
| `subscriptions_meta.value` as `to_jsonb($::text)` / read `value #>> '{}'` | treating that column as plain text SQL |
| vue-i18n literal `@` as `{'@'}` (e.g. `maya{'@'}example.com`) | raw `@` in locale strings (linked-message syntax; can blank the UI) |
| `PromptRepository` + `PromptVersionRepository` composed in module (no prompt/version SQL joins) | one mega-repo / joined `PromptVersion` DTO |
| Prompt resolve/write: `latest` (auto label) + `production` bool (+ version number); no custom labels on create/update | accepting arbitrary `labels` JSON on write / resolve-by-custom-label |
| inline response mapping at call sites | tiny `fooResponse(entity)` helpers |

## Locale

- Package `locale`: `.po` files, `T` / `TR`, Accept-Language.
- Keys are msgids (e.g. `not_authenticated`), not format strings.
- `getMsg` uses reflect so `go vet`/`printf` does not treat `Po.Get` as sprintf.
- Frontend `web/src/locales/*.json`: escape `@` with `{'@'}` so vue-i18n does not parse linked messages.

## Stack cheat sheet

- Module: `github.com/clivern/actx0`, Go **1.26.x**
- HTTP: chi; JSON/helpers in `pkg/util`
- Domain: `api` → `module` → `db` / `pkg/*` / `service/*`
- Config: viper (`app.*` keys)
- Common pkgs: `pkg/storage`, `pkg/qdrant`, `pkg/cache`, `pkg/stripe`, `pkg/resend`, `pkg/meteo`
- `pkg/chunk`: splitters are `recursive` and `markdown` only (no token splitter)

## Open notes

_Add scratch pad items here only if they must survive beyond one chat. Clear when done._
