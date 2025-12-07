# 🧱 DONELOGドメイン 初期設計まとめ（Draft）

## 🎯 ドメインの目的

このシステムは、
ユーザーが読んだ記事の数や、学習した問題の数など、
**「何かしらやったことの数」(count)** を記録し、
**グラフ・月ごとの比較・推移** といった形で振り返れるようにするためのドメインである。

---

# 🏗 ユビキタス言語

| 用語             | 意味                                    |
| -------------- | ------------------------------------- |
| **DONELOG**    | 1件の記録。「何を何件やったか」を表すエンティティ             |
| **Track**      | DONELOG が「何についての記録か」を示す対象（資格、本、テーマなど） |
| **TrackId**    | Track を一意に示す ID                       |
| **Category**   | DONELOG を分類する軸（例：資格、読書、技術学習…）         |
| **LOGSUMMARY** | DONELOG を元に算出される「集計結果」（期間内合計、推移など）    |

---

# 🔖 DONELOG（Entity / Aggregate Root）

```text
DONELOG {
    id: DONELOGId        // この記録固有のID
    title: Title         // 記録のタイトル
    trackId: TrackId     // 何についての記録か（基礎情報試験 / Clean Architectureなど）
    category: Category   // フィルタ用の分類
    count: Count         // 何件やったか
    occurredOn: Date     // いつの出来事か（集計の基準）
}
```

### DONELOG の役割

* 1件分の「やったこと」を表す
* count や occurredOn によって、期間集計・月比較・推移の材料になる
* 集計結果(LOGSUMMARY)の元データとなる
* Track / Category など関連 VO の整合性を 1 集約内で保証する役割を持つ
* Command 側では DoneLogRepository を通じて永続化・再構築され、アプリケーションサービスがトランザクション境界を管理する

---

# 📦 Value Object（VO）候補

DONELOG の属性のうち、下記はすべて VO として扱える：

* **DONELOGId**
* **Title**
* **TrackId**
* **Category**
* **Count**
* **Date（OccurredOn）**

## DONELOGId

- フォーマット：ULID などの時系列ソートが可能な 26 文字英数字。
- 生成：サーバー側で発行し、クライアント入力は不可。
- 目的：集約内の一意性と作成順の推測を両立する。

## Title

- UTF-8 文字列で 1〜120 文字。改行や制御文字は拒否する。
- 入力時に前後の空白をトリムし、空文字ならエラー。
- 表示用途のため、XSS 対策と絵文字許可範囲はアプリケーション層で判断する。

## TrackId

- Track 集約の ID。`track_{slug}` のようにスネークケースで統一する。
- 英数字とハイフンのみ許可し、先頭は英字で始まる。
- DONELOG 作成時に存在確認と Active 状態をセットで検証する。

## Category（CategoryId）

- Category 集約の ID。TrackId と同じ命名規則を採用する（例：`cat_{slug}`）。
- Track に紐づく Category のみ選択可。非アクティブ ID は拒否する。
- 未分類を許容する場合は `null` または特別な ID（例：`cat_none`）で扱う。

## Count

- 1 以上の整数。通常は 32bit 範囲に収め、1〜9999 を目安にする。
- VO 内で加算・減算を提供し、0 以下になりそうな操作は例外を投げる。
- 小数や時間量が必要なら別 VO（Duration など）を導入する。

## Date（OccurredOn）

- `YYYY-MM-DD` 形式の日付。ユーザーのローカルタイムゾーン基準。
- DONELOG 登録日ではなく「実施した日」。未来日を禁止するかはドメインで決める。
- 並び順・集計期間の境界判断はこの VO を基準に行う。

### VO の役割

- すべての値をコンストラクタ経由で生成し、必ず不変条件を満たした状態だけを外部に渡す。
- 文字列/ID/数値/日付ごとに責務を分け、ドメイン層のバリデーションが散らばらないようにする。
- VO 外からは `String()` や `Int()` などの読み取りメソッドのみ公開し、直接の書き換えを防ぐ。
- 実装では private フィールドを持たせ、必ず `NewXXX` を通る構造にしている（Go 実装済み）。

---

# 📊 LOGSUMMARY（Value Object）

DONELOG の集合から導出される「集計結果」。

```text
LOGSUMMARY {
    category: CategoryId?       // nullなら全カテゴリ
    period: Period              // いつからいつまで
    totalCount: Count           // 期間内合計
    points: List<SummaryPoint>  // 推移を見るためのデータ
}

SummaryPoint {
    label: string               // 日 / 週 / 月などの軸に使う値
    count: Count                // その区切り単位の合計
}
```

