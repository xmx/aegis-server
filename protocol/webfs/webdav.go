package webfs

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/webdav"
)

func DAV(prefix, dir string) http.Handler {
	prefix = strings.TrimRight(prefix, "/")
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	dir = filepath.Clean(dir)

	dav := &webdav.Handler{
		Prefix:     prefix,
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
	}
	d := http.Dir(dir)

	return &davFS{
		dir: dir,
		dav: dav,
		hfs: d,
		han: http.FileServer(d),
	}
}

type davFS struct {
	dir string
	dav *webdav.Handler
	hfs http.FileSystem
	han http.Handler
}

func (f *davFS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		f.dav.ServeHTTP(w, r)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, f.dav.Prefix)
	accept := r.Header.Get("Accept")
	var processed bool
	if accept == "application/json" {
		info, ok, err := f.readdir(path)
		if err == nil && ok {
			processed = true
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_ = json.NewEncoder(w).Encode(info)
		}
	}
	if !processed {
		r.URL.Path = path
		f.han.ServeHTTP(w, r)
	}
}

func (f *davFS) readdir(path string) (*DirInfo, bool, error) {
	info := &DirInfo{
		Path:  path,
		Files: make(FileInfos, 0, 100),
	}

	file, err := f.hfs.Open(path)
	if err != nil {
		return nil, false, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, false, err
	}
	if !stat.IsDir() {
		return nil, false, nil
	}

	infos, err := file.Readdir(4096)
	for _, in := range infos {
		name, mode := in.Name(), in.Mode()
		inf := &FileInfo{
			Name:      name,
			Size:      in.Size(),
			Mode:      mode.String(),
			ModTime:   in.ModTime(),
			Directory: in.IsDir(),
		}
		info.Files = append(info.Files, inf)
		if mode&os.ModeSymlink == 0 {
			continue
		}

		link, err := os.Readlink(filepath.Join(path, name))
		if err != nil {
			continue
		}
		inf.Symlink = link
		if !filepath.IsAbs(link) {
			link = filepath.Join(path, link)
		}
		if fi, _ := os.Stat(link); fi != nil {
			inf.Directory = fi.IsDir()
		}
	}
	sort.Sort(info.Files)

	return info, true, nil
}

type DirInfo struct {
	Path  string    `json:"path"`
	Files FileInfos `json:"files"`
}

type FileInfo struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Mode      string    `json:"mode"`
	ModTime   time.Time `json:"mod_time"`
	Directory bool      `json:"directory"`
	Symlink   string    `json:"symlink,omitempty"`
}

type FileInfos []*FileInfo

func (fis FileInfos) Len() int {
	return len(fis)
}

func (fis FileInfos) Less(i, j int) bool {
	fi, fj := fis[i], fis[j]
	if fi.Directory == fj.Directory {
		return fi.Name < fj.Name
	} else if fi.Directory {
		return true
	} else {
		return false
	}
}

func (fis FileInfos) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}
