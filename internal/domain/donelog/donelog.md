# DONELOG Aggregate

## 役割
- 1 件の DONELOG（何を何件やったか）を完全な状態で保持する Aggregate Root。
- `id`, `title`, `trackID`, `categoryID`, `count`, `occurredOn` を VO として保持し、VO の不変条件によって整合性を担保する。
- Track/Category など他集約との結合は ID 参照のみ。表示名等は Query 側で解決する。

## 操作
- `NewDoneLog` で必須 VO を全て受け取り、ゼロ値を拒否する。
- `Update` で Title/Category/Count/OccurredOn を一括更新し、VO 経由で常にバリデーション後の値のみを保持する。

## Command/Query との関係
- Command 側 Application サービスから DoneLogRepository を通して永続化・復元され、トランザクション境界を定義する。
- Query 側では DONELOG から派生したプロジェクション（一覧、LOGSUMMARY 等）を利用し、Aggregate を直接返さない。

## TODO
- Repository インターフェースとインフラ実装（PostgreSQL 等）。
- ドメインイベント（DoneLogRecorded など）の定義と発行箇所。
