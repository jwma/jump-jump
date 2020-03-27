<h1 align="center">
  <br>Jump Jump<br>
</h1>

<p align="center"><em>开箱即用，Go 语言开发的一个功能完善的短链接系统。</em></p>
<p align="center">
  <a href="https://github.com/jwma/jump-jump/workflows/Go/badge.svg" target="_blank">
    <img src="https://github.com/jwma/jump-jump/workflows/Go/badge.svg" alt="ci">
  </a>
  <a href="https://img.shields.io/github/license/mashape/apistatus.svg" target="_blank">
      <img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="license">
  </a>
</p>

---

* [快速体验](#快速体验)
* [功能与使用](#功能与使用)
    * [短链接管理](#短链接管理)
* [本地启动](#本地启动)
* [如何访问短链接？](#如何访问短链接)
    * [设置短链接域名](#设置短链接域名)
    * [获取完整短链接](#获取完整短链接)
* [部署到服务器](#部署到服务器)
* [关注我的公众号了解更多](#关注我的公众号了解更多)

---

## 快速体验

[访问这里](http://anmuji.com/t/7pcu75)，来体验一下 Jump Jump 吧！

## 功能与使用

短链接基础功能已经开发完毕，后续的功能可以查看[版本规划](http://anmuji.com/t/h7ua8j)：

![Jump Jump 功能模块](http://rs.majiawei.com/jumpjump/features.png)

### 短链接管理

![短链接](http://rs.majiawei.com/jumpjump/v1.1.0copyshortlink.png)

## 本地启动

使用 `docker-compose` 启动，能够快速帮你启动 `redis`, `apiserver`, `landingserver`，使用如下命令：

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

## 如何访问短链接？

### 设置短链接域名

登入到管理后台之后，你能够看到短链接域名设置（只有管理员有权修改），在这里设置好你部署的域名/IP:Port，如：
`http://127.0.0.1:8000/` 或者 `http://anmuji.com/t/`，这里有一点需要注意的是，需要以 `/` 结尾。

### 获取完整短链接

访问短链接列表页面，如果你已经创建了短链接，那么可以在列表的第一个字段，悬停一下，会出现一个带有域名的完整短链接，点击就可以自动拷贝到
剪切板，你可以到需要使用的地方进行粘贴或者使用浏览器访问。

## 部署到服务器

这里提供了使用 docker-compose 的部署方案，[点击查看](http://anmuji.com/t/fk1ta3)。

## 关注我的公众号了解更多
![码极工作室](http://rs.majiawei.com/mjstudio/qrcode.png)
