package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...

type nullReader struct{}
type nullWriter struct{}

func (nullReader) Read(p []byte) (n int, err error)  { return len(p), nil }
func (nullWriter) Write(p []byte) (n int, err error) { return len(p), nil }

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	//fmt.Println("Logs from your program will appear here!")

	chroot_path := "./tmp"

	err := create_chroot_jail(chroot_path)
	if err != nil {
		fmt.Printf("Chroot Error %v", err)
		os.Exit(1)
	}

	copy_docker_explore(chroot_path)

	command := os.Args[3]
	args := os.Args[4:len(os.Args)]
	// arg := strings.Join(args, " ")

	//
	cmd := exec.Command(command, args...)
	cmd.Stdin = nullReader{}
	cmd.Stdout = nullWriter{}
	cmd.Stderr = nullWriter{}
	cmd.Dir = chroot_path

	err = cmd.Run()

	exit_code := cmd.ProcessState.ExitCode()

	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(exit_code)
	}

	os.Exit(exit_code)
}

func create_chroot_jail(path string) error {
	err := os.MkdirAll(path, 0750)
	if err != nil {
		return err
	}

	err = unix.Chroot(path)
	if err != nil {
		return err
	}

	return nil
}

func copy_docker_explore(chroot_path string) error {
	src_path := "/usr/local/bin/docker-explorer"
	dst_path := "/usr/local/bin/docker-explorer"

	src_file, err := os.Open(src_path)
	if err != nil {
		return err
	}

	defer src_file.Close()

	dst_file, err := os.Open(dst_path)
	if err != nil {
		return err
	}

	defer dst_file.Close()

	err = os.Chdir(chroot_path)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst_file, src_file)
	if err != nil {
		return err
	}

	return nil

}
