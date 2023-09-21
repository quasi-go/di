
# Quasi DI

A tiny, elegant dependency injection library for golang.

## Example Tests

See [example/example_test.go](example/example_test.go) to run the example tests below.

## Instance Bindings

This first example showcases basic usage via a couple commonly used functions.

`BindInstance()` stores the given instance and will automatically provide
it back when attempting to `Resolve[T]()` an instance of the same type. You can
choose to fully spell this call as `di.BindInstance[SimpleStruct](&dep)`
if you prefer further clarity.

`Resolve[T]()` will first look and see if it has been given any rules for
providing a pointer to an instance of type `T` (its returned value is of
type `*T`). If so, it will follow that rule to provide an instance. In this
case, we bound the type `ExampleStruct` to `&dep` above, so `resolved` points
to `dep`.

`Type[T]()` is a utility function that will return the `reflect.Type`
associated with `T`, a `*SimpleStruct` in this case. This function uses
a trick so that it works with interfaces also.

```go
package example

import (
	"fmt"
	"reflect"
	"testing"
)

import (
	"github.com/quasi-go/di"
)

type SimpleStruct struct {
	Name string
}

func TestSimpleStruct(t *testing.T) {
	dep := SimpleStruct{
		Name: "This is a test",
	}

	di.BindInstance(&dep)
	resolved, _ := di.Resolve[SimpleStruct]()

	// `resolved` should be a pointer back to `dep`

	if *resolved == dep {
		fmt.Println("resolved points to dep, as expected.")
	} else {
		t.Error("resolved and dep should point to the same instance of SimpleStruct.")
	}

	if reflect.TypeOf(resolved) == di.Type[*SimpleStruct]() {
		fmt.Println("resolved is a pointer to SimpleStruct, as expected.")
	} else {
		t.Error("resolved should point to an instance of SimpleStruct.")
	}

	// `resolved.Name` should equal the message we set above.

	if resolved.Name == "This is a test" {
		fmt.Println("resolved.Name is \"This is a test\", as expected")
	} else {
		t.Error("resolved.Name should equal \"This is a test\", as set above.")
	}
}
```

## Implicit Resolution

When no rule is set for a struct type passed to `Resolve[T]()`, it will
attempt to implicitly build a new instance of type `T` and return a pointer `*T`
to the new instance. It cycles through each member of the struct and will
either execute an associated rule if one has been bound for the type
(such as `di.BindInstance(&dep)` above) or it will attempt to recursively
`Resolve[T]()` the child, if it is a struct.

Observe the process below.

```go
type ExampleStruct struct {
	Dependency        SimpleStruct
	PtrToDependency   *SimpleStruct
	privateDependency *SimpleStruct
}

func TestExampleStruct(t *testing.T) {
	resolved, _ := di.Resolve[ExampleStruct]()

	// `resolved` should be a pointer to a new instance of `ExampleStruct`

	if reflect.TypeOf(resolved) == di.Type[*ExampleStruct]() {
		fmt.Println("resolved is a pointer to ExampleStruct, as expected.")
	} else {
		t.Error("resolved should point to an instance of SimpleStruct.")
	}

	// Because we used `BindInstance[T]()` in `TestSimpleStruct` above,
	// each time we `Resolve[SimpleStruct]()` it will return a pointer to the
	// same instance as &dep in the first test.

	dep, _ := di.Resolve[SimpleStruct]()

	// Behind the scenes, `Resolve[ExampleStruct]()` made a call to
	// `di.Resolve[SimpleStruct]()` to populate `resolved.dependency`. We can see
	// that they are identical.

	if resolved.Dependency == *dep {
		fmt.Println("resolved.dependency is the same SimpleStruct as dep points to, as expected.")
	} else {
		t.Error("resolved.dependency should be the same SimpleStruct as dep.")
	}

	// `Resolve[T]()` is smart enough that it know to inject a `*T` or `T` based on the
	// type of the struct member encountered.

	if resolved.PtrToDependency == dep {
		fmt.Println("resolved.ptrToDependency is a pointer to the same SimpleStruct and dep, as expected.")
	} else {
		t.Error("resolved.dependency should point to the same SimpleStruct as dep.")
	}

	// Of course, our Name is still set.

	if resolved.PtrToDependency.Name == "This is a test" {
		fmt.Println("Our Name is still set, as expected.")
	} else {
		t.Error("resolved.PtrToDependency.Name should be the same as when set above.")
	}

	// `Resolve[T]()` cannot resolve private members for you. In this case,
	// `resolved.privateDependency` is `nil`.

	if resolved.privateDependency == nil {
		fmt.Println("resolved.privateDependency is nil, as expected.")
	} else {
		t.Error("resolved.privateDependency should not have been resolved.")
	}
}
```

