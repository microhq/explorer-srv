package main

import (
	log "github.com/golang/glog"
	p "github.com/myodc/explorer-srv/profile/handler"
	s "github.com/myodc/explorer-srv/service/handler"
	u "github.com/myodc/explorer-srv/user/handler"
	"github.com/myodc/go-micro/cmd"
	"github.com/myodc/go-micro/server"
)

func main() {
	cmd.Init()

	server.Init(
		server.Name("go.micro.srv.explorer"),
	)

	server.Handle(
		server.NewHandler(
			new(s.Service),
		),
	)

	server.Handle(
		server.NewHandler(
			new(p.Profile),
		),
	)

	server.Handle(
		server.NewHandler(
			new(u.User),
		),
	)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
