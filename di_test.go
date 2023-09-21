package di

import (
	"reflect"
	"testing"
)

type Thing1 struct {
	name string
}

type Thing1Alt struct {
	Thing1
	subname string
}

func (t *Thing1Alt) test() string {
	return t.subname
}

type Thing2 struct {
	Thing1m Thing1
}

type Embed1 struct {
	Thing1m  Thing1
	Thing1m2 Thing1
}

type Thing3 struct {
	Embed1
	Thing1p2        *Thing1
	Thing2m         Thing2
	Thing1NoInjectM Thing1
	SomeNumberm     int
	Itestm          ITest
}

type ITest interface {
	test() string
}

func TestResolve(t *testing.T) {
	Reset()

	thing1, err := Resolve[Thing1]()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(thing1) != Type[*Thing1]() {
		t.Error("Failed asserting Thing1")
	}

	thing3, err := Resolve[Thing3]()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(thing3) != Type[*Thing3]() {
		t.Error("Failed asserting Thing3")
	}

	if reflect.TypeOf(thing3.Thing1m2) != Type[Thing1]() {
		t.Error("Failed asserting member to Thing1")
	}

	if reflect.TypeOf(thing3.Thing1p2) != Type[*Thing1]() {
		t.Error("Failed asserting member pointer to Thing1")
	}

	if thing1 != thing3.Thing1p2 {
		t.Error("Failed asserting pointers to Thing1 are same instance")
	}
}

func TestInstance(t *testing.T) {
	Reset()

	inst := Instance[Thing1]()
	resolved, _ := Resolve[Thing1]()

	if inst != resolved {
		t.Error("Instance[T]() and Resolve[]() should return the same object")
	}
}

func TestImpl(t *testing.T) {
	Reset()

	BindType[ITest, Thing1Alt]()
	impl := Impl[ITest]()
	resolved, _ := ResolveImpl[ITest]()

	if impl != resolved {
		t.Error("Instance[T]() and Resolve[]() should return the same object")
	}
}

func TestBindInstance(t *testing.T) {
	Reset()

	thing1 := &Thing1{name: "THING1"}
	BindInstance(thing1)
	thing3, err := Resolve[Thing3]()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(thing3) != Type[*Thing3]() {
		t.Error("Failed asserting Thing3")
	}

	if reflect.TypeOf(thing3.Thing1m2) != Type[Thing1]() {
		t.Error("Failed asserting member to Thing1")
	}

	if reflect.TypeOf(thing3.Thing1p2) != Type[*Thing1]() {
		t.Error("Failed asserting member pointer to Thing1")
	}

	if thing1 != thing3.Thing1p2 {
		t.Error("Failed asserting pointers to Thing1 are same instance")
	}

	if thing3.Thing1m2.name != "THING1" || thing3.Thing1p2.name != "THING1" {
		t.Error("Failed recalling correct Thing1 instance")
	}
}

func TestBindImpl(t *testing.T) {
	Reset()

	BindImpl[ITest](&Thing1Alt{subname: "subname"})
	object1, err := ResolveImpl[ITest]()

	if err != nil {
		t.Error(err)
	}

	if !reflect.TypeOf(object1).Implements(Type[ITest]()) {
		t.Error("Failed asserting type of implementation")
	}

	if object1.test() != "subname" {
		t.Error("Invalid test message")
	}
}

func TestBindImplFails(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	BindImpl[Thing1](&Thing1Alt{})
	t.Error("expected panic")
}

func TestBindType(t *testing.T) {
	Reset()

	thing1 := &Thing1Alt{subname: "embedded"}
	BindInstance(thing1)

	BindType[ITest, Thing1Alt]()
	object1, err := ResolveImpl[ITest]()

	if err != nil {
		t.Error(err)
	}

	if !reflect.TypeOf(object1).Implements(Type[ITest]()) {
		t.Error("Failed asserting type of implementation")
	}

	if object1.test() != "embedded" {
		t.Error("Invalid name", object1.test())
	}
}

func TestBindTypeFails(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	BindType[Thing1, Thing1Alt]()
	t.Error("expected panic")
}

func TestBindFactory(t *testing.T) {
	Reset()

	thing1 := &Thing1{name: "initial"}
	BindInstance(thing1)

	BindFactory(func(thing1 Thing1) (*Thing1Alt, error) {
		return &Thing1Alt{subname: thing1.name}, nil
	})

	alt, err := Resolve[Thing1Alt]()

	if err != nil {
		t.Error(err)
	}

	if alt.subname != "initial" {
		t.Error("Failed asserting created object")
	}

	alt2, err := Resolve[Thing1Alt]()

	if err != nil {
		t.Error(err)
	}

	if alt == alt2 {
		t.Error("Factory should return different values each time")
	}

	BindFactory(func(thing1 *Thing1) (ITest, error) {
		return &Thing1Alt{subname: thing1.name}, nil
	})

	iTest, err := ResolveImpl[ITest]()

	if err != nil {
		t.Error(err)
	}

	if iTest.test() != "initial" {
		t.Error("Failed asserting created object")
	}
}

func TestBindProvider(t *testing.T) {
	Reset()

	thing1 := &Thing1{name: "initial"}
	BindInstance(thing1)

	BindProvider(func(thing1 Thing1) (*Thing1Alt, error) {
		return &Thing1Alt{subname: thing1.name}, nil
	})

	alt, err := Resolve[Thing1Alt]()

	if err != nil {
		t.Error(err)
	}

	if alt.subname != "initial" {
		t.Error("Failed asserting created object")
	}

	alt2, err := Resolve[Thing1Alt]()

	if err != nil {
		t.Error(err)
	}

	if alt != alt2 {
		t.Error("Provider should return the same value each time")
	}
}
