# Goboot 配置化Web开发

## 简介
- 启动日志预览
```shell script
________  ________  ________  ________  ________  _________   
|\   ____\|\   __  \|\   __  \|\   __  \|\   __  \|\___   ___\ 
\ \  \___|\ \  \|\  \ \  \|\ /\ \  \|\  \ \  \|\  \|___ \  \_| 
 \ \  \  __\ \  \\\  \ \   __  \ \  \\\  \ \  \\\  \   \ \  \  
  \ \  \|\  \ \  \\\  \ \  \|\  \ \  \\\  \ \  \\\  \   \ \  \ 
   \ \_______\ \_______\ \_______\ \_______\ \_______\   \ \__\
    \|_______|\|_______|\|_______|\|_______|\|_______|    \|__|

[INFO] app [ go-server ] on [ dev ] run at port [ 8080 ] on time 2023-02-10 16:48:51
[INFO] local: http://localhost:8080/
```
- Goboot
- 是基于 Golang 的 web 开发框架 gin 实现的
- 一个二次包装的配置化开发模板
- 结合和部分Springboot的开发特性
- 在 Go 支持的情况下
- 实现配置化开发
- 实现常见配置支持
- 实现自动路由匹配

## 快速入门
- 第一步，安装GO环境
- 第二步，使用gomod初始化自己的项目
- 下面以项目名hello为例
- 新建项目 hello
```shell script
mkdir hello
```
- 进入项目
```shell script
cd hello
```
- 初始化gomod项目
- 注意，项目名要和创建的文件夹名称一致
```shell script
go mod init hello
```
- 第三步，拷贝goboot到项目中
- 这是目前的项目结构
```shell script
hello
|---goboot
    |---goboot.go
|---go.mod
```
- 第四步，编写自己的入口文件
- 也就是main.go
```shell script
vi main.go
```
```go
package main

import (
	"fmt"
	"hello/goboot"
	"time"

	"github.com/gin-gonic/gin"
)

// 定义API处理结构体
type Api struct {
}

// 给结构体绑定函数
// 这里的函数名为Hello
// 在后面会仔细讲解函数名的问题
func (api *Api) Hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"hello": "hello",
	})
}

// 定义主函数
func main() {
	// 获得默认的应用对象
    app := goboot.GetDefaultApplication()
    // 添加mapping处理对象，注意是指针
	app.AddHandlers(&Api{})
    // 添加默认主页的GET请求处理
	app.App.GET("/", func(c *gin.Context) {
		stime := time.Now().Format("2006-01-02 15:04:05")
		c.String(200, fmt.Sprintf("现在是北京时间：%v",stime))
	})
    // 运行应用
	app.Run()
}
```
- 第四步，编写配置文件
- 也就是 goboot.yml
```shell script
vi goboot.yml
```
```yaml
goboot:
  application:
    name: go-server
  profiles:
    active: dev
  server:
    port: 8080
    bannerPath: ./banner.txt
    staticResources:
      enable: true
      urlPath: /static
      filePath: ./static
    templateResources:
      enable: true
      filePath: ./templates/**/*.html
    session:
      enable: true
      # cookie/redis
      impl: cookie
      secretKey: 123456
      sessionKey: go-session
    redis:
      enable: true
      host: 127.0.0.1
      port: 6379
      password: ltb12315
      database: 0
    datasource:
      enable: true
      # mysql/postgres
      driver: mysql
      host: 127.0.0.1
      port: 6379
      url: user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
      username: root
      password: 123456
      database: test_db
    gorm:
      enable: true
    https:
      enable: false
      pemPath: ./https/server.pem
      keyPath: ./https/server.key
    gzip:
      enable: false
      # BestCompression/BestSpeed/DefaultCompression/NoCompression
      level: DefaultCompression
      excludeExtensions:
        - .mp4
      excludePaths:
        - /api/
      excludePathRegexes:
        - "*download"
    proxy:
      - name: github.com
        path: /github/
        redirect: http://github.com/
    mapping:
      - /api/
    cors:
      enable: true
      allowAllOrigins: true
      allowOrigins:
        - http://localhost/
      allowMethods:
        - GET
        - PUT
        - DELETE
        - POST
        - PATCH
        - OPTIONS
      allowHeaders:
        - token
        - Origin
        - secure
        - Auth
      exposeHeaders:
        - Content-Length
      allowCredentials: true
      maxAgeMinutes: 0
```
- 第五步，下载依赖
```shell script
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/gzip
go get github.com/gin-contrib/cors
go get github.com/gin-contrib/sessions
go get github.com/go-yaml/yaml
go get github.com/redis/go-redis/v9
go get github.com/gin-contrib/sessions/redis@v0.0.5
go get github.com/google/uuid
go get github.com/go-sql-driver/mysql
go get github.com/lib/pq
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get gorm.io/driver/postgres
```
- 第六步，启动运行
```shell script
go run main.go
```
- 第七步，浏览器访问查看
```shell script
http://localhost:8080/
http://localhost:8080/hello
```
- 这是最终的文件结构
```shell script
hello
|---goboot
    |---goboot.go
|---go.mod
|---goboot.yml
|---main.go
```


