package executor

import (
	"testing"
	"strconv"
	"os"
	"fmt"
)

var d DB

func TestNewDB(t *testing.T) {
	port,_ := strconv.ParseInt(os.Getenv("MYSQL_PORT"),10,32)
	config := DBConfig{
		Host: os.Getenv("MYSQL_HOST"),
		Port: int(port),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DBName: os.Getenv("MYSQL_DB"),
		User: os.Getenv("MYSQL_USER"),
	}
	fmt.Println(config)

	db,err := NewDB(config)
	d = db
	if err != nil {
		t.Error(err)
	}
}

func TestDb_Migrate(t *testing.T) {
	d.Migrate()
}