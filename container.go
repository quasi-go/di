package di

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

const (
	LOG_LEVEL_NONE    = 0
	LOG_LEVEL_ERROR   = 1
	LOG_LEVEL_WARNING = 2
	LOG_LEVEL_NOTICE  = 4
	LOG_LEVEL_DEFAULT = LOG_LEVEL_ERROR & LOG_LEVEL_WARNING & LOG_LEVEL_NOTICE
	LOG_LEVEL_TRACE   = 8
	LOG_LEVEL_ALL     = 15
)

type Rule interface {
	Resolve(c *Container) (reflect.Value, error)
}

type typeRule struct {
	typeTo reflect.Type
}

func (r *typeRule) Resolve(c *Container) (reflect.Value, error) {
	return c.ResolveType(r.typeTo)
}

type instanceRule struct {
	instance reflect.Value
}

func (r *instanceRule) Resolve(_ *Container) (reflect.Value, error) {
	return r.instance, nil
}

type factoryRule struct {
	callback any
}

func (r *factoryRule) Resolve(c *Container) (reflect.Value, error) {
	_, err := validateFactoryCallback(r.callback)

	if err != nil {
		return reflect.Value{}, err
	}

	returnValue, err := c.Call(r.callback)

	if err != nil {
		return reflect.Value{}, err
	}

	if len(returnValue) != 2 {
		return reflect.Value{}, errors.New("callback must return one value and an error")
	}

	return returnValue[0], nil
}

type providerRule struct {
	factoryRule
	instance *reflect.Value
}

func (r *providerRule) Resolve(c *Container) (reflect.Value, error) {
	if r.instance == nil {
		instance, err := r.factoryRule.Resolve(c)

		if err != nil {
			return reflect.Value{}, err
		}

		r.instance = &instance
	}

	return *r.instance, nil
}

type ruleStore map[Id]Rule

type Container struct {
	rules    ruleStore
	logger   *log.Logger
	logLevel int
}

var currentContainer = NewContainer()

func SetContainer(c *Container) {
	currentContainer = c
}

func GetContainer() *Container {
	return currentContainer
}

func NewContainer() *Container {
	return &Container{
		rules:    make(ruleStore),
		logLevel: LOG_LEVEL_DEFAULT,
	}
}

func resetContainer() {
	SetContainer(NewContainer())
}

func (c *Container) SetLogger(logger *log.Logger) {
	c.logger = logger
}

func (c *Container) SetLogLevel(logLevel int) {
	c.logLevel = logLevel
}

func (c *Container) SetRule(key Id, rule Rule) {
	if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_TRACE) {
		c.logger.Printf("Setting %s as (%s): %+v", key, reflect.TypeOf(rule).String(), rule)
	}

	c.rules[key] = rule
}

func (c *Container) HasRule(key Id) bool {
	_, exists := c.rules[key]
	return exists
}

func (c *Container) GetRule(key Id) Rule {
	value, _ := c.rules[key]
	return value
}

func (c *Container) ResolveType(typeInfo reflect.Type) (reflect.Value, error) {
	if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_TRACE) {
		c.logger.Printf("Resolving %s", typeInfo.String())
	}

	if typeInfo.Kind() == reflect.Pointer {
		typeInfo = typeInfo.Elem()
	}

	typeId := Id(typeInfo.String())

	if !c.HasRule(typeId) {
		if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_TRACE) {
			c.logger.Printf("Rule %s not found", typeInfo.String())
		}

		built, err := c.BuildType(typeInfo)

		if err != nil {
			return reflect.Zero(typeInfo), err
		}

		c.SetRule(typeId, &instanceRule{built})
	}

	return c.GetRule(typeId).Resolve(c)
}

