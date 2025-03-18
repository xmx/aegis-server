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

func New(prefix, dir string) http.Handler {
	dir = filepath.Clean(dir)
	dav := &webdav.Handler{
		Prefix:     prefix,
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
	}
	d := http.Dir(dir)

	return &davFS{
		dav: dav,
		hfs: d,
		han: http.StripPrefix(prefix, http.FileServer(d)),
	}
}

type davFS struct {
	dav *webdav.Handler
	hfs http.FileSystem
	han http.Handler // GET webui
}

func (f *davFS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		f.dav.ServeHTTP(w, r)
		return
	}

	path := r.URL.Path
	if !f.tryServeJSONDir(w, r, path) {
		f.han.ServeHTTP(w, r)
	}
}

func (f *davFS) tryServeJSONDir(w http.ResponseWriter, r *http.Request, path string) bool {
	accept := r.Header.Get("Accept")
	if accept != "application/json" {
		return false
	}

	path = strings.TrimPrefix(path, f.dav.Prefix)
	file, err := f.hfs.Open(path)
	if err != nil {
		return false
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}

	// 最多读取 maxsize 个文件，防止目录文件过多，影响内存安全。
	const maxsize = 10000
	infos, _ := file.Readdir(maxsize)
	files := make(fileInfos, 0, len(infos))
	for _, fi := range infos {
		st := readStat(fi.Sys())
		name, mode := fi.Name(), fi.Mode()
		info := &fileInfo{
			Name:       name,
			Size:       fi.Size(),
			Mode:       mode.String(),
			Directory:  fi.IsDir(),
			UpdatedAt:  fi.ModTime(),
			CreatedAt:  st.CreatedAt,
			AccessedAt: st.AccessedAt,
			User:       st.User,
			Group:      st.Group,
		}

		files = append(files, info)
		if mode&os.ModeSymlink == 0 {
			continue
		}
		// 链接类型文件
		if link, exx := os.Readlink(filepath.Join(path, name)); exx == nil {
			info.Symlink = link
		}
	}
	sort.Sort(files)
	ret := &dirInfo{Path: filepath.Clean(path), Files: files}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(ret)

	return true
}

type dirInfo struct {
	Path  string    `json:"path"`
	Files fileInfos `json:"files"`
}

type fileInfo struct {
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	Mode       string    `json:"mode"`
	Directory  bool      `json:"directory,omitempty"`
	Symlink    string    `json:"symlink,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	AccessedAt time.Time `json:"accessed_at,omitempty"`
	User       string    `json:"user,omitempty"`
	Group      string    `json:"group,omitempty"`
}

type fileInfos []*fileInfo

func (fis fileInfos) Len() int {
	return len(fis)
}

// Less 目录靠前显示，按照文件名升序。
func (fis fileInfos) Less(i, j int) bool {
	fi, fj := fis[i], fis[j]
	if fi.Directory == fj.Directory {
		return fi.Name < fj.Name
	} else if fi.Directory {
		return true
	} else {
		return false
	}
}

func (fis fileInfos) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

type sysStat struct {
	AccessedAt time.Time `json:"accessed_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       string    `json:"user,omitempty"`
	Group      string    `json:"group,omitempty"`
}
