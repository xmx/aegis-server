package launch

import (
	"context"
	"os"
)

func Run(ctx context.Context, cfgFile string) error {
	file, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
