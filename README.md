# stage

gin+gorm 面向接口（api）简易小demo，初学者也能快速上手

RESTful架构，结构清晰，传参灵活

### 项目结构

```
|-app
    |-controller 控制器
    |-dao 负责curd的
    |-middleware 中间件
    |-model 模型
    |-plugin 常用工具类
        |-driver 持久化工具类
        |-scene 统一输出工具类
        |-function.go 公共函数
    |-service 核心业务处理
|-config 配置文件和统一路由管理
    |-conf
        |-conf.go 配置方法    
    |-route
        |-route.go 路由配置文件 
    |-config.yml 统一环境配置文件   
    |-message.yml 统一状态码配置文件
|-main.go 程序执行入口
```

### RESTful API curl

```
添加用户：
curl --location --request POST 'http://localhost:8080/v1/user' \
--header 'action: add' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name":"tout",
    "password":"123"
}'

用户列表：
curl --location --request GET 'http://localhost:8080/v1/user?page=1&limit=10' \
--header 'action: list'

用户详情：
curl --location --request GET 'http://localhost:8080/v1/user?id=1' \
--header 'action: info'

修改用户:
curl --location --request PUT 'http://localhost:8080/v1/user' \
--header 'action: update' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id":1,
    "password":"123456"
}'

删除用户：
curl --location --request DELETE 'http://localhost:8080/v1/user' \
--header 'action: delete' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id":1
}'
```

注：一个router下的header action需与controller方法同名，联调时可以通过action快速定位接口

### 返参示例

```
{
    "status": 10000,
    "msg": "请求成功",
    "body": null
}
```

### 测试user表结构

```
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```