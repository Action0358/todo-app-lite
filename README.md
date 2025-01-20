# Todo App (軽量版)

このリポジトリは、フロントエンド、バックエンド、およびデータベース（SQLite）を使用した軽量なTodoアプリケーションの実装です。このアプリケーションは、ユーザーがタスクを管理し、追加、削除、更新できるシンプルなTodoリストを提供します。

## 概要

このTodoアプリは、以下の要素で構成されています：

- **フロントエンド**: HTML, CSS, JavaScriptを使用し、シンプルなインターフェースを提供します。バックエンドが静的コンテンツも提供します。
- **バックエンド**: Go（Golang）を使用してRESTful APIを提供し、クライアントとデータベースとのやり取りを担当します。GoがWebサーバーとAPIサーバーの両方の役割を担い、静的コンテンツもサーバー側で提供します。
- **データベース**: SQLiteを使用し、タスクデータを永続化します。

## 機能

- **タスクの作成**: 新しいTodoタスクを追加できます。
- **タスクの表示**: 追加されたタスクを一覧表示します。
- **タスクの更新**: 既存のタスクを更新できます。
- **タスクの削除**: 不要なタスクを削除できます。
- **データ永続化**: SQLiteデータベースにタスクを保存します。
- **静的コンテンツの提供**: Goサーバーが静的コンテンツ（HTML, CSS, JavaScript）を提供し、フロントエンドの表示を行います。

## システム構成

### フロントエンド
- **HTML**: タスクの入力フォームとリスト表示を担当します。静的ファイルはバックエンド（Goサーバー）から提供されます。
- **CSS**: シンプルで直感的なUIを実現します。
- **JavaScript**: フロントエンドからAPIを呼び出し、タスクの追加・削除・更新を行います。

### バックエンド
**Go (Golang)**: RESTful APIを提供するサーバーを構築します。また、Goサーバーは静的コンテンツ（HTML, CSS, JS）も提供し、フロントエンドとバックエンドが単一のサーバーで完結します。
- 使用するGoパッケージ:
  - `net/http` - HTTPサーバーの作成。
  - `encoding/json` - JSON形式のリクエスト/レスポンス処理。
  - `github.com/jinzhu/gorm` - ORM（Object Relational Mapping）ライブラリを使用してSQLiteと接続。

### データベース
**SQLite**: 軽量で高速なデータベース。タスクの情報（タイトル、状態など）を保存します。
- **テーブル構造**
- `todos`:
  - `id` (INTEGER, 主キー)
  - `title` (TEXT, タスク名)
  - `completed` (BOOLEAN, 完了状態)

## セットアップ方法

### 必要なもの

- Go (Golang)がインストールされていること。
- SQLiteがインストールされていること（またはGoのSQLiteライブラリがインストールされていれば、SQLiteのセットアップは不要）。
- `docker`および`docker-compose`（オプション）。

### インストール手順

1. このリポジトリをクローンします
```bash
git clone https://github.com/Action0358/todo-app-lite.git
cd todo-app-lite
```

2. Dockerコンテナのビルドと起動
```bash
docker-compose up --build
docker-compose up
```

3. コンテナにアクセスしてSQLiteをインストール
```bash
docker exec -it server /bin/bash
apt-get update
apt-get install sqlite3
```

4. アプリケーションへのアクセス
```bash
http://localhost:8080
```