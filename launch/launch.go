package launch

import (
	"context"
	"fmt"
	"os"

	"github.com/xmx/aegis-server/business/service"
)

func Run(ctx context.Context, cfgFile string) error {
	file, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	defer file.Close()

	initialService := service.Initial()

	initCfg, err := initialService.Wait(ctx)
	if err != nil {
		return err
	}
	fmt.Println(initCfg)

	return nil
}
