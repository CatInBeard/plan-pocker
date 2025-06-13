package main

import (
	"os"
	"sync"
)

type Setting struct {
	DefaultValue string
	EnvVar       string
}

const STAY_CONNECTED_PLAYER_SETTING string = "STAY_CONNECTED_PLAYER"
const REPOSITORY_CACHE_TIMEOUT_SETTING string = "REPOSITORY_CACHE_TIMEOUT"
const GAME_STATUS_CACHE_LIVE_TIMEOUT_SETTING string = "GAME_STATUS_CACHE_TIMEOUT"

var settings = map[string]Setting{
	STAY_CONNECTED_PLAYER_SETTING:          {"20", "STAY_CONNECTED_PLAYER"},
	REPOSITORY_CACHE_TIMEOUT_SETTING:       {"0", "REPOSITORY_CACHE_TIMEOUT"},
	GAME_STATUS_CACHE_LIVE_TIMEOUT_SETTING: {"0", "GAME_STATUS_CACHE_TIMEOUT"},
}

var settingsCache = make(map[string]string)
var mu sync.RWMutex

func GetSetting(key string) string {
	mu.RLock()
	defer mu.RUnlock()

	if value, found := settingsCache[key]; found {
		return value
	}

	setting, exists := settings[key]
	if !exists {
		return ""
	}

	if envValue, exists := os.LookupEnv(setting.EnvVar); exists {
		if _, found := settingsCache[key]; !found {
			mu.Lock()
			settingsCache[key] = envValue
			mu.Unlock()
		}
		return envValue
	}

	return setting.DefaultValue
}
