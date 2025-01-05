# 使用するベースイメージ
FROM python:3.9-slim

# 必要なパッケージをインストール
RUN apt-get update && apt-get install -y \
    sqlite3 \
    && pip install sqlite-web

# 作業ディレクトリを設定
WORKDIR /app

# SQLite Webの実行コマンド
CMD ["sqlite_web", "-H", "0.0.0.0", "/app/database/todos.db"]
