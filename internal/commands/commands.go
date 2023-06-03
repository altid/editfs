package commands

import "github.com/altid/libs/service/commander"

var Commands = []*commander.Command{
	{
		Name:        "open",
		Args:        []string{"<path>"},
		Description: "Open file at path",
		Heading:     commander.DefaultGroup,
	},
	{
		Name:        "close",
		Args:        []string{"<path>"},
		Description: "Close file at path",
		Heading:     commander.DefaultGroup,
	},
	{
		Name:        "x",
		Args:        []string{""},
		Description: "Execute Sam command",
		Heading:     commander.DefaultGroup,
	},
}