func (c *Container) BuildType(typeInfo reflect.Type) (reflect.Value, error) {
	if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_TRACE) {
		c.logger.Printf("Building %s", typeInfo)
	}

	structPtr := reflect.New(typeInfo)

	if typeInfo.Kind() != reflect.Struct {
		if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_TRACE) {
			c.logger.Printf("%s is not a struct, returning without resolving children", typeInfo.String())
		}

		return structPtr, nil
	}

	structElem := structPtr.Elem()

	for i := 0; i < typeInfo.NumField(); i++ {
		typeField := typeInfo.Field(i)
		structField := structElem.Field(i)
		inject, err := c.shouldInject(typeField)

		if err != nil {
			if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_ERROR) {
				c.logger.Printf("ERROR: %s", err)
			}

			return reflect.Zero(typeInfo), err
		}

		if !inject {
			if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_TRACE) {
				c.logger.Printf("Tagged as @noinject, skipping %s", typeField.Name)
			}

			continue
		}

		if !structField.CanSet() {
			if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_WARNING) {
				c.logger.Printf("WARNING: Can't set private member `%s` of `%s`. You need to make this member public to "+
					"inject it or add the tag `inject:\"@noinject\"` to mark that the field is skipped", typeField.Name, typeInfo.String())
			}

			continue
		}

		if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_TRACE) {
			c.logger.Printf("Resolving child %s", typeField.Name)
		}

		childType := typeField.Type
		isInterface := childType.Kind() == reflect.Interface
		isPointer := childType.Kind() == reflect.Pointer

		if isInterface && !c.HasRule(Id(childType.String())) {
			if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_TRACE) {
				c.logger.Printf("%s is an interface but has no rule set, skipping", childType.String())
			}
			continue
		}

		if isPointer {
			childType = childType.Elem()
		}

		builtChild, err := c.ResolveType(childType)

		if err != nil {
			if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_TRACE) {
				c.logger.Printf("ERROR: %s", err)
			}

			return reflect.Zero(typeInfo), err
		}

		var elem reflect.Value

		if isPointer || isInterface {
			elem = builtChild
		} else {
			elem = builtChild.Elem()
		}

		if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_TRACE) {
			c.logger.Printf("Setting child %#v", elem.Interface())
		}

		structField.Set(elem)
	}

	if structPtr.Type().Implements(Type[Initializable]()) {
		if c.logger != nil && hasLogLevel(c.logLevel, LOG_LEVEL_TRACE) {
			c.logger.Printf("Initializing %s", structPtr.Type().String())
		}

		structPtr.MethodByName("Initialize").Call([]reflect.Value{})
	}

	return structPtr, nil
}

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct ||
		(t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct)
}

func isInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Interface
}

func isStructOrInterface(t reflect.Type) bool {
	return isStruct(t) || isInterface(t)
}

func (c *Container) shouldInject(field reflect.StructField) (bool, error) {
	injectTag := field.Tag.Get("inject")
	t := field.Type

	var inject bool

	switch injectTag {
	case "":
		inject = true
	case "@noinject":
		return false, nil
	default:
		errorMessage := fmt.Sprintf("Invalid `inject` tag value \"%s\" on member %s",
			injectTag, t.String())
		return false, errors.New(errorMessage)
	}

	if !inject {
		return false, nil
	}

	canConstruct := isStructOrInterface(t)

	if !canConstruct {
		return false, nil
	}

	return true, nil
}

func (c *Container) Call(callback any) ([]reflect.Value, error) {
	funcType := reflect.TypeOf(callback)
	funcValue := reflect.ValueOf(callback)

	var args []reflect.Value

	for i := 0; i < funcType.NumIn(); i++ {
		argType := funcType.In(i)
		arg, err := c.ResolveType(argType)

		if err != nil {
			message := fmt.Sprintf("could not resolve argument #%d of callback; type %s could not be resolved: %s", i, argType, err)
			return nil, errors.New(message)
		}

		if argType.Kind() != reflect.Pointer && argType.Kind() != reflect.Interface {
			arg = arg.Elem()
		}

		args = append(args, arg)
	}

	return funcValue.Call(args), nil
}

func validateFactoryCallback(callback any) (reflect.Type, error) {
	typeInfo := reflect.TypeOf(callback)

	if typeInfo.Kind() != reflect.Func {
		return reflect.TypeOf(nil), errors.New("callback must be a function")
	}

	if typeInfo.NumOut() != 2 {
		return reflect.TypeOf(nil), errors.New("callback must have only one return value and one error value")
	}

	errorType := typeInfo.Out(1)

	if errorType != Type[error]() {
		return reflect.TypeOf(nil), errors.New("the second return value must be an error")
	}

	returnType := typeInfo.Out(0)

	if returnType.Kind() == reflect.Pointer {
		return returnType.Elem(), nil
	} else if returnType.Kind() == reflect.Interface {
		return returnType, nil
	}

	return reflect.TypeOf(nil), errors.New("callback must return an interface or a pointer to the constructed value")
}

func hasLogLevel(value int, test int) bool {
	return (value & test) != 0
}
