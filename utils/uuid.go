package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func NewUUIDv7() (string, error) {
	var uuid [16]byte
	timestamp := uint64(time.Now().UnixMilli())

	uuid[0] = byte(timestamp >> 40)
	uuid[1] = byte(timestamp >> 32)
	uuid[2] = byte(timestamp >> 24)
	uuid[3] = byte(timestamp >> 16)
	uuid[4] = byte(timestamp >> 8)
	uuid[5] = byte(timestamp)

	if _, err := rand.Read(uuid[6:]); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x70
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return formatUUID(uuid), nil
}

func formatUUID(uuid [16]byte) string {
	var dst [36]byte

	hex.Encode(dst[0:8], uuid[0:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:36], uuid[10:16])

	return string(dst[:])
}
