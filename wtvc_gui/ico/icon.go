package ico

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

var (
	//go:embed src/favicon.ico
	Favicon []byte
	//go:embed src/add.png
	AddIcon []byte
)

func LoadIcon() fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "favicon.ico",
		StaticContent: Favicon,
	}
}

func LoadAddIcon() fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "add.png",
		StaticContent: AddIcon,
	}
}
