<!-- BEGIN:nextjs-agent-rules -->
# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` before writing any code. Heed deprecation notices.
<!-- END:nextjs-agent-rules -->

# Project Continuity Rules

- Before any substantial change, read `docs/project-continuity.md`, `docs/project-map.md`, `docs/current-state.md`, and `docs/development-rules.md`.
- Treat `idc-finance/` as the only active workspace unless the user explicitly says otherwise.
- Do not rely on editor-local memory, temporary chat context, or private notes for project state. Persist reusable state into tracked repo docs.
- For non-trivial business logic, write comments that explain business intent, external API assumptions, state transitions, and compatibility constraints.
- When architecture, external integration behavior, workflow, or next-step priorities change, update the continuity docs in the same change set.
- Never commit secrets, credentials, local machine paths, or editor-private state. Keep those in local env files that are ignored by git.
