package gee

import (
	"net/http"
	"strings"
)

// 定义一个存放路由的结构体

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
type router struct {
	roots    map[string]*node
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
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由到- route map 内
/**
method 方法
pattern /hello 路由地址
handler  就是 Handler interface  的 ServeHTTP 方法
*/
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// 处理运行--路由绑定的func handle 函数中，将从路由匹配得到的 Handler 添加到 c.handlers列表中，执行c.Next()。
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
