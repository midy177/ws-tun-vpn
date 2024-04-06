package app_context

import (
	"fyne.io/fyne/v2"
	jsoniter "github.com/json-iterator/go"
	"log"
	"ws-tun-vpn/types"
)

const (
	MainPage         NavEventType = 0x01
	AddConfigPage    NavEventType = 0x02
	ModifyConfigPage NavEventType = 0x03
	ConnectTo        NavEventType = 0x04
)

type NavEventType rune

type AppConfigs struct {
	Label string `json:"label"`
	types.ClientConfig
}

type NavEvent struct {
	TargetPage NavEventType
}

type AppContext struct {
	Window      fyne.Window
	Preferences fyne.Preferences
	AppConfigs  []AppConfigs
	// Channel for navigation events
	NavChannel     chan NavEvent
	SelectedItemID int64
}

func NewAppContext(w fyne.Window, p fyne.Preferences) *AppContext {
	return &AppContext{
		Window:         w,
		Preferences:    p,
		NavChannel:     make(chan NavEvent, 1),
		SelectedItemID: -1,
	}
}

func (ctx *AppContext) LoadAppConfigs() {
	// Create your settings content here
	var appConfigs []AppConfigs
	settingsJSON := ctx.Preferences.String("AppConfigs")
	if settingsJSON != "" {
		err := jsoniter.Unmarshal([]byte(settingsJSON), &appConfigs)
		if err != nil {
			log.Println("Error loading settings:", err)
		}
		ctx.AppConfigs = appConfigs
	} else {
		// Set default settings if no saved settings are found
		ctx.AppConfigs = nil
	}
}

func (ctx *AppContext) UpdateAppConfigs() {
	// Serialize settings to JSON
	settingsJSON, err := jsoniter.Marshal(ctx.AppConfigs)
	if err != nil {
		log.Println("Error marshaling settings:", err)
		return
	}
	// Save JSON string to app preferences
	ctx.Preferences.SetString("AppConfigs", string(settingsJSON))
}
