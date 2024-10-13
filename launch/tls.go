package launch

import (
	"context"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
	"os"
)

func saveTLS(ctx context.Context, svc service.ConfigCertificate) {
	key, _ := os.ReadFile("resources/temp/lo.zzu.wiki.key")
	pem, _ := os.ReadFile("resources/temp/lo.zzu.wiki.pem")

	req := &request.ConfigCertificateCreate{
		PublicKey:  string(pem),
		PrivateKey: string(key),
		Enabled:    true,
	}
	svc.Create(ctx, req)
}
