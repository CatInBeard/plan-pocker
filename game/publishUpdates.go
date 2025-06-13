package main

import (
	"fmt"
	"reflect"
	"shared/cache"
	"shared/logger"
	"shared/pubsub"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const PUBLISH_GAME_STATE_REDIS_CHANEL = "GAME_UPDATE"
const GAME_LAST_UPDATE_PREFIX = "game_last_update_"

type GameStatusUpdate struct {
	GameId          string   `json:"gameId"`
	Players         []Player `json:"players"`
	Deck            []int    `json:"deck"`
	AllowCustomDeck bool     `json:"allowCustomDeck"`
	Vote            float64  `json:"vote"`
}

func (g *GameStatusUpdate) prepare() {
	if g.Vote == 0 {
		for i := range g.Players {
			if g.Players[i].Vote > 0 {
				g.Players[i].Vote = 1
			}
		}
	}
}

func (g *GameStatusUpdate) Equals(other GameStatusUpdate) bool {
	return g.GameId == other.GameId &&
		reflect.DeepEqual(g.Players, other.Players) &&
		reflect.DeepEqual(g.Deck, other.Deck) &&
		g.AllowCustomDeck == other.AllowCustomDeck &&
		g.Vote == other.Vote
}

func UpdateGameState(gameStatusUpdate GameStatusUpdate) {
	pubSubClient := pubsub.GetPubSubClient()
	cacheClient := cache.GetCacheClient()

	gameStatusUpdate.prepare()

	var gameStatusUpdateCached GameStatusUpdate

	err := cacheClient.GetStructValue(GAME_LAST_UPDATE_PREFIX+gameStatusUpdate.GameId, &gameStatusUpdateCached)

	if err == nil {
		if gameStatusUpdate.Equals(gameStatusUpdateCached) {
			return
		}
	} else if err != redis.Nil {
		logger.Log(logger.ERROR, "[PUC-001] Failed to get value from cache", fmt.Sprintf("Key: %s, Error: %s", GAME_LAST_UPDATE_PREFIX+gameStatusUpdate.GameId, err.Error()))
	}

	updateDelay := GetSetting(GAME_RESEND_WITHOUT_UPDATE_SETTING)
	updateDelaySeconds, _ := strconv.Atoi(updateDelay)

	cacheClient.SetStructValue(GAME_LAST_UPDATE_PREFIX+gameStatusUpdate.GameId, gameStatusUpdate, time.Duration(updateDelaySeconds)*time.Second)

	logger.Log(logger.DEBUG, "[PUP-001] Publish game update id: "+gameStatusUpdate.GameId, fmt.Sprintf("Game: %+v", gameStatusUpdate))
	pubSubClient.Publish(PUBLISH_GAME_STATE_REDIS_CHANEL, gameStatusUpdate)
}

func CalculateGameStateByGameId(gameId string) {

	playerRepository := NewPlayerRepository()

	players, err := playerRepository.GetCachedPlayers(gameId)

	if err != nil {
		logger.Log(logger.ERROR, "[PUPR-001] Failed to get players from game "+gameId, fmt.Sprintf("Error: %s", err.Error()))
	}

	gameRepository := NewGameRepository()

	game, err := gameRepository.SelectByShortLink(gameId)
	if err != nil {
		logger.Log(logger.ERROR, "[PUGR-001] Failed to get game "+gameId, fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	gameStateRepository := NewGameStateRepository()

	gameState, err := gameStateRepository.GetGameState(gameId)
	var vote float64
	if err != nil {
		vote = 0
		if err != redis.Nil {
			logger.Log(logger.ERROR, "[PUGSR-001] Failed to get game state "+gameId, fmt.Sprintf("Error: %s", err.Error()))
		} else {
			logger.Log(logger.WARNING, "[PUGSR-002] Game state "+gameId+" not found", fmt.Sprintf("Error: %s", err.Error()))
		}
	} else {
		vote = gameState.VoteStatus
	}

	gameStatusUpdate := GameStatusUpdate{
		GameId:          gameId,
		Players:         players,
		Deck:            game.Settings.Deck,
		AllowCustomDeck: game.Settings.AllowCustomDeck,
		Vote:            vote,
	}

	UpdateGameState(gameStatusUpdate)
}