---

# 📈 集計ユースケース（LOGSUMMARY が扱う要求）

1. **ある期間の日ごとの合計 count を見たい**
   → 推移グラフ用

2. **ある期間の月ごとの total count の変化を見たい**
   → 月比較・成長推移の把握

3. **今月と先月の count の合計をカテゴリー別に比較したい**
   → Category フィルター前提

---

# 🧪 集計ユースケース詳細

## 1. 期間 × 日別推移

- 入力：`categoryId?`, `period (startDate, endDate)`
- 出力：`LOGSUMMARY(period, points=日単位, totalCount)`
- ビジネスルール
  - `period` は最大 90 日程度を上限にし、長期は月別に切り替える。
  - 日付ごとのブランクも `count=0` で出力し、グラフが途切れないようにする。
  - Category 未指定の場合は全 DONELOG を対象。
- エッジケース
  - Track / Category が非アクティブでも期間内の DONELOG は集計対象。
  - 同じ日に複数 DONELOG がある場合は `count` を合算する。

## 2. 期間 × 月別推移

- 入力：`categoryId?`, `period (startMonth, endMonth)`
- 出力：`LOGSUMMARY(points=月単位, totalCount)`
- ビジネスルール
  - 期間は最長 24 ヶ月程度を目安にする。
  - 月ラベルは `YYYY-MM`。月途中のデータもその月に丸める。
  - 前月比など比率計算はアプリケーション層で行う。
- エッジケース
  - データが存在しない月も `count=0` で埋める。

## 3. 今月 vs 先月 Category 比較

- 入力：`period (thisMonthRange, lastMonthRange)`
- 出力：`List<{categoryId, thisMonthCount, lastMonthCount, diff}>`
- ビジネスルール
  - 指定 Category のみ比較するか、全 Category を `SortOrder` 順に表示する。
  - 比較対象の月は「当月」「前月」と固定だが、将来は任意の 2 期間を比較できる拡張を想定。
- エッジケース
  - 今月または先月どちらかのみ存在する Category は片方の `count` を 0 とする。
  - 非アクティブ Category でも過去データがあれば比較対象に含める。

---

# ⚙ 集計の責務（Domain Service想定）

LOGSUMMARY 自体は「結果」だけを持ち、
計算は Domain Service が行う想定。

例：

```text
LogSummaryService {
    summarizeByDay(categoryId, period): LOGSUMMARY
    summarizeByMonth(categoryId, period): LOGSUMMARY
}
```

---

# 📐 集約境界（初期案）

## DONELOG 集約

- Aggregate Root。単独で完結した 1 件の事実を表現する。
- Track / Category は ID 参照のみで、表示用の値は保持しない。
- 不変条件の例
  - Count > 0。
  - TrackId / Category は存在している集約のみ参照する。
  - 同じ TrackId + OccurredOn の重複登録をどう扱うか（禁止 or 許容）を明示しておく。
- Track / Category の表示名が変わっても DONELOG の履歴は壊れず、参照時に最新情報を解決する。

## Track 集約

- DONELOG が「何を記録したか」を紐付ける対象。
- 属性例：TrackId, Name, DefaultCategoryId?, Color, SortOrder, ActiveFlag。
- 不変条件の例
  - Name はユニーク。
  - ActiveFlag = false の Track には新規 DONELOG を紐付けない。
- Category への参照は片方向（Category → Track には逆参照を持たせない）。
- DONELOG が存在する Track は削除せず、アーカイブ（非アクティブ化）のみ許可する。

## Category 集約

- DONELOG / Track を横断的に分類する軸。
- 属性例：CategoryId, Name, Color, SortOrder, ActiveFlag。
- 不変条件の例
  - Name はユニーク。
  - 存在しない Category を DONELOG / Track が参照できない。
- 削除ポリシー：紐づく DONELOG がある場合は削除せず非アクティブ化で代替する。
- Track との関係は一方向。Track が Category を参照しても Category 側は Track を直接知らない。

## 集約間の整合性

- DONELOG 作成・編集時に Track / Category の存在確認と Active 状態を検証する。
- Track / Category 変更時の影響範囲はアプリケーションサービスで調停し、必要に応じて同一トランザクションで扱う。
- Track → Category は一方向参照だが、Category 切替による DONELOG 更新はドメインサービスで吸収できるように設計する。
- LOGSUMMARY は DONELOG 集合の読み取り専用 VO で、副作用を発生させない。

