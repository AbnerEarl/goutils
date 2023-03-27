package gins

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Context struct {
	*gin.Context
}
type Server struct {
	*gin.Engine
}
type RouterGroup struct {
	*gin.RouterGroup
}

func HandleFunc(handler func(c *Context)) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		handler(&Context{Context: c})
	}
}
func NewServer() *Server {
	server := &Server{Engine: gin.Default()}
	return server
}

func (c *Context) ArgsObject(key string) interface{} {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		return data
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsObjectDefault(key string, val interface{}) interface{} {
	data, exists := c.Get(key)
	if exists && data != nil {
		return data
	} else {
		return val
	}

}

func (c *Context) ArgsInt64(key string) int64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.ParseInt(fmt.Sprint(data), 10, 64)
		if e == nil {
			return result
		} else {
			err = e
		}
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsInt64Default(key string, val int64) int64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.ParseInt(fmt.Sprint(data), 10, 64)
		if e == nil {
			return result
		} else {
			err = e
		}
	} else {
		return val
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsString(key string) string {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		return fmt.Sprint(data)
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsStringDefault(key string, val string) string {
	data, exists := c.Get(key)
	if exists && data != nil {
		return fmt.Sprint(data)
	} else {
		return val
	}
}

func (c *Context) ArgsInt64Array(key string) []int64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []int64
		tag := true
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				v, e := strconv.ParseInt(fmt.Sprint(item), 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				v, e := strconv.ParseInt(item, 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		}
		if tag {
			return result
		}
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsInt64ArrayDefault(key string, val []int64) []int64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []int64
		tag := true
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				v, e := strconv.ParseInt(fmt.Sprint(item), 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				v, e := strconv.ParseInt(item, 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		}
		if tag {
			return result
		}
	} else {
		return val
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsUint64Array(key string) []uint64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []uint64
		tag := true
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				v, e := strconv.ParseUint(fmt.Sprint(item), 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				v, e := strconv.ParseUint(item, 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		}
		if tag {
			return result
		}
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsUint64ArrayDefault(key string, val []uint64) []uint64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []uint64
		tag := true
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				v, e := strconv.ParseUint(fmt.Sprint(item), 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				v, e := strconv.ParseUint(item, 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		}
		if tag {
			return result
		}
	} else {
		return val
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsUintArray(key string) []uint {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []uint
		tag := true
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				v, e := strconv.ParseUint(fmt.Sprint(item), 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, uint(v))
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				v, e := strconv.ParseUint(item, 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, uint(v))
			}
		}
		if tag {
			return result
		}
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsUintArrayDefault(key string, val []uint) []uint {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []uint
		tag := true
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				v, e := strconv.ParseUint(fmt.Sprint(item), 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, uint(v))
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				v, e := strconv.ParseUint(item, 10, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, uint(v))
			}
		}
		if tag {
			return result
		}
	} else {
		return val
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsIntArray(key string) []int {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []int
		tag := true
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				v, e := strconv.Atoi(fmt.Sprint(item))
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				v, e := strconv.Atoi(item)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		}
		if tag {
			return result
		}
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsIntArrayDefault(key string, val []int) []int {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []int
		tag := true
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				v, e := strconv.Atoi(fmt.Sprint(item))
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				v, e := strconv.Atoi(item)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		}
		if tag {
			return result
		}
	} else {
		return val
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsFile(key string) (*multipart.FileHeader, error) {
	return c.FormFile(key)
}

func (c *Context) ArgsFileDefault(key string, file *multipart.FileHeader) (*multipart.FileHeader, error) {
	f, e := c.FormFile(key)
	if e != nil {
		return file, nil
	}
	return f, nil
}

func (c *Context) ArgsFloat64Array(key string) []float64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []float64
		tag := true
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				v, e := strconv.ParseFloat(fmt.Sprint(item), 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				v, e := strconv.ParseFloat(item, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		}
		if tag {
			return result
		}
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsFloat64ArrayDefault(key string, val []float64) []float64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []float64
		tag := true
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				v, e := strconv.ParseFloat(fmt.Sprint(item), 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				v, e := strconv.ParseFloat(item, 64)
				if e != nil {
					err = e
					tag = false
					break
				}
				result = append(result, v)
			}
		}
		if tag {
			return result
		}
	} else {
		return val
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsStringArray(key string) []string {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []string
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				result = append(result, fmt.Sprint(item))
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				result = append(result, fmt.Sprint(item))
			}
		}
		return result
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsStringArrayDefault(key string, val []string) []string {
	data, exists := c.Get(key)
	if exists && data != nil {
		var result []string
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			arr := data.([]interface{})
			for _, item := range arr {
				result = append(result, fmt.Sprint(item))
			}
		} else {
			arr := strings.Split(fmt.Sprint(data), ",")
			for _, item := range arr {
				result = append(result, fmt.Sprint(item))
			}
		}
		return result
	} else {
		return val
	}
}

func (c *Context) ArgsInt(key string) int {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.Atoi(fmt.Sprint(data))
		if e == nil {
			return result
		} else {
			err = e
		}
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsIntDefault(key string, val int) int {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.Atoi(fmt.Sprint(data))
		if e == nil {
			return result
		} else {
			err = e
		}
	} else {
		return val
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsUint64(key string) uint64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.ParseUint(fmt.Sprint(data), 10, 64)
		if e == nil {
			return result
		} else {
			err = e
		}
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsUint64Default(key string, val uint64) uint64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.ParseUint(fmt.Sprint(data), 10, 64)
		if e == nil {
			return result
		} else {
			err = e
		}
	} else {
		return val
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsFloat64(key string) float64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.ParseFloat(fmt.Sprint(data), 64)
		if e == nil {
			return result
		} else {
			err = e
		}
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsFloat64Default(key string, val float64) float64 {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.ParseFloat(fmt.Sprint(data), 64)
		if e == nil {
			return result
		} else {
			err = e
		}
	} else {
		return val
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsUint(key string) uint {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.ParseUint(fmt.Sprint(data), 10, 64)
		if e == nil {
			return uint(result)
		} else {
			err = e
		}
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsUintDefault(key string, val uint) uint {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.ParseUint(fmt.Sprint(data), 10, 64)
		if e == nil {
			return uint(result)
		} else {
			err = e
		}
	} else {
		return val
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsBool(key string) bool {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.ParseBool(fmt.Sprint(data))
		if e == nil {
			return result
		} else {
			err = e
		}
	} else {
		err = fmt.Errorf("the param '%s' is not exsit", key)
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (c *Context) ArgsBoolDefault(key string, val bool) bool {
	var err error
	data, exists := c.Get(key)
	if exists && data != nil {
		result, e := strconv.ParseBool(fmt.Sprint(data))
		if e == nil {
			return result
		} else {
			err = e
		}
	} else {
		return val
	}
	c.Abort()
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "参数解析异常",
		"data":    err,
	})
	panic(err)
}

func (s *Server) Group(relativePath string, handlers ...func(c *Context)) *RouterGroup {
	//Group 重写路由组注册
	sHandlers := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		sHandlers = append(sHandlers, HandleFunc(handle))
	}
	return &RouterGroup{s.Engine.Group(relativePath, sHandlers...)}
}

func (s *Server) GET(relativePath string, handlers ...func(c *Context)) gin.IRoutes {
	//GET 拓展Get请求（根）
	sHandlers := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		sHandlers = append(sHandlers, HandleFunc(handle))
	}
	return s.Engine.GET(relativePath, sHandlers...)
}

func (s *Server) POST(relativePath string, handlers ...func(c *Context)) gin.IRoutes {
	//POST 拓展POST请求（根）
	sHandlers := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		sHandlers = append(sHandlers, HandleFunc(handle))
	}
	return s.Engine.POST(relativePath, sHandlers...)
}

func (s *Server) PUT(relativePath string, handlers ...func(c *Context)) gin.IRoutes {
	//POST 拓展PUT请求（根）
	sHandlers := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		sHandlers = append(sHandlers, HandleFunc(handle))
	}
	return s.Engine.PUT(relativePath, sHandlers...)
}

func (s *Server) DELETE(relativePath string, handlers ...func(c *Context)) gin.IRoutes {
	//POST 拓展DELETE请求（根）
	sHandlers := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		sHandlers = append(sHandlers, HandleFunc(handle))
	}
	return s.Engine.DELETE(relativePath, sHandlers...)
}

func (s *Server) ANY(relativePath string, handlers ...func(c *Context)) gin.IRoutes {
	//POST 拓展DELETE请求（根）
	sHandlers := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		sHandlers = append(sHandlers, HandleFunc(handle))
	}
	return s.Engine.Any(relativePath, sHandlers...)
}

