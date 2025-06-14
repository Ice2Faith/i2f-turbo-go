package goboot

// /////////////////////////////////////////////////////////
// goboot 简介
// 提供配置化的常用功能的提供能力
// 以及基于反射的URL自动映射处理函数能力
// 简单使用 见 readme.md 详细说明
// 依赖包下载
// go get github.com/gin-gonic/gin
// go get github.com/gin-contrib/gzip
// go get github.com/gin-contrib/cors
// go get github.com/gin-contrib/sessions
// go get gopkg.in/yaml.v3
// go get github.com/redis/go-redis/v9
// go get github.com/gin-contrib/sessions/redis@v0.0.5
// go get github.com/google/uuid
// go get github.com/go-sql-driver/mysql
// go get github.com/lib/pq
// go get gorm.io/gorm
// go get gorm.io/driver/mysql
// go get gorm.io/driver/postgres
// /////////////////////////////////////////////////////////
import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	goredis "github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// /////////////////////////////////////////////////////////
// goboot 默认配置区
// /////////////////////////////////////////////////////////

// 默认配置文件名称
const DefaultConfigFile string = "./goboot.yml"

// 默认 Banner 内容
// Banner在线生成三方网址
// https://www.bootschool.net/ascii
const DefaultBannerText string = `
________  ________  ________  ________  ________  _________   
|\   ____\|\   __  \|\   __  \|\   __  \|\   __  \|\___   ___\ 
\ \  \___|\ \  \|\  \ \  \|\ /\ \  \|\  \ \  \|\  \|___ \  \_| 
 \ \  \  __\ \  \\\  \ \   __  \ \  \\\  \ \  \\\  \   \ \  \  
  \ \  \|\  \ \  \\\  \ \  \|\  \ \  \\\  \ \  \\\  \   \ \  \ 
   \ \_______\ \_______\ \_______\ \_______\ \_______\   \ \__\
    \|_______|\|_______|\|_______|\|_______|\|_______|    \|__|
`

// /////////////////////////////////////////////////////////
// goboot API区
// /////////////////////////////////////////////////////////

// 定义标准响应码
const (
	ApiCodeOk  int = 200
	APiCodeErr int = 500
)

// 定义具有JSON的TAG的响应结构
type ApiResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 绑定基础函数
func (api *ApiResp) Ret(code int, msg string, data interface{}) *ApiResp {
	api.Code = code
	api.Msg = msg
	api.Data = data
	return api
}

func (api *ApiResp) Ok(data interface{}) *ApiResp {
	api.Code = ApiCodeOk
	api.Msg = ""
	api.Data = data
	return api
}

func (api *ApiResp) Err(msg string) *ApiResp {
	api.Code = APiCodeErr
	api.Msg = msg
	api.Data = nil
	return api
}

func (api *ApiResp) Error(code int, msg string) *ApiResp {
	api.Code = code
	api.Msg = msg
	api.Data = nil
	return api
}

// 绑定具有响应能力的函数
func (api *ApiResp) GinRet(c *gin.Context, code int, msg string, data interface{}) *ApiResp {
	api.Code = code
	api.Msg = msg
	api.Data = data
	c.JSON(200, api)
	return api
}

func (api *ApiResp) GinOk(c *gin.Context, data interface{}) *ApiResp {
	api.Code = ApiCodeOk
	api.Msg = ""
	api.Data = data
	c.JSON(200, api)
	return api
}

func (api *ApiResp) GinErr(c *gin.Context, msg string) *ApiResp {
	api.Code = APiCodeErr
	api.Msg = msg
	api.Data = nil
	c.JSON(200, api)
	return api
}

func (api *ApiResp) GinError(c *gin.Context, code int, msg string) *ApiResp {
	api.Code = code
	api.Msg = msg
	api.Data = nil
	c.JSON(200, api)
	return api
}

