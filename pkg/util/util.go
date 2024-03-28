package util

import (
	validate "github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
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
	return "tun" + string(b)
}
