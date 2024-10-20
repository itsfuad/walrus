section .data
    head dd 0                 ; Pointer to the head of the list (initially null)
    current dd 0              ; Pointer used for traversal
    format db "Value: %d", 10, 0 ; Format string for printing

    ; Define a couple of nodes in the data section
    node1 dd 5, 0             ; First node: value = 5, next = null (0)
    node2 dd 10, 0            ; Second node: value = 10, next = null (0)

section .bss
    ; Space for dynamically allocated nodes can be defined here, if needed

section .text
    global _start

extern printf                ; Declare printf as an external function

_start:
    ; Add node2 to the list
    mov eax, node2            ; Load address of node2
    call list_add             ; Add it to the list

    ; Add node1 to the list
    mov eax, node1            ; Load address of node1
    call list_add             ; Add it to the list

    ; Traverse the list and print values
    call list_traverse

    ; Exit the program
    mov eax, 1                ; sys_exit
    xor ebx, ebx              ; Status 0
    int 0x80

; Add a node to the head of the list
list_add:
    mov ebx, [head]           ; Load current head
    mov [eax + 4], ebx        ; Set the next pointer of the new node to the current head
    mov [head], eax           ; Update the head to the new node
    ret

; Traverse the list and print each node's value
list_traverse:
    mov eax, [head]           ; Start at the head
    .loop:
        cmp eax, 0            ; If current node is null, we're done
        je .done
        push eax              ; Save the current node pointer

        mov ebx, [eax]        ; Load node value
        push ebx              ; Push value as argument to printf
        push format           ; Push format string
        call printf           ; Call printf
        add esp, 8            ; Clean up the stack (2 arguments)

        pop eax               ; Restore the current node pointer
        mov eax, [eax + 4]     ; Move to the next node
        jmp .loop

    .done:
        ret