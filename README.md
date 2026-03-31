# repo-hand-off

> **This project is archived and no longer maintained.**
> The approach of a standalone CLI tool turned out to be less practical than implementing handoff as a Claude Code skill with a stop hook. This repository is kept for reference only.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[日本語](README_ja.md)

A CLI tool that captures and shares the "state" of development — enabling seamless handoffs between humans and AI agents.

## Why Archived?

The core idea — capturing project state for handoffs — is valid, but building it as a separate CLI tool adds unnecessary friction. A better approach is to implement handoff as a **Claude Code skill** triggered by a **stop hook**, which integrates directly into the development workflow without requiring a separate tool.

## What This Was

`handoff` was a Go CLI that could:

- **`handoff export`** — Capture project state (Git info, project files, directory structure) into `HANDOFF.md`
- **`handoff prompt`** — Generate AI-ready context prompts (Markdown/XML)
- **`handoff import`** — Restore project files from a `HANDOFF.md`
- **`handoff diff`** — Compare saved state against current files
- **`handoff doctor`** — Diagnose handoff readiness (with `--strict` for CI)

## License

MIT
