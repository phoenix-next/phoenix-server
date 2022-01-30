# PhoeniX-server

## Introduction

- PhoeniX 后台服务器，用于集中管理、存储用户信息
- 项目使用 Golang + Gin 制作

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
```

最后，将可执行文件和配置文件置于相同目录下，并执行可执行文件即可运行服务器

## Credits

- 项目的结构参考了[Slime 学术分享平台](https://github.com/BFlameSwift/SlimeScholar-Go)，以及[Gin-Vue 代码框架](https://github.com/flipped-aurora/gin-vue-admin)
