package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/library/sqldb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	args := os.Args
	set := flag.NewFlagSet(args[0], flag.ExitOnError)
	ver := set.Bool("v", false, "打印版本")
	cfg := set.String("c", "resources/config/application.json", "配置文件")
	_ = set.Parse(args[1:])
	if *ver {
		return
	}

	//dsn := "2LqbbPiusZmES3w.root:DuX1CBVcI93FnF9u@tcp(gateway01.ap-southeast-1.prod.aws.tidbcloud.com:4000)/test?tls=tidb"
	//conn, err := sqldb.TiDB(dsn, time.Minute)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer conn.Close()
	//db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn}))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = db.AutoMigrate(&model.Certificate{})
	//log.Println(err)
}
