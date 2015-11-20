package handler

import (
	"golang.org/x/net/context"

	srv "github.com/micro/explorer-srv/proto/service"
	"github.com/micro/explorer-srv/service/db"
)

func (s *Service) CreateVersion(ctx context.Context, req *srv.CreateVersionRequest, rsp *srv.CreateVersionResponse) error {
	return db.CreateVersion(req.Version)
}

func (s *Service) ReadVersion(ctx context.Context, req *srv.ReadVersionRequest, rsp *srv.ReadVersionResponse) error {
	version, err := db.ReadVersion(req.Id)
	if err != nil {
		return err
	}
	rsp.Version = version
	return nil
}

func (s *Service) UpdateVersion(ctx context.Context, req *srv.UpdateVersionRequest, rsp *srv.UpdateVersionResponse) error {
	return db.UpdateVersion(req.Version)
}

func (s *Service) DeleteVersion(ctx context.Context, req *srv.DeleteVersionRequest, rsp *srv.DeleteVersionResponse) error {
	return db.DeleteVersion(req.Id)
}

func (s *Service) SearchVersion(ctx context.Context, req *srv.SearchVersionRequest, rsp *srv.SearchVersionResponse) error {
	versions, err := db.SearchVersion(req.ServiceId, req.Version, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Versions = versions
	return nil
}
