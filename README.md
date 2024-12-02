## 短链接生成器

短链接生成器是将一个长链接例如

- (https://www.google.com/search?q=sad&sca_esv=e058c774e46ad98b&sxsrf=ADLYWIJR81hSYagdglnIMd2c2qHxZFsmMg%3A1733131120255&ei=cHtNZ6-VD92C2roP_uT9iAo&ved=0ahUKEwivjZbG4IiKAxVdgVYBHX5yH6EQ4dUDCA8&uact=5&oq=sad&gs_lp=Egxnd3Mtd2l6LXNlcnAiA3NhZDIFEC4YgAQyBRAuGIAEMgUQABiABDIFEC4YgAQyBRAuGIAEMgUQABiABDIFEAAYgAQyCBAuGIAEGNQCMgUQABiABDIFEC4YgAQyFBAuGIAEGJcFGNwEGN4EGN8E2AEBSOsTUIYJWJwScAJ4AZABAJgBwwGgAckIqgEDMC42uAEDyAEA-AEBmAIHoAK0B6gCE8ICChAAGLADGNYEGEfCAgQQIxgnwgIKEAAYgAQYQxiKBcICCxAAGIAEGLEDGIMBwgIIEAAYgAQYsQPCAgcQIxgnGOoCwgITEAAYgAQYQxi0AhiKBRjqAtgBAcICFhAuGIAEGEMYtAIYyAMYigUY6gLYAQHCAgoQIxiABBgnGIoFwgILEAAYgAQYkQIYigXCAgsQLhiABBiRAhiKBcICCxAuGIAEGNEDGMcBwgILEC4YgAQYxwEYrwHCAhQQLhiABBiXBRjcBBjeBBjgBNgBAZgDBogGAZAGAroGBggBEAEYAZIHAzIuNaAHx1M&sclient=gws-wiz-serp)

转为短链接 

- http://bit.ly/3Z9T0Em 

这在一些需要分享URL的场景非常有用，例如在限制URL长度的twitter, 和各大评论区，电子书等媒介中。

### 短链接生成器非常适合初学者入门Web开发

1. 该项目相对简单，接口只有两个：
    - POST /api/url 接口, 接受长URL
    - GET /:code 接口: 把短URL重定向到长URL

2. 该项目非常实用，bitly和tinyURL公司就是以此为主要业务。

### 项目难点

这是一个读多写少的项目。

1. 难点1： GET请求需要服务时延低，响应速度快，使重定向的用户没有痛感。如果请求都要访问数据库，涉及到磁盘IO会增加响应时间，同时大量的请求会给数据库很多压力。所以需要使用redis进行缓存。

2. 难点2： 短URL的id如何生成，短id是可能重复的，需要使用重试机制提高成功率。

### 项目特点

1. 追求最佳实践，按照依赖倒置的原则，使用接口对重要的依赖进行解耦合，方便未来进行重构升级。
2. web框架使用echo,这纯粹是自己的喜好，因为echo使用装饰器的设计模式，我能看懂，同时天生支持全局错误处理。
3. 使用sqlc而不是orm, 因为我更偏爱直接使用sql与数据库打交道，而不是再学其他orm的语法。orm不适用于复杂的sql
4. 使用viper加载项目配置, postgres数据库， 并使用redis进行缓存


### 开发环境

1. 下载golang migrate[https://github.com/golang-migrate/migrate], 数据库迁移工具

```sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

2. 下载sqlc[https://github.com/sqlc-dev/sqlc], 将sql转为go代码
```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

3. 启动postgres 数据库 (使用docker)
```sh
docker run --name postgres-url \
	-e POSTGRES_USER=lang \
	-e POSTGRES_PASSWORD=password \
	-e POSTGRES_DB=urldb \
	-p 5432:5432 \
	-d postgres   
```

4. 启动redis (使用docker)
```sh
docker run --name redis \
	-p 6379:6379 \
	-d redis
```

### MCV架构

1. View（视图）： 人机交互接口，就是浏览器看到的界面。

2. Controller（控制器）： 负责处理用户请求，将请求中的数据反序列化，验证数据的格式是否正确，调用模型层处理，并把结果序列化返回给用户。

3. Model（模型）：负责程序的数据逻辑和业务规则，通常包含数据结构、数据库交互以及与应用逻辑相关的功能。

该项目是前后端分离的项目，后端不涉及V层，只提供接口返回json数据，只包含C和M层。M层太大，通常为了方便开发，又把M层分为repository和service层。
后端包含:

- repository: 持久层，与数据库交互。
- service: 逻辑层，汇总业务逻辑。依赖repository
- controller: 控制层，功能不变，如上。 依赖service

各个层级职责清晰，使用接口进行解耦合；controller层不应调用repository层的方法，而是通过service层进行调用。

### 运行项目

后端:
```sh
make migrate_up
sqlc generate
go mod tidy
go run main.go
```

前端:
```sh
cd frontend
pnpm install
npm run dev
```