## 配置文件
- 配置文件，直接在配置文件中
- 加注释给出说明
- 注意，在此处的配置环境中
- 不同于springboot配置
- 这里是区分大小写，严格匹配的
- 这个java同学需要注意
```yaml
# 配置根节点
goboot:
  # 应用配置      
  application:
    # 应用名称
    name: go-server
  # 多环境配置
  profiles:
    # 激活的环境，找不到指定配置就是默认配置文件
    # 查找规则：goboot.yml goboot-${goboot.prfiles.active}.yml
    # 比如这里，就查找goboot-dev.yml
    active: dev
  # 服务配置
  server:
    # 服务的启动端口    
    port: 8080
    # 也可以配置自己的启动banner
    bannerPath: ./banner.txt
    # 静态资源配置  
    staticResources:
      # 是否启用
      enable: true
      # url中的路径
      urlPath: /static
      # 解析为静态资源的路径
      filePath: ./static
    # 模板文件配置
    templateResources:
      # 是否启用    
      enable: true
      # 模板文件的匹配规则  
      filePath: ./templates/**/*.html
    # session 配置部分
    session:
      # 是否开启session
      enable: true
      # 使用的session存储类型，目前有以下两种可选
      # 当选redis时，必须配置redis
      # cookie/redis
      impl: cookie
      # session存储的加密秘钥
      secretKey: 123456
      # session在客户端的cookie键名称
      sessionKey: go-session
    # redis 配置
    redis:
      # 是否开启redis
      enable: true
      # redis 主机
      host: 127.0.0.1
      # redis 端口
      port: 6379
      # redis 访问密码
      password: ltb12315
      # redis 使用的数据库
      database: 0
    # 数据源配置
    datasource:
      # 是否启用数据源
      enable: true
      # 数据源驱动类型，支持以下类型
      # mysql/postgres
      driver: mysql
      # 数据源主机
      host: 127.0.0.1
      # 数据源端口
      port: 3306
      # 当url有配置时，按照url配置进行，其他数据源参数无效
      # 没有配置时，使用其他参数解析
      url: user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
      # 数据源用户名
      username: root
      # 数据源密码
      password: 123456
      # 数据源数据库
      database: test_db
    # ORM 配置
    gorm:
      # 是否启用 ORM
      enable: true
    # HTTPS的配置部分
    https:
      # 是否启用
      enable: false
      # 分别配置HTTPS的pem文件和key文件
      pemPath: ./https/server.pem
      keyPath: ./https/server.key
    # gzip响应压缩配置
    gzip:
      # 是否启用
      enable: false
      # 压缩级别：BestCompression，BestSpeed，DefaultCompression，NoCompression
      level: DefaultCompression
      # 排除的后缀列表
      excludeExtensions:
        - .mp4
      # 排除的路径前缀列表
      excludePaths:
        - /api/
      # 排除的路径匹配正则列表
      excludePathRegexes:
        - "*download"
    # 代理配置
    proxy:
      # 可以配置多个代理配置
      # 代理的名称，可以随意  
      - name: github.com
        # 代理的路径
        path: /github/
        # 目标跳转路径
        redirect: http://github.com/
    # 自动路径映射配置
    mapping:
      # 可以配置多个进行按照匹配规则自动路由  
      - /api/
    # 跨域配置
    cors:
      # 是否启用
      enable: true
      allowAllOrigins: true
      allowOrigins:
        - http://localhost/
      allowMethods:
        - GET
        - PUT
        - DELETE
        - POST
        - PATCH
        - OPTIONS
      allowHeaders:
        - token
        - Origin
        - secure
        - Auth
      exposeHeaders:
        - Content-Length
      allowCredentials: true
      maxAgeMinutes: 0
```

