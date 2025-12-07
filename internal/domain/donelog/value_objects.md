# Value Objects (VO)

## 目的
- DONELOG 集約が扱う ID/文字列/数値/日付を型レベルで制約し、不正値を排除する。
- コンストラクタ (`NewXXX`) 経由でのみ生成することで、不変条件（Count > 0, Title 1〜120 文字など）を強制する。
- private フィールド ＋ 読み取りメソッド（`String()`, `Int()` など）により、外部からの直接変更を防ぐ。

## 主な VO

| VO | 役割と制約 |
| --- | --- |
| `DoneLogID` | ULID 形式 26 文字。サーバ生成のみ。 |
| `TrackID` / `CategoryID` | `track_{slug}`, `cat_{slug}` のように slug 形式。英数字＋`_-`、先頭は英字。 |
| `Title` | UTF-8 文字列、1〜120 文字。改行・制御文字不可。前後の空白はトリム。 |
| `Count` | 1 以上の整数。加減算は VO メソッドのみ。 |
| `OccurredOn` | `YYYY-MM-DD`。ユーザーのローカルタイムゾーン基準で、未来日可否はドメインで判断。 |
| `Period` | `OccurredOn` の閉区間。`Contains` 判定を提供。 |

## 実装メモ (Go)
- `internal/domain/donelog/value_objects.go` に実装。
- VO ごとに `String()` / `Int()` などの読み取り API を提供し、Aggregate やサービスから扱いやすくしている。
