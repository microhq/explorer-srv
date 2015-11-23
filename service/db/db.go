package db

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/golang/glog"
	srv "github.com/micro/explorer-srv/proto/service"
)

var (
	db            *sql.DB
	url           = "explorer:explorer@tcp(127.0.0.1:3306)/"
	serviceSchema = `CREATE TABLE IF NOT EXISTS services (
id varchar(36) primary key,
name varchar(255),
owner varchar(255),
description varchar(255),
created integer,
updated integer,
url varchar(255),
readme text,
metadata text,
private boolean,
unique (name, owner))`
	versionSchema = `CREATE TABLE IF NOT EXISTS service_versions (
id varchar(36) primary key,
service_id varchar(36),
version varchar(255),
api text,
sources text,
dependencies text,
metadata text,
created integer,
updated integer,
private boolean,
unique (service_id, version))`
	q = map[string]string{
		"create":             "INSERT into explorer.services (id, name, owner, description, created, updated, url, readme, metadata, private) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"delete":             "DELETE from explorer.services where id = ?",
		"update":             "UPDATE explorer.services set name = ?, owner = ?, description = ?, updated = ?, url = ?, readme = ?, metadata = ?, private = ? where id = ?",
		"read":               "SELECT * from explorer.services where id = ?",
		"listAsc":            "SELECT * from explorer.services order by id asc limit ? offset ?",
		"listDesc":           "SELECT * from explorer.services order by id desc limit ? offset ?",
		"searchName":         "SELECT * from explorer.services where name = ? limit ? offset ?",
		"searchOwner":        "SELECT * from explorer.services where owner = ? limit ? offset ?",
		"searchNameAndOwner": "SELECT * from explorer.services where name = ? and owner = ? limit ? offset ?",

		// version
		"createVersion":  "INSERT into explorer.service_versions (id, service_id, version, api, sources, dependencies, metadata, created, updated, private) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"deleteVersion":  "DELETE from explorer.service_versions where id = ?",
		"updateVersion":  "UPDATE explorer.service_versions set version = ?, api = ?, sources = ?, dependencies = ?, metadata = ?, updated = ?, private = ? where id = ?",
		"readVersion":    "SELECT * from explorer.service_versions where id = ?",
		"searchVersion":  "SELECT * from explorer.service_versions where service_id = ? and version = ? limit ? offset ?",
		"searchVersions": "SELECT * from explorer.service_versions where service_id = ? limit ? offset ?",
	}
	st = map[string]*sql.Stmt{}
)

func init() {
	var d *sql.DB
	var err error

	if d, err = sql.Open("mysql", url); err != nil {
		log.Fatal(err)
	}
	if _, err := d.Exec("CREATE DATABASE IF NOT EXISTS explorer"); err != nil {
		log.Fatal(err)
	}
	d.Close()
	if d, err = sql.Open("mysql", url+"explorer"); err != nil {
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
	_, err = st["create"].Exec(service.Id, service.Name, service.Owner, service.Description, service.Created, service.Updated, service.Url, service.Readme, string(md), service.Private)
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
	_, err = st["update"].Exec(service.Name, service.Owner, service.Description, service.Updated, service.Url, service.Readme, string(md), service.Private, service.Id)
	return err
}

func Read(id string) (*srv.Service, error) {
	var md string
	service := &srv.Service{}

	r := st["read"].QueryRow(id)
	if err := r.Scan(&service.Id, &service.Name, &service.Owner, &service.Description, &service.Created, &service.Updated, &service.Url, &service.Readme, &md, &service.Private); err != nil {
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
		if err := r.Scan(&service.Id, &service.Name, &service.Owner, &service.Description, &service.Created, &service.Updated, &service.Url, &service.Readme, &md, &service.Private); err != nil {
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
