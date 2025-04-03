# Lesson 2: Go (Golang) Language Overview

Go (or Golang) is an open-source programming language developed at Google. It is designed for simplicity, high performance, and ease of use in building scalable and reliable software. Its clean syntax and powerful standard library make it ideal for backend services, distributed systems, and CLI tools.

## Documentation

### General go

- https://go.dev/tour/welcome/1
- https://go.dev/doc/tutorial/
- https://gobyexample.com/

### Error handling

- https://www.youtube.com/watch?v=d2DySHJ7oVk&themeRefresh=1
- https://go.dev/blog/error-handling-and-go
- https://earthly.dev/blog/golang-errors/

### Misc

- https://go101.org/article/type-embedding.html
- https://go.dev/blog/slog

## 1. Interfaces

### Description
Interfaces in Go define a set of method signatures. A type implements an interface simply by implementing its methods—there’s no explicit declaration of intent. This enables a form of implicit polymorphism.

### Basic Example
```go
type Speaker interface {
    Speak() string
}

type Dog struct{}

func (d Dog) Speak() string {
    return "Woof!"
}

func talk(s Speaker) {
    fmt.Println(s.Speak())
}
```

### Function Interface Example (http.HandlerFunc-like)
```go
type HandlerFunc func(string) string

func (f HandlerFunc) Serve(input string) string {
    return f(input)
}

func HelloHandler(name string) string {
    return "Hello " + name
}

var f HandlerFunc = HelloHandler
fmt.Println(f.Serve("Gopher"))
```

### Multiple Interfaces Implementation
```go
type Reader interface {
    Read() string
}

type Writer interface {
    Write(string)
}

type ReadWriter struct {
    data string
}

func (rw *ReadWriter) Read() string {
    return rw.data
}

func (rw *ReadWriter) Write(s string) {
    rw.data = s
}

func useReader(r Reader) {
    fmt.Println("Read:", r.Read())
}

func useWriter(w Writer) {
    w.Write("new content")
}
```

---

## 2. Receivers

### Description
Receivers allow functions to be associated with types. There are two types:
- Value receiver: works on a copy of the type.
- Pointer receiver: works on the original instance.

Receivers are mostly syntactic sugar to organize functions by their associated type.

> Note: Go does **not** support generic receivers.

### Code Example
```go
type Person struct {
    Name string
}

func (p Person) Greet() {
    fmt.Println("Hello,", p.Name)
}

func (p *Person) Rename(newName string) {
    p.Name = newName
}
```

---

## 3. Value Types & Pointers

### Description
Go is a value-oriented language. Pointers allow you to modify the original value and manage memory efficiently.

### Code Example
```go
a := 10
b := &a
*b = 20
fmt.Println(a) // Output: 20
```

### Struct Pointer Creation
```go
type User struct {
    Name string
}

u1 := &User{Name: "Ana"}   // via address-of operator
u2 := new(User)             // via built-in new function
u2.Name = "Bob"
```

---

## 4. Arrays, Slices, Maps, `make()`

### Description
- **Arrays**: fixed size.
- **Slices**: dynamic size.
- **Maps**: key-value pairs.
- **make()**: built-in function to initialize slices, maps, and channels.

### Code Example
```go
arr := [3]int{1, 2, 3}
slc := make([]int, 0, 10) // len=0, cap=10
m := make(map[string]int)
m["age"] = 30
```

### Explanation
- `make([]int, len, cap)` allocates a slice with a given length and capacity.
- `make(map[string]int)` allocates a map.
- Arrays are rarely used directly; slices are more idiomatic.

---

## 5. Enum & `iota` 🤮

### Description
Go lacks real enums. `iota` helps generate incremented constants.

### Code Example
```go
type Status int

const (
    Pending Status = iota
    Approved
    Rejected
)
```

### Caveats
Avoid overusing `iota`:
- You can’t easily trace values
- Adding a value at the top shifts all others, causing bugs in protocols/data formats.

Instead, consider using explicit constants.

---

## 6. Type Alias & Struct Composition

### Description
- `type A = B`: alias.
- `type A B`: new, distinct type.
- Struct composition allows reusing fields and methods via embedding.

### Code Example
```go
type ID = int       // alias

type UserID int     // new type, distinct from int

func PrintID(id UserID) {
    fmt.Println(id)
}

// Struct Composition

type Animal struct {
    Name string
}

func (a Animal) Speak() string {
    return "I'm " + a.Name
}

type Dog struct {
    Animal
    Breed string
}

d := Dog{Animal: Animal{Name: "Fido"}, Breed: "Labrador"}
fmt.Println(d.Speak()) // Calls embedded Animal.Speak
fmt.Println(d.Name)    // Accesses embedded field
```

---

## 7. Errors & Handling Errors

### Description
Errors in Go are values implementing the `error` interface.

### Basic Error Handling
```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

### Error Wrapping & Unwrapping
```go
if err := doSomething(); err != nil {
    return fmt.Errorf("doSomething failed: %w", err)
}

// Check wrapped errors
if errors.Is(err, fs.ErrNotExist) { ... }
if customErr, ok := err.(*MyError); ok { ... }
```

### Custom Error Type
```go
type MyError struct {
    Msg string
    Err error
}

func (e *MyError) Error() string {
    return e.Msg + ": " + e.Err.Error()
}

func (e *MyError) Unwrap() error {
    return e.Err
}
```

---

## 8. Generics, Interface Constraints, Type Casting & Assertion

### Description
Generics allow reusable, type-safe functions and structs. Type assertions help extract concrete types from interfaces.

### Code Example
```go
func Print[T any](val T) {
    fmt.Println(val)
}

// With constraints
func Sum[T int | float64](a, b T) T {
    return a + b
}

// Assertion
var x any = "hello"
s, ok := x.(string) // s is "hello", ok is true

// Type conversion
var i int = 5
var f float64 = float64(i) // explicit conversion
```

---

## 9. `for` & `range`

### Description
Go has one looping construct: `for`. It can be used like `while`, traditional `for`, or `foreach` via `range`.

### Code Example
```go
nums := []int{1, 2, 3}
for i, val := range nums {
    fmt.Println(i, val)
}

for i := 0; i < 10; i++ { ... }
for condition { ... }
for { ... } // Infinite loop
```

---

## 10. Default Values

### Description
Go assigns zero-values to all variables at initialization.

### Common Defaults
```go
var i int       // 0
var f float64   // 0.0
var s string    // ""
var b bool      // false
var m map[string]int // nil
var sl []int         // nil

type User struct {
    Name string
    Age int
}

var u User // Name: "", Age: 0
```

Zero values ensure every variable has a valid value, reducing null reference issues.

---

Feel free to extend this doc with concurrency (goroutines, channels), packages, testing, or context!


