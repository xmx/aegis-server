package launch

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go/http3"
	"os"

	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/memconf"
	"gorm.io/gorm"
)

func Run(ctx context.Context, db *gorm.DB, cfgFile string) error {
	file, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	defer file.Close()

	{
		tables := []any{
			model.ConfigCertificate{},
			model.ConfigLogger{},
			model.ConfigServer{},
		}
		if err = db.AutoMigrate(tables...); err != nil {
			return err
		}
	}
	qry := query.Use(db)

	configCertificateRepository := repository.ConfigCertificate(qry)
	configCertificateConfigurer := memconf.ConfigCertificate(configCertificateRepository)

	srv := &http3.Server{
		TLSConfig: &tls.Config{GetConfigForClient: configCertificateConfigurer.Certificate}
	}

	return nil
}
