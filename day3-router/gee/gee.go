package gee

import (
	"log"
	"net/http"
	"strings"
)

/**
首先我们定义 HandlerFunc 这是框架处理HTTP请求的类型  就是 Handler interface  的 ServeHTTP 方法
*/

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine 实现ServeHTTP的接口
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups
}

// RouterGroup 需要有访问Router的能力，为了方便，我们可以在Group中，保存一个指针，指向Engine，
//  整个框架的所有资源都是由Engine统一协调的，那么就可以通过Engine间接地访问各种接口了。
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

// New is the constructor of gee.Engine 使用 new() 实例化
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group 被定义为创建一个新的 RouterGroup
// 记住所有组共享同一个 Engine 实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	//添加路由
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 实现ServeHTTP 方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	// 获取 分组里面的所有 middlewares 的HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := NewContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

// Use 将中间件添加到组中
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// gee.go
// Default use Logger() & Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}
