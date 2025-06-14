package example10_1

import (
	"crypto/sha256"
	"io"
	"os"
)

func FileHash(filepath string) []byte {
	// opens file from passed in filepath
	file, _ := os.Open(filepath)

	// calculating hashcode for the file by using its content as hash input
	sha := sha256.New()
	io.Copy(sha, file)

	// returning the hash result as a slice of bytes and not the sha256 hash object
	return sha.Sum(nil)
}
