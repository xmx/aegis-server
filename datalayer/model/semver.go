package model

import (
	"encoding/binary"
	"strconv"
	"strings"
)

type Semver struct {
	Version string `bson:"version" json:"version"`
	Number  uint64 `bson:"number"  json:"number"`
}

// ParseSemver https://semver.org/
//
// https://semver.org/#backusnaur-form-grammar-for-valid-semver-versions
// <valid semver> ::= <version core>
//
//	| <version core> "-" <pre-release>
//	| <version core> "+" <build>
//	| <version core> "-" <pre-release> "+" <build>
func ParseSemver(version string) Semver {
	sem := Semver{Version: version}

	core := version
	var prerelease string
	if strings.Contains(version, "-") {
		core, prerelease, _ = strings.Cut(core, "-")
	} else if strings.Contains(version, "+") {
		core, prerelease, _ = strings.Cut(version, "+")
	}

	slots := make([]byte, 8)
	cores := strings.Split(core, ".")
	for i := range 3 {
		num, _ := strconv.ParseUint(cores[i], 10, 64)
		slots[i] = byte(num)
	}
	copy(slots[4:], prerelease)

	num := binary.BigEndian.Uint64(slots)
	sem.Number = num

	return sem
}

func FormatSemver(num uint64) Semver {
	sem := Semver{Number: num}
	slots := make([]byte, 8)
	binary.BigEndian.PutUint64(slots, num)

	var cores []string
	for i := range 3 {
		str := strconv.FormatInt(int64(slots[i]), 10)
		cores = append(cores, str)
	}
	core := strings.Join(cores, ".")
	var prerelease []byte
	for i, v := range slots[3:] {
		if i != 0 {
			prerelease = append(prerelease, v)
		}
	}
	if len(prerelease) != 0 {
		elems := []string{core, string(prerelease)}
		sem.Version = strings.Join(elems, "-")
	} else {
		sem.Version = core
	}

	return sem
}
