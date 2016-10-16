package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mgutz/dat/v1"
	"github.com/mgutz/dat/v1/sqlx-runner"
)

var Conn *runner.Connection

type Msg struct {
	Id         int64  `db:"id"`
	Content    string `db:"content"`
	CreateUser string `db:"create_user"`
	CreateTime int64  `db:"create_time"`
	UpdateTime int64  `db:"update_time"`
}

func UpdateMsg(id int64, content string) error {
	var msg = Msg{}
	msg.Content = content
	msg.UpdateTime = time.Now().Unix()
	_, err := Conn.Update("msg").SetWhitelist(msg, "content", "update_time").Where("id = $1", id).Exec()
	return err
}
func DeleteMsg(id int64) error {
	_, err := Conn.DeleteFrom("msg").Where("id = $1", id).Exec()
	return err
}

func InsertMsg(content string, user string) error {
	var msg = Msg{
		Content:    content,
		CreateUser: user,
	}
	msg.CreateTime = time.Now().Unix()
	msg.UpdateTime = msg.CreateTime
	_, err := Conn.InsertInto("msg").Record(msg).Blacklist("id").Exec()
	return err
}
func SelectAllMsg() (msgs []Msg, err error) {
	err = Conn.SQL(`SELECT * FROM msg ORDER BY update_time DESC`).QueryStructs(&msgs)
	return
}
func SelectIdMsg(id int) (msg Msg, err error) {
	err = Conn.SQL(`SELECT * FROM msg WHERE id = $1`, id).QueryStruct(&msg)
	return
}
func SelectNewestMsg() (msg Msg, err error) {
	err = Conn.SQL(`SELECT * FROM msg ORDER BY update_time DESC LIMIT 1`).QueryStruct(&msg)
	return
}
func SelectPageMsg(filter string, pageSize int, nth int) (msgs []Msg, err error) {
	if nth <= 0 || pageSize <= 0 {
		err = errors.New("pageSize or nth error")
	}

	err = Conn.SQL(`
		SELECT * FROM msg where content  like $1  ORDER BY update_time DESC LIMIT $2 OFFSET $3;
		`, "%%"+filter+"%%", pageSize, pageSize*(nth-1)).QueryStructs(&msgs)
	return
}
func GetMsgCount(filter string) (count int, err error) {
	err = Conn.SQL(`
		SELECT COUNT(*) FROM msg where content  like $1;
		`, "%%"+filter+"%%").QueryStruct(&count)
	return
}
func createSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS public.msg(
            Id SERIAL PRIMARY KEY,
            content character varying,
            create_user character varying,
            create_time bigint,
	        update_time bigint
        );`,
	}
	for _, q := range queries {
		fmt.Println(q)
		_, err := Conn.SQL(q).Exec()
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
func DBInit() {

	constr := "dbname=msg user=postgres password=amber host=2cats.xyz sslmode=disable"
	fmt.Println("DB CONF:", constr)
	db, err := sql.Open("postgres", constr)

	if err != nil {
		panic(err)
	}

	// set to reasonable values for production
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(16)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	//runner.LogQueriesThreshold = 10 * time.Millisecond

	Conn = runner.NewConnection(db, "postgres")
	createSchema()
}
