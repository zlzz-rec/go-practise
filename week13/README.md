## 一.工程结构

+ 该工程为巨石架构
+ 路由层使用iris框架.iris本身包含依赖注入功能, 但必须依赖iris app对象提供http服务. 为了支持job和task服务放弃iris自带依赖注入功能, 使用google的wire来实现
+ 服务模块包含myap、admin、job和task
+ 代码分层包含controller、service、biz和data. 为解决路由定义混乱问题可能在controller层前设计router层区分业务
+ 服务配置采用启动时接收启动参数的方式

```
.
├── README.md
├── cmd
│   ├── admin
│   ├── job
│   ├── myapp
│   │   ├── Makefile
│   │   ├── main.go
│   │   ├── run.sh
│   │   ├── wire.go
│   │   └── wire_gen.go
│   └── task
├── go.mod
├── go.sum
├── internal
│   ├── admin
│   ├── job
│   ├── myapp
│   │   ├── biz
│   │   ├── config
│   │   ├── controller
│   │   ├── data
│   │   └── service
│   └── task
├── pkg
│   ├── mysql
│   │   └── mysql.go
│   └── prometheus
│       └── prometheus.go
└── test
```



## 二.构建启动

```
cd cmd/myapp
wire && make && ./myapp \
    --server_address=0.0.0.0:8080\
    --mysql_address=rm-2ze8cu948171215g6.mysql.rds.aliyuncs.com:3306 \
    --mysql_user=feedmaster \
    --mysql_password=Feed_2018 \
    --mysql_database=testing \
    --redis_address=r-2zee75baa6194c14.redis.rds.aliyuncs.com:6379 \
    --redis_password=Feed@2018 \
    --redis_db=0 \
    --debug
```

## 三.说明

+ 1）微服务架构（BFF、Service、Admin、Job、Task 分模块）
  + 本项目做模块拆分(只有myapp有代码定义, 其它模块为空目录)
+ 2）API 设计（包括 API 定义、错误码规范、Error 的使用）
  + api定义在controller目录的app.go文件
  + 错误码在pkg的zerror包的error_define.go定义 
  + 错误码根因处理在pkg的zerror包的error_handler.go处理
  + 错误堆栈包装示例, 可以参考internal/data/hello.go的21行  
+ 3）gRPC 的使用
  + 本项目暂为使用gRPC
+ 4）Go 项目工程化（项目结构、DI、代码分层、ORM 框架）
  + 项目结构根据controller,service,biz,data做的分层
  + data层支持mysql和redis  
+ 5）并发的使用（errgroup 的并行链路请求)
  + main.go函数入口中使用errgroup进行了主app和Prometheus的并行链路处理
+ 6）微服务中间件的使用（ELK、Opentracing、Prometheus、Kafka）
  + 本项目在main.go中启用了Prometheus
  + 本项目在internal/biz/hello.go的30行发送了kafka消息   
+ 7）缓存的使用优化（一致性处理、Pipeline 优化）
  + pipeline示例暂未添加



