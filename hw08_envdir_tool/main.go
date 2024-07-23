package main

import (
	"log"
	"os"
)

var ErrNotEnoughArgs = "not enough arguments"

func main() {
	// Place your code here.
	args := os.Args
	if len(args) < 3 {
		log.Fatalln(ErrNotEnoughArgs)
		return
	}

	env, err := ReadDir(args[1])
	if err != nil {
		log.Println(err)
		return
	}

	code := RunCmd(args[2:], env)
	log.Println(code)
}