## 接口开发
- 接口开发，可以使用gin框架自己的方式
- 也可以使用配置中的mapping自动映射两种模式
- 也可以实现GobootController接口定义分组路由
- 三种种模式，都是基于封装的goboot

### gin模式接口开发
- gin模式，就是通过应用实例，获取App属性，得到gin.Engine实现
- 得到 boot 对象
```go 
boot := goboot.GetDefaultApplication()
```
- 拿到 gin.Engine 对象
```go
engine := boot.App
```
- 然后，就可以和原始的gin开发一样开发了
```go
engine.GET("/", func(c *gin.Context) {
    stime := time.Now().Format("2006-01-02 15:04:05")
    c.String(200, fmt.Sprintf("现在是北京时间：%v",stime))
})
```


### mapping自动映射模式接口开发
- 此模式，首先，你得知道工作原理
- 工作原理
    - 首先，需要在配置中 goboot.server.mapping 配置上自动映射的路径
    - 例如，举例配置中的 /api/ 这个路径
    - 则，如果请求路径为：http://localhost:8080/api/hello/go
    - 则 /api/hello/go 就是符合mapping配置的一个路径
    - 则，去除 /api/ 这一层之后，得到的路径为 hello/go
    - 在对boot的代码中，将一个带有方法的结构体，添加到 handlers 中
    - 则，表示这些对象的方法，都具备可以自动映射的能力
    - 假设，结构体为 Api , 具有一个 Hello_Go 方法
    - 那么，对于路径 hello/go 就被映射到 Hello_Go 方法上
    - 具体的映射规则如下：
        - URL路径: /gin-web/hello-goboot
        - 将URL路径按照每层路径分隔
        - 得到：gin-web，hello-goboot
        - 对每一级，按照横向分隔
        - 得到：gin,web和hello,goboot
        - 按照每一级内使用大驼峰（Capital）组合
        - 得到：GinWeb，HelloGoboot
        - 对每一级使用下划线组合
        - 得到：GinWeb_HelloGoboot
        - 这个就是这个路径对应要映射的方法名
        - 那么，将会在注册的 handlers 中，查找名称为这样的一个方法来处理请求
        - 下面给出一些映射案例：
        - 一： hello ---> Hello
        - 二： helloWorld --> HelloWorld
        - 三： hello-world --> HelloWorld
        - 四： hello-go/hello-gin --> HelloGo_HelloGin
    - 针对restful类型接口，限制请求方式的适配
        - 在映射函数的定义上
        - 如果是X*_开头的函数，将限定为指定规则对应的请求方式
        - 具体的关系如下
            - XG_ --> GET
            - XU_ --> PUT
            - XP_ --> POST
            - XD_ --> DELETE
            - XH_ --> PATCH
            - XA_ --> ANY
        - 对于包含这些前缀的函数
        - 映射的路径匹配的函数名
        - 将去除这些后缀后进行匹配
        - 举例：
            - 方法名：XP_Get_User
            - 则对应的请求：POST /get/user
            - 当使用其他请求类型时，将404：GET /get/user
    - 映射函数的要求：
        - 入参可以有多个
        - 顺序可以任意
        - 也可以无参数
        - 支持的参数如下
            - c *gin.Context
            - boot * goboot.GobootApplication
            - engine * gin.Engine
            - request * http.Request
            - resp * goboot.ApiResp
            - ctxResp * goboot.CtxResp
            - redis * redis.Client
            - redisCli * goboot.RedisCli
            - session sessions.Session
            - db *sql.DB
            - gormDb * gorm.DB
            - 自定义绑定请求参数的结构体
                - 注意，必须是结构体类型
                - 结构体支持值类型或指针类型
        - 方法案例
            - 一：func (api *Api) Hello()
            - 二：func (api *Api) Hello(c *gin.Context)
            - 三：func (api *Api) Hello(boot *goboot.GobootApplication,c *gin.Context)
            - 四：func (api *Api) Hello(c *gin.Context, post User, boot *goboot.GobootApplication, engine *gin.Engine, request *http.Request)
- 使用代码示例
```go
package main

import (
	"hello/goboot"

	"github.com/gin-gonic/gin"
)

// 定义API处理结构体
type Api struct {
}

// 给结构体绑定函数
// 这里的函数名为Hello
// 在后面会仔细讲解函数名的问题
func (api *Api) Hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"hello": "hello",
	})
}

// 定义主函数
func main() {
	// 获得默认的应用对象
    app := goboot.GetDefaultApplication()
    // 添加mapping处理对象，注意是指针
	app.AddHandlers(&Api{})
    // 运行应用
	app.Run()
}
```
- 因此，如果使用mapping模式开发
- 分三步走
- 第一步，确认配置文件中的mapping有添加
- 第二步，编写一个符合mapping要求的结构体，也就是具有方法
- 第三步，调用boot对象的AddHandles方法，添加处理的所有结构体

