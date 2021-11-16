package main

import (
	"fmt"
	"io"
	"lgmontenegro/brainfuck"
	"strings"
)

func main() {
	teste := brainfuck.NewCompiler()

	r := strings.NewReader(`++++++++++
[
>+++++++
>++++++++++
>+++
<<<-
]
>++.
>+.
+++++++..
+++.
>++.
<<+++++++++++++++.
>.
+++.
------.
--------.
>+.
`)

	bufferReadSize := make([]byte, 1)
	for {
		_, err := r.Read(bufferReadSize)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = teste.Compile(bufferReadSize[0])
		if err != nil {
			panic(err.Error())
		}
	}
}
