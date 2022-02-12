package main

import (
	"crypto/rand"
	"math/big"

	"github.com/segmentio/ksuid"
)

// Alternative alphabet for further usage e.g. pw generation
// const l = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

// b64a contains the base64 character alphabet
const b64a = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// getStrLength returns the amount of valid characters instead of the number of bytes (like len())
func getStrLength(s string) int {
	return len([]rune(s))
}

// createRandomTokenB64 ... builds a string based on alphabet b64a
func createRandomTokenB64(l int) (r string, err error) {
	for {
		if l <= getStrLength(r) {
			return r, nil
		}
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(b64a))))
		if err != nil {
			return "", err
		}
		r = r + string(b64a[n.Int64()])
	}
}

// GetRandomToken32 ... generates a 32 byte token based upon base64 encoding
func GetRandomToken32() (string, error) {

	s, err := createRandomTokenB64(32)

	return s, err
}

func (app *Controller) customGenerator() string {
	if app.stats.NewRequest || app.stats.RequestID == "" {
		id := ksuid.New()
		app.stats.NewRequest = false
		app.stats.RequestID = id.String()
		return app.stats.RequestID
	}
	return app.stats.RequestID
}