## Type Bindings

`BindType[I, C]()` binds an interface type `I` to a concrete struct type `C`. Above
we implemented the method `GetName()` for `SimpleStruct`, which satisfies
`TestInterface`. Now we can bind to `TestInterface` so that a `*SimpleStruct`
will be provided when `TestInterface` is resolved.

When we `ResolveImpl[TestInterface]()`, it will attempt to resolve `SimpleStruct`
with `Resolve[SimpleStruct]()`. Same as any other call to `Resolve[T]()` this
operation will follow any previous bindings to `SimpleStruct` if they exist
(in this case, the same as the `dep` variables above), or will attempt to
build a new instance if not. The call to `ResolveImpl[I]()` will return an instance
of the provided interface `I`, as opposed to `Resolve[T]()`, which returns a `*T`.

```go
type TestInterface interface {
	GetName() string
}

func (s *SimpleStruct) GetName() string {
	return s.Name
}

func TestBindType(t *testing.T) {
	di.BindType[TestInterface, SimpleStruct]()
	
	fromInterface, _ := di.ResolveImpl[TestInterface]() // Note the call to ResolveImpl
	fromConcrete, _ := di.Resolve[SimpleStruct]()

	// These two should be pointers to the same object.

	if fromInterface == fromConcrete {
		fmt.Println("resolving the interface TestInterface returns the same as resolving SimpleStruct, as expected.")
	} else {
		t.Error("resolved.privateDependency should not have been resolved.")
	}

	// And in fact, they are pointers to the same `SimpleStruct` from our examples above.
	// The `GetName()` method returns the `.Name` member.

	if fromInterface.GetName() == "This is a test" {
		fmt.Println("GetName returns the same string we set above, as expected.")
	} else {
		t.Error("resolved.PtrToDependency.Name should be the same as when set above.")
	}
}
```

## Binding Interfaces

In some cases you would like to bind an interface to an instance you've already
constructed. This is similar to `BindInstance`, but it associates an interface
with the instance instead of its concrete type. Note below that we're also
overwriting the previous call to `BindType[TestInterface, SimpleStruct]`.
This is allowed.

Notice below we only update the binding to `TestInterface`. `Resolve[SimpleStruct]()`
is unaffected. `di.ResolveImpl[TestInterface]()` was previously the same
as calling `di.Resolve[SimpleStruct]()` because of `BindType[TestInterface, SimpleStruct]()`.
Now `BindImpl[TestInterface](impl)` changed what `TestInterface` resolves to, but
it left the binding to `SimpleStruct` unaffected.

```go
func TestBindImpl(t *testing.T) {
	impl := &SimpleStruct{
		Name: "Not the same as our first test",
	}

	di.BindImpl[TestInterface](impl)

	// Now we can call `ResolveImpl[TestInterface]()` to get back our impl.

	resolved, _ := di.ResolveImpl[TestInterface]()

	// These should be the same.

	if resolved == impl {
		fmt.Println("resolved and impl are the same, as expected.")
	} else {
		t.Error("resolved and impl should be the same.")
	}

	// Note that this is not the same result as out call to
	// `ResolveImpl[TestInterface]()` above.

	if resolved.GetName() == "Not the same as our first test" {
		fmt.Println("We're getting back the new message, as expected.")
	} else {
		t.Error("We should be getting back the new message here.")
	}

	original, _ := di.Resolve[SimpleStruct]()

	if original.GetName() == "This is a test" {
		fmt.Println("We're getting back the original message, as expected.")
	} else {
		t.Error("We should be getting back the new original here.")
	}
}
```

## Automatic Initialization

