# ベースイメージとして nginx を使用
FROM nginx:stable-alpine

# Nginxの設定ファイルをコピー
COPY nginx.conf /etc/nginx/conf.d/default.conf

# publicディレクトリの内容をコンテナ内の適切な場所にコピー
COPY public/ /usr/share/nginx/html/public/

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]