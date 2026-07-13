package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	timeCost    uint32 = 3
	memoryCost  uint32 = 32 * 1024
	parallelism uint8  = 4
	keyLen             = 32
	saltLen            = 16
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	rand.Read(salt)
	key := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, parallelism, keyLen)

	b64 := base64.RawStdEncoding
	saltB64 := b64.EncodeToString(salt)
	keyB64 := b64.EncodeToString(key)

	return fmt.Sprintf("$argon2id$v=19$m=%v,t=%v,p=%v$%v$%v", strconv.FormatUint(uint64(memoryCost), 10), strconv.FormatUint(uint64(timeCost), 10), strconv.FormatUint(uint64(parallelism), 10), saltB64, keyB64), nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	// Expect: ["", "argon2id", "v=19", "m=..,t=..,p=..", "<saltB64>", "<hashB64>"]
	parts := strings.Split(hash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return false, errors.New("invalid argon2id PHC")
	}

	otherParts := strings.Split(parts[3], ",")
	if len(otherParts) != 3 {
		return false, errors.New("invalid argon2id PHC")
	}

	m, err := strconv.ParseUint(otherParts[0][2:], 10, 32)
	if err != nil {
		return false, errors.New("invalid memoryCost")
	}
	t, err := strconv.ParseUint(otherParts[1][2:], 10, 32)
	if err != nil {
		return false, errors.New("invalid timeCost")
	}
	p, err := strconv.ParseUint(otherParts[2][2:], 10, 8)
	if err != nil {
		return false, errors.New("invalid parallelism")
	}

	b64 := base64.RawStdEncoding
	salt, err := b64.DecodeString(parts[4])
	if err != nil {
		return false, errors.New("invalid salt encoding")
	}
	want, err := b64.DecodeString(parts[5])
	if err != nil {
		return false, errors.New("invalid key encoding")
	}

	got := argon2.IDKey([]byte(password), salt, uint32(t), uint32(m), uint8(p), uint32(len(want)))
	// Avoid timing leaks:
	return subtle.ConstantTimeCompare(got, want) == 1, nil
}
