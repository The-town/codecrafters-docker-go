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
	output, err := cmd.CombinedOutput()
	process := cmd.Process
	process_state, _ := process.Wait()

	if err != nil {
		fmt.Printf("Err: %v", err)

		exit_code := process_state.ExitCode()
		fmt.Println(exit_code)
		os.Exit(exit_code)
	}

	fmt.Printf(string(output))
	fmt.Fprint(os.Stderr, string(output))
}
