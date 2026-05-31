---
name: basltrouter-cicd
description: >
  CI/CD conventions, Git practices, branch naming, commit format, PR standards,
  and release workflow for the BasaltRouter project. Use this skill whenever the user
  mentions creating a branch, writing a commit message, opening a PR, tagging a
  release, setting up a pipeline, or asks about naming anything in the repo
  (branches, tags, environments). Also trigger for any question about what to do
  before pushing code, how to structure a PR, or what the CI checks are.
---

# BasaltRouter CI/CD Practices

## 1. Branch Naming

All branches must follow this pattern:

```
<type>/<scope>/<short-description>
```

### Types

| Type | When to use |
|---|---|
| `feat` | New feature or capability |
| `fix` | Bug fix |
| `chore` | Tooling, deps, config — no production code change |
| `docs` | Documentation only |
| `refactor` | Code restructure with no behavior change |
| `perf` | Performance improvement |
| `test` | Adding or fixing tests only |
| `ci` | CI/CD pipeline changes |
| `hotfix` | Urgent production fix branched from `main` |
| `release` | Release preparation branch |

### Scopes (BasaltRouter-specific)

```
gateway       adapter       keystore      vault
budget        ratelimit     router        cache
events        audit         notifier      worker
api           dashboard     auth          db
deploy        docs          config
```

### Rules

- All lowercase, words separated by hyphens
- Short description: 2–5 words max, no articles (a/the)
- No special characters except `/` and `-`
- Max 60 characters total

### Examples

```
feat/gateway/streaming-support
fix/budget/redis-window-reset
chore/deps/upgrade-pgx-v5
docs/adapter/anthropic-setup-guide
refactor/router/circuit-breaker-cleanup
hotfix/keystore/timing-attack-fix
release/v1.2.0
```

---

## 2. Commit Message Format

BasaltRouter follows **Conventional Commits 1.0.0**.

### Structure

```
<type>(<scope>): <subject>

[optional body]

[optional footer(s)]
```

### Subject Line Rules

- Type and scope are lowercase
- Subject is lowercase, imperative mood ("add" not "added", "fix" not "fixes")
- No period at the end
- Max 72 characters total (type + scope + subject)
- Must complete the sentence: *"If applied, this commit will..."*

### Body Rules

- Wrap at 100 characters per line
- Explain **what** and **why**, not how
- Separate from subject with a blank line
- Use present tense

### Footer Rules

- `BREAKING CHANGE:` prefix for breaking API changes (triggers major version bump)
- `Closes #<issue>` or `Fixes #<issue>` to link issues
- `Co-authored-by: Name <email>` for pair work
- Multiple footers each on their own line

### Examples

**Simple fix:**
```
fix(budget): correct rolling window reset calculation
```

**Feature with body:**
```
feat(adapter): add Anthropic streaming support

Implements SSE streaming for the Anthropic Messages API adapter.
Token counts are accumulated from stream chunks since Anthropic only
reports usage on the final chunk.

Closes #42
```

**Breaking change:**
```
feat(api): rename virtual_key to vkey in all responses

BREAKING CHANGE: the `virtual_key` field in all key management API
responses is now `vkey`. Clients must update their field references.

Closes #88
```

**Chore with no body:**
```
chore(deps): upgrade golang-migrate to v4.18.1
```

**Multi-scope (use the most specific):**
```
refactor(router): extract circuit breaker into separate package
```

### What NOT to do

```
# Too vague
fix: bug fix

# Past tense
feat(gateway): added streaming

# Period at end
fix(budget): correct window reset.

# No scope when scope is obvious
feat: add anthropic adapter

# WIP commits (squash before merging)
wip: halfway through adapter
```

---

## 3. Branch Lifecycle

### Standard Feature Flow

```
main
 └── feat/gateway/streaming-support   ← branch from main
      └── [commits]
      └── PR → main                   ← squash merge
```

### Hotfix Flow

```
main
 └── hotfix/keystore/timing-attack-fix  ← branch from main
      └── [fix commit]
      └── PR → main                     ← merge commit (preserve history)
      └── tag: v1.x.y
```

### Release Flow

```
main
 └── release/v1.2.0   ← branch from main when ready
      └── chore(release): bump version to v1.2.0
      └── [only bugfixes allowed here]
      └── PR → main
      └── tag: v1.2.0
```

### Rules

- **Never commit directly to `main`** — all changes via PR
- **Never rebase shared branches** — only rebase local, unpushed branches
- Delete branches after merge
- `main` is always deployable
- Branch from `main` always (not from another feature branch) unless explicitly building on another feature

---

## 4. Commit Hygiene

### Before Every Commit

- `go fmt ./...` — format all Go code
- `golangci-lint run` — must pass with zero errors
- `go test ./...` — all tests must pass locally
- No debug logging (`fmt.Println`, `log.Printf` added for debugging) left in
- No commented-out code blocks
- No `.env` or secrets in the commit

### Atomic Commits

Each commit must be a single logical change. If you're writing "and" in your subject line, split it into two commits.

```
# Wrong — two things in one commit
feat(budget): add soft limit check and fix hard limit redis key

# Right — two separate commits
feat(budget): add soft limit threshold check
fix(budget): correct hard limit redis key format
```

### Squashing

Before opening a PR, squash WIP/fixup commits into logical units. It's fine to have multiple commits in a PR — each should tell a coherent story.

```bash
# Interactive rebase to squash last N commits
git rebase -i HEAD~N
```

---

## 5. Pull Request Standards

### PR Title

Must follow the same Conventional Commits format as a commit subject:

```
feat(adapter): add Groq provider support
fix(ratelimit): handle concurrent request counter overflow
```

### PR Description Template

Every PR must include this structure (add as GitHub PR template at `.github/pull_request_template.md`):