## LOGSUMMARY

- DONELOG から派生する集計結果を表す Value Object。
- Aggregate ではなくリードモデル扱い。計算は `LogSummaryService` などのドメインサービスが担当し、DONELOG に副作用を与えない。
- LogSummaryService は DONELOG 集約を横断して日次・月次などの集計を計算する Domain Service であり、集約に属さない派生ロジックを担う。
- Query 側では LogSummary を DTO に変換して API 経由で返す。Command 側の更新と切り離しておける。

---

# 🧩 現在の状態まとめ

* DONELOG：形が確定
* LOGSUMMARY：目的・構造が明確
* VO：どれをVOにするか認識済み
* 集計ユースケース：固まった
* 集約境界：初期案として妥当

---

# 📡 アプリケーションレイヤーまとめ

| 分類 | 内容 |
| ---- | ---- |
| Command API | DONELOG/Track/Category の作成・更新・削除。ドメインイベントを発行。 |
| Query API | DONELOG 一覧、Track/Category カタログ、日次/月次/カテゴリ比較集計。 |
| CQRS 方針 | Command と Query を別パッケージ + プロジェクションで分離。 |
| TODO | Domain/UseCase/Projection 実装、Query モデルのストレージ設計、監視整備。 |

---

# 🛠 Go Backend 実装方針（DDD）

## レイヤー構成

1. **Domain（`internal/domain/donelog`）**
   - エンティティ / VO / ドメインサービス（`LogSummaryService` など）。
   - リポジトリインターフェース（`DoneLogRepository`, `TrackRepository`, `CategoryRepository`）。
   - ドメインイベントを定義する場合もここに置く。
2. **Application（`internal/app`）**
   - UseCase / ApplicationService をコマンド・クエリ単位で実装。
   - トランザクション境界を握り、リポジトリやドメインサービスをオーケストレーション。
   - DTO や入力バリデーションロジックをここでまとめる。
3. **Interface / Adapter（`internal/interface/http`, `internal/interface/presenter`）**
   - HTTP ハンドラ、GraphQL Resolver など I/O ごとのアダプタ。
   - Application レイヤーの入出力 DTO とのマッピング。
4. **Infrastructure（`internal/infrastructure/...`）**
   - DB 実装（`postgres`, `sqlite`）、外部 API、時計 (`Clock`)、ID 生成器など。
   - Domain レイヤーのインターフェースを満たす。
5. **Cmd（`cmd/api`）**
   - DI 初期化、設定読み込み、サーバ起動。

## ディレクトリ案

```text
cmd/
  api/
    main.go
internal/
  domain/
    donelog/
      donelog.go          // Aggregate / Entity
      value_object.go     // Title, Count, etc.
      log_summary.go      // Value Object
      repository.go       // Repository interfaces
      service.go          // LogSummaryService
  app/
    donelog/
      command/
        create_donelog.go
        update_donelog.go
      query/
        list_tracks.go
        summarize_by_day.go
  interface/
    http/
      handler.go
      middleware.go
    presenter/
      donelog_response.go
  infrastructure/
    persistence/
      postgres/
        donelog_repository.go
    clock/
      system_clock.go
    id/
      ulid_generator.go
pkg/
  config/
  logger/
```

## リポジトリ / サービス契約

- `DoneLogRepository`
  - `Save(ctx, *DONELOG) error`
  - `FindByID(ctx, DONELOGId) (*DONELOG, error)`
  - `ListByPeriod(ctx, Period, TrackId?, CategoryId?) ([]*DONELOG, error)`
- `TrackRepository`
  - `FindActiveByID(ctx, TrackId) (*Track, error)`
  - `ListActive(ctx) ([]*Track, error)`
- `CategoryRepository`
  - `FindActiveByID(ctx, CategoryId) (*Category, error)`
  - `ListActive(ctx) ([]*Category, error)`
- `LogSummaryService`
  - 既存の `summarizeByDay`, `summarizeByMonth` を実装し、Application 層から利用。

## Go 実装 TODO

1. Domain パッケージに DONELOG / VO / リポジトリ interface を定義する。
2. Application レイヤーで Create / Update / Delete / Summaries の UseCase を書き下ろす。
3. Infrastructure 層で DB スキーマを決め、リポジトリ実装を用意。
4. HTTP Handler から Application を呼び出すアダプタを準備（REST なら `/donelogs`, `/tracks`, `/summaries` など）。
5. DI / 設定読み込みを `cmd/api` で構築。テストは Application / Domain をユニットテスト、Interface は E2E でカバー。

