package main

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func TestError(t *testing.T) {

	err := errors.New("whoops")
	e2 := errors.Wrap(err, "inner")
	fmt.Printf("%+v", e2)

}
