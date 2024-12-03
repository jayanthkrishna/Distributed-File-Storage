package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"

	"log"
	"os"
	"strings"
)

const ROOT = "jayanthCAS"

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	FileName string
}

func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.FileName)
}

func (p *PathKey) FirstPathName() string {
	return strings.Split(p.FullPath(), "/")[0]
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		FileName: key,
	}
}

type StoreOpts struct {
	Root              string
	PathTransformFunc PathTransformFunc
}
type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = ROOT
	}
	return &Store{
		StoreOpts: opts,
	}
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))

	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := range paths {
		from, to := i*blockSize, i*blockSize+blockSize

		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		FileName: hashStr,
	}

}

func (s *Store) Has(key string) bool {

	pathKey := s.PathTransformFunc(key)

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)

	if err == err.(*os.PathError) {
		return false
	}

	return true
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}
func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	firstPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())
	fmt.Println(firstPathNameWithRoot)
	defer func() {
		log.Printf("deleted %s from the disk", pathKey.FileName)
	}()
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Store) Read(key string) (io.Reader, error) {

	f, err := s.readStream(key)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathkey := s.PathTransformFunc(key)
	fullpathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathkey.FullPath())
	// fullPath := pathkey.FullPath()

	f, err := os.Open(fullpathWithRoot)

	if err != nil {
		return nil, err
	}

	return f, nil

}

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
}

func (s *Store) writeStream(key string, r io.Reader) error {

	pathkey := s.PathTransformFunc(key)

	pathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathkey.Pathname)
	if err := os.MkdirAll(pathWithRoot, os.ModePerm); err != nil {
		return err
	}

	pathAndFileName := fmt.Sprintf("%s/%s", s.Root, pathkey.FullPath())

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