---

# 🧾 Command API（書き込み系）

## Create DONELOG

- **Endpoint**: `POST /api/donelogs`
- **Request**

```json
{
  "title": "Clean Architecture 1章",
  "trackId": "track_clean_architecture",
  "categoryId": "cat_reading",
  "count": 3,
  "occurredOn": "2024-05-01"
}
```

- **Validation**
  - Title/TrackId/CategoryId/Count/OccurredOn は VO のルールを順守。
  - Track / Category が存在し Active であることを確認。
- **Behavior**
  - DONELOG を生成し `DoneLogRepository.Save`。
  - `DoneLogRecorded` イベントを発行（Projection 更新用）。
- **Response**: `201 Created`, body に新規 DONELOG の ID を返す。

## Update DONELOG

- **Endpoint**: `PUT /api/donelogs/{id}`
- **Request**: Create と同じフィールド。部分更新ではなく全フィールド更新を想定。
- **Validation / Behavior**
  - 指定 ID の DONELOG が存在することを確認。
  - VO で再検証し、変更差分があれば集約を更新。
  - `DoneLogUpdated` イベントを発行。
- **Response**: `200 OK` or `204 No Content`。

## Delete DONELOG

- **Endpoint**: `DELETE /api/donelogs/{id}`
- **Behavior**
  - 集約を取得し、削除ポリシー（論理削除 or 物理削除）に従う。
  - Projection との整合性のため `DoneLogDeleted` イベントを発行。
- **Response**: `204 No Content`。

## Manage Track

- **Endpoint**
  - `POST /api/tracks`
  - `PUT /api/tracks/{id}`
  - `PATCH /api/tracks/{id}/archive` などのアクティブ状態変更
- **Request**

```json
{
  "id": "track_clean_architecture",
  "name": "Clean Architecture",
  "defaultCategoryId": "cat_reading",
  "color": "#00AA88",
  "sortOrder": 100
}
```

- **Behavior**
  - Track 集約を作成/更新し、Active 状態を管理。
  - Track 変更に応じた DONELOG の整合性はドメインサービスで評価（例：DefaultCategory 変更時）。
- **Response**: `201` or `200`。

## Manage Category

- **Endpoint**
  - `POST /api/categories`
  - `PUT /api/categories/{id}`
  - `PATCH /api/categories/{id}/archive`
- **Request**: Track と同様。`name`, `color`, `sortOrder`, `active` など。
- **Behavior**
  - Category 集約を作成/更新。非アクティブ化時に DONELOG との参照を調整（必要なら `uncategorized` へ再マップ）。
- **Response**: `201` or `200`。

## Command エラーハンドリング

- 400: VO バリデーションエラー。
- 404: Track / Category / DONELOG が見つからない。
- 409: 処理中の競合（同一 Track + OccurredOn のユニーク制約違反など）。
- 500: 予期しないエラー。Application 層でログを取り、クライアントにはトレーサブルなエラーコードを返す。

---

# ♻️ CQRS 導入方針

## なぜ CQRS にするのか

- DONELOG の書き込み（コマンド）は単純な事実の保存が中心だが、読み取り（クエリ）はグラフや期間比較など重い集計を要求するため、責務を分離することで設計がクリアになる。
- Command 側は厳格なドメインモデルとトランザクション制御に集中でき、Query 側は読みやすいデータ構造（集計済みビューやキャッシュ）を自由に使える。
- 将来的に読み取りのスケールが増えた場合、読み取り専用 DB やキャッシュ層を追加する判断が容易になる。
- 変更履歴や監査ログを Command 側で記録しつつ、Query 側で非同期に投影することで、ユーザー体験とデータ整合性のバランスを取りやすい。

## Command モデル（書き込み側）

- 対象：Create / Update / Delete DONELOG、Track / Category の管理。
- レイヤー構成：Domain（Aggregate, VO）、Application（Command UseCase）、Infrastructure（`DoneLogRepository` など）。
- ストレージ：`donelog` や `track` テーブルを正規化し、トランザクション境界を明確にする。
- API：REST なら POST/PUT/DELETE、GraphQL なら Mutation として公開。
- ドメインイベント：DONELOG 作成/更新時に `DoneLogRecorded` などのイベントを発行し、Query 側の投影処理をトリガーする。

## Query モデル（読み取り側）

