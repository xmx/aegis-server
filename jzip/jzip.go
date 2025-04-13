package jzip

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
)

type Manifest struct {
	Version     int         `json:"version"`
	Application Application `json:"application"`
}

type Application struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Main    string `json:"main"`
	Version string `json:"version"`
}

func Open(name string) (*JZip, error) {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}

	manifestFile, err := zr.Open("manifest.json")
	if err != nil {
		_ = zr.Close()
		return nil, fmt.Errorf("缺少 manifest.json: %w", err)
	}
	defer manifestFile.Close()

	manifest := new(Manifest)
	if err = json.NewDecoder(manifestFile).Decode(manifest); err != nil {
		_ = zr.Close()
		return nil, fmt.Errorf("解析 manifest.json 错误: %w", err)
	}
	app := manifest.Application
	if app.ID == "" {
		_ = zr.Close()
		return nil, fmt.Errorf("软件ID不能为空")
	}
	if app.Name == "" {
		_ = zr.Close()
		return nil, fmt.Errorf("软件名不能为空")
	}

	return &JZip{Manifest: manifest, zipFile: zr}, nil
}

type JZip struct {
	Manifest *Manifest
	zipFile  *zip.ReadCloser
}

func (jz *JZip) Open(name string) (io.ReadCloser, error) {
	return jz.zipFile.Open(name)
}

func (jz *JZip) Close() error {
	return jz.zipFile.Close()
}

func (jz *JZip) String() string {
	app := jz.Manifest.Application
	return "ID: " + app.ID + " , Name: " + app.Name + " , Version: " + app.Version
}
