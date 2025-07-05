package container

import "fmt"

type Command struct {
	Action  string
	Program string
	Args    []string
}

func ParseArgs(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{Action: "help"}, nil
	}

	switch args[0] {
	case "help":
		return Command{Action: "help"}, nil
	case "run":
		if len(args) < 2 {
			return Command{}, fmt.Errorf("run command requires program argument")
		}
		return Command{
			Action:  "run",
			Program: args[1],
			Args:    args[2:],
		}, nil
	default:
		return Command{}, fmt.Errorf("unknown command: %s", args[0])
	}
}