- 対象：一覧取得、Track/Category カタログ、`LOGSUMMARY` 系集計。
- 読み取り専用のプロジェクションを用意する。
  - 例：`donelog_summary_daily`, `donelog_summary_monthly`, `category_totals` といったマテリアライズドビュー / テーブル。
  - 実装初期は 1 つの RDB でも問題ないが、将来的にリードレプリカや専用ストレージ（Elasticsearch, BigQuery）に切り出せるようインターフェースを分離する。
- アプリケーション層では Command 側とは別パッケージ（`internal/app/donelog/query`）を維持し、DTO とプレゼンターも Query 用に分ける。
- API：GET または GraphQL Query。Query 側はキャッシュ（Redis 等）を挟んでも Command に影響しない。

## 投影（Projection）戦略

- Command 側で発行したイベントを購読し、Query モデルを更新する投影処理を `internal/app/projection` などに配置する。
- 初期段階では同一トランザクション内で同期更新（Command → Query テーブル更新）でもよいが、イベントキュー（e.g. PostgreSQL NOTIFY, メッセージブローカー）を挟む設計も視野に入れる。
- 投影処理は冪等に実装し、失敗時は再実行できるようにする。

## Go 実装 TODO（CQRS 対応）

1. Command / Query 用に Application パッケージを明確に分離し、インターフェースを交差させない。
2. Query 専用の読み取りモデル（構造体, DTO, SQL）を定義する。Domain Aggregate をそのまま返さない。
3. DONELOG 作成/更新で発火するドメインイベントと、これを処理する Projection 用 Application Service を設計する。
4. Infrastructure 層で Projection の永続化先（テーブル or ビュー）を用意し、必要ならマイグレーションに追加する。
5. 運用面：Query モデルが Command モデルと最終的に一致しているか確認するヘルスチェックやメトリクスを整える。

---

# 🔍 Query API（読み取り系）

## GET /api/donelogs

- **用途**: 最新の DONELOG 一覧をページング取得。
- **Query Params**
  - `trackId?`, `categoryId?`, `occurredOnFrom?`, `occurredOnTo?`
  - `page`, `limit`（デフォルト 20）
- **Response**

```json
{
  "items": [
    {
      "id": "01HYR1X5C9XM9P6H7K71M9QAHX",
      "title": "Clean Architecture 1章",
      "track": { "id": "track_clean_architecture", "name": "Clean Architecture" },
      "category": { "id": "cat_reading", "name": "読書" },
      "count": 3,
      "occurredOn": "2024-05-01"
    }
  ],
  "page": 1,
  "totalCount": 123
}
```

- **データソース**: Command テーブルを直接参照してもよいが、Query 専用のビュー（JOIN 済み）を利用するとクエリが簡潔になる。

## GET /api/tracks, /api/categories

- **用途**: UI 初期表示用のカタログ。Active のみ返す。
- **Response**: `[{ "id": "...", "name": "...", "color": "...", "sortOrder": 100 }]`
- **キャッシュ**: ETag/Cache-Control、またはアプリ内メモリキャッシュを使用。

## GET /api/summaries/daily

- **Query Params**
  - `categoryId?`, `startDate`, `endDate`
- **Response**

```json
{
  "period": { "startDate": "2024-05-01", "endDate": "2024-05-31" },
  "totalCount": 42,
  "points": [
    { "label": "2024-05-01", "count": 3 },
    { "label": "2024-05-02", "count": 0 }
  ]
}
```

- **データソース**: `donelog_summary_daily` プロジェクションまたは集計ビュー。

## GET /api/summaries/monthly

- **Query Params**: `categoryId?`, `startMonth`, `endMonth`
- **Response**: `{"points":[{"label":"2024-05","count":42}], "totalCount": 120}`
- **データソース**: `donelog_summary_monthly` プロジェクション。

## GET /api/summaries/category-comparison

- **Query Params**: `month`（当月）、`previousMonth?`（省略時は month-1）
- **Response**

```json
{
  "items": [
    {
      "category": { "id": "cat_reading", "name": "読書" },
      "thisMonthCount": 20,
      "lastMonthCount": 15,
      "diff": 5
    }
  ]
}
```

- **データソース**: Category 集計プロジェクション。SortOrder 順に並べる。

## GraphQL との両立

- REST に加えて GraphQL を導入する場合は Query 用 Resolver を `internal/interface/graphql/query` に追加し、Application Query を呼び出すだけにする。
- Command は Mutation で公開し、CQRS の概念が崩れないよう Query Resolver から Command UseCase を呼ばない。

---
