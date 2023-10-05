package di

import (
	"fmt"
	"reflect"
	"sync"
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
