package utils

import "fmt"

func ConvertSteamId(steamId64 int64) string {
	universe := (steamId64 >> 56) & 0xFF
	lowerBit := (steamId64 - 76561197960265728) & 1
	higherBits := (steamId64 - 76561197960265728 - lowerBit) / 2

	return fmt.Sprintf("STEAM_%d:%d:%d", universe, lowerBit, higherBits)
}
