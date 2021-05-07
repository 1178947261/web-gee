package gee

import "net/http"

//将路由(router)独立出来，方便之后增强。
//设计上下文(Context)，封装 Request 和 Response ，提供对 JSON、HTML 等返回类型的支持。

// H 是我们定义响应数据的结构
type H map[string]interface{}

// Context 定义上下文结构体
type Context struct {
	// 原始对象
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求的信息
	// 请求URL 地址
	Path string
	// 请求方法
	Method string
	// 响应的状态码
	StatusCode int
}
