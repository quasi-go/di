package example

import (
	"fmt"
	"github.com/quasi-go/di"
	"reflect"
	"testing"
)

type A struct {
	Name string
}

type B struct {
	Type       string
	PtrDep     *A
	ValueDep   A
	privateDep *A
}

// Initialize() is called immediately after created

func (b *B) Initialize() {
	b.Type = "None"
}

type C struct {
	B
}

func (c *C) Initialize() {
	c.Type = "Embedded"
}

// All structs implement I

type I interface{}

func TestAll(t *testing.T) {
	// `BindInstance[T](inst)` bind type `T` to the passed `inst`.

	a := &A{Name: "Alice"}
	di.BindInstance[A](a)

	// Instance[A]() retrieves the same instance `a` from above

	resolvedA1 := di.Instance[A]()

	// `resolvedA` is a `*A`
	// `resolvedA` == `a`

	// Even without explicitly binding a type, the library can implicitly build new structs

	resolvedB := di.Instance[B]()

	// Private members are not set by the library, so `resolvedB.privateDep` === nil
	// `resolvedB` is a `*B`
	// `resolvedB.Type` == "None" from the call to `B.initialize()`
	// `resolvedB.PtrDep` === `a`
	// `resolvedB.ValueDep` == `*a`

	// You can also bind an interface `I` to an instance that implements it with `BindImpl[I]()`

	di.BindImpl[I](resolvedB)

	// Note that when resolving an interface, use `Impl[I]()` instead of `Instance[I]()`

	resolvedI1 := di.Impl[I]()

	// `resolvedI` === `resolvedB`

	// We can also bind an interface to a type
	// Note that we are overwriting our previous binding ot `I`

	di.BindType[I, C]()
	resolvedI2 := di.Impl[I]()

	// di.Impl[I]() === di.Instance[C]()

	// `BindProvider(func)` will bind a type to a provider function. The parameters  are `Instanced()`-ed
	// before being injected to the function. The bound type is the return type of the function.

	di.BindProvider(func(injectedA *A) (I, error) {
		// injectedA === resolvedA

		// Providers can be used for things that cannot be automatically resolved, such as private
		// instance variables
		return &B{
			privateDep: injectedA,
		}, nil
	})

	resolvedI3 := di.Impl[I]()

	// di.Impl[I]() is now generated by the callback we defined.
	// The generated struct is only constructed once; the callback is not invoked multiple times.
	// di.Impl[I]() === di.Impl[I]()

	// `BindFactory(func)` works the same as `BindProvider(func)` except that it will be invoked
	// to return a new instance each time it is `Instance()`-ed

	di.BindFactory(func() (*A, error) {
		return &A{
			Name: "Created from factory",
		}, nil
	})

	resolvedA2 := di.Instance[A]()

	// di.Instance[A]() != di.Instance(A)[]

	fmt.Println("resolvedA1: (", reflect.TypeOf(resolvedA1), ")", resolvedA1)
	fmt.Println("resolvedB:  (", reflect.TypeOf(resolvedB), ")", resolvedB)
	fmt.Println("resolvedI1: (", reflect.TypeOf(resolvedI1), ")", resolvedI1)
	fmt.Println("resolvedI2: (", reflect.TypeOf(resolvedI2), ")", resolvedI2)
	fmt.Println("resolvedI3: (", reflect.TypeOf(resolvedI3), ")", resolvedI3)
	fmt.Println("resolvedA2: (", reflect.TypeOf(resolvedA2), ")", resolvedA2)
}
