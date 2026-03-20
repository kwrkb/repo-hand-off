# LESSONS.md

## 設計

- Markdown パーサーでセクション分割する際、ファイル内容に `##` ヘッダーが含まれると誤ってセクション境界と認識される。対策: Extra ファイルセクションに `## Extra: ` プレフィックスを付与し、既知セクション名またはプレフィックスのみを境界として認識する。`filepath.Ext` 等のヒューリスティックは `## v1.2` のような一般ヘッダーで誤認する
- XML レンダリングを共有する場合、文字列操作（`TrimPrefix`/`TrimSuffix`）でタグを剥がすのは脆い。共通の body 生成関数を抽出し、各呼び出し元がラッパーを付与する方式にする
- デフォルト値は一箇所（config パッケージ）で定義し、利用側で二重にフォールバックしない。`BuildDirTree` の `maxDepth <= 0` ガードは安全ネットとして残すが、正規のデフォルトは config が持つ

## セキュリティ

- `handoff import` のように外部ファイルからパスを取り込むコマンドでは、パストラバーサル（`../../outside.md`、`/tmp/out.md`）を防ぐ。`filepath.Abs` で解決後、`strings.HasPrefix(resolved, workDir+sep)` で検証する
- YAML 設定で glob パターンを受け取る場合、`filepath.Match` の `ErrBadPattern` をロード時にバリデーションする。不正パターンが実行時に静かに無視されるのを防ぐ

## Go

- `os.OpenFile` で `O_CREATE|O_EXCL` を使うと、ファイル存在チェックと作成をアトミックに行える。`os.Stat` → `os.WriteFile` の TOCTOU パターンを避ける
- `f.Close()` のエラーは必ずチェックする。書き込みバッファのフラッシュ失敗でデータロスの可能性がある
- map のイテレーション順序は非決定的。ユーザー向け出力やテストで安定した順序が必要な場合は `sort.Strings(keys)` でソートする

## テスト

- cobra の `PersistentPreRunE` で設定をロードする場合、cmd パッケージ内のグローバル変数（`cfg`, `workDir`）に依存する。単体テストより統合テスト（ビルド後の CLI 実行）で動作確認するのが実用的
