# Branch Protection Setup

After the first CI run produces status check contexts, configure branch protection on `main`:

```bash
gh api repos/LevantateLabs/basaltrouter/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["lint","format","vet","test","coverage","build","migrate","security","secrets"]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"required_approving_review_count":1}' \
  --field required_signatures=true \
  --field required_linear_history=true
```

Requirements:
- Never commit directly to `main`
- All PRs require 1 approval (2 for `internal/vault`, `internal/keystore`)
- Signed commits enforced
- Squash merge for `chore/*`, `feat/*`, `fix/*` branches
