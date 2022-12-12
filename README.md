# skyee
一个游戏服务器，利用go的特性实现完成无锁化，使用actor模式

# Test
运行一个服务，获取配置
```sh
cd ./examples/web
go build -o ./web/webserver ./main.go

cd ./web
./webserver

```

请求 Post 获取配置
```sh
POST http://127.0.0.1:1004/ClientPack?key=123 HTTP/1.1
content-type: application/json

{
    "ConfigVersion": "default",
    "GameID": "test",
    "Head": 123
}
```