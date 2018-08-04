package syringe

import (
	"reflect"
)

func NewContainer() *Container {
	return &Container{
		definitions: make(map[string]*definition, 0),
	}
}

type Container struct {
	definitions map[string]*definition
	parent      *Container
}

func (c *Container) RegisterType(name string, module interface{}, lc Lifecycle) {
	c.definitions[name] = &definition{
		module:    reflect.ValueOf(module),
		lifecycle: lc,
		resolved:  false,
	}
}

func (c *Container) RegisterValue(name string, value interface{}) {
	c.definitions[name] = &definition{
		module:    reflect.ValueOf(value),
		lifecycle: Singleton,
		resolved:  true,
	}
}

func (c *Container) Resolve(name string) interface{} {
	instance := resolveDependency(c, name)
	return instance.Interface()
}

func (c *Container) Scope() *Container {
	nc := NewContainer()
	nc.parent = c
	return nc
}
