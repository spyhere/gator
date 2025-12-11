package main

import "fmt"

func extractCommand(args []string) (command, error) {
	if len(args) < 2 {
		return command{}, fmt.Errorf("No commands with arguments were provided, exiting...")
	}
	commandName := args[1]
	var commandArgs []string
	if len(args) > 2 {
		commandArgs = args[2:]
	}
	return command{
		name: commandName,
		args: commandArgs,
	}, nil
}
