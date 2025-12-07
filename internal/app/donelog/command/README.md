# Application Commands (DONELOG)

- `CreateDoneLog`, `UpdateDoneLog`, `DeleteDoneLog`: DONELOG 集約の基本的な作成/更新/削除ユースケース。
- 依存するリポジトリ: `DoneLogRepository`, `TrackRepository`, `CategoryRepository`。
- 入力 DTO（Command）でバリデーション後、Domain の VO/Entity へ変換する。
- 将来的に Track/Category 管理の Command もこのパッケージに追加する。
