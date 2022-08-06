package main

import (
	"os"

	"watzek.dev/apps/dokumentatorin/commands"
)

func main() {
	err := commands.Execute(os.Args[1:])
	if err != nil {
		panic(err)
	}
}
