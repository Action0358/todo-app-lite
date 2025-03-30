# ベースイメージとしてGoの公式イメージを使用
FROM golang:1.23

# 作業ディレクトリを設定
WORKDIR /app

# Goモジュールを初期化するため、go.modとgo.sumをコンテナ内で作成
RUN go mod init todo-app-lite

# 依存関係をダウンロード
RUN go mod tidy

# Airをインストール
RUN go install github.com/air-verse/air@latest

# ソースコードをコピー（serverディレクトリ内のファイルをコンテナの/appにコピー）
COPY . .

# Airを実行するコマンドを設定
CMD ["air"]