```markdown
## What
<!-- One paragraph: what does this PR do? -->

## Why
<!-- Why is this change needed? Link the issue. -->
Closes #

## How
<!-- Brief explanation of the approach. Not a line-by-line walkthrough. -->

## Testing
<!-- How was this tested? Unit tests? Manual? Which providers? -->

## Breaking Changes
<!-- List any breaking API or config changes. Write NONE if not applicable. -->

## Checklist
- [ ] `go fmt` and `golangci-lint` pass
- [ ] Tests added or updated
- [ ] Migrations reversible (down migration included)
- [ ] No secrets or credentials in code
- [ ] Documentation updated if needed
- [ ] CHANGELOG.md updated if user-facing change
```

### PR Rules

- Min 1 approval required to merge (2 for `internal/vault`, `internal/keystore`, security-related changes)
- All CI checks must be green before merge
- No merging your own PR without review (except hotfixes under time pressure — document why in PR)
- Keep PRs focused — under 400 lines changed is ideal, 800 lines is the soft ceiling. Split larger changes.
- Reviewers should be assigned, not just requested
- Draft PRs are encouraged for early feedback — prefix title with `[WIP]` if not using GitHub Draft mode

### Merge Strategy

| Branch type | Merge strategy | Reason |
|---|---|---|
| `feat/*`, `fix/*`, `chore/*`, `docs/*`, `refactor/*`, `perf/*`, `test/*`, `ci/*` | **Squash merge** | Clean linear history on main |
| `hotfix/*` | **Merge commit** | Preserve fix history, easier to audit |
| `release/*` | **Merge commit** | Preserve release history |

---

## 6. Versioning

BasaltRouter follows **Semantic Versioning 2.0.0 (SemVer)**.

```
v<MAJOR>.<MINOR>.<PATCH>
```

| Change | Version bump |
|---|---|
| Breaking API or config change (`BREAKING CHANGE:` footer) | MAJOR |
| New feature, backward compatible (`feat:`) | MINOR |
| Bug fix, backward compatible (`fix:`) | PATCH |
| Hotfix on released version | PATCH |

### Pre-release Tags

```
v1.0.0-alpha.1    ← early unstable
v1.0.0-beta.1     ← feature complete, testing
v1.0.0-rc.1       ← release candidate
v1.0.0            ← stable release
```

### Git Tags

```bash
# Annotated tags only — never lightweight tags for releases
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

---

## 7. CI Pipeline (GitHub Actions)

### Triggers

| Pipeline | Triggers on |
|---|---|
| `ci.yml` (PR checks) | Every PR, every push to PR branch |
| `release.yml` | Push of `v*` tag |
| `nightly.yml` | Cron: 02:00 UTC daily on `main` |

### PR Pipeline (`ci.yml`) — Must All Pass

```
lint        → golangci-lint run
format      → gofmt -l (fail if any files need formatting)
vet         → go vet ./...
test        → go test -race -coverprofile=coverage.out ./...
coverage    → fail if coverage drops below 70%
build       → go build ./cmd/gateway ./cmd/api ./cmd/worker
migrate     → run migrations against test Postgres, verify rollback
security    → govulncheck ./...
secrets     → gitleaks detect --no-git
```

### Release Pipeline (`release.yml`)

```
[all PR pipeline steps]
    ↓
docker build (gateway, api, worker)
    ↓
docker push → ghcr.io/basltrouter/<image>:<version>
    ↓
helm package → deploy/helm/basltrouter/
    ↓
github release → attach binaries, helm chart, checksums
    ↓
update CHANGELOG.md from conventional commits
```

### Nightly Pipeline

```
go test -race ./...              ← full test suite with race detector
govulncheck ./...                ← fresh CVE scan
docker build --no-cache          ← verify clean build
notify on failure → Discord/Slack webhook
```

---

## 8. Environment & Secret Handling

### Environments

| Name | Branch | Purpose |
|---|---|---|
| `local` | any | Developer machine, Docker Compose |
| `ci` | PR branches | Ephemeral, GitHub Actions |
| `staging` | `main` | Auto-deployed on merge to main |
| `production` | `v*` tag | Manual trigger after tag |

### Rules

- **Never hardcode credentials** — all secrets via environment variables
- Use `.env.example` with placeholder values — always keep it updated
- `.env` is in `.gitignore` — never commit it
- CI secrets stored in GitHub Actions Secrets, never in workflow YAML
- Production secrets in a secrets manager (Vault, AWS Secrets Manager) — not in environment files

---

## 9. CHANGELOG

BasaltRouter maintains a `CHANGELOG.md` following [Keep a Changelog](https://keepachangelog.com) format.

### Structure

```markdown
# Changelog

## [Unreleased]
### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security

## [1.2.0] - 2025-06-15
### Added
- Groq provider adapter (#55)
### Fixed
- Budget window not resetting on month boundary (#61)
```

### Rules

- Every user-facing change gets a CHANGELOG entry in the PR
- Internal refactors, test additions, and CI changes do not need entries
- Security fixes always get a `### Security` entry regardless of scope
- The `[Unreleased]` section is converted to a version heading at release time

---

## 10. Quick Reference Card

```
Branch:   <type>/<scope>/<short-description>
Commit:   <type>(<scope>): <subject>
PR title: same format as commit subject
Tag:      v<MAJOR>.<MINOR>.<PATCH>  (annotated only)

Types:    feat fix chore docs refactor perf test ci hotfix release
Scopes:   gateway adapter keystore vault budget ratelimit router
          cache events audit notifier worker api dashboard auth db
          deploy docs config

Merge:    feat/fix/chore/docs → squash merge into main
          hotfix/release      → merge commit into main

Before PR: fmt → lint → vet → test → no secrets → atomic commits
```
