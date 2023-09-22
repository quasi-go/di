package di

import (
	"errors"
	"fmt"
	"reflect"
)

type Initializable interface {
	Initialize()
}

func Convert[T any](value reflect.Value) (*T, error) {
	if value.Kind() != reflect.Pointer {
		errorMessage := fmt.Sprintf("expected pointer, got %s (%s)", value.Type(), value.Kind())
		return nil, errors.New(errorMessage)
	}

	objectPtr := value.Interface()
	typeInfo := Type[T]()

	if typeInfo.Kind() != reflect.Interface {
		if result, ok := objectPtr.(*T); ok {
			return result, nil
		}
	} else {
		errorMessage := fmt.Sprintf("invalid type %s (%s), Convert must not be an interface", value.Type(), value.Kind())
		return nil, errors.New(errorMessage)
	}

	errorMessage := fmt.Sprintf("%s object could not be converted to %s", value.Type(), Type[T]())
	return nil, errors.New(errorMessage)
}

func Resolve[T any]() (*T, error) {
	c := GetContainer()

	typeInfo := Type[T]()
	built, err := c.ResolveType(typeInfo)

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
	built, err := c.ResolveType(typeInfo)

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
	validateImpl[T, U]()

	GetContainer().SetRule(
		TypeId[T](),
		&instanceRule{reflect.ValueOf(impl)},
	)
}

func BindType[T any, U any]() {
	validateImpl[T, U]()

	GetContainer().SetRule(
		TypeId[T](),
		&typeRule{Type[U]()},
	)
}

func BindFactory(callback any) {
	returnType, err := validateFactoryCallback(callback)

	if err != nil {
		panic(err.Error())
	}

	GetContainer().SetRule(
		Id(returnType.String()),
		&factoryRule{callback},
	)
}

func BindProvider(callback any) {
	returnType, err := validateFactoryCallback(callback)

	if err != nil {
		panic(err.Error())
	}

	GetContainer().SetRule(
		Id(returnType.String()),
		&providerRule{factoryRule: factoryRule{callback}},
	)
}

func Instance[T any]() *T {
	inst, err := Resolve[T]()

	if err != nil {
		panic(err)
	}

	return inst
}

func Impl[T any]() T {
	inst, err := ResolveImpl[T]()

	if err != nil {
		panic(err)
	}

	return inst
}

func Reset() {
	resetContainer()
}

func Call[T any](callback any) (*T, error) {
	callbackType := reflect.TypeOf(callback)
	nReturn := callbackType.NumOut()

	if nReturn != 2 || callbackType.Out(0) != Type[*T]() || callbackType.Out(1) != Type[error]() {
		return new(T), errors.New("the callback must return a *" + Type[T]().String() + "and an error")
	}

	returnValues, err := GetContainer().Call(callback)

	if err != nil {
		return new(T), err
	}

	convertedT, err := Convert[T](returnValues[0])

	if err != nil {
		return new(T), err
	}

	returnedErr := returnValues[1].Interface()

	if returnedErr == nil {
		return convertedT, nil
	}

	if convertedErr, ok := returnedErr.(error); ok {
		return convertedT, convertedErr
	}

	return new(T), errors.New("returned value could not be converted to error")
}

func Invoke(callback any) {
	_, err := GetContainer().Call(callback)

	if err != nil {
		panic(err)
	}
}

func validateImpl[T any, U any]() {
	if Type[T]().Kind() != reflect.Interface {
		message := fmt.Sprintf("*%s must be an interface", Type[T]())
		panic(message)
	}

	if !Type[*U]().Implements(Type[T]()) {
		message := fmt.Sprintf("*%s does not implement %s", Type[U](), Type[T]())
		panic(message)
	}
}
