package db

import (
	"errors"
	"database/sql"
	"time"

	log "github.com/golang/glog"
	_ "github.com/cockroachdb/cockroach/sql/driver"
	srv "github.com/myodc/explorer-srv/proto/profile"
)

var (
	db *sql.DB
	url = "http://root@192.168.99.100:26257"
	profileSchema = `CREATE TABLE IF NOT EXISTS profiles (
id varchar(255) primary key,
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
	q = map[string]string {
		"delete": "DELETE from explorer.profiles where id = $1",
		"create": `INSERT into explorer.profiles (
				id, name, owner, type, display_name, blurb, url, location, created, updated) 
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		"update": "UPDATE explorer.profiles set name = $2, owner = $3, type = $4, display_name = $5, blurb = $6, url = $7, location = $8, updated = $9 where id = $1",
		"read": "SELECT * from explorer.profiles where id = $1",
		"list": "SELECT * from explorer.profiles limit $1 offset $2",
		"searchName": "SELECT * from explorer.profiles where name = $1",
		"searchOwner": "SELECT * from explorer.profiles where owner = $1",
		"searchNameAndOwner": "SELECT * from explorer.profiles where name = $1 and owner = $2",

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
	_, err := st["update"].Exec(profile.Id, profile.Name, profile.Owner, profile.Type, profile.DisplayName,
					profile.Blurb, profile.Url, profile.Location, profile.Updated)
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
