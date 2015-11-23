package db

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/golang/glog"
	srv "github.com/micro/explorer-srv/proto/profile"
)

var (
	db            *sql.DB
	url           = "explorer:explorer@tcp(127.0.0.1:3306)/"
	profileSchema = `CREATE TABLE IF NOT EXISTS profiles (
id varchar(36) primary key,
name varchar(255),
owner varchar(255),
type integer,
display_name varchar(255),
blurb varchar(255),
url varchar(255),
location varchar(255),
created integer,
updated integer,
unique (name));`
	q = map[string]string{
		"delete": "DELETE from explorer.profiles where id = ?",
		"create": `INSERT into explorer.profiles (
				id, name, owner, type, display_name, blurb, url, location, created, updated) 
				values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"update":             "UPDATE explorer.profiles set name = ?, owner = ?, type = ?, display_name = ?, blurb = ?, url = ?, location = ?, updated = ? where id = ?",
		"read":               "SELECT * from explorer.profiles where id = ?",
		"list":               "SELECT * from explorer.profiles limit ? offset ?",
		"searchName":         "SELECT * from explorer.profiles where name = ? limit ? offset ?",
		"searchOwner":        "SELECT * from explorer.profiles where owner = ? limit ? offset ?",
		"searchNameAndOwner": "SELECT * from explorer.profiles where name = ? and owner = ? limit ? offset ?",
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
	if _, err = d.Exec(profileSchema); err != nil {
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

func Create(profile *srv.Profile) error {
	profile.Created = time.Now().Unix()
	profile.Updated = time.Now().Unix()
	_, err := st["create"].Exec(profile.Id, profile.Name, profile.Owner, profile.Type, profile.DisplayName,
		profile.Blurb, profile.Url, profile.Location, profile.Created, profile.Updated)
	return err
}

func Delete(id string) error {
	_, err := st["delete"].Exec(id)
	return err
}

func Update(profile *srv.Profile) error {
	profile.Updated = time.Now().Unix()
	_, err := st["update"].Exec(profile.Name, profile.Owner, profile.Type, profile.DisplayName,
		profile.Blurb, profile.Url, profile.Location, profile.Updated, profile.Id)
	return err
}

func Read(id string) (*srv.Profile, error) {
	profile := &srv.Profile{}

	r := st["read"].QueryRow(id)
	if err := r.Scan(&profile.Id, &profile.Name, &profile.Owner, &profile.Type, &profile.DisplayName,
		&profile.Blurb, &profile.Url, &profile.Location,
		&profile.Created, &profile.Updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return profile, nil
}

func Search(name, owner string, limit, offset int64) ([]*srv.Profile, error) {
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

	var profiles []*srv.Profile

	for r.Next() {
		profile := &srv.Profile{}
		if err := r.Scan(&profile.Id, &profile.Name, &profile.Owner, &profile.Type, &profile.DisplayName,
			&profile.Blurb, &profile.Url, &profile.Location,
			&profile.Created, &profile.Updated); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		profiles = append(profiles, profile)

	}
	if r.Err() != nil {
		return nil, err
	}

	return profiles, nil
}
