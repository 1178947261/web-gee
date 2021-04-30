package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc定义路由进来的请求方法-

type HandlerFunc func(*Context)

// Engine  实现了HTTP ServeHTTP 的接口
// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router
}

// New 这是构造函数 gee.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET 定义添加GET请求的方法
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST 定义添加POST请求的方法
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run 定义启动http服务器的方法
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// handler = type HandlerFunc func(*Context)
	c := newContext(w, req)
	engine.router.handle(c)
}
