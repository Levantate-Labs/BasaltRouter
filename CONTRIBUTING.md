# Contributing to BasaltRouter

Thank you for your interest in contributing. This document summarizes the
workflow and conventions enforced in this repository. For full CI/CD details,
see [`.cursor/skills/basaltrouter-cicd/SKILL.md`](.cursor/skills/basaltrouter-cicd/SKILL.md).

## Quick Reference

```
Branch:   <type>/<scope>/<short-description>
Commit:   <type>(<scope>): <subject>
PR title: same format as commit subject
Tag:      v<MAJOR>.<MINOR>.<PATCH>  (annotated only)
```

**Types:** `feat`, `fix`, `chore`, `docs`, `refactor`, `perf`, `test`, `ci`, `hotfix`, `release`

**Scopes:** `gateway`, `adapter`, `keystore`, `vault`, `budget`, `ratelimit`, `router`, `cache`, `events`, `audit`, `notifier`, `worker`, `api`, `dashboard`, `auth`, `db`, `deploy`, `docs`, `config`

## Branch Naming

```
<type>/<scope>/<short-description>
```

- All lowercase, words separated by hyphens
- Short description: 2–5 words, no articles
- Max 60 characters total
- Example: `feat/gateway/streaming-support`

## Commit Messages

BasaltRouter follows [Conventional Commits 1.0.0](https://www.conventionalcommits.org/).

```
<type>(<scope>): <subject>

[optional body]

[optional footer(s)]
```

- Subject: lowercase, imperative mood, no period, max 72 characters
- Use `BREAKING CHANGE:` footer for breaking changes
- Use `Closes #<issue>` to link issues

## Before Every Commit

```bash
go fmt ./...
golangci-lint run
go vet ./...
go test ./...
```

Or run `make check`.

- No debug logging left in
- No commented-out code blocks
- No `.env` or secrets in commits
- Signed commits are required

## Pull Requests

- **Never commit directly to `main`** — all changes via PR
- PR title follows Conventional Commits format
- Fill out the PR template completely
- Min 1 approval required (2 for `internal/vault`, `internal/keystore`, security changes)
- All CI checks must pass before merge
- Keep PRs focused — under 400 lines changed when possible

### Merge Strategy

| Branch type | Merge strategy |
|---|---|
| `feat/*`, `fix/*`, `chore/*`, `docs/*`, `refactor/*`, `perf/*`, `test/*`, `ci/*` | Squash merge |
| `hotfix/*`, `release/*` | Merge commit |

## Local Development

```bash
cp .env.example .env
make docker-up
make migrate
make test
make lint
```

## License

By contributing, you agree that your contributions will be licensed under the
Business Source License 1.1, converting to Apache License 2.0 on the Change Date
specified in [LICENSE](LICENSE).
