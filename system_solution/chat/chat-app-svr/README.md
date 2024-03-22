# chat-app-svr

## 数据库
docker 启动一个 postgres SQL 
```shell
docker run --name postgres -e POSTGRES_PASSWORD=123456 -v postgres:/var/lib/postgresql/data -p 5432:5432 -d postgres:13
```

访问数据库和缓存的逻辑则自己进行编写，通过 orm 框架来实现dao层的代码。

## api协议
```shell
goctl api go -api chat-app-svr.api -dir .
```

## gRPC 服务
通过 go-zero 生成模板代码 ，代码逻辑在 internal/logic 内编写

增加 Client则需要再 internal/svc 下加入对应的Client.