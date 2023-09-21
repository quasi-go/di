package di

import "reflect"

type Id string

type tId[T any] struct {
	_type T
}

func Type[T any]() reflect.Type {
	tid := tId[T]{}
	iType, _ := reflect.TypeOf(tid).FieldByName("_type")
	return iType.Type
}

func TypeId[T any]() Id {
	iType := Type[T]()
	return Id(iType.String())
}

func ObjectTypeId(object any) Id {
	value := reflect.ValueOf(object)

	if value.Kind() == reflect.Pointer {
		value = value.Elem()

	}
	return Id(value.Type().String())
}