### GobootController 路由分组模式
- 这个模式，其实也是自动路由的一种变体
- 之前的mapping模式相当于全局自动映射
- 而controller模式，则是分组的自动映射模式
- 使用上和mapping一样，只不过处理的结构体
- 也就是说，定义的处理函数，和mapping模式一样定义即可
- 需要实现接口 GobootController
    - 关于实现接口，在Golang中，接口的实现，不需要什么implements/extends等关键字
    - 只需要将接口中的每个方法在结构体中实现即可
- 下面就以一个示例来说明
- 首先，定义自己的mapping结构体
- 实现接口中定义个path方法
```go
// 定义处理结构
type AdminController struct {
}
// 实现接口方法，返回这个分组路由为 /admin/
func (con *AdminController) Path() string {
	return "/admin/"
}
// 添加自己的路径映射处理函数
func (admin *AdminController) XP_Get(ctx *goboot.CtxResp) any {
	return ctx.ApiJsonOk("ok")
}
```
- 将controller添加到路由中
```go
// 拿到应用对象
app := goboot.GetDefaultApplication()
// 添加controller
app.AddControllers(&AdminController{})
// 运行应用
app.Run()
```

## 自动映射函数
- 上面说了mapping模式的自动映射函数
- 只是简单的介绍了映射函数
- 下面就来详细的说明映射函数
- 以及自动注入的入参的作用或设计初衷
- 下面讲解映射函数，绑定的结构体，都以 Api 讲解
- 配置的mapping 为 /api/
```go
type Api struct{
}
```
### 需要手动指定路径
- 缺少像springmvc的注解声明方式
- 则取而代之，使用函数名作为路径匹配规则
- 原始写法
```go
engine.GET("/api/hello",func(c *gin.Context){
  c.JSON(200,gin.H{
    "data":"hello",
  })
})
```
- 使用Goboot之后
- 则可以改写为如下方式
```go
func (api * Api) Hello(c * gin.Context){
  c.JSON(200,gin.H{
    "data":"hello",
  })
}
```
- 按照匹配规则，Hello函数名进行匹配请求路径
- 这样，避免了两个问题
- 直接使用engine对象
- 写明请求路径
- 这两个问题，都极大的增大了耦合性

### 数据响应之后，必须自行return
- 在使用gin进行响应数据时
- 需要明确的指定return
- 即时是abort也必须进行返回
- 否则如果后续有其他响应逻辑
- 则会连带执行其他响应
- 导致两个响应结合再一起
- 问题如下
```go
engine.GET("/api/hello", func(c *gin.Context) {
  c.AbortWithStatusJSON(200, gin.H{
    "data": "error",
  })
  // 如果此处没有return
  // 则下面的正常响应将会继续响应
  // return
  c.JSON(200, gin.H{
    "data": "hello",
  })
})
```
- 这里演示的这种情况实际中是很常见的
- 只不过，实际中对于abort是有条件的
- 但是依旧不能避免return
- 下面，在goboot中
- 封装了两个结构，来解决此问题
- 这都是基于自动映射实现的
- 因为自动映射，不关心返回值，返回值不会被处理
- 方法一，使用goboot.ApiResp结合gin.Context实现直接返回
```go
func (api *Api) Hello(resp *goboot.ApiResp, c *gin.Context) *goboot.ApiResp {
	return resp.GinOk(c, "hello")
}
```
- 这种方式，通过自动注入goboot.ApiResp结合gin.Context实现
- 方式二，和方法一一致，只不过自己实例化goboot.ApiResp指针
```go
func (api *Api) Hello(c *gin.Context) *goboot.ApiResp {
	return (&goboot.ApiResp{}).GinOk(c, "hello")
}
```
- 这种方式，只需要注入gin.Context即可
- 方式三，和方法二类似，只不过通过方法实例化goboot.ApiResp指针
```go
func (api *Api) Hello(c *gin.Context) *goboot.ApiResp {
	return goboot.ApiOk(nil).GinOk(c, "hello")
}
```
- 这种方式，比较起来容易接受
- 方式四，推荐方式，直接使用goboot.CtxResp实现
```go
func (api *Api) Hello(resp *goboot.CtxResp) *goboot.CtxResp {
	return resp.ApiJsonOk("hello")
}
```
- 这种方式，最为简单实用
- 一般业务场景中，这种模式，在加上自动解析请求参数注入
- 就是一般的使用模式
- 如下
```go
func (api *Api) Login(resp *goboot.CtxResp,user * User) *goboot.CtxResp {
	return resp.ApiJsonOk("ok")
}
```

