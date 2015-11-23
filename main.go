package main

import (
	log "github.com/golang/glog"
	org "github.com/micro/explorer-srv/organization/handler"
	p "github.com/micro/explorer-srv/profile/handler"
	se "github.com/micro/explorer-srv/search/handler"
	s "github.com/micro/explorer-srv/service/handler"
	t "github.com/micro/explorer-srv/token/handler"
	u "github.com/micro/explorer-srv/user/handler"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
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

	server.Handle(
		server.NewHandler(
			new(t.Token),
		),
	)

	server.Handle(
		server.NewHandler(
			new(se.Search),
		),
	)

	server.Handle(
		server.NewHandler(
			new(org.Organization),
		),
	)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
