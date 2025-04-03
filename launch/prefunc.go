package launch

import (
	"fmt"
	"log/slog"

	"github.com/xmx/aegis-server/business/customvalid"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/validation"
)

func registerValidator(valid *validation.Validate, repo repository.All, log *slog.Logger) {
	validDB := customvalid.NewValidDB(repo, log)
	err := valid.RegisterCustomValidation(validDB.Password)
	fmt.Println(err)
}
