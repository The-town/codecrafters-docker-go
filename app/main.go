package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	//fmt.Println("Logs from your program will appear here!")

	chroot_path := "./tmp"

	dirwalk("/usr/bin")

	err := create_chroot_jail(chroot_path)
	if err != nil {
		fmt.Printf("Chroot Error %v", err)
		os.Exit(1)
	}

	err = copy_docker_explore(chroot_path)
	if err != nil {
		fmt.Printf("Copy Docker Explore Error %v", err)
		os.Exit(1)
	}

	command := os.Args[3]
	args := os.Args[4:len(os.Args)]
	// arg := strings.Join(args, " ")

	//
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
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

	// Go言語の仕様としてchrootする際は/dev/nullが必要
	// https://rohitpaulk.com/articles/cmd-run-dev-null
	err = os.MkdirAll("/dev/null", 0750)
	if err != nil {
		return err
	}

	return nil
}

func copy_docker_explore(chroot_path string) error {
	src_path := "/usr/local/bin/docker-explorer"
	dst_path := "./docker-explorer"

	src_file, err := os.Open(src_path)
	if err != nil {
		return err
	}

	defer src_file.Close()

	err = os.Chdir(chroot_path)
	if err != nil {
		return err
	}

	os.MkdirAll("./usr/local/bin", 0750)
	os.Chdir("./usr/local/bin")

	dst_file, err := os.Create(dst_path)
	if err != nil {
		return err
	}

	defer dst_file.Close()

	_, err = io.Copy(dst_file, src_file)
	if err != nil {
		return err
	}

	return nil

}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}
