package main // import "crawshaw.io/ltboss"

import (
	"fmt"
	"os"
)

var cmdname = "ltboss"

func main() {
	cmdname = os.Args[0]

	if len(os.Args) == 1 {
		help(nil)
	} else if os.Args[1] == "help" {
		help(os.Args[1:])
	}

	for _, cmd := range commands {
		if cmd.name == os.Args[1] {
			cmd.run(os.Args[1:])
			panic("run should exit")
		}
	}
	fmt.Fprintf(os.Stderr, "%s: unknown subcommand %q\nRun '%s help' for usage.\n", cmdname, os.Args[1], cmdname)
	os.Exit(2)
}

type command struct {
	name     string
	oneLiner string
	usage    string
	docs     string
	run      func(args []string)
}

var commands = []command{
	{
		name:     "start",
		oneLiner: "create a new service",
		usage:    `start [-name servicename] [service flags] binpath [binflags]`,
		docs:     "TODO",
		run:      func(args []string) { fmt.Printf("TODO start\n") },
	},
	{
		name:     "stop",
		oneLiner: "shut down a running service",
		usage:    `stop servicename`,
		docs:     "TODO",
		run:      func(args []string) { fmt.Printf("TODO stop\n") },
	},
	{
		name:     "reload",
		oneLiner: "replace a running service with a new process",
		usage:    `reload [-timeout duration] [service flags] servicename`,
		docs:     "TODO",
		run:      func(args []string) { fmt.Printf("TODO reload\n") },
	},
	{
		name:     "show",
		oneLiner: "details about a running service",
		usage:    `show [service name]`,
		docs:     "TODO",
		run:      func(args []string) { fmt.Printf("TODO show\n") },
	},
	{
		name:     "ls",
		oneLiner: "list services",
		usage:    `ls [pattern]`,
		docs:     "TODO",
		run:      func(args []string) { fmt.Printf("TODO ls\n") },
	},
}
