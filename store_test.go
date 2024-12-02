package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mypicture"
	pathname := CASPathTransformFunc(key)

	fmt.Println(pathname)
}

func TestStore(t *testing.T) {

	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}

	s := NewStore(opts)

	data := bytes.NewReader([]byte("Some file bytes"))

	if err := s.writeStream("somekey", data); err != nil {
		t.Error(err)
	}

}