// 直接方法构建Api响应
func ApiRet(code int, msg string, data interface{}) *ApiResp {
	return &ApiResp{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func ApiOk(data interface{}) *ApiResp {
	return &ApiResp{
		Code: ApiCodeOk,
		Data: data,
	}
}

func ApiErr(msg string) *ApiResp {
	return &ApiResp{
		Code: APiCodeErr,
		Msg:  msg,
	}
}

func ApiError(code int, msg string) *ApiResp {
	return &ApiResp{
		Code: code,
		Msg:  msg,
	}
}

// /////////////////////////////////////////////////////////
// goboot 令牌验证区
// /////////////////////////////////////////////////////////
type Tokens struct {
}

func (tk Tokens) MakeUUID() string {
	uid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return uid.String()
}
func (tk Tokens) MakeNumberUUID() uint32 {
	uid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return uid.ID()
}
func (tk Tokens) MakeToken() string {
	str := tk.MakeUUID()
	return strings.ReplaceAll(str, "-", "")
}

func (tk Tokens) FindToken(c *gin.Context, tokenName string) string {
	token := c.GetHeader(tokenName)
	if token == "" {
		token = c.Query(tokenName)
	}
	if token == "" {
		token = c.PostForm(tokenName)
	}
	if token == "" {
		token = c.Param(tokenName)
	}
	return token
}

// /////////////////////////////////////////////////////////
// goboot 上下文响应区
// /////////////////////////////////////////////////////////

// 定义一个直接响应的上下文包装结构
type CtxResp struct {
	Context *gin.Context
	Session sessions.Session
	App     *GobootApplication
}

// 定义响应ApiResp的JSON函数
func (api *CtxResp) ApiJsonRet(code int, msg string, data interface{}) *CtxResp {
	api.Context.JSON(200, ApiRet(code, msg, data))
	return api
}
func (api *CtxResp) ApiJsonOk(data interface{}) *CtxResp {
	api.Context.JSON(200, ApiOk(data))
	return api
}
func (api *CtxResp) ApiJsonErr(msg string) *CtxResp {
	api.Context.JSON(200, ApiErr(msg))
	return api
}
func (api *CtxResp) ApiJsonError(code int, msg string) *CtxResp {
	api.Context.JSON(200, ApiError(code, msg))
	return api
}

// 定义普通结构JSON响应函数
func (api *CtxResp) Json(obj interface{}) *CtxResp {
	api.Context.JSON(200, obj)
	return api
}

// 定义普通字符串响应函数
func (api *CtxResp) String(format string, args ...interface{}) *CtxResp {
	api.Context.String(200, format, args...)
	return api
}

// 定义普通模板响应函数
func (api *CtxResp) Html(template string, obj interface{}) *CtxResp {
	api.Context.HTML(200, template, obj)
	return api
}

// 获取session
func (api *CtxResp) SessionGet(key interface{}) interface{} {
	return api.Session.Get(key)
}

// 设置session
func (api *CtxResp) SessionSet(key interface{}, val interface{}) *CtxResp {
	api.Session.Set(key, val)
	api.Session.Save()
	return api
}

// /////////////////////////////////////////////////////////
// goboot Log区
// /////////////////////////////////////////////////////////
// 控制台日志输出
func Log(level string, format string, args ...interface{}) {
	time := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(fmt.Sprintf("[%v] [%v]", time, level), fmt.Sprintf(format, args...))
}
func LogInfo(format string, args ...interface{}) {
	Log("INFO ", format, args...)
}
func LogWarn(format string, args ...interface{}) {
	Log("WARN ", format, args...)
}
func LogError(format string, args ...interface{}) {
	Log("ERROR", format, args...)
}

// /////////////////////////////////////////////////////////
// goboot 配置区
// /////////////////////////////////////////////////////////

// 配置结构
type GobootConfig struct {
	Goboot     Goboot `yaml:"goboot"`
	ConfigFile string
}

// 配置根节点
type Goboot struct {
	Application Application `yaml:"application"`
	Profiles    Profiles    `yaml:"profiles"`
	Server      Server      `yaml:"server"`
}

// 应用配置
type Application struct {
	Name string `yaml:"name"`
}

// 环境配置
type Profiles struct {
	Active string `yaml:"active"`
}

// 服务器配置
type Server struct {
	Port       int    `yaml:"port"`
	BannerPath string `yaml:"bannerPath"`

	StaticResources   StaticResources   `yaml:"staticResources"`
	TemplateResources TemplateResources `yaml:"templateResources"`
	Https             Https             `yaml:"https"`
	Gzip              Gzip              `yaml:"gzip"`
	Cors              Cors              `yaml:"cors"`
	Proxy             Proxy             `yaml:"proxy"`
	Mapping           Mapping           `yaml:"mapping"`
	Session           Session           `yaml:"session"`
	Redis             Redis             `yaml:"redis"`
	Datasource        Datasource        `yaml:"datasource"`
	Gorm              Gorm              `yaml:"gorm"`

	FileServer FileServer `yaml:"fileServer"`
}

// 静态资源项配置
type StaticResourcesItem struct {
	UrlPath  string `yaml:"urlPath"`
	FilePath string `yaml:"filePath"`
	TryFiles string `yaml:"tryFiles"`
}

// 静态资源配置
type StaticResources struct {
	Enable bool                  `yaml:"enable"`
	Items  []StaticResourcesItem `yaml:"items"`
}

// 模板渲染配置
type TemplateResources struct {
	Enable   bool   `yaml:"enable"`
	FilePath string `yaml:"filePath"`
}

// Session 配置
type Session struct {
	Enable     bool   `yaml:"enable"`
	Impl       string `yaml:"impl"`
	SecretKey  string `yaml:"secretKey"`
	SessionKey string `yaml:"sessionKey"`
}

// Redis 配置
type Redis struct {
	Enable   bool   `yaml:"enable"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

// Datasource 配置
type Datasource struct {
	Enable   bool   `yaml:"enable"`
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// Gorm 配置
type Gorm struct {
	Enable bool `yaml:"enable"`
}

// HTTPS配置
type Https struct {
	Enable  bool   `yaml:"enable"`
	PemPath string `yaml:"pemPath"`
	KeyPath string `yaml:"keyPath"`
}

// 代理配置
type ProxyItem struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Redirect string `yaml:"redirect"`
}

type Proxy struct {
	Enable bool        `yaml:"enable"`
	Items  []ProxyItem `yaml:"items"`
}

// 映射配置
type Mapping struct {
	Enable bool     `yaml:"enable"`
	Items  []string `yaml:"items"`
}

// GZIP配置
type Gzip struct {
	Enable             bool     `yaml:"enable"`
	Level              string   `yaml:"level"`
	ExcludeExtensions  []string `yaml:"excludeExtensions"`
	ExcludePaths       []string `yaml:"excludePaths"`
	ExcludePathRegexes []string `yaml:"excludePathRegexes"`
}

// 跨域配置
type Cors struct {
	Enable           bool     `yaml:"enable"`
	AllowAllOrigins  bool     `yaml:"allowAllOrigins"`
	AllowOrigins     []string `yaml:"allowOrigins"`
	AllowMethods     []string `yaml:"allowMethods"`
	AllowHeaders     []string `yaml:"allowHeaders"`
	ExposeHeaders    []string `yaml:"exposeHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	MaxAgeMinutes    int      `yaml:"maxAgeMinutes"`
}

// /////////////////////////////////////////////////////////
// goboot 应用区
// /////////////////////////////////////////////////////////

// 应用实例
type GobootApplication struct {
	App         *gin.Engine
	Config      *GobootConfig
	Handlers    []interface{}
	Listeners   *GobootLifecycleListener
	Redis       *RedisCli
	Controllers []GobootController
	Db          *sql.DB
	GormDb      *gorm.DB
}

// 控制器，需要提供基础路径
type GobootController interface {
	Path() string
}

func (app *GobootApplication) AddControllers(handlers ...GobootController) *GobootApplication {
	app.Controllers = append(app.Controllers, handlers...)
	return app
}

// redis 客户端封装
type RedisCli struct {
	Redis   *goredis.Client
	Context context.Context
}

func (redis *RedisCli) Set(key string, val interface{}) *goredis.StatusCmd {
	return redis.Redis.Set(redis.Context, key, val, 0)
}
func (redis *RedisCli) SetExpire(key string, val interface{}, expire time.Duration) *goredis.StatusCmd {
	return redis.Redis.Set(redis.Context, key, val, expire)
}
func (redis *RedisCli) GetCheck(key string) (string, error) {
	cmd := redis.Redis.Get(redis.Context, key)
	return cmd.Result()
}
func (redis *RedisCli) Get(key string) string {
	cmd := redis.Redis.Get(redis.Context, key)
	val, _ := cmd.Result()
	return val
}

// 应用监听器
type GobootListener func(boot *GobootApplication)

// 应用声明周期监听器
type GobootLifecycleListener struct {
	OnConfiged                 []GobootListener
	OnBeforeUse                []GobootListener
	OnBeforeStaticResources    []GobootListener
	OnBeforeTemplatesResources []GobootListener
	OnBeforeProxy              []GobootListener
	OnBeforeMapping            []GobootListener
	OnPrepared                 []GobootListener
	OnBeforeBanner             []GobootListener
	OnBeforeRun                []GobootListener
}

// 处理器必须是struct类型的指针
// 因为需要使用反射
// URL路径与处理函数的匹配规则
//
// url: user/info
//
// 1. URL中的每个部分都会被处理为Capital大驼峰格式
// User Info
//
// 2. 然后使用下划线连接，就是对应的函数名称
// User_Info
//
// 举例:
//
// user --> User
//
// user/info --> User_Info
//
// user/fav-icon --> User_FavIcon
//
// find-user/geo-range --> FindUser_GeoRange
func (app *GobootApplication) AddHandlers(handlers ...interface{}) *GobootApplication {
	app.Handlers = append(app.Handlers, handlers...)
	return app
}

// 获取默认的应用实例
func GetDefaultApplication() *GobootApplication {
	LogInfo("default application initial...")
	return GetApplication(DefaultConfigFile, nil)
}

// 从指定文件读取应用配置
// 始终返回配置，第二个返回值表示是否正确读取了配置
// 不会处理Profiles
func ReadGobootConfig(cfgFile string) (config *GobootConfig, ok bool) {
	// 初始化默认配置
	config = &GobootConfig{
		Goboot: Goboot{
			Server: Server{
				Port: 8080,
			},
		},
	}

	// 读取配置文件
	bytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		LogWarn("read config file %v error of %v", cfgFile, err)
	}

	// 解析yaml到结构
	if err == nil {
		err = yaml.Unmarshal(bytes, config)
		if err != nil {
			LogWarn("parse yaml config file %v error of %v", cfgFile, err)
		} else {
			ok = true
			config.ConfigFile = cfgFile
		}
	}

	return
}

// 从指定的配置文件获取配置
// 会处理Profiles的重定向配置
// 至多重定向配置一次
func ResolveGobootConfig(cfgFile string) *GobootConfig {
	// 初始化默认配置
	config := &GobootConfig{
		Goboot: Goboot{
			Server: Server{
				Port: 8080,
			},
		},
	}

	// 从配置文件读取配置
	cfg, ok := ReadGobootConfig(cfgFile)
	if ok {
		LogInfo("load config yaml file %v", cfgFile)
		// 判断是否具有active激活配置
		if cfg.Goboot.Profiles.Active != "" {
			LogInfo("find profile active %v", cfg.Goboot.Profiles.Active)
			redirectCfgFile := fmt.Sprintf("goboot-%v.yml", cfg.Goboot.Profiles.Active)
			// 读取激活配置
			rcfg, rok := ReadGobootConfig(redirectCfgFile)
			// 正确读取激活配置，则使用激活配置，否则还是原配置
			if rok {
				LogInfo("override config yaml file %v", redirectCfgFile)
				config = rcfg
			} else {
				LogInfo("fallback use config yaml file %v", cfgFile)
				config = cfg
			}
		} else {
			config = cfg
		}
	}
	return config
}

// 使用配置文件获取应用，同时支持添加生命周期监听器
// 可以用在各个阶段进行定制化的调整
func GetApplication(cfgFile string, listener *GobootLifecycleListener) *GobootApplication {
	LogInfo("use config yaml %v initial application with listener %v", cfgFile, listener)
	config := ResolveGobootConfig(cfgFile)
	return GetConfigApplication(config, listener)
}

// 调用指定的监听器
func invokeListeners(boot *GobootApplication, listeners []GobootListener) {
	for _, item := range listeners {
		item(boot)
	}
}

// 使用结构和监听器获取应用实例
// 此方法作为最底层的调用，但依然开放出去使用
func GetConfigApplication(config *GobootConfig, listener *GobootLifecycleListener) *GobootApplication {
	// 由于是指针，避免空指针
	if config == nil {
		config = &GobootConfig{
			Goboot: Goboot{
				Server: Server{
					Port: 8080,
				},
			},
		}
	}
	if listener == nil {
		listener = &GobootLifecycleListener{}
	}
	// 实例化应用结构
	boot := &GobootApplication{
		App:       gin.Default(),
		Config:    config,
		Handlers:  []interface{}{},
		Listeners: listener,
	}

	// 调用监听器
	LogInfo("goboot configed.")
	invokeListeners(boot, boot.Listeners.OnConfiged)

	engine := boot.App

	server := boot.Config.Goboot.Server

	// 配置 redis
	if server.Redis.Enable {
		if server.Redis.Port == 0 {
			server.Redis.Port = 6379
		}
		if server.Redis.Host == "" {
			server.Redis.Host = "127.0.0.1"
		}
		redisAddr := fmt.Sprintf("%v:%v", server.Redis.Host, server.Redis.Port)
		boot.Redis = &RedisCli{
			Redis: goredis.NewClient(&goredis.Options{
				Addr:     redisAddr,
				Password: server.Redis.Password,
				DB:       server.Redis.Database,
			}),
			Context: context.Background(),
		}
		LogInfo("goboot redis config, connect to: %v", redisAddr)
	}

	// 数据源配置
	if server.Datasource.Enable {
		if server.Datasource.Host == "" {
			server.Datasource.Host = "127.0.0.1"
		}
		if server.Datasource.Driver == "mysql" {
			url := server.Datasource.Url
			if url == "" {
				// user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
				url = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?tls=skip-verify&autocommit=true", server.Datasource.Username, server.Datasource.Password, server.Datasource.Host, server.Datasource.Port, server.Datasource.Database)
			}
			db, err := sql.Open(server.Datasource.Driver, url)
			if err != nil {
				panic(err)
			}
			boot.Db = db

			LogInfo("goboot datasource(mysql) config, connect to: %v", url)

			// 提供gorm连接
			if server.Gorm.Enable {
				gormDB, err := gorm.Open(mysql.New(mysql.Config{
					Conn: boot.Db,
				}), &gorm.Config{})
				if err != nil {
					panic(err)
				}
				boot.GormDb = gormDB

				LogInfo("goboot gorm(mysql) config.")
			}

		} else if server.Datasource.Driver == "postgres" {
			if server.Datasource.Port == 0 {
				server.Datasource.Port = 5432
			}
			url := server.Datasource.Url
			if url == "" {
				// postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full
				url = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=verify-full", server.Datasource.Username, server.Datasource.Password, server.Datasource.Host, server.Datasource.Port, server.Datasource.Database)
			}
			db, err := sql.Open(server.Datasource.Driver, url)
			if err != nil {
				panic(err)
			}
			boot.Db = db

			LogInfo("goboot datasource(postgres) config, connect to: %v", url)

			// 提供gorm连接
			if server.Gorm.Enable {
				gormDB, err := gorm.Open(postgres.New(postgres.Config{
					Conn: boot.Db,
				}), &gorm.Config{})
				if err != nil {
					panic(err)
				}
				boot.GormDb = gormDB

				LogInfo("goboot gorm(postgres) config.")
			}
		} else {
			LogWarn("goboot datasource config not support auto config, type : %v", server.Datasource.Driver)
		}

	}

	LogInfo("goboot before use.")
	invokeListeners(boot, boot.Listeners.OnBeforeUse)

	// 配置跨域
	if server.Cors.Enable {
		LogInfo("goboot enable cors.")
		corsConfig := cors.Config{}

		corsConfig.AllowAllOrigins = server.Cors.AllowAllOrigins
		if !corsConfig.AllowAllOrigins {
			if len(server.Cors.AllowOrigins) > 0 {
				corsConfig.AllowOrigins = server.Cors.AllowOrigins
			}
			if len(server.Cors.AllowMethods) > 0 {
				corsConfig.AllowMethods = server.Cors.AllowMethods
			}
			if len(server.Cors.AllowHeaders) > 0 {
				corsConfig.AllowHeaders = server.Cors.AllowHeaders
			}
		}

		if len(server.Cors.ExposeHeaders) > 0 {
			corsConfig.ExposeHeaders = server.Cors.ExposeHeaders
		}
		corsConfig.AllowCredentials = server.Cors.AllowCredentials
		if server.Cors.MaxAgeMinutes > 0 {
			corsConfig.MaxAge = time.Minute * time.Duration(server.Cors.MaxAgeMinutes)
		}

		engine.Use(cors.New(corsConfig))
	}

	// 配置gzip
	if server.Gzip.Enable {
		LogInfo("goboot enable gzip.")
		gzipLevel := gzip.DefaultCompression
		levelStr := server.Gzip.Level
		if levelStr == "BestCompression" {
			gzipLevel = gzip.BestCompression
		} else if levelStr == "BestSpeed" {
			gzipLevel = gzip.BestSpeed
		} else if levelStr == "DefaultCompression" {
			gzipLevel = gzip.DefaultCompression
		} else if levelStr == "NoCompression" {
			gzipLevel = gzip.NoCompression
		}
		options := []gzip.Option{}
		if len(server.Gzip.ExcludeExtensions) > 0 {
			op := gzip.WithExcludedExtensions(server.Gzip.ExcludeExtensions)
			options = append(options, op)
		}
		if len(server.Gzip.ExcludePaths) > 0 {
			op := gzip.WithExcludedPaths(server.Gzip.ExcludePaths)
			options = append(options, op)
		}
		if len(server.Gzip.ExcludePathRegexes) > 0 {
			op := gzip.WithExcludedPathsRegexs(server.Gzip.ExcludePathRegexes)
			options = append(options, op)
		}
		engine.Use(gzip.Gzip(gzipLevel, options...))
	}

	// 配置 session
	if server.Session.Enable {
		LogInfo("goboot enable session.")
		if server.Session.SessionKey == "" {
			server.Session.SessionKey = "go-session"
		}

		// 判断是否使用redis作为session-store
		sessionImpl := server.Session.Impl
		if sessionImpl == "redis" {
			if !server.Redis.Enable {
				panic("redis session require enable redis config [goboot.server.redis.enable]")
			}
			if server.Redis.Port == 0 {
				server.Redis.Port = 6379
			}
			if server.Redis.Host == "" {
				server.Redis.Host = "127.0.0.1"
			}
			redisAddr := fmt.Sprintf("%v:%v", server.Redis.Host, server.Redis.Port)
			store, err := redis.NewStore(server.Redis.Database, "tcp", redisAddr, server.Redis.Password, []byte(server.Session.SecretKey))
			if err != nil {
				panic(err)
			}
			engine.Use(sessions.Sessions(server.Session.SessionKey, store))
			LogInfo("goboot enable session(redis).")
		} else {
			store := cookie.NewStore([]byte(server.Session.SecretKey))
			engine.Use(sessions.Sessions(server.Session.SessionKey, store))
			LogInfo("goboot enable session(cookie).")
		}
	}

	LogInfo("goboot before static resources.")
	invokeListeners(boot, boot.Listeners.OnBeforeStaticResources)

	// 配置静态资源
	if server.StaticResources.Enable {
		for _, staticItem := range server.StaticResources.Items {
			LogInfo("goboot enable static resources, mapping: %v --> %v", staticItem.UrlPath, staticItem.FilePath)
			if _, err := os.Stat(staticItem.FilePath); os.IsNotExist(err) {
				os.MkdirAll(staticItem.FilePath, 0777)
			}
			engine.Static(staticItem.UrlPath, staticItem.FilePath)
		}
		// 处理404时的资源的try files处理
		engine.NoRoute(func(c *gin.Context) {
			reqPath := c.Request.URL.Path
			if !strings.HasSuffix(reqPath, "/") {
				reqPath = reqPath + "/"
			}
			for _, staticItem := range server.StaticResources.Items {
				urlPath := staticItem.UrlPath
				if !strings.HasSuffix(urlPath, "/") {
					urlPath = urlPath + "/"
				}
				if !strings.HasPrefix(reqPath, urlPath) {
					continue
				}
				filesArr := strings.Split(staticItem.TryFiles, " ")
				for _, fileItem := range filesArr {
					if fileItem == "" {
						continue
					}
					tryFile := staticItem.FilePath + "/" + fileItem
					_, err := os.Stat(tryFile)
					if err == nil {
						LogInfo("[try files] url: %v try to %v", reqPath, urlPath+fileItem)
						c.File(tryFile)
					}
				}
			}
		})
	}

	LogInfo("goboot before templates resources.")
	invokeListeners(boot, boot.Listeners.OnBeforeTemplatesResources)

	// 配置模板
	if server.TemplateResources.Enable {
		LogInfo("goboot enable templates resources, path: %v", server.TemplateResources.FilePath)
		if _, err := os.Stat(server.TemplateResources.FilePath); os.IsNotExist(err) {
			os.MkdirAll(server.TemplateResources.FilePath, 0777)
		}

		engine.LoadHTMLGlob(server.TemplateResources.FilePath)
	}

	LogInfo("goboot before file-server.")

	// 配置文件服务器
	if server.FileServer.Enable {
		LogInfo("goboot enable file-server at rootPath: %v", server.FileServer.RootPath)
		engine.Use(FileServerMiddleware(server.FileServer))
	}

	LogInfo("goboot before proxy.")
	invokeListeners(boot, boot.Listeners.OnBeforeProxy)

	// 配置代理
	if server.Proxy.Enable {
		LogInfo("goboot enable %v proxy(s)", len(server.Proxy.Items))
		for _, item := range server.Proxy.Items {
			LogInfo("goboot proxy, path: %v", item)
		}
		engine.Use(ProxyMiddleware(server.Proxy))
	}

	LogInfo("goboot before mapping.")
	invokeListeners(boot, boot.Listeners.OnBeforeMapping)

	// 配置自动映射
	if server.Mapping.Enable {
		LogInfo("goboot enbale %v mapping(s)", len(server.Mapping.Items))
		for _, item := range server.Mapping.Items {
			LogInfo("goboot mapping, path: %v", item)
		}
		engine.Use(MappingMiddleware(server.Mapping, boot))
	}

	LogInfo("goboot prepared.")
	invokeListeners(boot, boot.Listeners.OnPrepared)

	return boot
}

// 映射请求中间件
func MappingMiddleware(mapping Mapping, boot *GobootApplication) gin.HandlerFunc {
	return func(c *gin.Context) {
		hasMatched := false
		// 检查路径前缀匹配
		urlPath := c.Request.URL.Path
		for _, item := range mapping.Items {
			if strings.HasPrefix(urlPath, item) {
				hasMatched = true
				proxyPath := urlPath[len(item):]
				LogInfo("goboot mapping, path: %v", item)
				MappingHandler(boot, c, proxyPath, boot.Handlers...)
			}
		}
		if !hasMatched {
			// 如果不匹配，继续执行
			c.Next()
		}

	}
}

// GET,PUT,POST,DELETE,PATCH
// XG_,XU_,XP_,XD_,XH_,XA_
// 处理mapping自动映射
// 将配置中的mapping自动按照URL路径映射
func MappingHandler(boot *GobootApplication, c *gin.Context, proxyPath string, handlers ...interface{}) {
	// 处理自动映射的异常为404
	errorMsg := "request not found."
	defer func() {
		err := recover()
		if err != nil {
			c.String(404, errorMsg)
		}
	}()

	if len(handlers) == 0 {
		errorMsg = "not found any handlers."
		panic(errorMsg)
	}

	requestMethod := c.Request.Method

	// 将路径转换为函数名
	paths := strings.Split(proxyPath, "/")
	methodName := ""

	for _, path := range paths {
		item := strings.Trim(path, " \t\n\r")
		if item == "" {
			continue
		}
		parts := strings.Split(item, "-")
		item = ""
		for _, part := range parts {
			if part == "" {
				continue
			}
			item += strings.ToUpper(part[:1]) + part[1:]
		}
		if item == "" {
			continue
		}

		item = strings.ToUpper(item[:1]) + item[1:]
		if methodName != "" {
			methodName += "_"
		}
		methodName += item
	}

	// 遍历处理器，查找符合映射规则的函数
	for _, handler := range handlers {
		// 类型
		htype := reflect.TypeOf(handler)
		// 实际类型
		rtype := htype
		// 当时指针类型时，需要拿到实际类型
		if htype.Kind() == reflect.Ptr && htype.Elem().Kind() == reflect.Struct {
			rtype = htype.Elem()
		}
		// 只有是结构体类型的处理器才算处理器
		if rtype.Kind() == reflect.Struct {
			// 得到函数个数
			mcnt := htype.NumMethod()
			for i := 0; i < mcnt; i++ {
				mm := htype.Method(i)
				funcName := mm.Name
				funcMethod := ""
				if strings.HasPrefix(funcName, "XG_") {
					funcMethod = "GET"
					funcName = funcName[3:]
				} else if strings.HasPrefix(funcName, "XP_") {
					funcMethod = "POST"
					funcName = funcName[3:]
				} else if strings.HasPrefix(funcName, "XU_") {
					funcMethod = "PUT"
					funcName = funcName[3:]
				} else if strings.HasPrefix(funcName, "XD_") {
					funcMethod = "DELETE"
					funcName = funcName[3:]
				} else if strings.HasPrefix(funcName, "XH_") {
					funcMethod = "PATCH"
					funcName = funcName[3:]
				} else if strings.HasPrefix(funcName, "XA_") {
					funcMethod = ""
					funcName = funcName[3:]
				}
				// 如果函数名匹配
				if funcName == methodName {
					if funcMethod != "" {
						if funcMethod != requestMethod {
							errorMsg = "request method allow, require only " + funcMethod
							panic(errorMsg)
						}
					}
					// 拿到函数对象
					method := reflect.ValueOf(handler).MethodByName(mm.Name)
					// 获取入参个数
					paramCnt := method.Type().NumIn()
					callArgs := []reflect.Value{}
					matchFlag := true
					// 为每个函数入参注入值
					for p := 0; p < paramCnt; p++ {
						arg := method.Type().In(p)
						val, ok := HandleMappingMethodArg(arg, boot, c)
						if ok {
							callArgs = append(callArgs, val)
						} else {
							matchFlag = false
							break
						}
					}
					// 如果出现无法注入的参数，则不匹配，进行继续匹配
					if !matchFlag {
						continue
					}
					// 匹配成功的函数，进行调用
					method.Call(callArgs)
					return
				}
			}

		}
	}

	// 执行到这里，说明没有任何函数匹配
	errorMsg = "not found any handler method in handlers"
	panic(errorMsg)
}

// 为自动映射的方法添加调用参数
// 实现对boot，ctx,engine,request的方法入参自动注入
// 对结构体的请求参数自动填充能力
func HandleMappingMethodArg(arg reflect.Type, boot *GobootApplication, c *gin.Context) (reflect.Value, bool) {
	engine := boot.App
	request := c.Request
	resp := ApiOk(nil)
	ctxResp := &CtxResp{
		Context: c,
		Session: nil,
		App:     boot,
	}
	redis := boot.Redis.Redis
	redisCli := boot.Redis
	gormDb := boot.GormDb
	if boot.Config.Goboot.Server.Session.Enable {
		ctxResp.Session = sessions.Default(c)
		if arg.String() == "sessions.Session" {
			return reflect.ValueOf(ctxResp.Session), true
		}
	}

	// 如果是指针类型的参数
	if arg.Kind() == reflect.Ptr {
		// 分别判断是否是支持注入的内置类型，如果是，直接注入
		if arg.Elem() == reflect.TypeOf(*c) {
			return reflect.ValueOf(c), true
		} else if arg.Elem() == reflect.TypeOf(*ctxResp) {
			return reflect.ValueOf(ctxResp), true
		} else if arg.Elem() == reflect.TypeOf(*resp) {
			return reflect.ValueOf(resp), true
		} else if arg.Elem() == reflect.TypeOf(*request) {
			return reflect.ValueOf(request), true
		} else if arg.Elem() == reflect.TypeOf(*boot) {
			return reflect.ValueOf(boot), true
		} else if arg.Elem() == reflect.TypeOf(*engine) {
			return reflect.ValueOf(engine), true
		} else if arg.Elem() == reflect.TypeOf(*redisCli) {
			return reflect.ValueOf(redisCli), true
		} else if arg.Elem() == reflect.TypeOf(*gormDb) {
			return reflect.ValueOf(gormDb), true
		} else if arg.Elem() == reflect.TypeOf(*redis) {
			return reflect.ValueOf(redis), true
		} else if arg.Elem().Kind() == reflect.Struct {
			// 如果不是预定义的，但是是结构体，则自动请求参数绑定注入
			bindParam := reflect.New(arg.Elem()).Interface()
			c.ShouldBind(bindParam)
			return reflect.ValueOf(bindParam), true
		}
	} else if arg.Kind() == reflect.Struct {
		// 如果直接是结构体类型，直接实例化，自动请求参数绑定注入
		bindParam := reflect.New(arg).Interface()
		c.ShouldBind(bindParam)
		return reflect.ValueOf(bindParam).Elem(), true
	}

	// 其他类型，则绑定参数失败
	return reflect.ValueOf(false), false
}

// 代理请求中间件
func ProxyMiddleware(proxy Proxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		hasMatched := false
		// 检查路径前缀匹配
		urlPath := c.Request.URL.Path
		for _, item := range proxy.Items {
			if strings.HasPrefix(urlPath, item.Path) {
				hasMatched = true
				redirect := item.Redirect
				proxyPath := urlPath[len(item.Path):]
				LogInfo("goboot proxy, path: %v", item)
				ProxyHandler(c, redirect, proxyPath)
			}
		}
		if !hasMatched {
			// 如果不匹配，继续执行
			c.Next()
		}

	}
}