func (r *RouterGroup) GET(relativePath string, handlers ...func(c *Context)) gin.IRoutes {
	//GET 拓展Get请求（子）
	rHandlers := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandlers = append(rHandlers, HandleFunc(handle))
	}
	return r.RouterGroup.GET(relativePath, rHandlers...)
}

func (r *RouterGroup) POST(relativePath string, handlers ...func(c *Context)) gin.IRoutes {
	//POST 拓展Post请求（子）
	rHandlers := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandlers = append(rHandlers, HandleFunc(handle))
	}
	return r.RouterGroup.POST(relativePath, rHandlers...)
}

func (r *RouterGroup) PUT(relativePath string, handlers ...func(c *Context)) gin.IRoutes {
	//GET 拓展Get请求（子）
	rHandlers := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandlers = append(rHandlers, HandleFunc(handle))
	}
	return r.RouterGroup.PUT(relativePath, rHandlers...)
}

func (r *RouterGroup) DELETE(relativePath string, handlers ...func(c *Context)) gin.IRoutes {
	//POST 拓展Post请求（子）
	rHandlers := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandlers = append(rHandlers, HandleFunc(handle))
	}
	return r.RouterGroup.DELETE(relativePath, rHandlers...)
}

func (r *RouterGroup) ANY(relativePath string, handlers ...func(c *Context)) gin.IRoutes {
	//POST 拓展Post请求（子）
	rHandlers := make([]gin.HandlerFunc, 0)
	for _, handle := range handlers {
		rHandlers = append(rHandlers, HandleFunc(handle))
	}
	return r.RouterGroup.Any(relativePath, rHandlers...)
}

func (r *RouterGroup) Use(middlewares ...func(c *Context)) gin.IRoutes {
	//Use 拓展中间件注册
	rMiddlewares := make([]gin.HandlerFunc, 0)
	for _, middleware := range middlewares {
		rMiddlewares = append(rMiddlewares, HandleFunc(middleware))
	}
	return r.RouterGroup.Use(rMiddlewares...)
}
