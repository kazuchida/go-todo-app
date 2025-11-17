# go-todo-app

軽量な Todo アプリケーションのサンプル（Go）です。
このリポジトリはクリーンアーキテクチャの考え方に沿って、
ドメイン・ユースケース・インフラストラクチャを分離して実装しています。

**主なディレクトリ構成**
- `internal/domain` : エンティティとビジネスルール（`Todo`）
- `internal/usecase` : ユースケース（アプリケーションの振る舞い）
- `internal/adapter/repository` : `TodoRepository` の SQL 実装
- `internal/infrastructure/database` : Postgres 接続とスキーマ初期化
- `cmd/api` : HTTP サーバのエントリポイント

**Prerequisites**
- Go 1.20 以上がインストールされていること
- Docker（ローカルで Postgres を使う場合）

**Quick Start (ローカル実行)**

1) Postgres を Docker で起動（デフォルト設定）

簡易で起動する場合:

```sh
docker run --name todoapp \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=todoapp \
  -p 5432:5432 -d postgres:15
```

環境変数を明示して DB 名などを設定する場合（推奨）:

```sh
docker run --name todo-postgres \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=todoapp \
	-p 5432:5432 -d postgres:15
```

2) 依存パッケージを取得してモジュールを整える

```sh
cd /path/to/go-todo-app
go get github.com/lib/pq@latest
go mod tidy
```

3) アプリを起動

```sh
go run ./cmd/api
```

`cmd/api/main.go` はデフォルトで以下の DB 設定を使います（必要に応じてソースを書き換えてください）：

- host: `localhost`
- port: `5432`
- user: `postgres`
- password: `postgres`
- dbname: `todoapp`

起動時にスキーマ初期化 (`internal/infrastructure/database.InitSchema`) を実行します。

**データベースへ直接スキーマを適用する（任意）**
アプリ側で自動初期化しない場合は `psql` で手動作成できます。`postgres.go` にあるスキーマ定義は以下のとおりです。

```sql
CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

例:
```sh
psql -h localhost -U postgres -d todoapp -c "CREATE TABLE IF NOT EXISTS todos (...);"
```

**テスト**

```sh
go test ./...
```

**設計メモ**
- クリーンアーキテクチャに基づき、`usecase` 層は `TodoRepository` インターフェースに依存します。実装（Postgres）は `internal/adapter/repository` にあります。
- 各メソッドは `context.Context` を受け取り、ハンドラでのキャンセル/タイムアウトが DB 操作に伝搬されます。

**補足**
- DB 設定やポート、パスワードは本番では環境変数やシークレットマネージャで安全に管理してください。
- 他の Postgres ドライバ（例: `pgx`）に差し替える場合は、`internal/infrastructure/database` と `go.mod` を更新してください。

---

**API 一覧**

以下はこのアプリが提供する主要なエンドポイント一覧です。`http://localhost:8080` でサーバーが起動している前提です。

| メソッド | パス | リクエストボディ (JSON) | 説明 | 期待ステータス |
|---|---|---|---:|---:|
| POST | `/api/todos` | `{ "title": "...", "description": "..." }` | 新しい Todo を作成する | `201 Created` / `200 OK` |
| GET | `/api/todos` | - | Todo の一覧を取得する | `200 OK` |
| GET | `/api/todos/{id}` | - | 指定 ID の Todo を取得する | `200 OK` / `404 Not Found` |
| PUT | `/api/todos/{id}` | `{ "title": "...", "description": "...", "completed": true|false }` | Todo を更新する（完全更新） | `200 OK` / `404 Not Found` |
| PATCH | `/api/todos/{id}/toggle` | - | Todo の完了状態を反転する | `200 OK` / `404 Not Found` |
| DELETE | `/api/todos/{id}` | - | Todo を削除する | `204 No Content` / `200 OK` / `404 Not Found` |

簡単な `curl` 例（サーバーが `localhost:8080` で起動している想定）：

```sh
# Create
curl -i -X POST http://localhost:8080/api/todos \
	-H "Content-Type: application/json" \
	-d '{"title":"買い物","description":"牛乳を買う"}'

# List
curl -i http://localhost:8080/api/todos

# Get
curl -i http://localhost:8080/api/todos/1

# Update
curl -i -X PUT http://localhost:8080/api/todos/1 \
	-H "Content-Type: application/json" \
	-d '{"title":"買い物（更新）","description":"牛乳、パン","completed":false}'

# Toggle
curl -i -X PATCH http://localhost:8080/api/todos/1/toggle

# Delete
curl -i -X DELETE http://localhost:8080/api/todos/1
```