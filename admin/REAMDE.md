# Docker Hub のイメージの自動 build について

Docker Hubの [automated-build](https://docs.docker.com/docker-hub/builds/#create-an-automated-build) 機能を使って、master ブランチに push されると、自動でイメージをビルドして Docker Hub に登録されるようになっている。

`Dockerfile` とイメージの対応は以下の通り

* `/docker/app/base/Dockerfile`: `showwin/ishocon2_app_base`
* `/docker/benchmarker/Dockerfile`: `showwin/ishocon2_bench`


`Dockerfile_app`: `ishocon2_app:latest` は言語毎に分けたことで使わなくなった。


## 手動でやる場合
### アプリケーションベースイメージ作成手順(メモ)
```
$ docker build -f ./docker/app/base/Dockerfile . -t showwin/ishocon2_app_base:$version
$ docker login
$ docker push showwin/ishocon2_app_base:$version
```


### ベンチマーカーイメージ作成手順(メモ)
```
$ docker build -f ./docker/benchmarker/Dockerfile . -t showwin/ishocon2_bench:$version
$ docker login
$ docker push showwin/ishocon2_bench:$version
```
