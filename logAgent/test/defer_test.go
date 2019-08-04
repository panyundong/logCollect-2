package test

import (
	"fmt"
	"testing"
)

func TestDefer(t *testing.T) {
	defer func() {
		fmt.Println("这是第二个,第三执行")
	}()
	defer func() {
		fmt.Println("这是第一个,第二执行")
	}()

	fmt.Println("第三个先执行")
}