## 主要函数或结构
- 常量：DefaultConfigFile ，指定了默认的配置文件的名称 为 ./goboot.yml
- 常量：DefaultBannerText ，指定了默认的应用banner的文本值
- 结构：ApiResp ，定义了标准的接口返回结构，code，msg，data
    - 以及包含了常用的填值结构方法
    - 以及包含了针对gin的JSON返回的结构方法Gin*系列
    - 以及全局静态方法Api*系列
- 常量：ApiCodeOk ，指定了默认的ApiResp返回正常时的code值
- 常量：APiCodeErr ，指定了默认的ApiResp异常返回时的code值
- 结构：Tokens ，定了了几个结构方法，用于获取UUID和从请求中获取token的结构方法
- 结构：CtxResp ，是最常用的mapping系列自动映射函数中最常用的一个入参，包含了context,session,app
    - 以及包含了对ApiResp结构响应JSON的ApiJson*系列结构函数
    - 以及包含了原始gin响应的Json/string/html函数
    - 以及包含了对session设置获取的Session*系列函数
- 函数：Log* 系列全局函数，使用自定义的控制台数据日志
- 结构：GobootConfig 定义了解析配置文件的根配置结构
    - 此结构包含了整个配置文件中的配置信息
    - 如有需要，可以进行获取
- 结构：GobootApplication 是封装的goboot的应用实例结构
    - 整个goboot的上下文，引擎等都在此结构中进行包含
    - 其中包含了，gin.Engine,GobootConfig,Handlers,GobootLifecycleListener,RedisCli,GobootController,sql.DB,gorm.DB
    - 此实例，通过Get*Application系列函数进行初始化获取
    - 最终设置完毕之后，使用结构函数 Run 来启动一个应用
- 接口：GobootController 是针对 GobootApplication 中Controllers定义的接口
    - 用于定义分组路由的自动映射
    - 其中包含一个 Path 方法，用于获取分组路由的路径
- 结构：RedisCli 是对 redis.Client 的简单封装
    - 主要是为了简化原来的redis.Client的使用
    - 目前提供了简单的GET和set方法
- 函数类型： GobootListener 定义了在应用初始化和启动的各个生命周期进行监听的接口函数
    - 可以用于监听对应周期应用的状态
    - 或者在对应的周期进行修改应用配置的目的
- 结构：GobootLifecycleListener 定义了一组声明周期各个环节的监听集合
    - 用来组装 GobootListener
- 函数：GetDefaultApplication 用来获取一个默认配置文件配置的应用实例
    - 实际上是使用默认配置 goboot.yml 调用 GetApplication 来获取应用实例
    - 这也是最常用的一个函数
- 函数：GetApplication 支持监听器的根据指定配置文件获取应用实例
    - 实际上是使用 ResolveGobootConfig 来获取配置结构，调用 GetConfigApplication 来获取应用实例
- 函数：GetConfigApplication 直接根据配置结构获取应用实例
- 函数：ReadGobootConfig 将指定的配置文件，解析为配置结构
- 函数：ResolveGobootConfig 读取指定的配置文件，并根据Profiles重定向读取配置
- 函数：MappingHandler 负责进行结构的路径自动映射，实现函数调用的处理方法
    - 这个方法服务于自动映射mapping和GobootController
    - 实现将请求按照规则，调用目标函数的过程
- 函数：HandleMappingMethodArg 负责实现参数类型的实际参数的自动绑定
    - 是为 MappingHandler 实现自动注入函数调用入参的核心函数调用
- 函数：ProxyHandler 负责进行实现proxy配置进行自动代理的处理函数

