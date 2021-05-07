package gee

import "net/http"

// 首先我们定义一个给大家科普一下:
//海王：吊着你的今晚吃啥
//炮王：约炮和吃饭一样的今晚吃啥
//直男：心直口快的低情商的今晚吃啥
//直女：心直口快的低情商女
//舔狗：爱而不得的今晚吃啥
//潮男：会打扮的今晚吃啥
//御姐：不爱搭理人的女的
//萝莉：大部分是装的
//渣男：今晚吃啥
//渣女：遍地都是
//好男孩：我

/**
首先我们定义 HandlerFunc 这是框架处理HTTP请求的类型  就是 Handler interface  的 ServeHTTP 方法
*/

type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine 我们定义一个结构体来实现- Handler interface  的 ServeHTTP 方法
type Engine struct {
	router map[string]HandlerFunc
}

// New 第一种，在Go语言中，可以直接以 var 的方式声明结构体即可完成实例化
//var t T
//t.a = 1
//t.b = 2
//
////第二种，使用 new() 实例化
//t := new(T)
//
////第三种，使用字面量初始化
//t := T{a, b}
//t := &T{} //等效于 new(T)/**
//  &Engine{router: make(map[string]HandlerFunc)}

func New() *Engine {
	engine := new(Engine)
	// map 需要使用make()函数来分配内存
	// 否则会出现 assignment to entry in nil map
	engine.router = make(map[string]HandlerFunc)
	return engine
}

// 添加路由到- route map 内
/**
method 方法
pattern /hello 路由地址
handler  就是 Handler interface  的 ServeHTTP 方法
*/
func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// http 行为 + 路由地址作为key
	key := method + "-" + pattern
	// 实现HandlerFunc 的方法作为值
	e.router[key] = handler
}

func (e *Engine) Run(address string) (err error) {

	// 因为Engine 实现- Handler interface
	return http.ListenAndServe(address, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 获取到-key  然后取出 实现HandlerFunc 的方法
	key := req.Method + "-" + req.URL.Path
	handler := e.router[key]
	handler(w, req)
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

func (e *Engine) PUT(pattern string, handler HandlerFunc) {
	e.addRoute("PUT", pattern, handler)
}
