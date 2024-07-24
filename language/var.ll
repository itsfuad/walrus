; Declare the main function
define i32 @main() {
entry:
    ; Allocate space for the variable 'a' on the stack
    %a = alloca i32, align 4

    ; Store the value 4 in a temporary variable
    %4val = add i32 4, 6

    ; Store the result in the variable 'a'
    store i32 %4val, i32* %a, align 4

    ; Load the value of 'a' for returning
    %result = load i32, i32* %a, align 4

    ; Return the result
    ret i32 %result
}