package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

	Params map[string]string
	// middleware
	handlers []HandlerFunc
	index    int
}

//使用 new() 实例化 初始化结构体里面的字段

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

//index是记录当前执行到第几个中间件，当在中间件中调用Next方法时，
//控制权交给了下一个中间件，直到调用到最后一个中间件，
//然后再从后往前，调用每个中间件在Next方法之后定义的部分。

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Json 返回JSON数据
func (c Context) Json(code int, obj interface{}) {

	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	//NewEncoder创建一个将数据写入w的*Encoder。 (约等于 new 了一个类) C.Writer 就是我们要写入的流
	encoder := json.NewEncoder(c.Writer)
	//Encode将v的json编码写入输出流，并会写入一个换行符   Encode 就是把OBJ 对象写入到  C.Writer
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.Json(code, H{"message": err})
}
