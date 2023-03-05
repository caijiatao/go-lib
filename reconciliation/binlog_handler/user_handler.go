package binlog_handler

import (
	"database/sql"

	"github.com/siddontang/go-log/log"
)

type User struct {
	Id         uint64
	UserInfo   string
	CreateTime int64
	UpdateTime int64
}

func (u *User) UserTableName() string {
	return "user"
}

type UserDataHandler struct {
	db *sql.DB
}

func (u *UserDataHandler) parseUserRows(rows [][]interface{}) []User {
	// 将interface 转化为具体的数据
	users := make([]User, 0, len(rows))
	return users
}

// 判断数据是否存在
func (u *UserDataHandler) isDataExists(user User) bool {
	return false
}

func (u *UserDataHandler) createData(user User) error {
	return nil
}

// 更新数据，通过updateTime来防止旧的更新覆盖新的更新
func (u *UserDataHandler) updateData(user User) error {
	return nil
}

func (u *UserDataHandler) deleteData(user User) error {
	return nil
}

func (u UserDataHandler) SyncBinLogData(action string, rows [][]interface{}, table string) error {
	var (
		err error
	)
	users := u.parseUserRows(rows)
	for _, user := range users {
		switch action {
		case "insert":
			if u.isDataExists(user) {
				err = u.updateData(user)
			} else {
				err = u.createData(user)
			}
		case "update":
			if u.isDataExists(user) {
				err = u.updateData(user)
			}
		case "delete":
			if u.isDataExists(user) {
				err = u.deleteData(user)
			}
		}
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func NewUserDataHandler() *UserDataHandler {
	return &UserDataHandler{}
}
