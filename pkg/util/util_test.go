package util

import (
	"fmt"
	"testing"
)

func TestUtil(t *testing.T) {
	name := LenToSubNetMask(27)
	fmt.Print(name + "\n")
}
