package db

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/cockroachdb/cockroach/sql/driver"
	log "github.com/golang/glog"
	token "github.com/micro/explorer-srv/proto/token"
)

var (
	db          *sql.DB
	url         = "http://root@192.168.99.100:26257"
	tokenSchema = `CREATE TABLE IF NOT EXISTS tokens (
id varchar(255) primary key,
namespace varchar(255),
name varchar(255),
created integer,
updated integer,
unique (namespace, name));`
	q = map[string]string{
		"delete": "DELETE from explorer.tokens where id = $1",
		"create": `INSERT into explorer.tokens (
				id, namespace, name, created, updated) 
				values ($1, $2, $3, $4, $5)`,
		"update":                 "UPDATE explorer.tokens set namespace = $2, name = $2, updated = $3 where id = $1",
		"read":                   "SELECT * from explorer.tokens where id = $1",
		"list":                   "SELECT * from explorer.tokens limit $1 offset $2",
		"searchNamespace":        "SELECT * from explorer.tokens where namespace = $1 limit 1",
		"searchName":             "SELECT * from explorer.tokens where name = $1 limit 1",
		"searchNamespaceAndName": "SELECT * from explorer.tokens where namespace = $1 and name = $2 limit 1",
	}
	st = map[string]*sql.Stmt{}
)

func init() {
	var d *sql.DB
	var err error

	if d, err = sql.Open("cockroach", url); err != nil {
		log.Fatal(err)
	}
	if _, err := d.Exec("CREATE DATABASE explorer"); err != nil && err.Error() != `database "explorer" already exists` {
		log.Fatal(err)
	}
	d.Close()
	if d, err = sql.Open("cockroach", url+"?database=explorer"); err != nil {
		log.Fatal(err)
	}
	if _, err = d.Exec(tokenSchema); err != nil {
		log.Fatal(err)
	}
	db = d

	for query, statement := range q {
		prepared, err := db.Prepare(statement)
		if err != nil {
			log.Fatal(err)
		}
		st[query] = prepared
	}
}

func Create(token *token.Token) error {
	token.Created = time.Now().Unix()
	token.Updated = time.Now().Unix()
	_, err := st["create"].Exec(token.Id, token.Namespace, token.Name, token.Created, token.Updated)
	return err
}

func Delete(id string) error {
	_, err := st["delete"].Exec(id)
	return err
}

func Update(token *token.Token) error {
	token.Updated = time.Now().Unix()
	_, err := st["update"].Exec(token.Id, token.Namespace, token.Name, token.Updated)
	return err
}

func Read(id string) (*token.Token, error) {
	token := &token.Token{}

	r := st["read"].QueryRow(id)
	if err := r.Scan(&token.Id, &token.Namespace, &token.Name, &token.Created, &token.Updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return token, nil
}

func Search(namespace, name string, limit, offset int64) ([]*token.Token, error) {
	var r *sql.Rows
	var err error

	if len(namespace) > 0 && len(name) > 0 {
		r, err = st["searchNamespaceAndName"].Query(namespace, name, limit, offset)
	} else if len(namespace) > 0 {
		r, err = st["searchNamespace"].Query(namespace, limit, offset)
	} else if len(name) > 0 {
		r, err = st["searchName"].Query(name, limit, offset)
	} else {
		r, err = st["list"].Query(limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var tokens []*token.Token

	for r.Next() {
		token := &token.Token{}
		if err := r.Scan(&token.Id, &token.Namespace, &token.Name, &token.Created, &token.Updated); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		tokens = append(tokens, token)

	}
	if r.Err() != nil {
		return nil, err
	}

	return tokens, nil
}
