package di

import (
	"errors"
	"fmt"
	"reflect"
)

type Initializeable interface {
	Initialize()
}

func Convert[T any](value reflect.Value) (*T, error) {
	if value.Kind() != reflect.Pointer {
		errorMessage := fmt.Sprintf("expected pointer, got %s (%s)", value.Type(), value.Kind())
		return nil, errors.New(errorMessage)
	}

	objectPtr := value.Interface()
	typeInfo := Type[T]()

	if typeInfo.Kind() == reflect.Struct {
		if result, ok := objectPtr.(*T); ok {
			return result, nil
		}
	} else {
		errorMessage := fmt.Sprintf("invalid type %s (%s), Convert expects a struct", value.Type(), value.Kind())
		return nil, errors.New(errorMessage)
	}

	errorMessage := fmt.Sprintf("%s object could not be converted to %s", value.Type(), Type[T]())
	return nil, errors.New(errorMessage)
}

func Resolve[T any]() (*T, error) {
	c := GetContainer()

	typeInfo := Type[T]()
	built, err := c.Resolve(typeInfo)

	if err != nil {
		return nil, err
	}

	return Convert[T](built)
}

func ConvertImpl[T any](value reflect.Value) (T, error) {
	objectPtr := value.Interface()
	typeInfo := Type[T]()

	if typeInfo.Kind() == reflect.Interface {
		if result, ok := objectPtr.(T); ok {
			return result, nil
		}
	} else {
		errorMessage := fmt.Sprintf("invalid type %s (%s), ConvertImpl expect an interface", value.Type(), value.Kind())
		return *new(T), errors.New(errorMessage)
	}

	errorMessage := fmt.Sprintf("%s object could not be converted to %s", value.Type(), Type[T]())
	return *new(T), errors.New(errorMessage)
}

func ResolveImpl[T any]() (T, error) {
	c := GetContainer()

	typeInfo := Type[T]()
	built, err := c.Resolve(typeInfo)

	if err != nil {
		return *new(T), err
	}

	return ConvertImpl[T](built)
}

func BindInstance[T any](instance *T) {
	GetContainer().SetRule(
		TypeId[T](),
		&instanceRule{reflect.ValueOf(instance)},
	)
}

func BindImpl[T any, U any](impl *U) {
	if !Type[*U]().Implements(Type[T]()) {
		message := fmt.Sprintf("*%s does not implement %s", Type[U](), Type[T]())
		panic(message)
	}

	GetContainer().SetRule(
		TypeId[T](),
		&instanceRule{reflect.ValueOf(impl)},
	)
}

func BindType[T any, U any]() {
	if !Type[*U]().Implements(Type[T]()) {
		message := fmt.Sprintf("*%s does not implement %s", Type[U](), Type[T]())
		panic(message)
	}

	GetContainer().SetRule(
		TypeId[T](),
		&typeRule{Type[U]()},
	)
}

func BindFactory(callback any) {
	returnType, err := GetContainer().validateFactoryCallback(callback)

	if err != nil {
		panic(err.Error())
	}

	GetContainer().SetRule(
		Id(returnType.String()),
		&factoryRule{callback},
	)
}

func BindProvider(callback any) {
	returnType, err := GetContainer().validateFactoryCallback(callback)

	if err != nil {
		panic(err.Error())
	}

	GetContainer().SetRule(
		Id(returnType.String()),
		&providerRule{factoryRule: factoryRule{callback}},
	)
}

