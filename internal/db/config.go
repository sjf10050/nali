package db

import "github.com/spf13/viper"

// Preloaded config values — read once at startup and used on the hot path,
// avoiding repeated viper.GetString calls (which involve mutex locks and
// multiple map lookups).
var (
	selectedLang string
	selectedIPv4 string
	selectedIPv6 string
	selectedCDN  string
)

// PreloadConfig reads the user-visible config values from viper once.
// Must be called after config.ReadConfig and before any Find() call.
func PreloadConfig() {
	selectedLang = viper.GetString("selected.lang")
	if selectedLang == "" {
		selectedLang = "zh-CN"
	}
	selectedIPv4 = viper.GetString("selected.ipv4")
	selectedIPv6 = viper.GetString("selected.ipv6")
	selectedCDN = viper.GetString("selected.cdn")
}
