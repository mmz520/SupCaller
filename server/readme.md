# 生成建议

## 目录结构
```
|── common
|   ├── config
|   ├── middleware
|   ├── models
|   ├── utils
|── internal
|   |── security
|   |   ├── controllers
|   |   ├── service
|   |   ├── dao
|   |   ├── middleware
|   |   ├── models
|   |   |   ├── user.go
|   |   |   ├── role.go
|   |   |   ├── auth.go
|   |── resource
|   |   ├── trunk
|   |   |   ├── controllers
|   |   |   ├── service
|   |   |   ├── dao
|   |   |   ├── middleware
|   |   |   ├── models
|   |   |   |   ├── trunk.go
|   |   ├── lineConfig
|   |   |   ├── controllers
|   |   |   ├── service
|   |   |   ├── dao
|   |   |   ├── middleware
|   |   |   ├── models
|   |   |   |   ├── lineConfig.go
|── config.yaml
|── main.go
|── readme.md
```


## 基础设置
- mysql
- redis
- jwt
- oauth
- 全局错误处理
- 热更配置
- 日志包括操作日志中间件
- 国际化，时区
- 监控
- 限流 熔断
