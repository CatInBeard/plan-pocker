package main

import (
	"encoding/json"
	"fmt"
	"shared/logger"
	"shared/pubsub"
	"time"
)

const PUBLISH_GAME_STATE_REDIS_CHANEL = "GAME_UPDATE"

type Player struct {
	Name   string `json:"userName"`
	UID    string `json:"uid"`
	GameId string `json:"gameId"`
	Vote   int    `json:"vote"`
}

type GameStatusUpdate struct {
	GameId          string   `json:"gameId"`
	Players         []Player `json:"players"`
	Deck            []int    `json:"deck"`
	AllowCustomDeck bool     `json:"allowCustomDeck"`
	Vote            float64  `json:"vote"`
}

type CbVoters struct {
	Action string   `json:"action"`
	Voters []Player `json:"voters"`
	Vote   float64  `json:"vote"`
}

type CbDeck struct {
	Action      string `json:"action"`
	Deck        []int  `json:"deck"`
	AllowCustom bool   `json:"allowCustom"`
}

func SubscribeGameUpdates(handler func(string, interface{})) {

	client := pubsub.GetPubSubClient()

	logger.Log(logger.INFO, "[WSS-005] Start read from channel "+PUBLISH_GAME_STATE_REDIS_CHANEL, fmt.Sprintf("Channel:%s", PUBLISH_GAME_STATE_REDIS_CHANEL))

	for {
		messageChannel, err := client.Subscribe(PUBLISH_GAME_STATE_REDIS_CHANEL)
		if err != nil {
			logger.Log(logger.ERROR, "[WSS-001] Failed to subscribe to channel", fmt.Sprintf("Failed to subscribe to channel %s: %+v", PUBLISH_GAME_STATE_REDIS_CHANEL, err))
			time.Sleep(5 * time.Second)
			continue
		}
		defer client.Unsubscribe(PUBLISH_GAME_STATE_REDIS_CHANEL)

		for msg := range messageChannel {
			var subscriptionMsg GameStatusUpdate
			if err := json.Unmarshal([]byte(msg.Payload), &subscriptionMsg); err != nil {
				logger.Log(logger.ERROR, "[WSS-003] Error unmarshalling message", fmt.Sprintf("Error unmarshalling message: %", err.Error()))
				continue
			}
			logger.Log(logger.DEBUG, "[WSS-004] Send mesaage to susbcribers", fmt.Sprintf("Message %+v", subscriptionMsg))
			processForSend(subscriptionMsg, handler)
		}

		logger.Log(logger.WARNING, "[WSS-002] channel closed, attempting to reconnect", "")
		time.Sleep(3 * time.Second)
	}
}

func processForSend(msg GameStatusUpdate, handler func(string, interface{})) {
	handler(msg.GameId, cbVotesrFromGameStatusUpdate(msg))
	handler(msg.GameId, cbDeckFromGameStatusUpdate(msg))
	if msg.Vote != 0 {
		handler(msg.GameId, cbDeckFromGameStatusUpdate(msg))
	}
}

func cbVotesrFromGameStatusUpdate(gs GameStatusUpdate) CbVoters {
	return CbVoters{
		Action: "voters",
		Vote:   gs.Vote,
		Voters: gs.Players,
	}
}

func cbDeckFromGameStatusUpdate(gs GameStatusUpdate) CbDeck {
	return CbDeck{
		Action:      "deck",
		Deck:        gs.Deck,
		AllowCustom: gs.AllowCustomDeck,
	}
}
