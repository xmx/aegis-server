package main

import (
	"log"
	"time"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/library/sqldb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "2LqbbPiusZmES3w.root:geDbAqK2ltgwuxbQ@tcp(gateway01.ap-southeast-1.prod.aws.tidbcloud.com:4000)/test?tls=tidb"
	conn, err := sqldb.TiDB(dsn, time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn}))
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Certificate{})
	log.Println(err)
}
