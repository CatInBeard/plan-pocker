package main

import (
	"fmt"
	"shared/cache"
	"shared/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type GameState struct {
	ShortLink  string
	VoteStatus float64
}

type GameStateRepository struct {
	cache        cache.CacheClient
	cacheTimeout time.Duration
}

const CACHE_GAME_STATUS_SHORT_LINK_PREFIX string = "game_status_by_short_"

func NewGameStateRepository() *GameStateRepository {

	cacheLiveDurationString := GetSetting(GAME_STATUS_CACHE_LIVE_TIMEOUT_SETTING)
	cacheLiveDuration, _ := strconv.Atoi(cacheLiveDurationString)

	return &GameStateRepository{
		cache:        cache.GetCacheClient(),
		cacheTimeout: time.Duration(cacheLiveDuration) * time.Second,
	}
}

func (r *GameStateRepository) SetGameState(gameState GameState) error {
	err := r.cache.SetStructValue(CACHE_GAME_STATUS_SHORT_LINK_PREFIX+gameState.ShortLink, gameState, r.cacheTimeout)
	logger.Log(logger.DEBUG, "[GSR-002] Set game state to cache", fmt.Sprintf("State: %+v", gameState))
	if err != nil {
		logger.Log(logger.ERROR, "[GSR-001] Failed to set value to cache", fmt.Sprintf("Shortlink: %s, vote: %f, Error: %s", gameState.ShortLink, gameState.VoteStatus, err.Error()))
	}
	return err
}

func (r *GameStateRepository) GetGameState(shortLink string) (*GameState, error) {
	var gameState GameState
	err := r.cache.GetStructValue(CACHE_GAME_STATUS_SHORT_LINK_PREFIX+shortLink, &gameState)
	if err != nil {
		if err != redis.Nil {
			logger.Log(logger.ERROR, "[GSR-002] Failed to get game state from cache", fmt.Sprintf("Shortlink: %s, Error: %s", shortLink, err.Error()))
		}
		return nil, err
	}
	return &gameState, err
}
