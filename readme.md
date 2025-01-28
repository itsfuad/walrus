# Walrus Programming Language
Walrus is a tiny and simple programming language designed for simplicity. Its syntax draws inspiration from Go, Rust, and TypeScript.

# Language Design Principles
1. **Expressive**: Just look at the code and you can tell what's going on without jumping around codes.
2. **Beginner friendly**
3. **Distinct behavior for every things** (most of the cases). Like in other language type cast often uses int(2.4) or (int)2.4; which looks like function calls. We have separate keyword 'as' for type casting like typescript keeping it as unambiguous as we can.
4. **Less over complications**: we removed OOP instead used structs which support methods, access modifiers and has interfaces with duck typing. We have simple approach. 

## Features

### Lexer
- Tokenizes the source code for parsing.

### Parser
- Handles syntactic analysis for:
  - **Builtin types**
    - Signed Integers: `i8`, `i16`, `i32`, `i64`, `i128`
    - Unsigned Intergers: `u8`, `u16`, `u32`, `u64`, `u128`
    - Floats: `f32`, `f64`
    - String: `str`
    - Null: `null`
    - Void: `void`
    - Map: `map[key]value`
    - Maybe: `maybe{type}`
  - **Variable Declaration and Assignment**
    - Mutable variables with `let`
    - Constant variables with `const`
    - Multiple variable declarations in one line
  - **Expressions**
    - Unary: `-`, `!`
    - Additive: `+`, `-`
    - Multiplicative: `*`, `/`, `%`, `^`
    - Grouping: `( )`
    - Type casting using `as`
  - **Data Structures**
    - Arrays: Indexing and assignment
    - Maps: Indexing and key-value assignments
  - **Conditionals**
    - `if`, `else if`, `else`
  - **Functions**
    - Declaration, calls, return values
    - Optional parameters
    - First-class functions and closures
  - **User-Defined Types**
    - Structs: Property access and assignment
    - Interfaces: Definition, implementation, and usage
  - **Operators**
    - Increment/Decrement: Prefix and Postfix
    - Assignment: `+=`, `-=`, `*=`, `/=`, `%=`
  - **Additional Constructs**
    - Switch statements
    - For loops (syntax under development)
  - **Rich Error Reporting**
    - Displays multiple errors during parsing and type checking

### Type Checking
- Ensures type safety across all constructs.
- Handles all parser-supported features except for loops and imports (in progress).

### Code Generation
- Planned for future releases.

# Installation and Usage

