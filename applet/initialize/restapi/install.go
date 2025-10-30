package restapi

import (
	"net"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/profile"
	"github.com/xmx/aegis-control/mongodb"
	"github.com/xmx/aegis-server/applet/initialize/request"
	"github.com/xmx/aegis-server/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewInstall(results chan<- *config.Config) *Install {
	return &Install{
		results: results,
	}
}

type Install struct {
	results chan<- *config.Config
}

func (inst *Install) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/install").POST(inst.setup)
	return nil
}

func (inst *Install) setup(c *ship.Context) error {
	req := new(request.InstallSetup)
	if err := c.Bind(req); err != nil {
		return err
	}

	addr := req.Server.Addr
	{
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		_ = ln.Close()
	}
	//{
	//	ln, err := net.Listen("udp", addr)
	//	if err != nil {
	//		return err
	//	}
	//	_ = ln.Close()
	//}
	{
		parent := c.Request().Context()
		uri := req.Database.URI
		db, err := mongodb.Open(uri)
		if err != nil {
			return err
		}
		cli := db.Client()
		defer cli.Disconnect(parent)

		if err = cli.Ping(parent, nil); err != nil {
			return err
		}
	}
	cfg := &config.Config{
		Server: config.Server{
			Addr:   addr,
			Static: req.Server.Static,
		},
		Database: config.Database{
			URI: req.Database.URI,
		},
		Logger: config.Logger{
			Level:   req.Logger.Level,
			Console: req.Logger.Console,
			Logger: &lumberjack.Logger{
				Filename:   config.LogFilename,
				MaxSize:    req.Logger.MaxSize,
				MaxAge:     req.Logger.MaxAge,
				MaxBackups: req.Logger.MaxBackups,
				LocalTime:  req.Logger.LocalTime,
				Compress:   req.Logger.Compress,
			},
		},
	}

	if err := profile.WriteFile(config.Filename, cfg); err != nil {
		return err
	}
	inst.results <- cfg

	return nil
}
