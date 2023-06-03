package editfs

import (
	"context"
	"fmt"

	"github.com/altid/editfs/internal/commands"
	"github.com/altid/editfs/internal/session"
	"github.com/altid/libs/config"
	"github.com/altid/libs/mdns"
	"github.com/altid/libs/service"
	"github.com/altid/libs/service/listener"
	"github.com/altid/libs/store"
)

type Editfs struct {
	run     func() error
	session *session.Session
	addr	string
	port	int
	name    string
	debug   bool
	mdns    *mdns.Entry
}

var defaults *session.Defaults = &session.Defaults{
	Address: "",
	Port:    564,
	SSL:     "simple",
	TLSCert: "",
	TLSKey:  "",
}

func CreateConfig(srv string, debug bool) error {
	return config.Create(defaults, srv, "", debug)
}

func Register(addr string, port int, srv string, debug bool) (*Editfs, error) {
	var err error
	if e := config.Marshal(defaults, srv, "", debug); e != nil {
		return nil, e
	}
	l, err := tolisten(defaults, addr, port, debug)
	if err != nil {
		return nil, err
	}
	session := &session.Session{
		Defaults: defaults,
		Verbose: debug,
	}
	ctx := context.Background()
	session.Parse(ctx)
	m := &Editfs{
		session: session,
		name:    srv,
		addr:	 addr,
		port:    port,
		debug:   debug,
	}
	c := service.New(srv, addr, debug)
	c.WithListener(l)
	c.WithStore(store.NewRamstore(debug))
	c.WithContext(ctx)
	c.WithCallbacks(session)
	c.WithRunner(session)
	c.SetCommands(commands.Commands)
	m.run = c.Listen
	return m, nil
}

func (edit *Editfs) Run() error {
	return edit.run()
}

func (edit *Editfs) Broadcast() error {
	url := fmt.Sprintf("%s:%d", edit.addr, edit.port)
	entry, err := mdns.ParseURL(url, edit.name)
	if err != nil {
		return err
	}
	if e := mdns.Register(entry); e != nil {
		return e
	}

	edit.mdns = entry
	return nil
}

func (edit *Editfs) Cleanup() {
	if edit.mdns != nil {
		edit.mdns.Cleanup()
	}
	edit.session.Quit()
}
func (edit *Editfs) Session() *session.Session {
	return edit.session
}

func tolisten(d *session.Defaults, addr string, port int, debug bool) (listener.Listener, error) {
	//if ssh {
	// 	return listener.NewListenSsh()
	//}

	dial := fmt.Sprintf("%s:%d", addr, port)
	if d.TLSKey == "none" && d.TLSCert == "none" {
		return listener.NewListen9p(dial, "", "", debug)
	}

	return listener.NewListen9p(dial, d.TLSCert, d.TLSKey, debug)
}
