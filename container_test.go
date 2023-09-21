package di

import (
	"fmt"
	"testing"
)

func TestGetSetContainer(t *testing.T) {
	resetContainer()
	defaultContainer := GetContainer()

	thing1 := Thing1{name: "THING1"}
	BindInstance(&thing1)

	object1, err := Resolve[Thing3]()

	if err != nil {
		fmt.Println(err)
	}

	if object1.Thing1m2.name != "THING1" {
		t.Error("Failed recalling correct Thing1 instance")
	}

	SetContainer(NewContainer())
	object2, _ := Resolve[Thing3]()

	if object2.Thing1m2.name != "" {
		t.Error("Failed recalling correct Thing1 instance after reset from default")
	}

	SetContainer(defaultContainer)
	object3, _ := Resolve[Thing3]()

	if object3.Thing1m2.name != "THING1" {
		t.Error("Failed recalling correct Thing1 instance after reset back to default")
	}
}
