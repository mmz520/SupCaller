package router

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	method  string
	path    string
	handler gin.HandlerFunc
}

var RouterPool map[string]RouterRegisterInterface = make(map[string]RouterRegisterInterface)

type RouterRegisterInterface interface {
	Register(group string, setControllers func(g *gin.RouterGroup))
	setGroup(g *gin.RouterGroup)
	InitGroup()
}

type RouterRegister struct {
	g              *gin.RouterGroup
	setControllers func(g *gin.RouterGroup)
}

func NewRouterRegister() *RouterRegister {
	return &RouterRegister{}
}

func (r *RouterRegister) Register(group string, setControllers func(g *gin.RouterGroup)) {
	r.setControllers = setControllers
	RouterPool[group] = r
}

func (r *RouterRegister) setGroup(g *gin.RouterGroup) {
	r.g = g
}

func (r *RouterRegister) InitGroup() {
	if r.setControllers != nil {
		r.setControllers(r.g)
	}
}

func AutoRegister(r *gin.RouterGroup) {
	for group, routes := range RouterPool {
		g := r.Group(group)
		routes.setGroup(g)
		routes.InitGroup()
	}
}
