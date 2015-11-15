package db

import (
	"errors"
	"encoding/json"
	"database/sql"
	"time"

	log "github.com/golang/glog"
	_ "github.com/cockroachdb/cockroach/sql/driver"
	srv "github.com/myodc/explorer-srv/proto/service"
)

var (
	db *sql.DB
	url = "http://root@192.168.99.100:26257"
	serviceSchema = `CREATE TABLE IF NOT EXISTS services (
id varchar(255) primary key,
name varchar(255),
owner varchar(255),
description varchar(255),
created integer,
updated integer,
metadata text,
unique (name, owner))`
	q = map[string]string {
		"create": "INSERT into explorer.services (id, name, owner, description, created, updated, metadata) values ($1, $2, $3, $4, $5, $6, $7)",
		"delete": "DELETE from explorer.services where id = $1",
		"update": "UPDATE explorer.services set name = $2, owner = $3, description = $4, updated = $5, metadata = $6 where id = $1",
		"read": "SELECT * from explorer.services where id = $1",
		"list": "SELECT * from explorer.services limit $1 offset $2",
                "searchName": "SELECT * from explorer.services where name = $1",
                "searchOwner": "SELECT * from explorer.services where owner = $1",
                "searchNameAndOwner": "SELECT * from explorer.services where name = $1 and owner = $2",
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
        if _, err = d.Exec(serviceSchema); err != nil {
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

func Create(service *srv.Service) error {
	md, err := json.Marshal(service.Metadata)
	if err != nil {
		return err
	}
	service.Created = time.Now().Unix()
	service.Updated = time.Now().Unix()
	_, err = st["create"].Exec(service.Id, service.Name, service.Owner, service.Description, service.Created, service.Updated, string(md))
	return err
}

func Delete(id string) error {
	_, err := st["delete"].Exec(id)
	return err
}

func Update(service *srv.Service) error {
	md, err := json.Marshal(service.Metadata)
	if err != nil {
		return err
	}
	service.Updated = time.Now().Unix()
	_, err = st["update"].Exec(service.Id, service.Name, service.Owner, service.Description, service.Updated, string(md))
	return err
}

func Read(id string) (*srv.Service, error) {
	var md string
	service := &srv.Service{}

	r := st["read"].QueryRow(id)
	if err := r.Scan(&service.Id, &service.Name, &service.Owner, &service.Description, &service.Created, &service.Updated, &md); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	var mdu map[string]string

	if err := json.Unmarshal([]byte(md), &mdu); err != nil {
		return nil, err
	}

	service.Metadata = mdu
	return service, nil
}

func Search(name, owner string, limit, offset int64) ([]*srv.Service, error) {
	var r *sql.Rows
	var err error

	if len(name) > 0 && len(owner) > 0 {
		r, err = st["searchNameAndOwner"].Query(name, owner, limit, offset)
	} else if len(name) > 0 {
		r, err = st["searchName"].Query(name, limit, offset)
	} else if len(owner) > 0 {
		r, err = st["searchOwner"].Query(owner, limit, offset)
	} else {
		r, err = st["list"].Query(limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var services []*srv.Service

	for r.Next() {
		var md string
		service := &srv.Service{}
		if err := r.Scan(&service.Id, &service.Name, &service.Owner, &service.Description, &service.Created, &service.Updated, &md); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		var mdu map[string]string
		if err := json.Unmarshal([]byte(md), &mdu); err != nil {
			return nil, err
		}
		service.Metadata = mdu
		services = append(services, service)

	}
	if r.Err() != nil {
		return nil, err
	}

	return services, nil
}
