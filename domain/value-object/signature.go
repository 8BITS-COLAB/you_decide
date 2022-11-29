package valueobject

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

func NewSignature(raw string, diff int) (string, int) {
	h := sha1.New()
	nonce := 0

	for {
		h.Write([]byte(fmt.Sprintf("%s:%d", raw, nonce)))
		buffer := h.Sum(nil)
		hash := fmt.Sprintf("%x", buffer)

		if strings.HasPrefix(hash, strings.Repeat("0", diff)) {
			return hash, nonce
		}
		nonce++
	}
}
