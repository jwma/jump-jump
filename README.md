Jump Jump
===
![Go](https://github.com/jwma/jump-jump/workflows/Go/badge.svg) ![license](https://img.shields.io/github/license/mashape/apistatus.svg)

使用 Go 开发的一个功能完善的短链接系统。旧版源码可查看 [prev 分支](https://github.com/jwma/jump-jump/tree/prev)。

## 功能/使用
短链接基础功能已经开发完毕，其他功能正在开发中，完成后，将会支持如下的功能：

![Jump Jump 功能模块](j2module.png?raw=true "Jump Jump 功能模块")

### 短链接管理
![短链接](shortlinklist.png?raw=true "短链接")

## 本地启动
最快速的启动方式，就是使用 docker-compose 启动，使用如下命令：
```shell script
# 克隆或下载项目源码到本地
git clone https://github.com/jwma/jump-jump.git

# 进入项目源码目录
cd jump-jump/

# 在本地构建容器镜像
make dockerimage

# 启动
docker-compose -f deployments/docker-compose.yaml -p jumpjump up -d

# 查看服务运行状态
docker-compose -f deployments/docker-compose.yaml -p jumpjump ps

# 如果看到 apiserver/landingserver 未启动成功，重启一下就好
docker-compose -f deployments/docker-compose.yaml -p jumpjump restart

# 创建用户，在服务正常运行的情况，运行 createuser 可以创建用户，使用如下
docker-compose -f deployments/docker-compose.yaml -p jumpjump exec apiserver ./createuser --help

Usage of ./createuser:
  -password string
        password.
  -role int
        role, 1: normal user, 2: administrator. (default 1)
  -username string
        username.

# 创建一个管理员角色的用户
docker-compose -f deployments/docker-compose.yaml -p jumpjump exec apiserver ./createuser -username=mj
 -password=12345 -role=2
```

在服务启动完毕且已经创建好用户之后，可以打开浏览器，访问 `http://localhost:8080` 进入管理后台进行短链接的管理工作。

### 如何访问短链接？
我们可以通过 `landingserver` 提供的服务来访问短链接，访问 `http://localhost:8081/{短链接ID}`，`{短链接ID}` 就是我们在
管理后台看见的 ID，尝试添加一个短链接进行访问吧~

## 部署到服务器
1. 通过源码构建镜像，推送到 Docker Hub 或者私有仓库；
2. 在服务器拉取镜像；
3. 参考源码中 `deployments/docker-compose.yaml` 文件编写属于你的 docker-compose 配置文件；
4. 通过 docker-compose 启动服务；
5. 通过 `createuser` 命令行工具创建用户； 
6. 直接暴露服务端口 / 使用 Nginx 转发请求到服务；
7. 搞定。

## 关注我的公众号了解更多
![qr code](qrcode.png?raw=true "qr code")
