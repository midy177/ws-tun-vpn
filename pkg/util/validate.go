package util

import (
	validate "github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"log"
)

func ValidateWithFatal(field any, tag string) {
	err := validate.New().Var(field, tag)
	if err != nil {
		// 校验失败，处理错误
		log.Fatal(err)
	}
}

func FlagRequiredWithFatal(cmd *cobra.Command, name string) {
	err := cmd.MarkFlagRequired(name)
	if err != nil {
		log.Fatal(err)
	}
}
