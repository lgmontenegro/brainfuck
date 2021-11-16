# My Brainfuck Interpreter

## An assessment exercise

### The accepted challenge:

This is a small assignment, while being able to show a bit of application design solutions.

For example requirements are built in a way to allow to mostly adhere to SOLID principles.
 
While not listed in requirements, if provided(highly advised) - tests will show testing techniques.

Application will also be manually tested for cases specific to codebase and if questions will arise.

 

#### Requirements:

* Make it stack-based.

> * This should mean that every loop must be in stack one level below current.
> * For example for script ‘+[-]-’ it could mean that operations would be as such: “add +1, run loop, sub 1”, where “run loop” includes necessary operations and actions inside.
> * Nested loops must be supported.

* Make it read input from io.Reader without knowing all input at once.

> * This requirement excludes look-ahead. For example to find closing bracket for loops.
> * As such next token should be read and processed only after current one is finished processing.
> * Loops preparation might be using the same process to fill them. This means whilst creating a loop – it can read necessary tokens from the input, while still adhering to this rule.

* Make it a library.

> * Other developers must be able to import this library and use it from inside their code.
> * This means that all necessary functionality MUST be exported and freely available to developer.

* Make it extensible.

> * Have functionality(functions/methods) to add and remove custom operations at runtime, before parsing the input.
> * No synchronization required: developers will add/remove operations only before parsing any script.
> * For example if developer that will use this library would want to add another operation(to square current cell value for example) to it – he MUST be able to do that. Same for removing operations.