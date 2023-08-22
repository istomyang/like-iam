# like-iam

技术描述：各组件App使用 Cobra/Viper/Pflag 作为命令行解析，存储层使用 GORM 连接 MySQL。apiserver 使用 Gin，严格采用 RESTful 标准，用户认证使用 JWT。authzserver 使用 gRPC 保证组件间传输高性能，使用 Redis的发布/订阅 保证一致性，使用 Redis作为 消息队列 保存授权日志。pumper 消费授权日志并分发到下游消费，比如 csv、sysloger、influx、kafka、es、mongo、prometheus。wathcer 是监控组件，类似于 k8s 的 CronJob，比如 负责一些定时清理工作。log 是基于 zap 开发的日志库，采用 std/log 的API风格，与 klog 也有集成，可以作为单独的库用于生产环境。iam-sdk-go/iam-sdk-rust 是 链式调用风格 的SDK库。error 是自己实现的错误包。component/shutdown 是自己实现的优雅关闭包。