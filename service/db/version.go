package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	srv "github.com/micro/explorer-srv/proto/service"
)

func CreateVersion(v *srv.Version) error {
	api, err := json.Marshal(v.Api)
	if err != nil {
		return err
	}
	src, err := json.Marshal(v.Sources)
	if err != nil {
		return err
	}
	dep, err := json.Marshal(v.Dependencies)
	if err != nil {
		return err
	}
	md, err := json.Marshal(v.Metadata)
	if err != nil {
		return err
	}
	v.Created = time.Now().Unix()
	v.Updated = time.Now().Unix()
	_, err = st["createVersion"].Exec(v.Id, v.ServiceId, v.Version, string(api), string(src), string(dep), string(md), v.Created, v.Updated, v.Private)
	return err
}

func DeleteVersion(id string) error {
	_, err := st["deleteVersion"].Exec(id)
	return err
}

func UpdateVersion(v *srv.Version) error {
	api, err := json.Marshal(v.Api)
	if err != nil {
		return err
	}
	src, err := json.Marshal(v.Sources)
	if err != nil {
		return err
	}
	dep, err := json.Marshal(v.Dependencies)
	if err != nil {
		return err
	}
	md, err := json.Marshal(v.Metadata)
	if err != nil {
		return err
	}
	v.Updated = time.Now().Unix()
	_, err = st["updateVersion"].Exec(v.Version, string(api), string(src), string(dep), string(md), v.Updated, v.Private, v.Id)
	return err
}

func ReadVersion(id string) (*srv.Version, error) {
	var md, src, dep, api string
	v := &srv.Version{}

	r := st["readVersion"].QueryRow(id)
	if err := r.Scan(&v.Id, &v.ServiceId, &v.Version, &api, &src, &dep, &md, &v.Created, &v.Updated, &v.Private); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(md), &v.Metadata); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(api), &v.Api); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(src), &v.Sources); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(dep), &v.Dependencies); err != nil {
		return nil, err
	}

	return v, nil
}

func SearchVersion(serviceId, version string, limit, offset int64) ([]*srv.Version, error) {
	var r *sql.Rows
	var err error

	if len(serviceId) > 0 && len(version) > 0 {
		r, err = st["searchVersion"].Query(serviceId, version, limit, offset)
	} else if len(serviceId) > 0 {
		r, err = st["searchVersions"].Query(serviceId, limit, offset)
	} else {
		return nil, errors.New("id and version blank")
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var versions []*srv.Version

	for r.Next() {
		var md, src, dep, api string
		v := &srv.Version{}
		if err := r.Scan(&v.Id, &v.ServiceId, &v.Version, &api, &src, &dep, &md, &v.Created, &v.Updated, &v.Private); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}

		if err := json.Unmarshal([]byte(md), &v.Metadata); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(api), &v.Api); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(src), &v.Sources); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(dep), &v.Dependencies); err != nil {
			return nil, err
		}

		versions = append(versions, v)

	}
	if r.Err() != nil {
		return nil, err
	}

	return versions, nil
}
