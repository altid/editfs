package main

import "github.com/altid/libs/fs"

type server struct{}

func (s *server) Run(ctrl *fs.Control, cmd *fs.Command) error {
	switch cmd.Name {
	case "open":
		return open(s, ctrl, cmd)
	case "save":
		//return save(s, ctrl, cmd)
	case "close":
		return ctrl.DeleteBuffer(cmd.Args[0], "document")
	}
	return nil
}

func (s *server) Quit() {

}
