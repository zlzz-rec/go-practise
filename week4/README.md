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
    --server_address=0.0.0.0:8080 \
    --mysql_address={mysql_address}:3306 \
    --mysql_user={mysql_user} \
    --mysql_password={mysql_password} \
    --mysql_database={mysql_database} \
    --debug
```