If we have setup that needs to be performed after the construction of the object,
we can implement `Initializeable`, an interface that consists of a single method
`Initialize()` that accepts no parameters and has no return value. This method
will be called immediately after an instance of the type is built, but not necessarily
each time `Resolve[T]()` or `ResolveImpl[I]()` is called.

```go
type InitializedStruct struct {
	Message string
}

func (i *InitializedStruct) Initialize() {
	i.Message = "This message was set from Initialize()"
}

func TestInitialize(t *testing.T) {
	// We resolve `InitializedStruct` implicitly (with no bindings).

	initialized, _ := di.Resolve[InitializedStruct]()

	// Here we see that the `Initialize()` method was called automatically to set
	// `initialized.Message`.
	if initialized.Message == "This message was set from Initialize()" {
		fmt.Println("The message was set correctly, as expected.")
	} else {
		t.Error("The message should have been set from the Initialize() method")
	}
}
```

## Binding to Provider Functions

`BindProvider(func)` binds a function that will construct our resolved type.
The library infers the type that we are binding from the callback's return type,
and can use `Resolve[T]()` and resolve inject objects as the parameters of the function.

The arguments injected into the provider are resolved the same as a direct call to
`Resolve[T]()` or `ResolveImpl[I]()`.

A provider can be used for things that the library cannot resolve itself, such as
setting private members like `.privateDependency` below.

The first time we resolve `ExampleStruct` it will invoke the func passed to create
an instance.

The second time we invoke `Resolve[ExampleStruct]()` it does not rerun the provider,
but instead returns the same `*ExampleStruct` at it constructed the first time.

```go
func TestBindProvider(t *testing.T) {
	di.BindProvider(func(dep1 *SimpleStruct, dep2 TestInterface) (*ExampleStruct, error) {
		resolved1, _ := di.Resolve[SimpleStruct]()

		if resolved1 == dep1 && resolved1.GetName() == "This is a test" {
			fmt.Println("resolved is the same as if we resolve SimpleStruct directly, as expected.")
		} else {
			t.Error("resolved should be the same as if we resolve SimpleStruct directly.")
		}

		resolved2, _ := di.ResolveImpl[TestInterface]()

		if resolved2 == dep2 && resolved2.GetName() == "Not the same as our first test" {
			fmt.Println("resolved is the same as if we resolve TestInterface directly, as expected.")
		} else {
			t.Error("resolved should be the same as if we resolve TestInterface directly.")
		}
		
		return &ExampleStruct{
			privateDependency: dep1,
		}, nil
	})
	
	resolved, _ := di.Resolve[ExampleStruct]()

	if reflect.TypeOf(resolved) == di.Type[*ExampleStruct]() {
		fmt.Println("resolved is a pointer to ExampleStruct, as expected.")
	} else {
		t.Error("resolved should point to an instance of SimpleStruct.")
	}

	if resolved.privateDependency.GetName() == "This is a test" {
		fmt.Println("resolved has the original message set, as expected.")
	} else {
		t.Error("resolved should have the original message set.")
	}
	
	again, _ := di.Resolve[ExampleStruct]()

	if resolved == again {
		fmt.Println("resolved and again point to the same ExampleStruct, as expected.")
	} else {
		t.Error("resolved and again should point to the same ExampleStruct.")
	}
}
```

## Non-Singleton Bindings

BindFactory(func) works exactly the same as BindProvider(func), with the exception that it runs
every time an instance is resolved.

BELOW, built2 and built2 point to different instances of `SimpleStruct`, each with its own incremented `.Name`.

```go
func TestBindFactory(t *testing.T) {
	count := 0

	di.BindFactory(func() (*SimpleStruct, error) {
		count++
		return &SimpleStruct{
			Name: fmt.Sprintf("%d times", count),
		}, nil
	})

	// 

	built1, _ := di.Resolve[SimpleStruct]()
	built2, _ := di.Resolve[SimpleStruct]()

	if built1 != built2 {
		fmt.Println("The two results don't point to the same instance, as expected.")
	} else {
		t.Error("The two results should be pointers to different instances.")
	}

	if built1.GetName() == "1 times" {
		fmt.Println("The first count works as expected.")
	} else {
		t.Error("Incorrect first count.")
	}

	if built2.GetName() == "2 times" {
		fmt.Println("The second count works as expected.")
	} else {
		t.Error("Incorrect first count.")
	}
}
```