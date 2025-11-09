# rest-grpc-todo-demo

一つの app に rest と grpc の２つの口を作成する。

## 前提

- Go 1.20 以上

## 1. セットアップ

```bash

go mod tidy

```

## 2.1. rest

```bash

go run ./cmd/rest-server

```

別ターミナル

```bash

# 一覧（最初は空）
curl http://localhost:8000/tasks

# 作成
curl -X POST http://localhost:8000/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"hoge"}'

# 更新(done=true)
curl -X PATCH http://localhost:8000/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"done":true}'

```

## 2.2. gRPC

```bash

go run ./cmd/grpc-server

```

Postman からの操作

1. Postman 起動
2. New
3. gRPC request
4. アドレスを localhost:50051 に指定
5. proto/task.proto を import
6. メソッドから CreateTask を指定(多分この段階で候補に出てくる)
7. 「サンプルメッセージを使用」を指定
8. 「呼び出す」