// 处理代理请求
func ProxyHandler(c *gin.Context, redirect string, proxyPath string) {
	remote, err := url.Parse(redirect)
	if err != nil {
		panic(err)
	}

	client := httputil.NewSingleHostReverseProxy(remote)

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = c.Request.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = path.Join(remote.Path, proxyPath)
		LogInfo("proxy req: %v %v", req.Method, req.URL)
	}
	client.ModifyResponse = func(resp *http.Response) error {
		LogInfo("proxy resp: %v | %v | %v", resp.StatusCode, resp.Status, resp.Request.URL)
		return nil
	}

	client.ServeHTTP(c.Writer, c.Request)
}

// 文件服务配置
type FileServer struct {
	Enable          bool   `yaml:"enable"`
	RootPath        string `yaml:"rootPath"` // 文件根路径
	UrlPath         string `yaml:"urlPath"`
	DisableUpload   bool   `yaml:"disableUpload"`   // 是否禁止上传
	DisableDownload bool   `yaml:"disableDownload"` // 是否禁止下载
	DisableList     bool   `yaml:"disableList"`     // 是否禁止举出文件
	DisableBrowser  bool   `yaml:"disableBrowser"`  // 是否禁止浏览文件
}

type FileInfoItem struct {
	Name       string `json:"name"`       // 文件类型
	Path       string `json:"path"`       // 相对路径
	Size       int64  `json:"size"`       // 文件大小（字节）
	SizeText   string `json:"sizeText"`   // 文件大小（带单位描述）
	IsDir      bool   `json:"isDir"`      // 是否为目录
	ModifyTime string `json:"modifyTime"` // 更新时间
}

