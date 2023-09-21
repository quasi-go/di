package di

import (
	"reflect"
	"testing"
)

func TestTypeId(t *testing.T) {
	if ObjectTypeId(&Thing1{}) != TypeId[Thing1]() {
		t.Error("Type IDs should match")
	}
}

func TestType(t *testing.T) {
	if reflect.TypeOf(Thing1{}) != Type[Thing1]() {
		t.Error("Types should match")
	}
}
