package session

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/altid/libs/service/commander"
)

func open(s *Session, cmd *commander.Command) error {

	// Create a ring buffer for edits in the server, and all modifications (Sam commands, etc) will be done via ctl messages
	// Clients modifying the underlying files should be wrapped in ctl messages outright, rather than allowing it
	stat, err := os.Stat(cmd.Args[0])
	if err != nil {
		return err
	}

	// TODO: Use status for the working directory
	switch stat.IsDir() {
	case true:
		return filepath.Walk(cmd.Args[0], func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if e := s.ctrl.CreateBuffer(p); e != nil {
				return e
			}

			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}

			mw, err := s.ctrl.MainWriter(p)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(mw, "%s\n", b)
			return err
		})
	case false:
		p := cmd.Args[0]
		if(!s.ctrl.HasBuffer(p)) {
			if e := s.ctrl.CreateBuffer(p); e != nil {
				return e
			}
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}

		mw, err := s.ctrl.MainWriter(p)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(mw, "%s\n", b)
		return err
	}

	return nil
}
