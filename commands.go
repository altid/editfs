package main

import "github.com/altid/libs/fs"

var Commands = []*fs.Command{
	{
		Name:        "x",
		Args:        []string{""},
		Description: "Execute Sam command",
		Heading:     fs.DefaultGroup,
	},
}
