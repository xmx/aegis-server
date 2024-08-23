package profile

import (
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func readdir(root, pattern string) iter.Seq2[string, error] {
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
				matched, err := filepath.Match(pattern, name)
				if err != nil {
					yield("", err)
					break
				}
				if !matched {
					continue
				}
			}

			if !yield(filepath.Join(root, name), nil) {
				break
			}
		}
	}
}
