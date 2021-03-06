# simple-restful-server
Userに対してHTTPメソッドを用いてCRUD操作を行う事ができるようなRESTfulなWeb Application

サーバにはGo、DBにはPostgreSQLを使用。docker-composeで起動

ポートはTCP/8081番を使用

## dir structure
#### `server`
serverの本体のソースコード（及びbuild後のバイナリファイル）

`main.go` : ハンドラの登録、ハンドラとする関数を記述

`db/db.go` : 環境変数からデータベースへの接続を記述

`model/model.go` : User構造体とそのメソッドを記述

#### `init.sql`
postgreSQLのコンテナ内で起動時に実行するsqlファイル

#### `Dockerfile`
server用のdocker imageをbuildするためのDockerfile

#### `docker-compose.yml`
serverとpostgreSQLのコンテナを協調させて起動させるためのdocker-composeファイル

#### `.env`
server/postgreSQL containerに読ませる環境変数を記述
データベースの認証情報などを記述

## HTTP methods
```
GET    /            # {"Message": "Hello, w"rld!"}を表示
GET    /users       # user の一覧を表示
GET    /users/:id   # 指定した id の user を表示
POST   /users       # user を追加
PUT    /users/:id   # 指定した id の user を更新
DELETE /users/:id   # 指定した id の user を削除
```

## start
```
$ docker-compose up -d
```
（portは8081をフォワーディングする設定）

## usage
### user登録
```
$ curl -XPOST -D - -H 'Content-Type:application/json' http://localhost:8081/users -d '{"name": "new_user", "email": "hoge@example.com" }'
```

### user表示
```
$ curl -XGET -D - -H 'Content-Type:application/json' http://localhost:8081/users/1
```

### user更新
```
$ curl -XPUT -D - -H 'Content-Type:application/json' http://localhost:8081/users/13 -d '{"name": "updated", "email": "updated@example.com" }'
```

### user一覧取得
```
$ curl -XGET -D - -H "Content-type: application/json" http://localhost:8081/users
```
ユーザが存在しない場合は`nil`を返却

### user削除
```
$ curl -XDELETE -D - -H 'Content-Type:application/json' http://localhost:8081/users/1
```

## other features
- 種々の内部エラーには500を返すように設定
- `GET /users/:id`で存在しないユーザを指定した場合には404を返すように設定
