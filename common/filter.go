package common

import (
	"net/http"
	"strings"
)

// FilterHandle 声明一个新的数据类型（函数类型）
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

// Filter 拦截器结构体
type Filter struct {
	filterMap map[string]FilterHandle //用来存储需要拦截的URI
}

// NewFilter Filter初始化函数
func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

// RegisterFilterUri 注册拦截器
func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

// GetFilterHandle 根据Uri获取对应的handle
func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

// WebHandle 声明新的函数类型
type WebHandle func(rw http.ResponseWriter, req *http.Request)

// Handle 执行拦截器，返回函数类型
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if strings.Contains(r.RequestURI, path) { //执行拦截业务逻辑
				err := handle(rw, r)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				break //跳出循环
			}
		}
		webHandle(rw, r) //执行正常注册的函数
	}
}
