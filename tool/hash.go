package tool

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"mime/multipart"
)

func GetHash(file multipart.File) (string, error) {
	h_ob := sha1.New()
	_, err := io.Copy(h_ob, file)
	if err != nil {
		return "", err
	}
	hash := base64.StdEncoding.EncodeToString(h_ob.Sum(nil))
	return hash, nil
}
