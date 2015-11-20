package handler

import (
	"golang.org/x/net/context"

	"github.com/micro/explorer-srv/profile/db"
	prf "github.com/micro/explorer-srv/proto/profile"
)

type Profile struct{}

func (s *Profile) Create(ctx context.Context, req *prf.CreateRequest, rsp *prf.CreateResponse) error {
	return db.Create(req.Profile)
}

func (s *Profile) Read(ctx context.Context, req *prf.ReadRequest, rsp *prf.ReadResponse) error {
	profile, err := db.Read(req.Id)
	if err != nil {
		return err
	}
	rsp.Profile = profile
	return nil
}

func (s *Profile) Update(ctx context.Context, req *prf.UpdateRequest, rsp *prf.UpdateResponse) error {
	return db.Update(req.Profile)
}

func (s *Profile) Delete(ctx context.Context, req *prf.DeleteRequest, rsp *prf.DeleteResponse) error {
	return db.Delete(req.Id)
}

func (s *Profile) Search(ctx context.Context, req *prf.SearchRequest, rsp *prf.SearchResponse) error {
	profiles, err := db.Search(req.Name, req.Owner, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Profiles = profiles
	return nil
}
