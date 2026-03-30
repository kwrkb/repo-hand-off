# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

repo-hand-off は、開発の「状態」（コード・計画・意図・学習・現況）を保存・共有し、人やAI間でシームレスに開発を引き継げるようにするCLIツール。詳細は VISION.md を参照。

## Status

MVP 実装済み。`handoff export` / `handoff prompt` が動作する。Phase 4 の README.md 作成が残タスク。

## Key Documents

- `VISION.md` — プロジェクトの目的・設計原則・スコープ
- `PLAN.md` — 実装計画・進捗管理
- `HANDOFF.md` — 現在の開発状態スナップショット（本ツールが生成する成果物、.gitignore 対象）

## Architecture

```
main.go              — エントリポイント
cmd/                  — cobra コマンド定義（root, export, prompt, import, diff, doctor）
internal/collector/   — Git情報・ファイル・ディレクトリ構造・TODO件数・CI検出の収集
internal/renderer/    — HANDOFF.md / AI向けプロンプト / doctor出力の生成
internal/doctor/      — handoff品質の診断（Rule interface + 組み込みルール）
internal/config/      — 設定管理
```

## Development Commands

```bash
go build -o handoff ./cmd/handoff  # ビルド
go test ./...               # 全テスト実行
./handoff export            # HANDOFF.md 生成
./handoff prompt            # AI向けプロンプトを stdout 出力
./handoff prompt --format xml  # XML形式で出力
./handoff doctor            # handoff品質の診断
./handoff doctor --format json  # JSON形式で出力
./handoff doctor --strict   # Error検出時 exit 1（CI用）
```

## Design Principles

- **状態ファースト**: コードではなく「状態」を中心に扱う
- **AIネイティブ**: AIにそのまま渡せる形式を前提とする
- **CLI中心**: 軽量で、どこでも使える
- **非侵襲**: 既存のGit/開発フローを壊さない
