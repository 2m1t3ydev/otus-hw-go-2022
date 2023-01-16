package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("args error count, need more than 2")
	}

	envPath := os.Args[1]
	commandAndParams := os.Args[2:]

	env, err := ReadDir(envPath)
	if err != nil {
		log.Fatalln(err)
	}

	exitCode := RunCmd(commandAndParams, env)
	os.Exit(exitCode)
}
