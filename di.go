package di

import (
	"errors"
	"fmt"
	"reflect"
)

type definitions struct {
	dependencies int
	constructor  reflect.Value
}

// Container ...
type Container struct {
	// Хранить созданные компоненты
	components map[reflect.Type]interface{}
	// Хранить определения на базе которых создаются компоненты
	definitions map[reflect.Type]definitions
}

// NewContainer ...
func NewContainer() *Container {
	return &Container{
		components:  make(map[reflect.Type]interface{}),
		definitions: make(map[reflect.Type]definitions),
	}
}

// Provide - регистрация компонентов + их связывание (wire)
func (c *Container) Provide(constructors ...interface{}) {
	c.register(constructors)
	c.wire()
}

// register of dependencies
func (c *Container) register(constructors []interface{}) {
	for _, constructor := range constructors {
		constructorType := reflect.TypeOf(constructor)
		if constructorType.Kind() != reflect.Func {
			panic(fmt.Errorf("%s must be constructor", constructorType.Name()))
		}

		outType := constructorType.Out(0) // constructor must return component

		if _, exists := c.definitions[outType]; exists {
			panic(fmt.Errorf("ambigious definitions %s already exists", constructorType.Name()))
		}

		paramsCount := constructorType.NumIn()
		c.definitions[outType] = definitions{
			dependencies: paramsCount,
			constructor:  reflect.ValueOf(constructor),
		}
	}
}

// wire -
func (c *Container) wire() {
	rest := make(map[reflect.Type]definitions, len(c.definitions))
	for key, value := range c.definitions {
		rest[key] = value
	}

	for {
		wired := 0

		for key, value := range rest {
			depsValues := make([]reflect.Value, 0)
			for i := 0; i < value.dependencies; i++ {
				depType := value.constructor.Type().In(i)
				if dep, exists := c.components[depType]; exists {
					depsValues = append(depsValues, reflect.ValueOf(dep))
				}
			}

			if len(depsValues) == value.dependencies {
				wired++
				component := value.constructor.Call(depsValues)[0].Interface()
				c.components[key] = component
				delete(rest, key)
				continue
			}
		}

		if len(rest) == 0 {
			return
		}

		if wired == 0 {
			// TODO: return component list (rest)
			panic(errors.New("some components has unmet dependencies"))
		}
	}
}

// Component ...
func (c *Container) Component(target interface{}) {
	if target == nil {
		panic("errors target cannot be nil")
	}

	targetValue := reflect.ValueOf(target)       // get value
	targetType := targetValue.Type()             // get value type
	targetValueType := targetValue.Elem().Type() // get value type type

	value, ok := c.components[targetValueType]
	if !ok {
		panic(errors.New("no such components"))
	}

	if targetType.Kind() != reflect.Ptr || targetValue.IsNil() {
		panic(errors.New("target must be a non-nil pointer"))
	}

	targetElemType := targetType.Elem()
	if !reflect.TypeOf(value).AssignableTo(targetElemType) {
		panic(errors.New("can't assign component to pointer"))
	}

	targetValue.Elem().Set(reflect.ValueOf(value))
}
