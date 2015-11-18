package db

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/cockroachdb/cockroach/sql/driver"
	log "github.com/golang/glog"
	srv "github.com/myodc/explorer-srv/proto/service"
)

var (
	db            *sql.DB
	url           = "http://root@192.168.99.100:26257"
	serviceSchema = `CREATE TABLE IF NOT EXISTS services (
id varchar(255) primary key,
name varchar(255),
owner varchar(255),
description varchar(255),
created integer,
updated integer,
url varchar(255),
readme text,
metadata text,
unique (name, owner))`
	versionSchema = `CREATE TABLE IF NOT EXISTS service_versions (
id varchar(255) primary key,
service_id varchar(255),
version varchar(255),
api text,
sources text,
dependencies text,
metadata text,
created integer,
updated integer,
unique (service_id, version))`
	q = map[string]string{
		"create":             "INSERT into explorer.services (id, name, owner, description, created, updated, url, readme, metadata) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		"delete":             "DELETE from explorer.services where id = $1",
		"update":             "UPDATE explorer.services set name = $2, owner = $3, description = $4, updated = $5, url = $6, readme = $7, metadata = $8 where id = $1",
		"read":               "SELECT * from explorer.services where id = $1",
		"listAsc":            "SELECT * from explorer.services order by id asc limit $1 offset $2",
		"listDesc":           "SELECT * from explorer.services order by id desc limit $1 offset $2",
		"searchName":         "SELECT * from explorer.services where name = $1",
		"searchOwner":        "SELECT * from explorer.services where owner = $1",
		"searchNameAndOwner": "SELECT * from explorer.services where name = $1 and owner = $2",

		// version
		"createVersion":  "INSERT into explorer.service_versions (id, service_id, version, api, sources, dependencies, metadata, created, updated) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		"deleteVersion":  "DELETE from explorer.service_versions where id = $1",
		"updateVersion":  "UPDATE explorer.service_versions set version = $2, api = $3, sources = $4, dependencies = $5, updated = $6, metadata = $7 where id = $1",
		"readVersion":    "SELECT * from explorer.service_versions where id = $1",
		"searchVersion":  "SELECT * from explorer.service_versions where service_id = $1 and version = $2",
		"searchVersions": "SELECT * from explorer.service_versions where service_id = $1",
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
	if _, err = d.Exec(versionSchema); err != nil {
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
	service.Readme = base64.StdEncoding.EncodeToString([]byte(service.Readme))
	_, err = st["create"].Exec(service.Id, service.Name, service.Owner, service.Description, service.Created, service.Updated, service.Url, service.Readme, string(md))
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
	service.Readme = base64.StdEncoding.EncodeToString([]byte(service.Readme))
	_, err = st["update"].Exec(service.Id, service.Name, service.Owner, service.Description, service.Updated, service.Url, service.Readme, string(md))
	return err
}

func Read(id string) (*srv.Service, error) {
	var md string
	service := &srv.Service{}

	r := st["read"].QueryRow(id)
	if err := r.Scan(&service.Id, &service.Name, &service.Owner, &service.Description, &service.Created, &service.Updated, &service.Url, &service.Readme, &md); err != nil {
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
	readme, _ := base64.StdEncoding.DecodeString(service.Readme)
	service.Readme = string(readme)
	return service, nil
}

func Search(name, owner, order string, limit, offset int64) ([]*srv.Service, error) {
	var r *sql.Rows
	var err error

	if len(name) > 0 && len(owner) > 0 {
		r, err = st["searchNameAndOwner"].Query(name, owner, limit, offset)
	} else if len(name) > 0 {
		r, err = st["searchName"].Query(name, limit, offset)
	} else if len(owner) > 0 {
		r, err = st["searchOwner"].Query(owner, limit, offset)
	} else {
		switch order {
		case "asc":
			r, err = st["listAsc"].Query(limit, offset)
		default:
			r, err = st["listDesc"].Query(limit, offset)
		}
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var services []*srv.Service

	for r.Next() {
		var md string
		service := &srv.Service{}
		if err := r.Scan(&service.Id, &service.Name, &service.Owner, &service.Description, &service.Created, &service.Updated, &service.Url, &service.Readme, &md); err != nil {
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
		readme, _ := base64.StdEncoding.DecodeString(service.Readme)
		service.Readme = string(readme)
		services = append(services, service)

	}
	if r.Err() != nil {
		return nil, err
	}

	return services, nil
}
