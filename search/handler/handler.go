package handler

import (
	"encoding/json"
	"fmt"

	es "github.com/mattbaird/elastigo/lib"
	search "github.com/micro/explorer-srv/proto/search"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

var (
	c = es.NewConn()
)

type Search struct{}

func (s *Search) Create(ctx context.Context, req *search.CreateRequest, rsp *search.CreateResponse) error {
	if len(req.Document.Index) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Create", "index cannot be blank")
	}
	if len(req.Document.Type) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Create", "type cannot be blank")
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(req.Document.Data), &data); err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Create", err.Error())
	}

	_, err := c.Index(req.Document.Index, req.Document.Type, req.Document.Id, nil, data)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Create", err.Error())
	}
	return nil
}

func (s *Search) Read(ctx context.Context, req *search.ReadRequest, rsp *search.ReadResponse) error {
	if len(req.Index) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Read", "index cannot be blank")
	}
	if len(req.Type) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Read", "type cannot be blank")
	}
	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Read", "id cannot be blank")
	}

	var data string
	err := c.GetSource(req.Index, req.Type, req.Id, nil, data)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Read", err.Error())
	}

	rsp.Document = &search.Document{
		Id:   req.Index,
		Type: req.Type,
		Data: data,
	}
	return nil
}

func (s *Search) Update(ctx context.Context, req *search.UpdateRequest, rsp *search.UpdateResponse) error {
	if len(req.Document.Index) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Update", "index cannot be blank")
	}
	if len(req.Document.Type) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Update", "type cannot be blank")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(req.Document.Data), &data); err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Update", err.Error())
	}
	_, err := c.Index(req.Document.Index, req.Document.Type, req.Document.Id, nil, data)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Update", err.Error())
	}
	return nil
}

func (s *Search) Delete(ctx context.Context, req *search.DeleteRequest, rsp *search.DeleteResponse) error {
	if len(req.Index) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Delete", "index cannot be blank")
	}
	if len(req.Type) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Delete", "type cannot be blank")
	}
	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Delete", "id cannot be blank")
	}

	_, err := c.Delete(req.Index, req.Type, req.Id, nil)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Delete", err.Error())
	}
	return nil
}

func (s *Search) Search(ctx context.Context, req *search.SearchRequest, rsp *search.SearchResponse) error {
	if len(req.Index) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Search", "index cannot be blank")
	}
	if len(req.Type) == 0 {
		return errors.BadRequest("go.micro.srv.explorer.Search", "type cannot be blank")
	}
	if len(req.Query) > 2048 {
		return errors.BadRequest("go.micro.srv.explorer.Search", fmt.Sprintf("%s %d", "query too long, length", len(req.Query)))
	}

	se := es.Search(req.Index).Type(req.Type)
	from := fmt.Sprintf("%d", req.Offset)
	size := fmt.Sprintf("%d", req.Limit)

	if len(req.Query) > 0 {
		var nq string
		for _, r := range req.Query {
			st := string(r)
			switch st {
			case "+", "-", "&&", "||", "!", "(", ")", "{", "}", "[", "]", "^", "\"", "~", "*", "?", ":", "\\", "/":
				nq += "\\"
			}
			nq += st
		}
		se = se.Search(nq)
	}
	out, err := se.From(from).Size(size).Result(c)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Search", err.Error())
	}

	for _, hit := range out.Hits.Hits {
		d, _ := hit.Source.MarshalJSON()
		rsp.Documents = append(rsp.Documents, &search.Document{
			Index: hit.Index,
			Type:  hit.Type,
			Id:    hit.Id,
			Data:  string(d),
		})
	}

	return nil
}
