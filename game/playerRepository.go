package main

import (
	"fmt"
	"shared/cache"
	"shared/logger"
	"strconv"
	"time"
)

type Player struct {
	Name   string `json:"userName"`
	UID    string `json:"uid"`
	GameId string `json:"gameId"`
	Vote   int    `json:"vote"`
}

type PlayerRepository struct {
	cache        cache.CacheClient
	cacheTimeout time.Duration
}

const PLAYER_CACHE_PATTERN string = "game_%s_player_%s"

func NewPlayerRepository() *PlayerRepository {

	cacheLiveDurationString := GetSetting(STAY_CONNECTED_PLAYER_SETTING)
	cacheLiveDuration, _ := strconv.Atoi(cacheLiveDurationString)

	return &PlayerRepository{
		cache:        cache.GetCacheClient(),
		cacheTimeout: time.Duration(cacheLiveDuration) * time.Second,
	}
}

func (r *PlayerRepository) SetPlayer(player Player) error {

	key := fmt.Sprintf(PLAYER_CACHE_PATTERN, player.GameId, player.UID)
	err := r.cache.SetStructValue(
		key,
		player,
		r.cacheTimeout,
	)

	if err != nil {
		logger.Log(logger.ERROR, "[PRE-001] Failed to set value to cache", fmt.Sprintf("Shortlink: %s, player: %+v, Error: %s", player.GameId, player, err.Error()))
	}

	logger.Log(logger.DEBUG, "[PRE-007] Add player to cache", fmt.Sprintf("Shortlink: %s, player: %+v, Key: %s", player.GameId, player, key))

	return err
}

func (r *PlayerRepository) GetPlayer(gameId, UID string) (*Player, error) {
	return r.GetPlayerFullKey(fmt.Sprintf(PLAYER_CACHE_PATTERN, gameId, UID))
}

func (r *PlayerRepository) GetPlayerFullKey(key string) (*Player, error) {
	var player Player
	err := r.cache.GetStructValue(key, &player)
	if err != nil {
		logger.Log(logger.ERROR, "[PRE-002] Failed to get value to cache", fmt.Sprintf("Key: %s, Error: %s", key, err.Error()))
		return nil, err
	}
	logger.Log(logger.DEBUG, "[PRE-005] Get player by full key", fmt.Sprintf("Key: %s, Player: %+v", key, player))
	return &player, err
}

func (r *PlayerRepository) GetCachedPlayerKeys(shortLink string) ([]string, error) {
	pattern := fmt.Sprintf("%s*", PLAYER_CACHE_PATTERN, shortLink)
	playerKeys, err := r.cache.GetKeysByPattern(pattern)

	if err != nil {
		logger.Log(logger.ERROR, "[PRE-003] Failed to get cached keys", fmt.Sprintf("Shortlink: %s, Error: %s", shortLink, err.Error()))
	}

	logger.Log(logger.DEBUG, "[PRE-006] Get player keys", fmt.Sprintf("shortLink: %s, Keys: %+v, pattern: %s", shortLink, playerKeys, pattern))

	return playerKeys, err
}

func (r *PlayerRepository) GetCachedPlayers(shortLink string) ([]Player, error) {

	playerKeys, err := r.GetCachedPlayerKeys(shortLink)
	if err != nil {
		return nil, err
	}

	var players []Player

	for _, key := range playerKeys {
		player, err := r.GetPlayerFullKey(key)
		if err != nil {
			logger.Log(logger.ERROR, "[PRE-004] Failed to get player data", fmt.Sprintf("Key: %s, Error: %s", key, err.Error()))
			continue
		}

		players = append(players, *player)
	}

	return players, nil
}
