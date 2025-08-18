1. 工具安装
```shell
go install github.com/gogf/gf/cmd/gf/v2@latest
gf -v
```

2. 创建项目模板
```shell
gf init my_goframe
```
会自动生成项目目录my_goframe，及其子目录，`cd my_goframe`  

3. 建表
```sql 
create table login(
    username varchar(20) not null comment '用户名',
    password char(32) not null comment 'md5之后的密码',
    primary key (username)
) comment '用户登录表';
```

4. 生成dao
在配置文件hack/config.yaml里指定正确的dao link。
```shell
gf gen dao
```  

5. 生成controller
在api下创建各个接口的Request和Response结构体。
```shell
gf gen ctrl
```  

6. 补全handler
在 internal/controller/user下的go文件里补全handler。  
在 internal/cmd/cmd.go里绑定路由。加一行`xxx.NewV1()`

7. 修改配置文件
修改manifest/config/config.yaml，修改server的address、配置logger、dao link等。  
注意:hack用于指导代码生成，而程序运行时使用的配置文件是manifest/config/config.yaml。

8. 启动运行
```shell
go run .
```  