### 测试Demo
- 文件结构
```shell script
hello
|---goboot
    |---goboot.go
|---templates
    |---index
        |---index.html
|---main.go
|---go.mod
|---goboot.yml
```
- 入口程序
- main.go
```go
package main

import (
	"hello/goboot"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Api struct {
}

func (api *Api) Hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"hello": "hello",
	})
}

func (api *Api) ApiResp(resp *goboot.ApiResp, c *gin.Context) {
	c.JSON(200, resp.Ok(gin.H{
		"hello": "hello",
	}))
}

func (api *Api) GinResp(resp *goboot.ApiResp, c *gin.Context) *goboot.ApiResp {
	return resp.GinOk(c, "hello")
}

func (api *Api) GinApiResp(c *gin.Context) *goboot.ApiResp {
	return goboot.ApiOk(nil).GinOk(c, "hello")
}

func (api *Api) CtxResp(resp *goboot.CtxResp) *goboot.CtxResp {
	return resp.ApiJsonOk("hello")
}

func (api *Api) Login(resp *goboot.CtxResp, user *User) *goboot.CtxResp {
	return resp.ApiJsonOk("ok")
}

type User struct {
	Username string `form:"username"`
}

func (user *User) User_Info(c *gin.Context, post User, boot *goboot.GobootApplication, engine *gin.Engine, request *http.Request) {
	c.JSON(200, gin.H{
		"user": post.Username,
		"boot": boot.Config.ConfigFile,
	})

}

func (user *User) Session_Set(ctx *goboot.CtxResp) any {
	ctx.SessionSet("user", "admin")
	return ctx.ApiJsonOk("ok")
}

func (user *User) Session_Get(ctx *goboot.CtxResp, session sessions.Session) any {
	val := ctx.SessionGet("user")
	val = session.Get("user")
	return ctx.ApiJsonOk(val)
}

func (user *User) Redis_Set(ctx *goboot.CtxResp, redis *goboot.RedisCli) any {
	redis.Set("user", "root")
	return ctx.ApiJsonOk("ok")
}

func (user *User) Redis_Get(ctx *goboot.CtxResp, redis *goboot.RedisCli) any {
	val := redis.Get("user")
	return ctx.ApiJsonOk(val)
}

type AdminController struct {
}

func (con *AdminController) Path() string {
	return "/admin/"
}

func (admin *AdminController) XP_Get(ctx *goboot.CtxResp) any {
	return ctx.ApiJsonOk("ok")
}

func main() {
	app := goboot.GetDefaultApplication()

	app.AddHandlers(&Api{}).
		AddHandlers(&User{})

	app.AddControllers(&AdminController{})

	app.App.GET("/", func(c *gin.Context) {
		stime := time.Now().Format("2006-01-02 15:04:05")
		c.HTML(200, "index/index.html", gin.H{
			"now": stime,
		})
		// c.String(200, stime)
	})

	app.Run()
}


```
- 模板文件
- templates/index/index.html
```html
{{ define "index/index.html" }}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>首页</title>
</head>
<body>
    <h2>现在是北京时间 {{.now}}</h2>
</body>
</html>
{{ end }}
```
- 配置文件
- goboot.yml
```yml
goboot:
  application:
    name: go-server
  profiles:
    active: dev
  server:
    port: 8080
    bannerPath: ./banner.txt
    staticResources:
      enable: true
      urlPath: /static
      filePath: ./static
    templateResources:
      enable: true
      filePath: ./templates/**/*.html
    session:
      enable: true
      # cookie/redis
      impl: cookie
      secretKey: 123456
      sessionKey: go-session
    redis:
      enable: true
      host: 127.0.0.1
      port: 6379
      password: ltb12315
      database: 0
    datasource:
      enable: true
      # mysql/postgres
      driver: mysql
      host: 127.0.0.1
      port: 6379
      url: user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
      username: root
      password: 123456
      database: test_db
    gorm:
      enable: true
    https:
      enable: false
      pemPath: ./https/server.pem
      keyPath: ./https/server.key
    gzip:
      enable: false
      # BestCompression/BestSpeed/DefaultCompression/NoCompression
      level: DefaultCompression
      excludeExtensions:
        - .mp4
      excludePaths:
        - /api/
      excludePathRegexes:
        - "*download"
    proxy:
      - name: github.com
        path: /github/
        redirect: http://github.com/
    mapping:
      - /api/
    cors:
      enable: true
      allowAllOrigins: true
      allowOrigins:
        - http://localhost/
      allowMethods:
        - GET
        - PUT
        - DELETE
        - POST
        - PATCH
        - OPTIONS
      allowHeaders:
        - token
        - Origin
        - secure
        - Auth
      exposeHeaders:
        - Content-Length
      allowCredentials: true
      maxAgeMinutes: 0
```