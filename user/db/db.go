package db

import (
	"errors"
	"database/sql"
	"time"

	log "github.com/golang/glog"
	_ "github.com/cockroachdb/cockroach/sql/driver"
	user "github.com/myodc/explorer-srv/proto/user"
)

var (
	db *sql.DB
	url = "http://root@192.168.99.100:26257"
	userSchema = `CREATE TABLE IF NOT EXISTS users (
id varchar(255) primary key,
username varchar(255),
email varchar(255),
salt varchar(16),
password text,
created integer,
updated integer,
unique (username),
unique (email));`
	sessionSchema = `CREATE TABLE IF NOT EXISTS sessions (
id varchar(255) primary key,
username varchar(255),
created integer,
expires integer);`
	q = map[string]string {
		"delete": "DELETE from explorer.users where id = $1",
		"create": `INSERT into explorer.users (
				id, username, email, salt, password, created, updated) 
				values ($1, $2, $3, $4, $5, $6, $7)`,
		"update": "UPDATE explorer.users set username = $2, email = $2, updated = $3 where id = $1",
		"read": "SELECT * from explorer.users where id = $1",
		"list": "SELECT * from explorer.users limit $1 offset $2",
		"searchUsername": "SELECT * from explorer.users where username = $1 limit 1",
		"searchEmail": "SELECT * from explorer.users where email = $1 limit 1",
		"searchUsernameAndEmail": "SELECT * from explorer.users where username = $1 and email = $2 limit 1",

		// users.sessions
		"createSession": "INSERT into explorer.sessions (id, username, created, expires) values ($1, $2, $3, $4)",
		"deleteSession": "DELETE from explorer.sessions where id = $1",
		"readSession": "SELECT * from explorer.sessions where id = $1",
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
        if _, err = d.Exec(userSchema); err != nil {
                log.Fatal(err)
        }
        if _, err = d.Exec(sessionSchema); err != nil {
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

func CreateSession(sess *user.Session) error {
	if sess.Created == 0 {
		sess.Created = time.Now().Unix()
	}

	if sess.Expires == 0 {
		sess.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	_, err := st["createSession"].Exec(sess.Id, sess.Username, sess.Created, sess.Expires)
	return err
}

func DeleteSession(id string) error {
	_, err := st["deleteSession"].Exec(id)
	return err
}

func ReadSession(id string) (*user.Session, error) {
	sess := &user.Session{}

	r := st["readSession"].QueryRow(id)
	if err := r.Scan(&sess.Id, &sess.Username, &sess.Created, &sess.Expires); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return sess, nil
}

func Create(user *user.User, salt string, password string) error {
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()
	_, err := st["create"].Exec(user.Id, user.Username, user.Email, salt, password, user.Created, user.Updated)
	return err
}

func Delete(id string) error {
	_, err := st["delete"].Exec(id)
	return err
}

func Update(user *user.User) error {
	user.Updated = time.Now().Unix()
	_, err := st["update"].Exec(user.Id, user.Username, user.Email, user.Updated)
	return err
}

func Read(id string) (*user.User, error) {
	user := &user.User{}

	r := st["read"].QueryRow(id)
	if err := r.Scan(&user.Id, &user.Username, &user.Email, nil, nil, &user.Created, &user.Updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return user, nil
}

func Search(username, email string, limit, offset int64) ([]*user.User, error) {
	var r *sql.Rows
	var err error

	if len(username) > 0 && len(email) > 0 {
		r, err = st["searchUsernameAndEmail"].Query(username, email, limit, offset)
	} else if len(username) > 0 {
		r, err = st["searchUsername"].Query(username, limit, offset)
	} else if len(email) > 0 {
		r, err = st["searchEmail"].Query(email, limit, offset)
	} else {
		r, err = st["list"].Query(limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var users []*user.User

	for r.Next() {
		user := &user.User{}
		if err := r.Scan(&user.Id, &user.Username, &user.Email, nil, nil, &user.Created, &user.Updated); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		users = append(users, user)

	}
	if r.Err() != nil {
		return nil, err
	}

	return users, nil
}

func SaltAndPassword(username, email string) (string, string, error) {
	var r *sql.Rows
	var err error

	if len(username) > 0 && len(email) > 0 {
		r, err = st["searchUsernameAndEmail"].Query(username, email, 1, 0)
	} else if len(username) > 0 {
		r, err = st["searchUsername"].Query(username, 1, 0)
	} else if len(email) > 0 {
		r, err = st["searchEmail"].Query(email, 1, 0)
	} else {
		return "", "", errors.New("username and email cannot be blank")
	}

	if err != nil {
		return "", "", err
	}
	defer r.Close()

	if !r.Next() {
		return "", "", errors.New("not found")
	}

	var salt, pass string
	user := &user.User{}
	if err := r.Scan(&user.Id, &user.Username, &user.Email, &salt, &pass, &user.Created, &user.Updated); err != nil {
		if err == sql.ErrNoRows {
			return "", "", errors.New("not found")
		}
		return "", "", err
	}
	if r.Err() != nil {
		return "", "", err
	}

	return salt, pass, nil
}
