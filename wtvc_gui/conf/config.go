package conf

import (
	jsoniter "github.com/json-iterator/go"
	"os"
	"ws-tun-vpn/pkg/syncmap"
	"ws-tun-vpn/types"
)

var CfgMap syncmap.Map[string, *types.ClientConfig]

var cfgFilename = "conf/config.json"

type Config []*types.ClientConfig

func Load() []string {
	file, err := os.ReadFile(cfgFilename)
	if err != nil {
		return nil
	}
	var cfg Config

	err = jsoniter.Unmarshal(file, cfg)
	if err != nil {
		return nil
	}
	//var names = make([]string, 0, len(cfg))
	//for _, v := range cfg {
	//CfgMap.Store(v.Name, v)
	//names = append(names, v.Name)
	//}
	return nil
}
func Reload() []string {
	var names = make([]string, 0, 10)
	CfgMap.Range(func(k string, v *types.ClientConfig) bool {
		names = append(names, k)
		return true
	})
	return names
}
