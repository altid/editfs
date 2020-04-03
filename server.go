package main

import "github.com/altid/libs/fs"

type server struct{}

func (s *server) Run(*fs.Control, *fs.Command) error {
	return nil
}

func (s *server) Quit() {

}
