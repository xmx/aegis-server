package main

import (
	"flag"
	"fmt"
	"github.com/xmx/aegis-server/library/credential"
	"github.com/xmx/aegis-server/quicsrv"
	"os"
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

	fmt.Println(*cfg)

	srv := quicsrv.Server{
		Cert: credential.Atomic(),
	}

	srv.ListenAndServe()
}
