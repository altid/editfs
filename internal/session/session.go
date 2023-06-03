package session

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/altid/libs/markup"
	"github.com/altid/libs/service/commander"
	"github.com/altid/libs/service/controller"
)

type ctlItem int

const (
	ctlCommand = iota
	ctlSucceed
	ctlStart
	ctlInput
	ctlErr
)

type Session struct {
	Defaults *Defaults
	Verbose  bool
	ctx      context.Context
	cancel   context.CancelFunc
	ctrl     controller.Controller
	debug    func(ctlItem, ...any)
}

type Defaults struct {
	Address string `altid:"address,no_prompt"`
	Port    int    `altid:"port,no_prompt"`
	SSL     string `altid:"ssl,prompt:SSL mode,pick:none|simple|certificate"`
	TLSCert string `altid:"tlscert,no_prompt"`
	TLSKey  string `altid:"tlskey,no_prompt"`
}

func (s *Session) Parse(ctx context.Context) {
	s.debug = func(ctlItem, ...any) {}
	s.ctx, s.cancel = context.WithCancel(ctx)
	if s.Verbose {
		s.debug = ctlLogging
	}
}

func (s *Session) Connect(Username string) error {
	return nil
}

func (s *Session) Run(ctrl controller.Controller, cmd *commander.Command) error {
	switch cmd.Name {
	case "x":
		//return x(s, ctrl, cmd)
		return nil
	case "open":
		return open(s, cmd)
	case "save":
		//return save(s, ctrl, cmd)
		return nil
	case "close":
		return ctrl.DeleteBuffer(cmd.Args[0])
	default:
		return fmt.Errorf("unsupported command %s", cmd.Name)
	}
}

func (s *Session) Quit() {
	s.cancel()
}

func (s *Session) Handle(bufname string, l *markup.Lexer) error {
	// We handle only commands and outright appends, which could be confusing
	return nil
}

func (s *Session) Start(c controller.Controller) error {
	// We don't really need to do anything here
	s.ctrl = c
	return nil
}

func (s *Session) Listen(c controller.Controller) {
	s.Start(c)
	<-s.ctx.Done()
}

func (s *Session) Command(cmd *commander.Command) error {
	return s.Run(s.ctrl, cmd)
}

func ctlLogging(ctl ctlItem, args ...any) {
	l := log.New(os.Stdout, "editfs ", 0)
	switch ctl {
	case ctlCommand:
		m := args[0].(*commander.Command)
		l.Printf("command name=\"%s\" heading=\"%d\" sender=\"%s\" args=\"%s\" from=\"%s\"", m.Name, m.Heading, m.Sender, m.Args, m.From)
	case ctlSucceed:
		l.Printf("%s succeeded\n", args[0])
	case ctlStart:
		l.Printf("start: addr=\"%s\", port=%d\n", args[0], args[1])
	case ctlInput:
		l.Printf("input: data=\"%s\" bufname=\"%s\"", args[0], args[1])
	case ctlErr:
		l.Printf("error: err=\"%v\"\n", args[0])
	}
}
