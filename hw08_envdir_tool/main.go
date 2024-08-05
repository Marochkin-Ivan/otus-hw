package main

import (
	"log"
	"os"
)

func main() {
	// Place your code here.
	args := os.Args
	if len(args) < 3 {
		log.Fatalln("not enough arguments")
		return
	}

	env, err := ReadDir(args[1])
	if err != nil {
		log.Println(err)
		return
	}

	os.Exit(RunCmd(args[2:], env))
}
