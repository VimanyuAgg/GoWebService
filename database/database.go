package database

import (
	"database/sql"
	"fmt"
	"github.com/goWebServices/one-o-one/config"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

var DbConn *sql.DB

func SetupDatabase() {
	var err error
	DbConn, err = sql.Open("mysql",fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/inventorydb",config.LOCAL_MYSQL_USERNAME, config.LOCAL_MYSQL_PASSWORD))
	if err != nil {
		log.Fatal("Error while making DBConnection...", err)
	}

}
