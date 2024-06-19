package printf

import (
	"fmt"

	"github.com/jwalton/gchalk"
)

func Red(f string, a ...any) {
	s := gchalk.Red(fmt.Sprintf(f, a...))
	fmt.Println(s)
}

func Green(f string, a ...any) {
	s := gchalk.Green(fmt.Sprintf(f, a...))
	fmt.Println(s)
}

func Blue(f string, a ...any) {
	s := gchalk.Blue(fmt.Sprintf(f, a...))
	fmt.Println(s)
}

func BlueS(s string, a ...any) {
	a = append([]any{gchalk.Blue(s)}, a...)
	fmt.Println(a...)
}
