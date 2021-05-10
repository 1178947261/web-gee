package gee

import "net/http"

/**
首先我们定义 HandlerFunc 这是框架处理HTTP请求的类型  就是 Handler interface  的 ServeHTTP 方法
*/

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine 实现ServeHTTP的接口
type Engine struct {
	router *router
}

// New is the constructor of gee.Engine 使用 new() 实例化
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 实现ServeHTTP 方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 初始化Context结构体
	c := NewContext(w, req)
	// 进入路由处理--运行绑定在路由上的 func
	engine.router.handle(c)
}
