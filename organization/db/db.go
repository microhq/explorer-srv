package db

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/golang/glog"
	org "github.com/micro/explorer-srv/proto/organization"
)

var (
	db                 *sql.DB
	url                = "explorer:explorer@tcp(127.0.0.1:3306)/"
	organizationSchema = `CREATE TABLE IF NOT EXISTS organizations (
id varchar(36) primary key,
name varchar(255),
email varchar(255),
owner varchar(255),
created integer,
updated integer,
unique (name));`
	membersSchema = `CREATE TABLE IF NOT EXISTS organization_members (
id varchar(36) primary key,
organization_id varchar(36) ,
username varchar(36),
roles text,
created integer,
updated integer,
unique (organization_id, username));`

	q = map[string]string{
		"delete": "DELETE from explorer.organizations where id = ?",
		"create": `INSERT into explorer.organizations (
				id, name, email, owner, created, updated) 
				values (?, ?, ?, ?, ?, ?)`,
		"update":             "UPDATE explorer.organizations set name = ?, email = ?, owner = ?, updated = ? where id = ?",
		"read":               "SELECT * from explorer.organizations where id = ?",
		"list":               "SELECT * from explorer.organizations limit ? offset ?",
		"searchName":         "SELECT * from explorer.organizations where name = ? limit ? offset ?",
		"searchOwner":        "SELECT * from explorer.organizations where owner = ? limit ? offset ?",
		"searchNameAndOwner": "SELECT * from explorer.organizations where name = ? and owner = ? limit ? offset ?",

		"deleteMember": "DELETE from explorer.organization_members where id = ?",
		"createMember": `INSERT into explorer.organization_members (
				id, organization_id, username, roles, created, updated) 
				values (?, ?, ?, ?, ?, ?)`,
		"updateMember":         "UPDATE explorer.organization_members set roles = ?, updated = ? where id = ?",
		"readMember":           "SELECT * from explorer.organization_members where id = ?",
		"searchUsername":       "SELECT * from explorer.organization_members where username = ? limit ? offset ?",
		"searchOrg":            "SELECT * from explorer.organization_members where organization_id = ? limit ? offset ?",
		"searchUsernameAndOrg": "SELECT * from explorer.organization_members where organization_id = ? and username = ? limit ? offset ?",
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
	if _, err = d.Exec(organizationSchema); err != nil {
		log.Fatal(err)
	}
	if _, err = d.Exec(membersSchema); err != nil {
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

func Create(organization *org.Organization) error {
	organization.Created = time.Now().Unix()
	organization.Updated = time.Now().Unix()
	_, err := st["create"].Exec(organization.Id, organization.Name, organization.Email, organization.Owner, organization.Created, organization.Updated)
	return err
}

func Delete(id string) error {
	_, err := st["delete"].Exec(id)
	return err
}

func Update(organization *org.Organization) error {
	organization.Updated = time.Now().Unix()
	_, err := st["update"].Exec(organization.Name, organization.Email, organization.Owner, organization.Updated, organization.Id)
	return err
}

func Read(id string) (*org.Organization, error) {
	organization := &org.Organization{}

	r := st["read"].QueryRow(id)
	if err := r.Scan(&organization.Id, &organization.Name, &organization.Email, &organization.Owner, &organization.Created, &organization.Updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return organization, nil
}

func Search(name, owner string, limit, offset int64) ([]*org.Organization, error) {
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

	var organizations []*org.Organization

	for r.Next() {
		organization := &org.Organization{}
		if err := r.Scan(&organization.Id, &organization.Name, &organization.Email, &organization.Owner, &organization.Created, &organization.Updated); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		organizations = append(organizations, organization)

	}
	if r.Err() != nil {
		return nil, err
	}

	return organizations, nil
}
