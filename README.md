jump jump
===
这是一个使用Go语言开发的一个短链接系统，将包含如下功能：
- [x] 生成短链接接口
- [x] 短链接跳转功能
- [x] 短链接数据统计功能
- [x] 获取单个短链接信息接口
- [x] 短链接更新接口（更新短链接设置，如目标链接、启用禁用）
- [ ] 后台用户验证
- [ ] 短链接数据报表

### 怎么用？
由于jump-jump还处于开发当中，数据报表和后台都还没完成，但已经实现了一个生成短链接的接口，下面是cURL的调用示例：
```
curl -d "url=https://github.com/jwma&isEnabled=true&description=mj github" \
-X POST http://localhost:8081/admin/link
```

### 开发环境
可以使用提供的`docker-compose-dev.yml`来启动jump-jump和依赖的Redis服务，并在开发过程中会自动编译新代码。
```
# 启动
docker-compose -f docker-compose-dev.yml up --build

# 停止并清除容器
docker-compose -f docker-compose-dev.yml down --volumes
```
启动成功后，可以打开`http://localhost:8081`访问jump-jump。

### 生产环境
可以使用提供的`docker-compose.yml`来启动jump-jump和依赖的Redis服务。
```
# 启动
docker-compose up --build

# 停止并清除容器
docker-compose down --volumes
```
启动成功后，可以打开`http://localhost:8080`访问jump-jump。
