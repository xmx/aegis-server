package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"time"
)

type UUID [16]byte

func (u UUID) String() string {
	return u.Hex()
}

func (u UUID) Hex() string {
	return hex.EncodeToString(u[:])
}

// V7 https://antonz.org/uuidv7/#go
func V7() UUID {
	// random bytes
	var value UUID
	_, _ = rand.Read(value[:])

	// current timestamp in ms
	timestamp := big.NewInt(time.Now().UnixMilli())
	timestamp.FillBytes(value[0:6])

	// version and variant
	value[6] = (value[6] & 0x0F) | 0x70
	value[8] = (value[8] & 0x3F) | 0x80

	return value
}
