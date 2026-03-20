# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

repo-hand-off は、開発の「状態」（コード・計画・意図・学習・現況）を保存・共有し、人やAI間でシームレスに開発を引き継げるようにするCLIツール。詳細は VISION.md を参照。

## Status

プロジェクトは初期段階。VISION.md のみ存在し、実装コードはまだない。

## Key Documents

- `VISION.md` — プロジェクトの目的・設計原則・スコープ
- `PLAN.md` — 実装計画（作成予定）
- `HANDOFF.md` — 現在の開発状態スナップショット（本ツールが生成する成果物）

## Design Principles

- **状態ファースト**: コードではなく「状態」を中心に扱う
- **AIネイティブ**: AIにそのまま渡せる形式を前提とする
- **CLI中心**: 軽量で、どこでも使える
- **非侵襲**: 既存のGit/開発フローを壊さない
