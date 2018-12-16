package processor

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
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


func base64Decode(text ...string) (string, error) {
	buff := bytes.NewBuffer(nil)
	for _, t := range text {
		if _, err := buff.WriteString(t); err != nil {
			return "", err
		}
	}
	data, err := base64.StdEncoding.DecodeString(buff.String())
	if err != nil {
		return "", err
	}
	return bytes.NewBuffer(data).String(), nil
}

func base64Encode(text ...string) (string, error) {
	buff := bytes.NewBuffer(nil)
	for _, t := range text {
		if _, err := buff.WriteString(t); err != nil {
			return "", err
		}
	}
	return base64.StdEncoding.EncodeToString(buff.Bytes()), nil
}