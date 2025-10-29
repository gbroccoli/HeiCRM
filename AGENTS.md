# Repository Guidelines

## Project Structure & Modules
- Root config: `config.yaml` (use `config.example.yaml` as a template).
- Shared packages: `pkg/` (`config`, `dbx`, `logx`).
- Services: `services/` (e.g., `services/auth/cmd/auth/main.go`).
- Infrastructure: `docker-compose.yml` (Postgres, NATS), `nats/` (JetStream data), `logs/`.
- Go workspace managed by `go.work` with per-package `go.mod` files.

## Build, Run, Test
- Start infra: `docker compose up -d db nats` (exposes Postgres 5432, NATS 4222/8222).
- Run auth service: `go run ./services/auth/cmd/auth` (listens on `:8080`).
- Build binary: `go build -o bin/auth ./services/auth/cmd/auth`.
- Run all tests: `go test ./...` (add `-v` or `-cover` as needed).

## Coding Style & Naming
- Formatting: `go fmt ./...`; validate: `go vet ./...`.
- `.editorconfig` enforced:
  - Go: tabs for indentation.
  - YAML/JSON/Shell: 2 spaces; line endings CRLF.
- Go conventions: packages are short lowercase (`config`, `dbx`); exported names use `CamelCase`; files use underscores (`login_test.go`).
- Errors: return `error`, wrap with context; prefer `context.Context` as first arg (`ctx`).

## Testing Guidelines
- Use Go’s standard `testing` package.
- Test files: `*_test.go` next to code (e.g., `services/auth/internal/handler/login_test.go`).
- Structure: table-driven tests for handlers and pkg utilities.
- Run subset: `go test ./services/auth/... -run TestLogin -v`.

## Commit & PR Guidelines
- Commits: imperative, scoped messages; recommended Conventional Commits (`feat:`, `fix:`, `chore:`, `refactor:`).
- PRs must include:
  - Summary of changes and rationale.
  - Linked issue (if any) and test plan/commands.
  - Screenshots/log snippets for API changes (e.g., `POST /auth/login`).

## Security & Config Tips
- Do NOT commit secrets. Keep real creds out of `config.yaml`; use env vars or an untracked override.
- Copy `config.example.yaml` and adjust locally; commit example-only changes.
- Rotate keys under `services/auth/internal/{jwt,keys}` outside of VCS; document how to reproduce locally.
- Review logs in `logs/` but avoid committing large or sensitive files.

