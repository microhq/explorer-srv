package handler

import (
	"crypto/rand"
	"time"

	token "github.com/micro/explorer-srv/proto/token"
	"github.com/micro/explorer-srv/token/db"
	"github.com/micro/go-micro/errors"
	uuid "github.com/streadway/simpleuuid"
	"golang.org/x/net/context"
)

var (
	alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func random(i int) string {
	bytes := make([]byte, i)
	for {
		rand.Read(bytes)
		for i, b := range bytes {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		}
		return string(bytes)
	}
	return "fuckyou"
}

type Token struct{}

func (s *Token) Generate(ctx context.Context, req *token.GenerateRequest, rsp *token.GenerateResponse) error {
	name := random(16)
	id, err := uuid.NewTime(time.Now())
	if err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Generate", err.Error())
	}
	namespace := "default"
	if len(req.Namespace) > 0 {
		namespace = req.Namespace
	}
	tk := &token.Token{
		Id:        id.String(),
		Namespace: namespace,
		Name:      name,
		Created:   time.Now().Unix(),
		Updated:   time.Now().Unix(),
	}
	if err := db.Create(tk); err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Generate", err.Error())
	}
	rsp.Token = tk
	return nil
}

func (s *Token) Create(ctx context.Context, req *token.CreateRequest, rsp *token.CreateResponse) error {
	if err := db.Create(req.Token); err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Create", err.Error())
	}
	return nil
}

func (s *Token) Read(ctx context.Context, req *token.ReadRequest, rsp *token.ReadResponse) error {
	token, err := db.Read(req.Id)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Read", err.Error())
	}
	rsp.Token = token
	return nil
}

func (s *Token) Update(ctx context.Context, req *token.UpdateRequest, rsp *token.UpdateResponse) error {
	if err := db.Update(req.Token); err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Update", err.Error())
	}
	return nil
}

func (s *Token) Delete(ctx context.Context, req *token.DeleteRequest, rsp *token.DeleteResponse) error {
	if err := db.Delete(req.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Delete", err.Error())
	}
	return nil
}

func (s *Token) Search(ctx context.Context, req *token.SearchRequest, rsp *token.SearchResponse) error {
	tokens, err := db.Search(req.Namespace, req.Name, req.Limit, req.Offset)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Search", err.Error())
	}
	rsp.Tokens = tokens
	return nil
}
