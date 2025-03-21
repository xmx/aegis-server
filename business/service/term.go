package service

import (
	"crypto/ed25519"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/socketerm"
	"github.com/xmx/aegis-server/protocol/asciicast"
	"github.com/xmx/aegis-server/protocol/vterminal"
	"golang.org/x/crypto/ssh"
)

func NewTerm(log *slog.Logger) *Term {
	return &Term{log: log}
}

type Term struct {
	log *slog.Logger
}

//goland:noinspection GoUnhandledErrorResult
func (svc *Term) PTY(conn *websocket.Conn, size *request.TermResize) error {
	if runtime.GOOS == "windows" {
		return errors.ErrUnsupported
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		for _, sh := range []string{
			"/usr/bin/fish", "/usr/bin/zsh", "/usr/bin/bash",
		} {
			if _, err := os.Lstat(sh); err == nil {
				shell = sh
			}
		}
	}
	if shell == "" {
		shell = "/bin/bash"
	}
	svc.log.Info("准备执行 SHELL", slog.Any("shell", shell))

	cmd := exec.Command(shell)
	tty, err := vterminal.NewPTMX(cmd, size.Cols, size.Rows)
	if err != nil {
		return err
	}
	defer tty.Close()

	castFile := filepath.Join("resources/static/play/", "pty.cast")
	err = svc.serveTerminal(conn, tty, castFile)

	return err
}

func (svc *Term) SSH(ws *websocket.Conn, req *request.TermSSH) error {
	_, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}
	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return err
	}

	auths := []ssh.AuthMethod{ssh.PublicKeys(signer)} // publickey 认证
	if req.Password != "" {
		auths = append(auths, ssh.Password(req.Password)) // password 认证，预先输入密码
	}
	auths = append(auths, ssh.RetryableAuthMethod(svc.keyboardInteractive(ws), 3))                    // keyboard-interactive 认证
	auths = append(auths, ssh.RetryableAuthMethod(ssh.PasswordCallback(svc.passwordCallback(ws)), 3)) // password 认证，后输入密码

	sshCfg := &ssh.ClientConfig{
		User:            req.Username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Minute,
	}
	tty, err := vterminal.NewSSH("tcp", req.Bastion, sshCfg, req.Cols, req.Rows)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer tty.Close()

	castFile := filepath.Join("resources/static/play/", "ssh.cast")
	err = svc.serveTerminal(ws, tty, castFile)

	return err
}

//goland:noinspection GoUnhandledErrorResult
func (svc *Term) serveTerminal(ws *websocket.Conn, tty vterminal.Typewriter, castFile string) error {
	conn := socketerm.New(ws, 10*time.Second)

	var cast asciicast.Writer
	if castFile != "" {
		file, err := os.Create(castFile)
		if err != nil {
			return err
		}
		defer file.Close()

		cols, rows, _ := tty.Size()
		header := asciicast.NewHeader(cols, rows)
		cast = header.NewWriter(file)
	}

	ch := make(chan struct{})
	go func() {
		var read io.Reader = tty
		if cast != nil {
			read = io.TeeReader(read, cast)
		}
		_, _ = io.Copy(conn, read)
		_ = ws.CloseNow()
		close(ch)
	}()

	var over bool
	for !over {
		code, data, err := conn.Recv()
		if err != nil {
			over = true
			break
		}
		switch code {
		case "i":
			if _, err = tty.Write([]byte(data)); err != nil {
				over = true
				break
			}
		case "r":
			cols, rows := asciicast.ParseResize(strings.ToLower(data))
			if cols <= 0 || rows <= 0 {
				break
			}

			if exx := tty.Resize(cols, rows); exx == nil && cast != nil {
				_ = cast.Resize(cols, rows)
			}
		}
	}
	_ = tty.Close()
	<-ch

	return nil
}

func (svc *Term) passwordCallback(ws *websocket.Conn) func() (string, error) {
	conn := socketerm.New(ws, 10*time.Second)
	return func() (string, error) {
		if _, err := conn.Write([]byte("请输入密码：")); err != nil {
			return "", err
		}
		answer, err := svc.readUntilCR(conn, false)
		if err != nil {
			return "", err
		}
		if _, err = conn.Write([]byte("\r\n")); err != nil {
			return "", err
		}

		return answer, nil
	}
}

func (svc *Term) keyboardInteractive(ws *websocket.Conn) ssh.KeyboardInteractiveChallenge {
	return func(name, instruction string, questions []string, echos []bool) ([]string, error) {
		conn := socketerm.New(ws, 10*time.Second)
		answers := make([]string, 0, len(questions))
		size := len(echos) - 1
		for i, question := range questions {
			if _, err := conn.Write([]byte(question)); err != nil {
				return nil, err
			}
			var echo bool
			if size >= i {
				echo = echos[i]
			}

			answer, err := svc.readUntilCR(conn, echo)
			if err != nil {
				return nil, err
			}
			answers = append(answers, answer)
			if _, err = conn.Write([]byte("\r\n")); err != nil {
				return nil, err
			}
		}

		return answers, nil
	}
}

func (*Term) readUntilCR(conn *socketerm.Conn, echo bool) (string, error) {
	var answer string
	for {
		_, data, err := conn.Recv()
		if err != nil {
			return "", err
		}
		if data == "\r" {
			break
		}
		answer += data
		if !echo {
			continue
		}
		if _, err = conn.Write([]byte(data)); err != nil {
			return "", err
		}
	}

	return answer, nil
}
