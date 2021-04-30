package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 返回的数据结构
type H map[string]interface{}

// 结构体- 可以看作是一个类 参数就是类属性
type Context struct {

	// HTTP 请求的原始对象
	Writer  http.ResponseWriter
	Request  *http.Request
	// HTTP 请求的信息
	//URL -
	Path string
	// 请求的方法类型列-POST -GET -PUT  等
	Method string
	// 响应请求的状态码 -
	StatusCode int

}

/**
 	开始封装 常用方法
 */

// 初始化-结构体

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Request:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}


// 获取POST参数

func (c Context) PostForm(key string) string {

	return c.Request.PostFormValue(key)
}

// 获取GET 参数

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// 设置返回的HTTP状态码

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}


// 设置响应的返回头

func (c Context) SetHeader(key string,value string)  {
	c.Writer.Header().Set(key,value)
}

// 返回字符串

func (c Context)String(code int, format string, values ...interface{})  {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 返回json 数据

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
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