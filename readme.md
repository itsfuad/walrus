# Walrus Programming language
A tiny simple programming language made for simplicity. It borrows syntax from 'go', 'rust' and 'typescript'

- [x] Lexer
- [x] Parser
    - [x] [Variable declare](#variable-declare-and-assign)
    - [x] [Variable assign](#variable-assign)
    - [x] [Expressions](#expressions)
        - [x] Unary (int, float, bool) `- !`
        - [x] Additive `+ -`
        - [x] Multiplicative `* / % ^`
        - [x] Grouping `( )`
    - [x] [Array](#array)
        - [x] Array indexing
    - [x] [Conditionals](#conditionals)
        - [x] if
        - [x] else
        - [x] else if
    - [x] [Functions](#functions)
        - [x] Function declaration
        - [x] Function call
        - [x] Function return
        - [x] Optional parameters
    - [x] [Closure](#functions)
    - [x] [User defined types](#user-defined-types)
        - [x] Struct
            - [x] Property access
            - [x] Property assignment
            - [x] Private property deny access
            - [x] Struct embedding
        - [x] Builtins (int, float, bool, string)
        - [x] Function
    - [x] [Increment/Decrement](#incrementdecrement)
        - [x] Prefix
        - [x] Postfix
    - [x] [Assignment operators](#assignment-operators)
        - [x] +=
        - [x] -=
        - [x] *=
        - [x] /=
        - [x] %=
    - [ ] [For loop](#for-loop)
        - [x] for [condition]
        - [ ] for [start] [condition] [end]
        - [ ] for [start] in [range] 
    - [x] [Traits](#traits)
        - [x] Implement
        - [x] Implement for struct
    - [ ] Generics
- [x] Analyzer
    - [x] Everything in parser except - 
        - [ ] For loop
        - [x] Traits
        - [ ] Implement
- [ ] Codegen

# Example

## Variable declare and assign
```rust
// Declare a variable with let or const keyword
let a := 10; // The variable is mutable and its type is inferred from the value e.g. int
const pi := 3.14; // constant variable with type float

// Declare a variable with type
let b: int = 20; // The variable is mutable and its type is int
const c: float = 3.14; // constant variable with type float
```

## Variable assign
```rust
let a := 10;
a = 20; // Assign a new value to a
```

## Expressions
```rust
let a := 10;
let b := 20;
let c := a + b; // c = 30
let d := a * b; // d = 200
let e := a / b; // e = 0.5
let f := a % b; // f = 10
let g := a ^ b; // g = 100000000000000000000
let h := -a; // h = -10
let i := !true; // i = false
```

## Array
```rust
let a := [1, 2, 3, 4, 5]; // Array of integers. One dimensional array
let b := a[0]; // b = 1
a[0] = 10; // a = [10, 2, 3, 4, 5]

// Array of arrays
let c := [[1, 2], [3, 4], [5, 6]]; // Array of arrays of integers
let d := c[0][0]; // d = 1
c[0][0] = 10; // c = [[10, 2], [3, 4], [5, 6]]
```

## Struct
```rust
type Person struct{
    name: str;
    age: int;
}

//Assign the type with @Name syntax. So we can distinguish between type and variable.
let p := @Person {
    name: "John",
    age: 20
};

// We could also assign the type with type inference
let p : Person = @Person { // Here 'Person' is the type, @Person is the type instance
    name: "John",
    age: 20
};
```

## Struct property access
```rust
let p := @Person {
    name: "John",
    priv age: 20 // Private property
};

let name := p.name; // name = "John"
let age := p.age; // Error: Cannot access private property
```

## Struct embedding
```rust
type Person struct{
    name: str;
    age: int;
}

type Employee struct{
    embed Person; // Embedding Person struct with 'embed' keyword
    salary: float;
}

let e := @Employee {
    name: "John", // Accessing name from Person
    age: 20,
    salary: 1000.0
};

let name := e.name; // name = "John"
let age := e.age; // age = 20
let salary := e.salary; // salary = 1000.0
```

## Conditionals
```rust
let a := 10;
let b := 20;

if a > b {
    print("a is greater than b");
} else if a < b {
    print("a is less than b");
} else {
    print("a is equal to b");
}
```

## Functions
```rust
fn add(a: int, b: int) -> int {
    return a + b;
}

let sum := add(10, 20); // sum = 30

// function with optional parameters
fn add(a: int, b: int, c?: int = 0) -> int {
    return a + b + c;
}

let sum := add(10, 20); // sum = 30

// functions are first class citizens so we can assign them to variables
let adder := add;

let sum := adder(10, 20); // sum = 30

// function with closure
fn add(a: int) -> fn(int) -> int {
    return fn(b: int) -> int {
        return a + b;
    };
}

let adder := add(10);
let sum := adder(20); // sum = 30
```

## User defined types
types are user defined data types. They can be structs, or a function signature or a wrapper around a built-in type.
```rust
type Circle struct {
    radius: float;
}

type FnType = fn(int, int) -> int;

type WrapperInt = int;

let c := @Circle {
    radius: 10.0
};

let f: FnType = fn(a: int, b: int) -> int {
    return a + b;
};

let w: WrapperInt = 10;
```

## Increment/Decrement
```rust
let a := 10;
let b := ++a; // b = 11, a = 11
let c := a++; // c = 11, a = 12
let d := --a; // d = 11, a = 11
let e := a--; // e = 11, a = 10
```

## Assignment operators
```rust
let a := 10;
a += 10; // a = 20
a -= 10; // a = 10
a *= 10; // a = 100
a /= 10; // a = 10
a %= 10; // a = 0
```

## For loop
```rust
for let i := 0; i < 10; i++ {
    print(i);
}
```

## Traits
Traits are shared behavior that can be implemented by different types. Traits are similar to interfaces in other languages.
```rust

type Circle struct {
    radius: float;
}

type Square struct {
    side: float;
}

trait Shape {
    fn area() -> float;
}

impl Shape for Circle {
    fn area() -> float {
        return 3.14 * this.radius * this.radius;
    }
}

impl Shape for Square {
    fn area() -> float {
        return this.side * this.side;
    }
}

let c := @Circle {
    radius: 10.0
};

let s := @Square {
    side: 10.0
};

print(c.area()); // 314.0
print(s.area()); // 100.0

```

## Implementing a trait for a struct
```rust
trait Printable {
    fn print();
}

type Person struct {
    name: str;
    age: int;
}

impl Printable for Person {
    fn print() {
        print("Name: ", this.name, " Age: ", this.age);
    }
}

let p := @Person {
    name: "John",
    age: 20
};

p.print(); // Name: John Age: 20

// Implementing a trait for a struct
impl Person {
    fn greet() {
        print("Hello, ", this.name);
    }
}

p.greet(); // Hello, John
```