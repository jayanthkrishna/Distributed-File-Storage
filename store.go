package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

type PathTransformFunc func(string) string

var DefaultPathTransformFunc = func(key string) string {
	return key
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}
type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func CASPathTransformFunc(key string) string {
	hash := sha1.Sum([]byte(key))

	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := range paths {
		from, to := i*blockSize, i*blockSize+blockSize

		paths[i] = hashStr[from:to]
	}

	return strings.Join(paths, "/")

}

func (s *Store) writeStream(key string, r io.Reader) error {

	pathName := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return err
	}

	filename := "somefile"
	pathAndFileName := pathName + "/" + filename

	f, err := os.Create(pathAndFileName)

	if err != nil {
		return err

	}

	n, err := io.Copy(f, r)

	if err != nil {
		return err
	}

	log.Printf("Written (%d) bytes to disk : %s \n", n, pathAndFileName)

	return nil
}
