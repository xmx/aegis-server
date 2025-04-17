package webdavfs

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/webdav"
)

func New(dir string) http.Handler {
	dir = filepath.Clean(dir)
	dav := &webdav.Handler{
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
	}
	d := http.Dir(dir)

	return &davFS{
		dav: dav,
		hfs: d,
		han: http.FileServer(d),
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

	if !f.tryServeJSONDir(w, r, r.URL.Path) {
		f.han.ServeHTTP(w, r)
	}
}

func (f *davFS) tryServeJSONDir(w http.ResponseWriter, r *http.Request, reqPath string) bool {
	accept := r.Header.Get("Accept")
	accepts := strings.Split(accept, ",")
	var allowed bool
	for _, str := range accepts {
		ct, _, _ := strings.Cut(str, ";")
		if allowed = ct == "application/json" || ct == "*/*"; allowed {
			break
		}
	}
	if !allowed {
		return false
	}

	file, err := f.hfs.Open(reqPath)
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
		if link, exx := os.Readlink(filepath.Join(reqPath, name)); exx == nil {
			info.Symlink = link
			if lstat, _ := os.Stat(filepath.Join(reqPath, link)); lstat != nil {
				info.Directory = lstat.IsDir()
			}
		}
	}
	sort.Sort(files)
	ret := &dirInfo{Path: path.Clean(reqPath), Files: files}

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
	Free       int64     `json:"free,omitempty"` // windows 磁盘
	Mode       string    `json:"mode,omitempty"`
	Directory  bool      `json:"directory,omitempty"`
	Symlink    string    `json:"symlink,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitzero"`
	CreatedAt  time.Time `json:"created_at,omitzero"`
	AccessedAt time.Time `json:"accessed_at,omitzero"`
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
