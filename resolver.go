package syringe

import (
	"log"
	"reflect"
	"unsafe"
)

func findDefinition(c *Container, name string) *definition {
	definition, ok := c.definitions[name]
	if ok {
		return definition
	}

	if c.parent != nil {
		return findDefinition(c.parent, name)
	}

	return nil
}

func resolveDependency(c *Container, name string) reflect.Value {
	definition := findDefinition(c, name)
	if definition.resolved {
		return definition.module
	}

	instancePtr := definition.module
	moduleType := instancePtr.Elem().Type()
	numFields := moduleType.NumField()

	if definition.lifecycle != Singleton {
		instancePtr = reflect.New(moduleType)
	}

	for i := 0; i < numFields; i++ {
		field := moduleType.Field(i)
		depName := field.Tag.Get(injectorTag)

		if depName != "" {
			dep := resolveDependency(c, depName)
			instanceField := instancePtr.Elem().FieldByName(field.Name)
			if !instanceField.CanSet() {
				instanceField = reflect.NewAt(instanceField.Type(), unsafe.Pointer(instanceField.UnsafeAddr())).Elem()
			}
			instanceField.Set(dep)
		} else {
			log.Fatalf("Dependency %s<%s> not found for %s", field.Name, field.Type.String(), instancePtr.String())
		}
	}

	if definition.lifecycle == Singleton && !definition.resolved {
		definition.resolved = true
	}

	return instancePtr
}
