package main

import (
	log "github.com/golang/glog"
	p "github.com/micro/explorer-srv/profile/handler"
	s "github.com/micro/explorer-srv/service/handler"
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

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
