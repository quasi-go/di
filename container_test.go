package di

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestGetSetContainerFails(t *testing.T) {
	resetContainer()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	BindAuto[Thing3]()
	object1, err := Resolve[Thing3]()

	if err != nil {
		t.Error(err)
	}

	// This fails because Thing1p2 is nil
	t.Log(object1.Thing1p2.name)
}

func TestGetSetContainer(t *testing.T) {
	resetContainer()

	thing1 := Thing1{name: "THING1"}
	BindInstance(&thing1)

	BindAuto[Thing3]()
	BindAuto[Embed1]()
	object1, err := Resolve[Thing3]()

	if err != nil {
		t.Fatal(err)
	}

	if object1.Thing1m2.name != "THING1" {
		t.Error("Failed recalling correct Thing1 instance")
	}

	SetContainer(NewContainer())
	_, err = Resolve[Thing3]()

	if err == nil {
		t.Error("Should fail after resetting to default")
	}
}

func TestSetConcurrent(t *testing.T) {
	resetContainer()
	defaultContainer := GetContainer()
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defaultContainer.SetRule(TypeId[Thing1](), &instanceRule{instance: reflect.ValueOf(&Thing1{})})
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestGetConcurrent(t *testing.T) {
	resetContainer()
	defaultContainer := GetContainer()
	defaultContainer.SetRule(TypeId[Thing1](), &instanceRule{instance: reflect.ValueOf(&Thing1{})})
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defaultContainer.GetRule(TypeId[Thing1]())
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestGetSetConcurrent(t *testing.T) {
	resetContainer()
	defaultContainer := GetContainer()
	defaultContainer.SetRule(TypeId[Thing1](), &instanceRule{instance: reflect.ValueOf(&Thing1{})})
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defaultContainer.SetRule(TypeId[Thing1](), &instanceRule{instance: reflect.ValueOf(&Thing1{})})
			defaultContainer.GetRule(TypeId[Thing1]())
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestSetHasConcurrent(t *testing.T) {
	resetContainer()
	defaultContainer := GetContainer()
	defaultContainer.SetRule(TypeId[Thing1](), &instanceRule{instance: reflect.ValueOf(&Thing1{})})
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defaultContainer.SetRule(TypeId[Thing1](), &instanceRule{instance: reflect.ValueOf(&Thing1{})})
			defaultContainer.HasRule(TypeId[Thing1]())
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestBuildTypeConcurrent(t *testing.T) {
	resetContainer()

	thing1 := Thing1{name: "THING1"}
	BindInstance(&thing1)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			GetContainer().BuildType(Type[Thing3]())
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestIsNil(t *testing.T) {
	var v reflect.Value
	fmt.Println(v == reflect.Value{})
}
