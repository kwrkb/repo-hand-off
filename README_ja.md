# repo-hand-off

> **このプロジェクトはアーカイブされており、メンテナンスは行われていません。**
> 単体の CLI ツールとして作るより、Claude Code のスキル + stop hook として実装する方が実用的であることが分かりました。このリポジトリは参考資料として残しています。

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[English](README.md)

開発の「状態」を保存・共有し、人やAI間でシームレスに引き継ぐCLIツール。

## アーカイブの理由

コアのアイデア（プロジェクトの状態をキャプチャして引き継ぐ）は妥当だが、独立した CLI ツールとして構築するのは余計な摩擦を生む。**Claude Code のスキル**として実装し、**stop hook** でトリガーすれば、開発フローに直接統合でき、別ツールを導入する必要がない。

## このツールでできたこと

- **`handoff export`** — プロジェクト状態（Git情報・プロジェクトファイル・ディレクトリ構造）を `HANDOFF.md` にキャプチャ
- **`handoff prompt`** — AI向けコンテキストプロンプトを生成（Markdown/XML）
- **`handoff import`** — `HANDOFF.md` からプロジェクトファイルを復元
- **`handoff diff`** — 保存した状態と現在のファイルを比較
- **`handoff doctor`** — handoff 品質の診断（`--strict` で CI 連携）

## License

MIT
