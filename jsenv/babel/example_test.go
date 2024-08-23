package babel_test

import (
	"testing"

	"github.com/xmx/aegis-server/jsenv/babel"
)

var codes = []string{
	`import { btn } from 'antd'; btn.click();`,
	`import http from 'k6/http'
import { check, sleep } from 'k6'

export default function () {
  const data = { username: 'username', password: 'password' }
  let res = http.post('https://myapi.com/login/', data)

  check(res, { 'success login': (r) => r.status === 200 })

  sleep(0.3)
}`,
}

func TestBabel(t *testing.T) {
	opt := map[string]any{
		"plugins": []string{
			"transform-modules-commonjs",
		},
	}
	for _, code := range codes {
		dest, err := babel.Transform(code, opt)
		t.Log(dest, err)
	}
}
