FROM golang:1.15

# 作業ディレクトリの作成・設定
WORKDIR /InfecShotAPI

# Go Modules を有効化
ENV GO111MODULE=on

COPY go.mod .

# go.mod 内のパッケージをダウンロード
RUN go mod download

EXPOSE 8080
