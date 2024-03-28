package util

import (
	"fmt"
	"testing"
)

func TestUtil(t *testing.T) {
	name := GenerateTunName(4)
	fmt.Print(name + "\n")
}
