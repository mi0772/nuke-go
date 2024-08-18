package main

import (
	"fmt"

	"github.com/mi0772/nuke-go/engine"
)

func main() {
	fmt.Println("nuke cache server")

	database, error := engine.InitializeDatabase("/fakepath", 10)
	if error != nil {
		panic("ciao")
	}

	e := database.Push("carlo", []byte("Here is a string...."))
	if e != nil {
		fmt.Println(e)
	}
	e = database.Push("carlo", []byte("secondo carlo"))
	if e != nil {
		fmt.Println(e)
	}
	e = database.Push("antonio", []byte("Here is a string...."))
	if e != nil {
		fmt.Println(e)
	}
	e = database.Push("banana", []byte("Here is a string...."))
	if e != nil {
		fmt.Println(e)
	}

	v, e := database.Pop("carlo")
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println("il valore trovato è : %s", v)
	}

}
