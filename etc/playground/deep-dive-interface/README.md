# Deep Dive Interfaces

> Interfaces in Go provide a way to specify **the behavior of an object**: if something can do this, then it can be used here.
https://go.dev/doc/effective_go#interfaces_and_types

There are two kinds of types in Go: **concrete** and **abstract**.
A concrete type defines both **shape** and **behavior**. It tells Go exactly how something is represented in memory and what operations it can perform.

An abstract type, on the other hand, defines **behavior only**. It says, "anything that can do this can be treated as that."
In Go, interfaces are abstract types. They don't care how a method works, only that it exists.

> A method exists if and only if its name and its input and output types match exactly.


## Thinking About Interfaces

Imagine a restaurant looking to hire a new Chef.
The job requirement says: **Can cook ramen**.

That's the entire job description. It doesn't say anything about age, gender, nationality, or moral background.
Anyone who can cook ramen qualifies. It doesn't matter how they cook it, as long as the result is ramen.

Even if one applicant can also bend fire like an avatar, that's fine. The requirement doesn't forbid extra skills; **it simply doesn't care**.

If we translate that into Go, it looks like this:

```go
type Chef interface {
    Cook(Ingredients, quantity) Ramen
}

func (Man) Cook(Ingredients, quantity) Ramen { ... }

func (Man) Shop(Money) Games { ... }

func (Woman) Cook(Ingredients, quantity) Ramen { ... }

func (WalterWhite) Cook(Ingredients, quantity) Crystal { ... }
```

In this context, both Man and Woman satisfy the Chef interface, since they both implement `Cook(Ingredients, quantity) Ramen`.
WalterWhite does not, even though he can cook—his "Shiny Crystal", not Ramen. He might be a genius in chemistry, but sorry Mr.Heisenberg you are not qualified for this kitchen.


## The `any` Type

Before Go 1.18, there was no any type. If you wanted a type that could hold anything, you used `interface{}`: an interface with no methods.
That's equivalent to saying, "If you can do nothing, you're still welcome." **With no requirements, everything qualifies**.

In Go 1.18, `any` was introduced as an alias for `interface{}`:

```go
type any = interface{}

// so this:
var x interface{}
// is the same as:
var x any
```

The meaning hasn't changed, only the readability has improved.

## Implicit 

In some languages, you must explicitly declare that a class implements an interface.
For example, in Java:

```java
class Man implements Chef {
    @Override
    public Ramen cook(Ingredients i, int quantity) { ... }
}
```

Go takes a simpler approach. You just define the method, and if it matches, it works:

```go
type Man struct{}
func (Man) Cook(Ingredients, quantity) Ramen { ... }
```

There is no `implements` keyword.
Go observes what you can do, then the Go compiler invited you to join.

> This implicit implementation reduces coupling between the abstraction and its implementation. The type doesn't need to know which interfaces it satisfies, and the interface doesn't need to know who implements it.

## How Interfaces Are Stored

You can read the full explanation here: https://research.swtch.com/interfaces

In short, an interface is stored as a pair: `(interface type, concrete value)`.
For example:
```go
var f io.Reader = &os.File{}
```
Internally, `f` holds both the interface type and the underlying value: `(io.Reader, *os.File)`.
This allows Go to recover the concrete type later:

```go
actualFile, ok := f.(*os.File)
```
or to inspect it dynamically:
```go
switch f.(type) {
case *os.File:
    // do something
}
```

> The real implementation uses internal structures called itables for efficiency. For a deep dive, see Russ Cox's article linked above.

## Best Practices

### Keep interfaces small

A good interface should do one thing, and do it well.
The smaller an interface is, the more types can implement it, and the less likely it'll break when you sneeze near it.

For example, Go's io package has this masterpiece of minimalism:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

That's it. One method. Yet this tiny interface powers files, sockets, buffers, HTTP bodies, and half the standard library.
Keep it small. Let composition do the heavy lifting, (e.g. `io.ReadCloser`).

###  Accept interfaces, return concretes

When writing functions, accept behavior, not implementation.
If all you need is something that can Read, don't demand a `*os.File`. Ask for an `io.Reader`.

```go
func Copy(r io.Reader, w io.Writer) error {
    buf := make([]byte, 1024)
    for {
        n, err := r.Read(buf)
        if n > 0 {
            if _, werr := w.Write(buf[:n]); werr != nil {
                return werr
            }
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
    }
    return nil
}
```
This function doesn't care who you are -- file, buffer, or interdimensional data stream -- as long as you can read and write.
But when returning values, prefer concrete types. It makes your API easier to reason about and avoids unnecessary indirection.

> Accept flexibility, return clarity.

The only exception is `error`, which is an interface by design. Returning an `error` is not only fine, it's expected.


### Avoid premature abstraction

Don't start your project by defining UserRepositoryInterface.
That's not architecture -- that's overengineering in disguise. Unless you already discovered the pattern in you mind.

Start concrete. Make something real.
Once you've discovered the pattern or added a second implementation, then you can introduce an interface.

Bad (fictional, but painfully real):

```go
type Repository interface {
    Save(user User)
    FindByID(id string) (User, error)
}
```

Better:

```go
type UserRepository struct {
    db *sql.DB
}
```

If you need to test the `UserRepository`, you can write mock for the `db` instead, go allow you to create custom driver.
Just like `slog.Logger`, you don't need to define interface for Logger just for testing, you can change the `io.Writer` to `io.Discard` or `t.Output()`.

Nowadays, tool like `testcontainers` exists, it's better than mocking. Used it if you ~~can~~ are allowed.  


### Define interfaces where they're used

Interfaces belong to the consumer, not the provider.
If package A depends on some behavior from package B, define the interface in A.
B doesn't need to know or care.

Example:

```go
// bad: defined in the provider
package storage

type Store interface {
    Save(data []byte) error
}

// better: defined in the consumer
package service

type Saver interface {
	Save(data []byte) error
}
```

That way, storage just provides a type that happens to satisfy Saver, without knowing it.
This keeps dependencies clean and the direction of ownership clear.

> The consumer writes the rules. The provider just happens to qualify for the job.

### Don't create interfaces just for testing

If the only reason your interface exists is because of a mocking framework,
you've accidentally invented a bureaucracy.

Keep your production code free of fake abstractions.
If you need mocks, you can mock like how Go team does for `httptest`, and `slog.Logger`. Instead of mocking
the entire logger, it accepts `io.Writer` that you can control.

### Abstraction has a cost

Interfaces use dynamic dispatch, meaning the actual method call is resolved at runtime.
This is fast enough for most use cases -- but not free.

If you're in a tight loop, this can lead to heap allocations and missed compiler optimizations.
That's why some standard library components, like `slog.Logger`, intentionally avoid interfaces internally.

Another cost is, maintain the code. Well-designed interface reduced coupling, but poorly designed interface will
increase headache.

### Abstraction should be earned

A well-designed interface is something you discover, not something you're architecting first.
When you notice multiple concrete types sharing the same pattern, extract it into an interface.
That's how `io.Reader`, `fmt.Stringer`, and `error` were born.


If you start with interfaces everywhere, you'll end up in more ceremony than hacking.

### Let the standard library be your north star

Every elegant Go package follows a quiet principle: keep it simple, and let the interfaces emerge naturally.
Each defines the minimal behavior necessary to express a concept.

