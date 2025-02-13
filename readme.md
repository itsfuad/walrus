# Walrus Programming language
A tiny simple programming language made for simplicity. It borrows syntax from 'go', 'rust' and 'typescript'

- [x] Lexer
- [x] Parser
    - [x] [Variable declare](#variable-declare-and-assign)
    - [x] [Variable assign](#variable-assign)
    - [x] [Expressions](#expressions)
        - [x] Unary (32, f32, bool) `- !`
        - [x] Additive `+ -`
        - [x] Multiplicative `* / % ^`
        - [x] Grouping `( )`
        - [x] Type cast `as`
    - [x] [Array](#array)
        - [x] Array indexing
    - [x] [Map](#map)
        - [x] Map indexing
        - [x] Map assignment
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
            - [x] Implement for struct
        - [x] Builtins (32, f32, bool, string)
        - [x] Function
        - [x] Interface
            - [x] Define
            - [x] Implement
            - [x] Use 
    - [x] [Increment/Decrement](#incrementdecrement)
        - [x] Prefix
        - [x] Postfix
    - [x] [Assignment operators](#assignment-operators)
        - [x] +=
        - [x] -=
        - [x] *=
        - [x] /=
        - [x] %=
    - [x] [For loop](#for-loop)
        - [x] for [condition]
        - [x] for [start] [condition] [end]
        - [x] for [start] in [range] 
    - [x] [Interaface](#interface)
    - [ ] Imports
    - [ ] Packages
    - [ ] Modules
    - [ ] Generics
- [x] Analyzer
    - [x] Everything in parser except - 
        - [ ] For loop
- [x] Rich multi error reporting system
- [ ] Codegen

# Example

## Variable declare and assign
```rust
// Declare a variable with let or const keyword
let a := 10; // The variable is mutable and its type is inferred from the value e.g. 32
const pi := 3.14; // constant variable with type f32

// Declare a variable with type
let b: 32 = 20; // The variable is mutable and its type is 32
const c: f32 = 3.14; // constant variable with type f32

let unsigned: u32 = 10; // Unsigned integer of 32 bits
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

## Grouping
```rust
let a := 10;
let b := 20;
let c := (a + b) * 2; // c = 60
```

## Type cast
```rust
let a := 10;
let b := a as f32; // b = 10.0 as float32
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

## Map
```rust
let myMap : map[str]i32 = map[str]i32 {
    "a": 10,
    "b": 20,
    "c": 30
};

let a := myMap["a"]; // a = 10
myMap["a"] = 20; // myMap = {"a": 20, "b": 20, "c": 30}
```

## Struct
```rust
type Person struct{
    name: str;
    age: 32;
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
fn add(a: 32, b: 32) -> 32 {
    ret a + b;
}

let sum := add(10, 20); // sum = 30

// function with optional parameters
fn add(a: 32, b: 32, c?: 32 = 0) -> 32 {
    ret a + b + c;
}

let sum := add(10, 20); // sum = 30

// functions are first class citizens so we can assign them to variables
let adder := add;

let sum := adder(10, 20); // sum = 30

// function with closure
fn add(a: 32) -> fn(32) -> 32 {
    ret fn(b: 32) -> 32 {
        ret a + b;
    };
}

let adder := add(10);
let sum := adder(20); // sum = 30
```

## User defined types
types are user defined data types. They can be structs, or a function signature or a wrapper around a built-in type.
```rust
type Circle struct {
    radius: f32;
}

type FnType fn(a: 32, b: 32) -> 32;

type WrapperInt 32;

let c := @Circle {
    radius: 10.0
};

let f: FnType = fn(a: 32, b: 32) -> 32 {
    ret a + b;
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

## Interface
Interfaces are a way to define a contract that a type must implement. It is a way to achieve polymorphism in the language.
```rust
type Shape interface {
    fn area() -> f32;
}

```

## Implementing a interface for a struct
```rust
type Printable interface {
    fn print();
}

type Person struct {
    name: str;
    age: 32;
}

impl Person {
    fn print() {
        print("Name: ", this.name, " Age: ", this.age);
    }
}

let p := @Person {
    name: "John",
    age: 20
};

fn printPerson(p: Printable) {
    p.print();
}

```