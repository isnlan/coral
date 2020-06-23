我们致力于提供完整的微服务研发体验，整合相关框架及工具后，微服务治理相关部分可对整体业务开发周期无感，从而更加聚焦于业务交付。

## Features
* Cache：优雅的接口化设计，非常方便的缓存序列化，推荐结合代理模式[overlord](https://github.com/bilibili/overlord)；
* DB：集成MySQL/Mongo/，添加熔断保护和统计支持，可快速发现数据层压力；
* Log：类似[zap](https://github.com/uber-go/zap)的field实现高性能日志库，并结合log-agent实现远程日志管理；
* Trace：基于opentracing，集成了全链路trace支持（gRPC/HTTP/）；
* HLF: 简洁的hyperleger fabric sdk
* Contract：智能合约的通用接口
## Quick start

### Requirments

Go version>=1.14

### Installation
```shell
GO111MODULE=on && go get -u github.com/snlansky/coral
