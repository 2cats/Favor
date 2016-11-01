package models

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Msg struct {
	Id         int64  `db:"id"`
	Content    string `db:"content"`
	CreateUser string `db:"create_user"`
	CreateTime int64  `db:"create_time"`
	UpdateTime int64  `db:"update_time"`
}

func UpdateMsg(id int64, content string) error {
	_, err := db.Exec(`UPDATE msg SET content=? ,update_time=? WHERE id=?`, content, time.Now().Unix(), id)
	return err
}
func DeleteMsg(id int64) error {
	_, err := db.Exec(`DELETE FROM msg WHERE id=?`, id)
	return err
}

func InsertMsg(content string, user string) error {
	now := time.Now().Unix()
	_, err := db.Exec(`INSERT INTO msg(content,create_user,create_time,update_time) VALUES(?,?,?,?)`,
		content, user, now, now,
	)
	return err
}
func SelectAllMsg() (msgs []Msg, err error) {
	records, err := db.Query(`SELECT * FROM msg ORDER BY update_time DESC`)
	if err != nil {
		return nil, err
	}
	defer records.Close()
	for records.Next() {
		var msg Msg
		records.Scan(&msg)
		msgs = append(msgs, msg)
	}
	err = records.Err()
	return
}
func SelectIdMsg(id int) (msg Msg, err error) {
	record := db.QueryRow(`SELECT id,content,create_user,create_time,update_time FROM msg ORDER BY update_time DESC LIMIT 1`)
	err = record.Scan(&msg.Id, &msg.Content, &msg.CreateUser, &msg.CreateTime, &msg.UpdateTime)
	return
}
func SelectNewestMsg() (msg Msg, err error) {
	record := db.QueryRow(`SELECT id,content,create_user,create_time,update_time FROM msg ORDER BY update_time DESC LIMIT 1`)
	if record == nil {
		err = errors.New("Nil Record")
		return
	}
	err = record.Scan(&msg.Id, &msg.Content, &msg.CreateUser, &msg.CreateTime, &msg.UpdateTime)
	return
}
func SelectPageMsg(filter string, pageSize int, nth int) (msgs []Msg, err error) {
	if nth <= 0 || pageSize <= 0 {
		err = errors.New("pageSize or nth error")
	}
	records, err := db.Query(`
		SELECT id,content,create_user,create_time,update_time FROM msg where content  like $1  ORDER BY update_time DESC LIMIT $2 OFFSET $3;
		`, "%%"+filter+"%%", pageSize, pageSize*(nth-1))
	if err != nil {
		return nil, err
	}
	defer records.Close()
	for records.Next() {
		var msg Msg
		records.Scan(&msg.Id, &msg.Content, &msg.CreateUser, &msg.CreateTime, &msg.UpdateTime)
		msgs = append(msgs, msg)
	}
	err = records.Err()
	return

}
func GetMsgCount(filter string) (count int, err error) {
	record := db.QueryRow(`SELECT COUNT(*) FROM msg where content  like ?;`, "%%"+filter+"%%")
	err = record.Scan(&count)
	return
}
func createSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS msg(
            Id INTEGER PRIMARY KEY ,
            content character varying,
            create_user character varying,
            create_time bigint,
	        update_time bigint
        );`,
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		panicWhenError(err)
	}
	return nil
}
func DBInit() {
	var err error
	db, err = sql.Open("sqlite3", Config.DBFile)
	panicWhenError(err)
	createSchema()
}
