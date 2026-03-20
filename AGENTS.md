# AGENTS.md

This file provides guidance to AI coding agents (Codex, Gemini, etc.) when working with code in this repository.

## Project Overview

repo-hand-off is a CLI tool that captures and shares development "state" (code, plans, intent, lessons, current status) for seamless handoffs between humans and AI agents. See VISION.md for details.

## Status

All MVP phases (1-8) complete. Export, prompt, import, diff commands are functional. AI config files (CLAUDE.md, AGENTS.md, GEMINI.md) are auto-detected.

## Key Documents

- `VISION.md` — Project purpose, design principles, scope
- `PLAN.md` — Implementation plan and progress
- `HANDOFF.md` — Development state snapshot (generated artifact, in .gitignore)

## Architecture

```
cmd/handoff/main.go   — Entrypoint
cmd/                   — Cobra command definitions (root, export, prompt, import, diff)
internal/collector/    — Git info, file, and directory structure collection
internal/renderer/     — HANDOFF.md / AI prompt generation
internal/parser/       — HANDOFF.md parsing
internal/differ/       — Section-level diff comparison
internal/config/       — Configuration management (.handoff.yaml)
```

## Development Commands

```bash
go build -o handoff ./cmd/handoff  # Build
go test ./...                      # Run all tests
go vet ./...                       # Static analysis
./handoff export                   # Generate HANDOFF.md
./handoff prompt                   # Output AI-ready prompt to stdout
./handoff prompt --format xml      # XML format output
./handoff import                   # Restore files from HANDOFF.md
./handoff diff                     # Compare HANDOFF.md with current state
```

## Design Principles

- **State-first**: Focuses on development "state", not just code
- **AI-native**: Output formats designed for direct AI consumption
- **CLI-centric**: Lightweight, usable anywhere
- **Non-invasive**: Never disrupts existing Git or development workflows
