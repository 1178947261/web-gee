package gee

import (
	"log"
	"net/http"
)

// 定义一个存放路由的结构体

type router struct {
	handlers map[string]HandlerFunc
}

// New 第一种，在Go语言中，可以直接以 var 的方式声明结构体即可完成实例化
// var t T
//t.a = 1
//t.b = 2
//
//第二种，使用 new() 实例化
//t := new(T)
//
//第三种，使用字面量初始化
//t := T{a, b}
//t := &T{} //等效于 new(T)/**
// &router{handlers: make(map[string]HandlerFunc)}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// 添加路由到- route map 内
/**
method 方法
pattern /hello 路由地址
handler  就是 Handler interface  的 ServeHTTP 方法
*/
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 处理运行--路由绑定的func
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