func Any2JsonString(obj any) string {
	// 将 struct 转为 JSON 字符串
	jsonBytes, err := json.Marshal(obj) // 返回 []byte
	if err != nil {
		panic(err)
	}

	// 转换为字符串
	jsonStr := string(jsonBytes)
	return jsonStr
}

func ConvertAsHumanSizeText(size int64) string {
	if size < 0 {
		return "invalid"
	}
	if size == 0 {
		return "0B"
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	f := float64(size)
	i := 0

	for ; i < len(units)-1; i++ {
		if f < 1024 {
			break
		}
		f /= 1024
	}

	if f == math.Floor(f) {
		return fmt.Sprintf("%d%s", int64(f), units[i])
	}
	return fmt.Sprintf("%.2f%s", f, units[i])
}

func ListFiles(fullPath string, rootPath string) ([]FileInfoItem, error) {

	var files []FileInfoItem

	items, err := os.ReadDir(fullPath)
	for _, item := range items {

		info, err := item.Info()
		if err != nil {
			return files, err
		}

		// 获取文件完整路径
		path := filepath.Join(fullPath, info.Name())

		// 获取相对路径
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return files, err
		}

		// 添加文件信息到列表
		files = append(files, FileInfoItem{
			Path:       strings.ReplaceAll(relPath, "\\", "/"),
			Name:       info.Name(),
			Size:       info.Size(),
			SizeText:   ConvertAsHumanSizeText(info.Size()),
			IsDir:      info.IsDir(),
			ModifyTime: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	// 排序逻辑：目录在前，文件在后，同类型按字母顺序
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir // 目录优先
		}
		return files[i].Name < files[j].Name // 同类型按字母顺序
	})
	return files, err
}

