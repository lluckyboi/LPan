package tool

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func GetHash(path string) (hash string) {
	file, err := os.Open(path)
	if err == nil {
		h_ob := sha1.New()
		_, err := io.Copy(h_ob, file)
		if err == nil {
			hash := h_ob.Sum(nil)
			hashvalue := hex.EncodeToString(hash)
			return hashvalue
		} else {
			return "something wrong when use sha256 interface..."
		}
	} else {
		log.Printf("failed to open %s\n", path)
	}
	defer file.Close()
	return
}
