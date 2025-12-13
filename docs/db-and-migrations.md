# DB 選定とマイグレーション方針

## DB: Supabase (Postgres)
- 採用理由: Postgres 互換で機能十分、無料枠で開始でき、将来有料プランへスケール可能。
- 接続: Supabase のサービスロール (service key) で発行される `postgres://...` を `DATABASE_URL` として利用する。
- ローカル開発: 必要に応じて `supabase start` でローカル Postgres を立て、同じスキーマを適用して動作確認できる。

## マイグレーション管理: goose
- 選定理由: Go プロジェクトで軽量に扱えて学習コストが低い、SQL ファイルで `up/down` を管理できる。
- ディレクトリ案: `migrations/` 直下に `0001_create_donelogs.sql` のように番号付きで配置。
- 運用イメージ:
  - 追加: `goose -dir ./migrations postgres "$DATABASE_URL" up`
  - 巻き戻し: `goose -dir ./migrations postgres "$DATABASE_URL" down`
  - 未適用のものだけを順番に適用する。CI でも同じコマンドを流す。
- 接続文字列: Supabase の service key から得られる `postgres://...` を `DATABASE_URL` として渡す。ローカルで Supabase を動かす場合は `supabase start` 後に発行される接続文字列を利用する。

## 今後の実装タスク
- `migrations/` 配下に `tracks`, `categories`, `donelogs` のスキーマを定義する `0001` を追加する。
- `internal/infrastructure/persistence/postgres/` に `DoneLogRepository` などの実装を作成する（`pgx` ドライバ想定）。
- `cmd/api` で `DATABASE_URL` を読み取り、DB 接続/DI を配線する。
- Supabase/DB に対して goose を流す手順を README か make タスクで共有する。
