
# Quasi DI

A tiny, elegant dependency injection library for golang.

## To The Point

```go
import "github.com/quasi-go/di"
```

`BindInstance[T](inst)` bind type `T` to the passed `inst`.

```go
a := &A{Name: "Alice"}
di.BindInstance[A](a)
```

`di.Resolve[A]()` retrieves the same instance `a` from above

```go
resolvedA1, err := di.Resolve[A]()
```

- `resolvedA` is a `*A`
- `resolvedA` == `a`.

Even without explicitly binding a type, the library can implicitly build new structs

```go
resolvedB, _ := di.Resolve[B]()
```

- Private members are not set by the library, so `resolvedB.privateDep` === nil
- `resolvedB` is a `*B`
- `resolvedB.Type` == "None" from the call to `B.initialize()`
- `resolvedB.PtrDep` === `a`
- `resolvedB.ValueDep` == `*a`

You can also bind an interface `I` to an instance that implements it with `BindImpl[I]()`

```go
di.BindImpl[I](resolvedB)
```

Note that when resolving an interface, use `ResolveImpl[I]()` instead of `Resolve[I]()`

```go
resolvedI1, _ := di.ResolveImpl[I]()
```
- `resolvedI` === `resolvedB`

We can also bind an interface to a type. Note that we are overwriting our previous binding to `I`.

```go
di.BindType[I, C]()
resolvedI2, _ := di.ResolveImpl[I]()
```

- `di.ResolveImpl[I]()` === `di.Resolve[C]()`

`BindProvider(func)` will bind a type to a provider function. The parameters  are `Resolved()`-ed
before being injected to the function. The bound type is the return type of the function.
Providers can be used for things that cannot be automatically resolved, such as private
instance variables.

```go
di.BindProvider(func(injectedA *A) (I, error) {
    return &B{
        privateDep: injectedA,
    }, nil
})

resolvedI3, _ := di.ResolveImpl[I]()
```

- `injectedA` === `resolvedA`
- `di.ResolveImpl[I]()` is now generated by the callback we defined.
- The generated struct is only constructed once; the callback is not invoked multiple times.
- `di.ResolveImpl[I]()` === `di.ResolveImpl[I]()`

`BindFactory(func)` works the same as `BindProvider(func)` except that it will be invoked
to return a new instance each time it is `Resolve()`-ed

```go
di.BindFactory(func() (*A, error) {
    return &A{
        Name: "Created from factory",
    }, nil
})

resolvedA2, _ := di.Resolve[A]()
```

- `di.Resolve[A]()` != `di.Resolve(A)[]`

## Example Tests

The example test above is implemented here: [example/tothepoint_test.go](example/tothepoint_test.go)

See [example/example_test.go](example/example_test.go) to see and run additional examples.

[example/README.md](example/README.md) provides the additional examples in a markdown document.

