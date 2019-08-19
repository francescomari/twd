package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

type arguments struct {
	self   string
	print  bool
	root   string
	prefix string
	cmd    []string
}

func parseArguments(args []string) (*arguments, error) {
	type state int

	const (
		stateStart state = iota
		stateArgs
		stateRoot
		statePrefix
		stateCmd
	)

	var a arguments

	s := stateStart

	for _, arg := range args {
		switch s {
		case stateStart:
			a.self = arg
			s = stateArgs
		case stateArgs:
			switch arg {
			case "-print":
				a.print = true
				s = stateArgs
			case "-root":
				s = stateRoot
			case "-prefix":
				s = statePrefix
			default:
				a.cmd = append(a.cmd, arg)
				s = stateCmd
			}
		case stateRoot:
			a.root = arg
			s = stateArgs
		case statePrefix:
			a.prefix = arg
			s = stateArgs
		case stateCmd:
			a.cmd = append(a.cmd, arg)
			s = stateCmd
		}
	}

	if s != stateCmd {
		return nil, fmt.Errorf("invalid arguments")
	}

	return &a, nil
}

func runCommand(args *arguments) error {
	root := args.root

	if root == "" {
		root = os.TempDir()
	}

	prefix := args.prefix

	if prefix == "" {
		prefix = "twd-"
	}

	dir, err := ioutil.TempDir(root, prefix)
	if err != nil {
		return fmt.Errorf("create temporary directory: %v", err)
	}

	cmd, err := exec.LookPath(args.cmd[0])
	if err != nil {
		return fmt.Errorf("lookup program: %v", err)
	}

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("change directory: %v", err)
	}

	if args.print {
		fmt.Println(dir)
	}

	if err := syscall.Exec(cmd, args.cmd, os.Environ()); err != nil {
		return fmt.Errorf("execute program: %v", err)
	}

	return nil
}

func main() {
	args, err := parseArguments(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "usage: %v [-print] [-root root] [-prefix prefix] command...\n", os.Args[0])
		os.Exit(1)
	}

	if err := runCommand(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
