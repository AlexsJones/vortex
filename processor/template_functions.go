package processor

import (
	"crypto/md5"
	"fmt"
	"io"
)

func hashMd5(text ...string) (string, error) {
	h := md5.New()
	for _, s := range text {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
