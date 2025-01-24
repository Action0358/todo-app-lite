# ベースイメージとして nginx を使用
FROM nginx:latest

# ローカルの public ディレクトリの内容を nginx が公開するディレクトリにコピー
COPY public /usr/share/nginx/html

# Nginxの設定ファイルをコピー
COPY nginx.conf /etc/nginx/conf.d/default.conf

# コンテナが使用するポート番号を指定
EXPOSE 80