package syringe

import "reflect"

type definition struct {
	module    reflect.Value
	lifecycle Lifecycle
	resolved  bool
}
