<!-- BEGIN:nextjs-agent-rules -->
# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` before writing any code. Heed deprecation notices.
<!-- END:nextjs-agent-rules -->

# Project Continuity Rules

- Before any substantial change, read `docs/project-continuity.md`, `docs/project-map.md`, `docs/current-state.md`, `docs/environment-setup.md`, `docs/development-rules.md`, `docs/cross-device-collaboration.md`, and `docs/active-handoff.md`.
- Treat `idc-finance/` as the only active workspace unless the user explicitly says otherwise.
- Do not rely on editor-local memory, temporary chat context, or private notes for project state. Persist reusable state into tracked repo docs.
- For non-trivial business logic, write comments that explain business intent, external API assumptions, state transitions, and compatibility constraints.
- When architecture, external integration behavior, workflow, or next-step priorities change, update the continuity docs in the same change set.
- Never commit secrets, credentials, local machine paths, or editor-private state. Keep those in local env files that are ignored by git.

# Cross-Client Execution Rules

- This repository is designed to support local IDEs, Codex desktop, Codex web, remote dev environments, and future collaborators with the same recovery path.
- Cross-device continuity only counts if the state is inside git-tracked files and pushed to the remote repository.
- If work is only described in chat or exists only in unpushed local changes, treat it as non-authoritative.
- Before starting work, check `git status` and read the latest handoff context from `docs/active-handoff.md`.
- Before ending work, update `docs/active-handoff.md` whenever the current focus, blockers, environment expectations, or next recommended task changed.
- If the change touches environment variables, also update `idc-finance/.env.example` and the setup docs.
- If the change touches schema, migrations, or demo data, also update the relevant migration/seed docs and note the verification status in `docs/active-handoff.md`.
- If verification could not be completed, record exactly what was not verified and why.

# Default Completion Standard

- End each meaningful work block in a state that another client can continue from after only `git pull`.
- Keep the repository free of stale references to removed systems, old backends, local-only machine setup, or deprecated implementation paths.
- Prefer a small number of authoritative docs over many overlapping notes; when a new process document is added, link it from the continuity entry points.
