package main

import (
	"fmt"
	"os"
)

func main() {
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "pull":
		pull(args)
	case "push":
		push(args)
	case "backup":
		backup(args)
	default:
		fmt.Printf("unknown command %q\n", cmd)
		os.Exit(1)
	}
}

// pull gets chords from server.
func pull(args []string) {
	fmt.Println("pulling", args)
}

// push updates chords on the server.
func push(args []string) {
	fmt.Println("pushing", args)
}

// backup makes a full backup of the database.
func backup(args []string) {
	fmt.Println("backing up", args)
}
