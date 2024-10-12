package directory

import (
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func Walk(root, pattern string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		dir, err := os.Open(root)
		if err != nil {
			yield("", err)
			return
		}
		names, err := dir.Readdirnames(-1)
		_ = dir.Close()
		if err != nil {
			yield("", err)
			return
		}

		slices.Sort(names)
		for _, name := range names {
			if strings.HasPrefix(name, ".") { // skip hiding file.
				continue
			}
			if pattern != "" {
				if matched, exx := filepath.Match(pattern, name); exx != nil {
					yield("", exx)
					break
				} else if !matched {
					continue
				}
			}

			path := filepath.Join(root, name)
			if !yield(path, nil) {
				break
			}
		}
	}
}
