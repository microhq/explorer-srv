package handler

import (
	"golang.org/x/net/context"

	"github.com/myodc/explorer-srv/service/db"
	srv "github.com/myodc/explorer-srv/proto/service"
)

type Service struct{}

func (s *Service) Create(ctx context.Context, req *srv.CreateRequest, rsp *srv.CreateResponse) error {
	return db.Create(req.Service)
}

func (s *Service) Read(ctx context.Context, req *srv.ReadRequest, rsp *srv.ReadResponse) error {
	service, err := db.Read(req.Id)
	if err != nil {
		return err
	}
	rsp.Service = service
        return nil
}

func (s *Service) Update(ctx context.Context, req *srv.UpdateRequest, rsp *srv.UpdateResponse) error {
        return db.Update(req.Service)
}

func (s *Service) Delete(ctx context.Context, req *srv.DeleteRequest, rsp *srv.DeleteResponse) error {
        return db.Delete(req.Id)
}

func (s *Service) Search(ctx context.Context, req *srv.SearchRequest, rsp *srv.SearchResponse) error {
	services, err := db.Search(req.Name, req.Owner, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Services = services
	return nil
}
