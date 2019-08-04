package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	var address string = "127.0.0.1:2379"

	split := strings.Split(address, ",")

	for _, value := range split {
		fmt.Println(value)
	}
}
