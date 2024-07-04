package main

import (
	"goboot/goboot"
	"net/http"

	// "time"

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

	// app.AddHandlers(&Api{}).
	// 	AddHandlers(&User{})

	// app.AddControllers(&AdminController{})

	// app.App.GET("/", func(c *gin.Context) {
	// 	stime := time.Now().Format("2006-01-02 15:04:05")
	// 	c.HTML(200, "index/index.html", gin.H{
	// 		"now": stime,
	// 	})
	// 	// c.String(200, stime)
	// })

	app.Run()
}
