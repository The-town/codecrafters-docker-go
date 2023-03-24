package main

import (
	"fmt"
	"os/exec"

	// Uncomment this block to pass the first stage!
	"os"
)

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	//fmt.Println("Logs from your program will appear here!")

	command := os.Args[3]
	args := os.Args[4:len(os.Args)]
	// arg := strings.Join(args, " ")

	//
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

	exit_code := cmd.ProcessState.ExitCode()
	fmt.Println(exit_code)
	os.Exit(exit_code)
}
