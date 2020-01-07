package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"tio-model/database/model"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type TDB_Postgres struct {
	db *sql.DB
}

func (p *TDB_Postgres) Init(addr string) error {
	//addr := os.Getenv("TIO_DB_POSTGRES_CONN")
	logrus.Debugf("Postgres Connstr: %s", addr)

	db, err := sql.Open("postgres", addr)
	if err != nil {
		logrus.Fatalf("Connect Postgres Error: %s", err.Error())
	}

	p.db = db

	if err := p.db.Ping(); err != nil {
		logrus.Fatalf("Ping Postgres Error: %s", err.Error())
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(30)
	db.SetConnMaxLifetime(5 * time.Minute)

	logrus.Info("Connect Postgres Success.")

	return nil
}

func (p *TDB_Postgres) Version() string {
	return "PosgresQL"
}

func (p *TDB_Postgres) SaveTioUser(user *model.User) error {
	sql := "INSERT INTO tio_user(name, passwd) VALUES ($1, $2)"
	logrus.Debugf("Save New User: [%s]", sql)

	_, err := p.db.Exec(sql, user.Name, user.Passwd)
	return err
}

func (p *TDB_Postgres) QueryTioUser(name string) (model.User, error) {

	u := model.User{}

	sql := "SELECT * FROM tio_user WHERE name=$1"
	logrus.Debugf("Query User: [%s]", sql)
	rows, err := p.db.Query(sql, name)
	if err != nil {
		return u, err
	}

	if rows.Next() {
		err = rows.Scan(&u.Id, &u.Name, &u.Passwd)
	} else {
		return u, errors.New("No Match User Record")
	}

	return u, nil
}
func (p *TDB_Postgres) UpdateTioUser(user *model.User) error {
	sql := "UPDATE tio_user SET passwd=$2 WHERE name=$1"
	logrus.Debugf("Update User: [%s]", sql)

	_, err := p.db.Exec(sql, user.Name, user.Passwd)

	return err
}

func (p *TDB_Postgres) DeleteTioUser(name string) error {
	sql := "DELETE tio_user WHERE name=$1"
	logrus.Debugf("Delete User: [%s]", sql)

	_, err := p.db.Exec(sql, name)

	return err
}

func (p *TDB_Postgres) SaveTioServer(s *model.Server) error {
	sql := "INSERT INTO server (name, version, uid, stype, domain, path, tversion, timestamp, status, image, raw) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	logrus.Debugf("Save New Server:[%s]", sql)

	_, err := p.db.Exec(sql, s.Name, s.Version, s.Uid, s.Stype, s.Domain, s.Path, s.TVersion, s.Timestamp, s.Status, s.Image, s.Raw)
	return err
}

func (p *TDB_Postgres) queryTioServerWithQuery(query string) ([]model.Server, error) {
	logrus.Debugf("Query Server:[%s]", query)

	var ss []model.Server

	rows, err := p.db.Query(query)
	if err != nil {
		return ss, err
	}

	for rows.Next() {
		s := model.Server{}
		err = rows.Scan(&s.Id, &s.Name, &s.Version, &s.Uid, &s.Stype, &s.Domain, &s.Path, &s.TVersion, &s.Timestamp, &s.Status, &s.Image, &s.Raw)
		if err != nil {
			logrus.Errorf("Scan Server Error. %s", err)
			continue
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (p *TDB_Postgres) QueryTioServer() ([]model.Server, error) {
	return p.queryTioServerWithQuery("SELECT * FROM server ORDER BY version desc, id asc")
}

func (p *TDB_Postgres) QueryTioServerByUser(uid, limit int, name string) ([]model.Server, error) {

	var sql string
	if limit > 0 {
		sql = fmt.Sprintf("SELECT * FROM server WHERE uid=$1 AND name=$2 ORDER BY version desc, id asc LIMIT %d", limit)
	} else {
		sql = fmt.Sprintf("SELECT * FROM server WHERE uid=$1 AND name=$2 ORDER BY version desc, id asc")
	}

	return p.queryTioServerWithQuery(sql)
}

func (p *TDB_Postgres) QueryTioServerById(sid int) (*model.Server, error) {
	sql := fmt.Sprintf("SELECT * FROM server WHERE id=$1")
	logrus.Debugf("Query Server:[%s] id: [%d]", sql, sid)
	s := model.Server{}

	rows, err := p.db.Query(sql, sid)
	if err != nil {
		return &s, err
	}

	for rows.Next() {
		err = rows.Scan(&s.Id, &s.Name, &s.Version, &s.Uid, &s.Stype, &s.Domain, &s.Path, &s.TVersion, &s.Timestamp, &s.Status, &s.Image, &s.Raw)
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func (p *TDB_Postgres) QueryTioServerByName(name string) (*model.Server, error) {
	sql := fmt.Sprintf("SELECT * FROM server WHERE name=$1 order by id desc")
	logrus.Debugf("Query Server:[%s]", sql)
	s := model.Server{}

	rows, err := p.db.Query(sql, name)
	if err != nil {
		return &s, err
	}

	for rows.Next() {
		err = rows.Scan(&s.Id, &s.Name, &s.Version, &s.Uid, &s.Stype, &s.Domain, &s.Path, &s.TVersion, &s.Timestamp, &s.Status, &s.Image, &s.Raw)
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func (p *TDB_Postgres) UpdateTioServer(s *model.Server) error {
	sql := "UPDATE server SET name=$2, version=$3, stype=$4, domain=$5, path=$6, tversion=$7, timestamp=$8, status=$9,image=$10, raw=$11 WHERE id=$1"
	logrus.Debugf("Update Server: [%s] [%v]", sql, s)

	_, err := p.db.Exec(sql, s.Id, s.Name, s.Version, s.Stype, s.Domain, s.Path, s.TVersion, s.Timestamp, s.Status, s.Image, s.Raw)

	return err
}

func (p *TDB_Postgres) DeleteTioServer(name string) error {
	sql := "DELETE server WHERE name=$1"
	logrus.Debugf("Delete server: [%s]", sql)

	_, err := p.db.Exec(sql, name)

	return err

}

//func (p *TDB_Postgres) QueryUserServer(uid int, name string) ([]model.Server, error) {
//	sql := fmt.Sprintf("SELECT * FROM server WHERE uid=$1 AND name=$2 order by version desc, id asc")
//	logrus.Debugf("Query User [%d] Serverless [%s] SQL [%s]", uid, name, sql)
//
//	var ss []model.Server
//
//	rows, err := p.db.Query(sql, uid)
//	if err != nil {
//		return ss, err
//	}
//
//	for rows.Next() {
//		s := model.Server{}
//		err = rows.Scan(&s.Id, &s.Name, &s.Version, &s.Uid, &s.Stype, &s.Domain, &s.Path, &s.TVersion, &s.Timestamp, &s.Status, &s.Image, &s.Raw)
//		if err != nil {
//			logrus.Errorf("Scan Server Error. %s", err)
//			continue
//		}
//		ss = append(ss, s)
//	}
//
//	return ss, nil
//}
