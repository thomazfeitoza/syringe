package syringe

import (
	"log"
	"reflect"
	"unsafe"
)

const LifetimeTransient = 1
const LifetimeSingleton = 2

const injectorTag = "di"

type definition struct {
	module   reflect.Value
	lifetime uint8
	resolved bool
}

func NewContainer() *Container {
	return &Container{
		definitions: make(map[string]*definition, 0),
	}
}

type Container struct {
	definitions map[string]*definition
	parent      *Container
}

func (c *Container) RegisterStruct(name string, module interface{}, lifetime uint8) {
	c.definitions[name] = &definition{
		module:   reflect.ValueOf(module),
		lifetime: lifetime,
		resolved: false,
	}
}

func (c *Container) RegisterValue(name string, value interface{}) {
	c.definitions[name] = &definition{
		module:   reflect.ValueOf(value),
		lifetime: LifetimeSingleton,
		resolved: true,
	}
}

func (c *Container) getValue(name string) *definition {
	definition, exists := c.definitions[name]
	if exists {
		return definition
	} else if c.parent != nil {
		return c.parent.getValue(name)
	}
	return nil
}

func (c *Container) resolveRecursive(name string) reflect.Value {
	definition := c.getValue(name)
	if definition.resolved {
		return definition.module
	}

	instancePtr := definition.module
	moduleType := instancePtr.Elem().Type()
	numFields := moduleType.NumField()

	if definition.lifetime != LifetimeSingleton {
		instancePtr = reflect.New(moduleType)
	}

	for i := 0; i < numFields; i++ {
		field := moduleType.Field(i)
		depName := field.Tag.Get(injectorTag)
		if depName != "" {
			dep := c.resolveRecursive(depName)
			instanceField := instancePtr.Elem().FieldByName(field.Name)
			if !instanceField.CanSet() {
				instanceField = reflect.NewAt(instanceField.Type(), unsafe.Pointer(instanceField.UnsafeAddr())).Elem()
			}
			instanceField.Set(dep)
		} else {
			log.Fatalf("Dependency %s<%s> not found for %s", field.Name, field.Type.String(), instancePtr.String())
		}
	}

	if definition.lifetime == LifetimeSingleton && !definition.resolved {
		definition.resolved = true
	}
	return instancePtr
}

func (c *Container) Resolve(name string) interface{} {
	instance := c.resolveRecursive(name)
	return instance.Interface()
}

func (c *Container) Scope() *Container {
	nc := NewContainer()
	nc.parent = c
	return nc
}