func SliceContains(list []string, elem string) bool {
	for _, item := range list {
		if item == elem {
			return true
		}
	}
	return false
}

func FileServerMiddleware(server FileServer) gin.HandlerFunc {

	pathBase := server.UrlPath
	if pathBase == "" {
		pathBase = "/file-server"
	}
	pathBase = strings.TrimSuffix(pathBase, "/")

	if !strings.HasPrefix(pathBase, "/") {
		pathBase = "/" + pathBase
	}
	if pathBase == "/" {
		pathBase = ""
	}
	pathList := pathBase + "/list"
	pathUpload := pathBase + "/upload"
	pathDownload := pathBase + "/download"
	pathBrowser := pathBase + "/browser"
	rootPath := server.RootPath
	if server.Enable {
		LogInfo("file-server enabled, url: %v --> path: %v", pathBase, rootPath)
		if !server.DisableBrowser {
			LogInfo("file-server enabled browser web-ui: %v, such %v/videos/cat to view /videos/cat files", pathBrowser, pathBrowser)
		}
		if !server.DisableList {
			LogInfo("file-server enabled list api: GET %v/{subPath} , such GET %v/videos/cat", pathList, pathList)
		}
		if !server.DisableUpload {
			LogInfo("file-server enabled upload api: POST %v/{subPath} with file={binary}, such POST %v/video/dog file=dog.mp4", pathUpload, pathUpload)
		}
		if !server.DisableDownload {
			LogInfo("file-server enabled download api: GET %v/{subPath}[?type=inline], such GET %v/video/dog/dog.mp4 to download file, GET %v/video/dog/dog.mp4?type=inline to preview in browser", pathDownload, pathDownload, pathDownload)
		}
	}
	return func(c *gin.Context) {
		// 如果未开启文件服务，直接跳过
		if !server.Enable {
			c.Next()
			return
		}
		// 检查路径前缀匹配
		urlPath := c.Request.URL.Path
		if !server.DisableBrowser && strings.HasPrefix(urlPath, pathBrowser) {
			filePath := urlPath[len(pathBrowser):]
			LogInfo("goboot file-server, browser path: %v", filePath)
			fullPath := filepath.Join(rootPath, filePath)
			regularFilePath := strings.ReplaceAll(filePath, "\\", "/")

			// 检查文件是否在允许的目录内
			allowedPath, _ := filepath.Abs(rootPath)
			absPath, _ := filepath.Abs(fullPath)
			if !strings.HasPrefix(absPath, allowedPath) {
				c.JSON(500, ApiError(500, filePath+" not allow access!"))
				return
			}

			// 检查文件是否存在
			info, err := os.Stat(fullPath)
			if os.IsNotExist(err) {
				c.JSON(404, ApiError(500, filePath+" not exists!"))
				return
			}

			// 检查是否是目录
			if !info.IsDir() {
				c.JSON(500, ApiError(500, filePath+" is not directory!"))
				return
			}

			files, err := ListFiles(fullPath, rootPath)

			if err != nil {
				c.JSON(500, ApiError(500, filePath+" list error!"))
				return
			}
			if info != nil {
				parentDir := filepath.Dir(fullPath)
				relPath, _ := filepath.Rel(rootPath, parentDir)
				parant, _ := os.Stat(parentDir)
				files = append([]FileInfoItem{
					FileInfoItem{
						Path:       strings.ReplaceAll(relPath, "\\", "/"),
						Name:       "..",
						Size:       parant.Size(),
						SizeText:   ConvertAsHumanSizeText(parant.Size()),
						IsDir:      parant.IsDir(),
						ModifyTime: parant.ModTime().Format("2006-01-02 15:04:05"),
					},
				}, files...)
			}

			if info != nil {
				curr := info
				relPath, _ := filepath.Rel(rootPath, fullPath)
				files = append([]FileInfoItem{
					FileInfoItem{
						Path:       strings.ReplaceAll(relPath, "\\", "/"),
						Name:       ".",
						Size:       curr.Size(),
						SizeText:   ConvertAsHumanSizeText(curr.Size()),
						IsDir:      curr.IsDir(),
						ModifyTime: curr.ModTime().Format("2006-01-02 15:04:05"),
					},
				}, files...)
			}

			html := `
			<html lang="zh">
    <head>
        <meta charset="UTF-8"/>
		<meta name="viewport" content="width=device-width, user-scalable=yes, initial-scale=1.0, minimum-scale=0.5, maximum-scale=2.0">
        <title>
			`
			html = html + regularFilePath
			html = html +
				`
			</title>
		<style>
        .file-page *{
            box-sizing: content-box;
            margin: 0;
            padding: 0;
            font-size: 16px;
        }
		@media screen and (max-width:320px) {
			.file-page *{
            	font-size: 6px;
        	}
		}
		@media screen and (min-width:320px) and (max-width:480px) {
			.file-page *{
            	font-size: 8px;
        	}
		}
		@media screen and (min-width:480px) and (max-width:640px) {
			.file-page *{
            	font-size: 10px;
        	}
		}
		@media screen and (min-width:640px) and (max-width:960px) {
			.file-page *{
            	font-size: 14px;
        	}
		}
        .file-page button{
            border: none;
            outline: none;
            padding: 3px 5px;
            background: dodgerblue;
            color: white;
            border-radius: 3px;
        }
        .file-page button:hover{
            background: deepskyblue;
        }
        .file-page button:active{
            background: orange;
            transform: scale(0.9);
        }
		.file-page button:disabled{
			filter: grayscale(0.6);
		}
        .file-path{
            font-size: 22px;
            font-weight: bold;
        }
        .file-divider{
            margin: 5px 3px;
        }
        .file-list{
            list-style-type: none;
        }
        .file-item{
            border-bottom: solid 1px #ddd;
            padding: 3px 2px;
			display: flex;
        }
        .file-item span{
            display: flex;
			align-items: center;
        }
        .file-type{
            width: 3%;
        }
        .file-name{
            width: 43%;
            color: dodgerblue;
        }
        .file-size{
            width: 10%;
        }
        .file-time{
            width: 15%;
        }
        .file-operation{
            width: 25%;
        }
        .file-operation button{
            display: inline-block;
        }
		.file-upload{
            float: right;
            background: limegreen !important;
        }
    </style>
    </head>
    <body>
        <div class="file-page">
            <div class="file-path">`
			html = html + regularFilePath
			if !server.DisableUpload {
				html = html + `
				<button class="file-upload" id="fileUploadButton" onclick="uploadFile()">upload</button>
                <input class="file-upload" type="file" id="fileInputDom" onchange="onFileChange(this)" style="display: none;"/>
				`
			}
			html = html +
				`</div>
            <hr class="file-divider"/>
            <ul class="file-list" id="fileListDom">
            </ul>
        </div>
        
    </body>
    <script>
        const pathBase = "`
			html = html + pathBase
			html = html + `"

        const pathList = "`
			html = html + pathList
			html = html + `"

        const pathUpload = "`
			html = html + pathUpload
			html = html + `"

        const pathDownload = "`
			html = html + pathDownload
			html = html + `"

        const pathBrowser = "`
			html = html + pathBrowser
			html = html + `"

		const fileList = `
			html = html + Any2JsonString(files)
			html = html + `

		function renderFileList(){
			let listDom = document.querySelector("#fileListDom");
			listDom.innerHTML=''
			let html=''
			for (let i = 0; i < fileList.length; i++) {
				let item=fileList[i]
				html+="<li class=\"file-item\">\n" +
					"            <span class=\"file-type\">"+(item.isDir?"+":"-")+"</span>\n" +
					"            <span class=\"file-name\" onclick=\"openFile("+(i)+")\">"+(item.name)+"</span>\n" +
					"            <span class=\"file-size\">"+(item.sizeText)+"</span>\n" +
					"            <span class=\"file-time\">"+(item.modifyTime)+"</span>\n" +
					"            <span class=\"file-operation\">\n"
			`
			if !server.DisableDownload {
				html = html + `
				if(!item.isDir){
							html+="   <button class=\"file-download\" onclick=\"downloadFile("+(i)+")\">download</button>\n" 
						}
				`
			}
			html = html + `
				html+="            </span>\n" +
					"        </li>"
			}
			
			listDom.innerHTML=html
		}

		renderFileList()

        function getBasePath(){
            let basePath=''
            let curPath=window.location.pathname
            let idx = curPath.indexOf(pathBrowser)
            if(idx>=0){
                basePath=curPath.substring(0,idx)
            }
            return basePath
        }
        function openFile(index) {
			let fileItem=fileList[index]
			if (fileItem.isDir) {
				let nextPath = '/' + pathBrowser + '/' + encodeURI(fileItem.path)
				nextPath = nextPath.replaceAll('//', '/')
				window.location.href = getBasePath() + nextPath
			} else {
			`
			if !server.DisableDownload {
				html = html + `
					let nextPath = '/' + pathDownload + '/' + encodeURI(fileItem.path)
					nextPath = nextPath.replaceAll('//', '/')
					window.location.href = getBasePath() + nextPath + "?type=inline"
				`
			} else {
				html = html + `
					debugger
				`
			}
			html = html + `
			}
		}

		function downloadFile(index) {
			let fileItem=fileList[index]
			if (fileItem.isDir) {
				openFile(index)
			} else {
			 `
			if !server.DisableDownload {
				html = html + `
					let nextPath = '/' + pathDownload + '/' + encodeURI(fileItem.path)
					nextPath = nextPath.replaceAll('//', '/')
					let url = getBasePath() + nextPath + "?type=attachment"
				
					let dom = document.createElement('a');
					dom.href = url
					dom.download = fileItem.name
					dom.style.display = 'none'
					document.body.append(dom)
					dom.click()
					document.body.removeChild(dom)
				`
			} else {
				html = html + `
					debugger
				`
			}
			html = html + `
			}

		}`
			if !server.DisableUpload {
				html = html + `
			function uploadFile(){
				var dom = document.querySelector('#fileInputDom');
				dom.click()
			}
			function onFileChange(dom){
				console.log(dom,dom.files,dom.files[0])
				let url=window.location.href
				url=url.replace(pathBrowser,pathUpload)
				let data =new FormData();
				data.append('file',dom.files[0]);

				const option = {
					method: 'post',
					mode: 'cors',
					body: data
				};
				let btnDom=document.querySelector("#fileUploadButton")
				btnDom.setAttribute('disabled','disabled')
				fetch(url, option)
					.then(res=>res.json())
					.then(function (data) {
					console.log('imgUrl', data);
					location.reload()
				}).finally(()=>{
					btnDom.removeAttribute("disabled")
				})

			}
				`
			}

			html = html + `
    </script>
    
</html>
			`

			c.Data(200, "text/html; charset=utf-8", []byte(html))
			return
		}
		if !server.DisableList && strings.HasPrefix(urlPath, pathList) {
			filePath := urlPath[len(pathList):]
			LogInfo("goboot file-server, list path: %v", filePath)
			fullPath := filepath.Join(rootPath, filePath)

			// 检查文件是否在允许的目录内
			allowedPath, _ := filepath.Abs(rootPath)
			absPath, _ := filepath.Abs(fullPath)
			if !strings.HasPrefix(absPath, allowedPath) {
				c.JSON(200, ApiError(500, filePath+" not allow access!"))
				return
			}

			// 检查文件是否存在
			info, err := os.Stat(fullPath)
			if os.IsNotExist(err) {
				c.JSON(200, ApiError(500, filePath+" not exists!"))
				return
			}

			// 检查是否是目录
			if !info.IsDir() {
				c.JSON(200, ApiError(500, filePath+" is not directory!"))
				return
			}

			files, err := ListFiles(fullPath, rootPath)

			if err != nil {
				c.JSON(200, ApiError(500, filePath+" list error!"))
				return
			}

			c.JSON(200, ApiOk(files))
			return
		}
		if !server.DisableUpload && strings.HasPrefix(urlPath, pathUpload) {
			filePath := urlPath[len(pathUpload):]
			LogInfo("goboot file-server, upload path: %v", filePath)
			fullPath := filepath.Join(rootPath, filePath)

			// 检查文件是否在允许的目录内
			allowedPath, _ := filepath.Abs(rootPath)
			absPath, _ := filepath.Abs(fullPath)
			if !strings.HasPrefix(absPath, allowedPath) {
				c.JSON(500, ApiError(500, filePath+" not allow access!"))
				return
			}

			// 检查文件是否存在
			info, err := os.Stat(fullPath)
			if os.IsNotExist(err) {
				c.JSON(404, ApiError(500, filePath+" not exists!"))
				return
			}

			// 检查是否是目录
			if !info.IsDir() {
				c.JSON(500, ApiError(500, filePath+" is directory!"))
				return
			}

			file, err := c.FormFile("file")

			savePath := filepath.Join(fullPath, file.Filename)
			err = c.SaveUploadedFile(file, savePath)

			relPath, _ := filepath.Rel(rootPath, savePath)
			c.JSON(200, ApiOk(relPath))
			return
		}
		if !server.DisableDownload && strings.HasPrefix(urlPath, pathDownload) {
			filePath := urlPath[len(pathDownload):]
			LogInfo("goboot file-server, download path: %v", filePath)
			fullPath := filepath.Join(rootPath, filePath)

			// 检查文件是否在允许的目录内
			allowedPath, _ := filepath.Abs(rootPath)
			absPath, _ := filepath.Abs(fullPath)
			if !strings.HasPrefix(absPath, allowedPath) {
				c.JSON(500, ApiError(500, filePath+" not allow access!"))
				return
			}

			// 检查文件是否存在
			info, err := os.Stat(fullPath)
			if os.IsNotExist(err) {
				c.JSON(404, ApiError(500, filePath+" not exists!"))
				return
			}

			// 检查是否是目录
			if info.IsDir() {
				c.JSON(500, ApiError(500, filePath+" is directory!"))
				return
			}

			// 打开文件
			file, err := os.Open(fullPath)
			if err != nil {
				c.JSON(500, ApiError(500, filePath+" open file error!"))
				return
			}
			defer file.Close()

			// 获取是否强制下载
			downloadType := c.Query("type")

			// 获取文件扩展名和MIME类型
			ext := strings.ToLower(filepath.Ext(fullPath))
			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				if ext == ".flv" {
					mimeType = "video/flv"
				}
			}
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}

			// 不是内联模式，默认下载
			if !SliceContains([]string{
				"inline", "preview", "view",
			}, downloadType) {
				c.Header("Content-Disposition", "attachment; filename=\""+info.Name()+"\"")
			}

			// 已知类型（图片、音视频等）由浏览器决定如何处理
			c.Header("Content-Type", mimeType)

			if SliceContains([]string{
				"inline", "preview", "view",
			}, downloadType) {
				suffix := strings.ToLower(filepath.Ext(info.Name()))
				if SliceContains([]string{
					".txt", ".log", ".md",
					".bat", ".sh",
					".css", ".js", ".ts", ".sass", ".less",
					".yml", ".yaml", ".properties",
					".xml", ".json", ".jsonl", ".jsonc",
					".sql",
					".c", ".h", ".cpp", ".hpp", ".hxx",
					".java", ".py", ".go", ".pl",
					".gitignore", ".gitattributes",
					".vue",
				}, suffix) {
					c.Header("Content-Type", "text/plain")
					c.Status(http.StatusOK)
					c.Header("Content-Length", strconv.FormatInt(info.Size(), 10))
					_, err = io.Copy(c.Writer, file)
					if err != nil {
						c.JSON(500, ApiError(500, filePath+" send file error!"))
						return
					}
					return
				}
			}

			// 处理Range请求
			rangeHeader := c.GetHeader("Range")
			if rangeHeader != "" {
				// 解析Range头（仅支持bytes单位）
				var start, end int64
				parseRange := false
				ranges := strings.Split(rangeHeader, "=")[1]
				parts := strings.Split(ranges, "-")

				if len(parts) == 2 {
					// 处理起始位置
					if parts[0] != "" {
						start, _ = strconv.ParseInt(parts[0], 10, 64)
					} else {
						// 处理后缀范围（如：-500 bytes）
						if parts[1] != "" {
							suffix, _ := strconv.ParseInt(parts[1], 10, 64)
							start = info.Size() - suffix
							end = info.Size() - 1
							parseRange = true
						}
					}

					// 处理结束位置
					if parts[1] != "" && parseRange {
						end, _ = strconv.ParseInt(parts[1], 10, 64)
					} else if parts[1] != "" {
						end, _ = strconv.ParseInt(parts[1], 10, 64)
						parseRange = true
					} else {
						end = info.Size() - 1
						parseRange = true
					}

					// 校验范围有效性
					if parseRange && start <= end && end < info.Size() {
						// 设置响应头
						c.Status(http.StatusPartialContent)
						c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, info.Size()))
						c.Header("Accept-Ranges", "bytes")
						c.Header("Content-Length", strconv.FormatInt(end-start+1, 10))

						// 使用SectionReader高效读取片段
						sectionReader := io.NewSectionReader(file, start, end-start+1)
						_, err = io.Copy(c.Writer, sectionReader)
						if err != nil {
							c.JSON(500, ApiError(500, filePath+" send section error!"))
							return
						}
						return
					}
				}
			}

			// 无Range请求，返回完整文件
			c.Status(http.StatusOK)
			c.Header("Content-Length", strconv.FormatInt(info.Size(), 10))
			_, err = io.Copy(c.Writer, file)
			if err != nil {
				c.JSON(500, ApiError(500, filePath+" send file error!"))
				return
			}

			return
		}
		// 如果不匹配，继续执行
		c.Next()

	}
}

