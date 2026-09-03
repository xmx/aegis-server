package restapi

import (
	"io"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
	"github.com/xmx/aegis-common/conpty"
)

type PTY struct {
	wsupg websocket.Upgrader
}

func NewPTY() *PTY {
	return &PTY{
		wsupg: websocket.Upgrader{
			HandshakeTimeout: 5 * time.Second,
			CheckOrigin: func(*http.Request) bool {
				return true
			},
		},
	}
}

func (pty *PTY) RegisterRoute(r *echo.Group) error {
	r.GET("/system/pty", pty.start)

	return nil
}

func (pty *PTY) start(c *echo.Context) error {
	w, r := c.Response(), c.Request()
	ws, err := pty.wsupg.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	var command string
	cmds := []string{"pwsh.exe", "C:\\Program Files\\PowerShell\\7\\pwsh.exe", "powershell.exe", "cmd.exe"}
	for _, cmd := range cmds {
		if command, _ = exec.LookPath(cmd); command != "" {
			break
		}
	}
	if command == "" {
		command = "C:\\WINDOWS\\system32\\cmd.exe"
	}
	cpty, err := conpty.Start(command)
	if err != nil {
		return err
	}
	defer cpty.Close()

	stdout, stdin := conpty.NewSocketIO(ws)
	go io.Copy(stdout, cpty)

	for {
		data, err1 := stdin.Read()
		if err1 != nil {
			return err1
		}

		switch data.Type {
		case "kill":
			break
		case "stdin":
			msg, err2 := data.AsString()
			if err2 != nil {
				return err2
			}
			if _, err2 = cpty.Write([]byte(msg)); err2 != nil {
				return err2
			}
		case "resize":
			msg, err2 := data.AsResize()
			if err2 != nil {
				return err2
			}
			if err2 = cpty.Resize(msg.Cols, msg.Rows); err2 != nil {
				return err2
			}
		default:
			return echo.ErrBadRequest
		}
	}
}
