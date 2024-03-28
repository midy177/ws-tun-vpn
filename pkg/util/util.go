package util

import (
	"bytes"
	"fmt"
	validate "github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"strconv"
)

func ValidateWithFatal(field any, tag, flag string) {
	err := validate.New().Var(field, tag)
	if err != nil {
		log.Println(err)
		log.Fatalf(" validation failed for flag: %s, failed on the '%s' tag", flag, tag)
	}
}

func FlagRequiredWithFatal(cmd *cobra.Command, name string) {
	err := cmd.MarkFlagRequired(name)
	if err != nil {
		log.Fatal(err)
	}
}

// GenerateTunName Randomly generate tun network card name
func GenerateTunName(n int) string {
	const charset = "0123456789"
	length := len(charset)
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(length)]
	}
	return "utun" + string(b)
}

// LenToSubNetMask 如 24 对应的子网掩码地址为 255.255.255.0
func LenToSubNetMask(subnet int) string {
	var buff bytes.Buffer
	for i := 0; i < subnet; i++ {
		buff.WriteString("1")
	}
	for i := subnet; i < 32; i++ {
		buff.WriteString("0")
	}
	masker := buff.String()
	a, _ := strconv.ParseUint(masker[:8], 2, 64)
	b, _ := strconv.ParseUint(masker[8:16], 2, 64)
	c, _ := strconv.ParseUint(masker[16:24], 2, 64)
	d, _ := strconv.ParseUint(masker[24:32], 2, 64)
	resultMask := fmt.Sprintf("%v.%v.%v.%v", a, b, c, d)
	return resultMask
}