// 启动应用
func (boot *GobootApplication) Run() {
	engine := boot.App
	app := boot.Config.Goboot.Application
	profiles := boot.Config.Goboot.Profiles
	server := boot.Config.Goboot.Server

	LogInfo("goboot run ...")

	// 配置控制器
	if len(boot.Controllers) > 0 {
		LogInfo("goboot enbale %v controllers(s)", len(boot.Controllers))
		for _, item := range boot.Controllers {
			groupPath := item.Path()
			groupRouters := engine.Group(groupPath)
			{
				groupRouters.Any("/*proxyPath", func(c *gin.Context) {
					proxyPath := c.Param("proxyPath")
					MappingHandler(boot, c, proxyPath, item)
				})
			}
		}
	}

	LogInfo("goboot brfore banner.")
	invokeListeners(boot, boot.Listeners.OnBeforeBanner)

	if server.BannerPath != "" {
		bytes, err := ioutil.ReadFile(server.BannerPath)
		if err != nil {
			LogInfo("goboot default banner.")
			fmt.Println(DefaultBannerText)
		} else {
			LogInfo("goboot read banner, file: %v", server.BannerPath)
			fmt.Println(string(bytes))
		}
	}

	LogInfo("app [%v] on [%v] run at port [%v]", app.Name, profiles.Active, server.Port)
	LogInfo("local: http://localhost:%v/", server.Port)

	iters, err := net.Interfaces()
	if err == nil {
		for _, iter := range iters {
			if (net.FlagUp & iter.Flags) != 0 {
				continue
			}
			if (net.FlagLoopback & iter.Flags) != 0 {
				continue
			}
			addrs, err2 := iter.Addrs()
			if err2 == nil {
				LogInfo("[net] %v :", iter.Name)
				for _, addr := range addrs {
					ipNet, ok := addr.(*net.IPNet)

					if ok && !ipNet.IP.IsLoopback() && !ipNet.IP.IsMulticast() {
						LogInfo("\thttp://%v:%v/", ipNet.IP, server.Port)
					}
				}
			}
		}

	}

	LogInfo("goboot before run.")
	invokeListeners(boot, boot.Listeners.OnBeforeRun)

	if len(boot.Handlers) > 0 {
		LogInfo("goboot loaded mapping %v hander(s)", len(boot.Handlers))
		for _, handler := range boot.Handlers {
			rtp := reflect.TypeOf(handler)
			if rtp.Kind() == reflect.Ptr && rtp.Elem().Kind() == reflect.Struct {
				rtp = rtp.Elem()
			}
			LogInfo("mapping handler, type: %v", rtp.Name())
		}
	}

	bindStr := fmt.Sprintf(":%v", server.Port)
	if server.Https.Enable {
		LogInfo("goboot run with https, pem: %v, key: %v", server.Https.PemPath, server.Https.KeyPath)
		engine.Run(bindStr, server.Https.PemPath, server.Https.KeyPath)
	} else {
		LogInfo("goboot run with http.")
		engine.Run(bindStr)
	}

}
