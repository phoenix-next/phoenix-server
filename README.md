# PhoeniX-server

## Introduction

- PhoeniX 后台服务器，用于集中管理、存储用户信息
- 项目使用 Golang 作为主要语言，使用 Gin 框架以及 Golang 的其他第三方库制作
- 本项目为服务端，Shell 客户端位于[这里](https://github.com/phoenix-next/phoenix-shell)，而桌面客户端位于[这里](https://github.com/phoenix-next/phoenix)

## Development

开发过程中，有一些常用的命令：

```sh
# 更新Swagger的API文档(本质上是重新生成文档)
swag init
# 格式化(format)API函数的Swagger注释
swag fmt
# 更新并安装项目依赖
go mod tidy
# 构建项目，生成可执行文件
go build .
```

## Deploy

首先，需要 clone 项目并进行编译（需要 Go 1.17 及以上），得到可执行文件 phoenix-server：

```sh
git clone https://github.com/phoenix-next/phoenix-server
cd phoenix-server
go mod tidy
go build .
```

其次，需要编写配置文件 phoenix-config.yml（配置文件必须命名为该名称），示例配置文件如下：

```yml
database:
  ip: '127.0.0.1' # MySQL的IP地址
  port: '3306' # MySQL的端口
  user: 'username' # MySQL的用户名
  password: 'password' # MySQL对应用户的密码
  database: 'phoenix' # 使用的MySQL数据库名称

email:
  user: 'user@126.com' # 发送注册邮件的邮箱
  password: 'password' # 上述邮箱的密码（或授权码）
  host: 'smtp.126.com' # 该邮箱所属的邮件服务器
  port: '465' # 该邮件服务器提供服务的端口

server:
  port: 8080 # 服务器运行的端口
  secret: 'secret' # 服务器用于签名的密钥
  docs: false # 是否允许访问API文档页面
  debug: false # 是否以Debug模式运行服务器
  cert: 'certFileName' # 仅在非Debug模式下有效，为使用的SSL证书的文件名
  key: 'keyFileName' # 仅在非Debug模式下有效，为使用的私钥的文件名
  backend_url: 'backend_url' # 服务器允许访问的后端网络地址
```

将可执行文件和配置文件置于**相同目录**下，并执行可执行文件即可运行服务器

P.S. 若以非Debug模式运行服务器，则服务器将使用HTTPS协议进行传输，SSL证书以及私钥也必须和可执行文件置于**相同目录**下

## Credits

- 项目的结构参考了[Slime 学术分享平台](https://github.com/BFlameSwift/SlimeScholar-Go)，以及[Gin-Vue 代码框架](https://github.com/flipped-aurora/gin-vue-admin)
- 感谢 [Golang](https://github.com/golang/go) 项目以及 [Gin](https://github.com/gin-gonic/gin) 框架，这是本项目的基石
- 另外，还要感谢 [viper](https://github.com/spf13/viper), [logrus](https://github.com/sirupsen/logrus), [swag](https://github.com/swaggo/swag) 等第三方库的作者，这加快了项目的开发