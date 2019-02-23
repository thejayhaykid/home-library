package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func Connect(connection string) (*xorm.Engine, error) {
	return xorm.NewEngine("mysql", connection)
}
