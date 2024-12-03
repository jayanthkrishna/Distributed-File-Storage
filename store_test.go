package main

import (
	"bytes"
	"fmt"
	"io"

	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mypicture"
	pathname := CASPathTransformFunc(key)

	fmt.Println(pathname)
}

func TestStore(t *testing.T) {

	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	f := []byte("Some file bytes")
	data := bytes.NewReader(f)
	key := "somekey"
	if err := s.writeStream(key, data); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("Expected to have a key %s", key)
	}

	r, err := s.Read(key)

	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)

	if string(b) != string(f) {
		t.Errorf("want %s have %s", f, b)
	}

	s.Delete(key)
}

func TestDelete(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	key := "somekey"
	s.Delete(key)
	// os.RemoveAll("jayanthCAS/f16b0")
}
