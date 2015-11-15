package handler

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/crypto/bcrypt"
	"github.com/myodc/go-micro/errors"
	"github.com/myodc/explorer-srv/user/db"
	user "github.com/myodc/explorer-srv/proto/user"
)

const (
	x = "cruft123"
)

var (
	alphanum      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
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

type User struct{}

func (s *User) Create(ctx context.Context, req *user.CreateRequest, rsp *user.CreateResponse) error {
	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Password), 10)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Create", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)
	return db.Create(req.User, salt, pp)
}

func (s *User) Read(ctx context.Context, req *user.ReadRequest, rsp *user.ReadResponse) error {
	user, err := db.Read(req.Id)
	if err != nil {
		return err
	}
	rsp.User = user
        return nil
}

func (s *User) Update(ctx context.Context, req *user.UpdateRequest, rsp *user.UpdateResponse) error {
        return db.Update(req.User)
}

func (s *User) Delete(ctx context.Context, req *user.DeleteRequest, rsp *user.DeleteResponse) error {
        return db.Delete(req.Id)
}

func (s *User) Search(ctx context.Context, req *user.SearchRequest, rsp *user.SearchResponse) error {
	users, err := db.Search(req.Username, req.Email, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Users = users
	return nil
}

func (s *User) Login(ctx context.Context, req *user.LoginRequest, rsp *user.LoginResponse) error {
	salt, hashed, err := db.SaltAndPassword(req.Username, req.Email)
	if err != nil {
		return err
	}

	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Login", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.Password)); err != nil {
		return errors.Unauthorized("go.micro.srv.explorer.login", err.Error())
	}
	// save session
	sess := &user.Session{
		Id: random(128),
		Username: req.Username,
		Created: time.Now().Unix(),
		Expires: time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	if err := db.CreateSession(sess); err != nil {
		return errors.InternalServerError("go.micro.srv.explorer.Login", err.Error())
	}
	rsp.Session = sess
	return nil
}

func (s *User) Logout(ctx context.Context, req *user.LogoutRequest, rsp *user.LogoutResponse) error {
	return db.DeleteSession(req.SessionId)
}

func (s *User) ReadSession(ctx context.Context, req *user.ReadSessionRequest, rsp *user.ReadSessionResponse) error {
	sess, err := db.ReadSession(req.SessionId)
	if err != nil {
		return err
	}
	rsp.Session = sess
        return nil
}