## Install Go
To run the compiler, you need to have go installed. You can download it from [here](https://golang.org/dl/)

## Testing a walrus file
To test a walrus file, you need to open `main.go` and change the `filePath` variable to the path of the file you want to test.
```go
func main() {
    filePath := "filename.wal"
    // ... rest of the code
}
```
Then run the following command
```sh
go run main.go
```
Or, if you're on windows then you can run the batch file `run.bat`
```sh
./run
```

# Running the tests
To run the tests, run the following command
```sh
go test ./...
```
Or, if you're on windows then you can run the batch file `test.bat`
```sh
./test
```

# Installing syntax highlighting for vscode
To install local build, go to `language-support` directory and run the following command
```sh
npm install -g vsce
```
## Build the extension
If you don't have node installed, you can download it from [here](https://nodejs.org/en/download/)
Then run the following command
```sh
vsce package
code --install-extension walrus-<full-output-name>.vsix
```
Or open vscode, go to extensions and click on the three dots on the top right corner and click on 'Install from VSIX' and select the generated vsix file.

## Install from marketplace
Or, you can install the extension from the marketplace [here](https://marketplace.visualstudio.com/items?itemName=Walrus.walrus)
Or, Search for 'Walrus' in the vscode extensions marketplace.

# Example

## Variable declare and assign
```rs
// Declare a variable with let or const keyword
let a := 10; // The variable is mutable and its type is inferred from the value e.g. i32
const pi := 3.14; // constant variable with type f32

// Declare a variable with type
let b: i32 = 20; // The variable is mutable and its type is i32
const c: f32 = 3.14; // constant variable with type f32

let unsigned: u32 = 10; // Unsigned integer of 32 bits
```
You can also declare multiple variables in a single line
```rs
//multiple variable declaration in one line
// with type
let t1: i32 = 43, t2: f32 = 3.5, t3: str;
// without type
let t4 := 43, t5 := 3.14, t6 := "hello";
```

## Variable assign
```rs
let a := 10;
a = 20; // Assign a new value to a
```

## Expressions
```rs
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
```rs
let a := 10;
let b := 20;
let c := (a + b) * 2; // c = 60
```

## Type cast
Type must be explicit in walrus. Type cast is done using 'as' keyword. For example, we cannot assign to a 8 bit integer to a 32 bit integer and vice versa. We need to cast the type. 
But when we declare a variable with explicit type, the value is implicitly casted to the type.
```rs
let a := 10;
let b := a as f32; // b = 10.0 as float32
```

## Array
```rs
let a := [1, 2, 3, 4, 5]; // Array of integers. One dimensional array
let b := a[0]; // b = 1
a[0] = 10; // a = [10, 2, 3, 4, 5]

// Array of arrays
let c := [[1, 2], [3, 4], [5, 6]]; // Array of arrays of integers
let d := c[0][0]; // d = 1
c[0][0] = 10; // c = [[10, 2], [3, 4], [5, 6]]
```

## Map
```rs
let myMap : map[str]i32 = map[str]i32 {
    "a" => 10,
    "b" => 20,
    "c" => 30
};

let a := myMap["a"]; // a = 10
myMap["a"] = 20; // myMap = {"a": 20, "b": 20, "c": 30}
```

## Maybe
Maybe is a type which can be null or the type specified. It is used to handle null values.
```rs
let a: maybe{i32} = null; // a can be null or an integer

a = 10; // a can be assigned an integer but the value is still maybe{i32}
a = null; // a can be assigned null

let b := 12;

//Now to use the value of a, we need to check if it is null or not
safe a {
    // when a is not null, this block will be executed
    b = 23 + a; // b = 33
} otherwise {
    // when a is null, this block will be executed
    b = 23;
}
```

## Struct
```rs
type Person struct{
    name: str;
    age: i32;
}

//Assign the type with @Name syntax. So we can distinguish between type and variable.
let p := Person {
    name: "John",
    age: 20
};

// We could also assign the type with type inference
let p : Person = Person { // Here 'Person' is the type, @Person is the type instance
    name: "John",
    age: 20
};
```

## Struct property access
```rs
let p := Person {
    name: "John",
    priv age: 20 // Private property
};

let name := p.name; // name = "John"
let age := p.age; // Error: Cannot access private property
```

## Conditionals
```rs
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
```rs
fn add(a: i32, b: i32) -> i32 {
    ret a + b;
}

let sum := add(10, 20); // sum = 30

// function with optional parameters
fn add(a: i32, b: i32, c?: i32 = 0) -> i32 {
    ret a + b + c;
}

let sum := add(10, 20); // sum = 30

// functions are first class citizens so we can assign them to variables
let adder := add;

let sum := adder(10, 20); // sum = 30

// function with closure
fn add(a: i32) -> fn(i32) -> i32 {
    ret fn(b: i32) -> i32 {
        ret a + b;
    };
}

let adder := add(10);
let sum := adder(20); // sum = 30
```

## Closure
Closures in simple terms are functions which are defined inside another function. They can access the variables of the parent function.
So when you return a function from a function, it is called a closure.
```rs
fn closure(a: i32) -> fn (b: i32) -> i32 {
    //return a function which takes an integer and returns the sum of a and b
    ret fn (b: i32) -> i32 {
        ret a + b;
    };
}

const addRes := add(1);

const closureRes1 := closure(1); // returns a function which returns a + b
const closureRes2 := closureRes1(2); // returns a + b
const closureRes3 := closure(1)(2); // one liner version of the above two lines
```

## User defined types
types are user defined data types. They can be structs, or a function signature or a wrapper around a built-in type.
```rs
type Circle struct {
    radius: f32;
}

type FnType fn(a: i32, b: i32) -> i32;

type WrapperInt i32;

let c := Circle {
    radius: 10.0
};

let f: FnType = fn(a: i32, b: i32) -> i32 {
    ret a + b;
};

let w: WrapperInt = 10;
```

## Increment/Decrement
```rs
let a := 10;
let b := ++a; // b = 11, a = 11
let c := a++; // c = 11, a = 12
let d := --a; // d = 11, a = 11
let e := a--; // e = 11, a = 10
```

## Assignment operators
```rs
let a := 10;
a += 10; // a = 20
a -= 10; // a = 10
a *= 10; // a = 100
a /= 10; // a = 10
a %= 10; // a = 0
```

## For loop
Syntax is not finalized yet


## Switch
```rs
let a := 10;

switch a {
    case 10: {
        print("a is 10");
    }
    case 20: {
        print("a is 20");
    }
    default: {
        print("a is neither 10 nor 20");
    }
}
```

## Interface
Interfaces are a way to define a contract that a type must implement. It is a way to achieve polymorphism in the language.
```rs
type Shape interface {
    fn area() -> f32;
}

```

## Implementing a interface for a struct
```rs
type Printable interface {
    fn print();
}

type Person struct {
    name: str;
    age: i32;
}

impl Person {
    fn print() {
        print("Name: ", this.name, " Age: ", this.age);
    }
}

let p := Person {
    name: "John",
    age: 20
};

fn printPerson(p: Printable) {
    p.print();
}

```

## Roadmap
- [ ] For loops
- [ ] Imports and modules
- [ ] Generics
- [ ] Advanced code generation

Stay tuned for updates and contribute to the project!
