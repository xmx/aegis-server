package temp_test

import (
	"archive/zip"
	"github.com/xmx/aegis-server/jsrun/jsmod"
	"github.com/xmx/aegis-server/jsrun/jsvm"
	"github.com/xmx/aegis-server/jzip"
	"os"
	"testing"
)

func TestGoja(t *testing.T) {
	mods := []jsvm.ModuleRegister{
		jsmod.NewConsole(os.Stdout, os.Stdout),
		jsmod.NewOS(),
	}
	eng, err := jsvm.New(mods)
	if err != nil {
		t.Fatal(err)
	}

	zr, err := zip.OpenReader("main.zip")
	if err != nil {
		t.Fatal(err)
	}
	_, err = eng.RunZip(zr)
	t.Logf("是否错误：%v", err)
}

func TestJZip(t *testing.T) {
	soft, err := jzip.Open("main.zip")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(soft)
}
