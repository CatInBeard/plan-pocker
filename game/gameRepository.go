package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"shared/cache"
	"shared/db"
	"shared/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Game struct {
	ID        int      `json:"id"`
	Shortlink string   `json:"shortlink"`
	Settings  Settings `json:"settings"`
}

type Settings struct {
	Deck            []int `json:"deck"`
	AllowCustomDeck bool  `json:"allowCustomDeck"`
}

type GameRepository struct {
	db           *db.DbClient
	cache        cache.CacheClient
	useReplica   bool
	cacheTimeout time.Duration
}

const CACHE_BY_SHORT_LINK_PREFIX string = "game_by_short_"

func NewGameRepository() *GameRepository {
	dbClient, replicaErr := db.GetDbClient()
	useReplica := true
	if replicaErr != nil {
		useReplica = false
	}

	cacheLiveDurationString := GetSetting(STAY_CONNECTED_PLAYER_SETTING)
	cacheLiveDuration, _ := strconv.Atoi(cacheLiveDurationString)

	return &GameRepository{
		db:           dbClient,
		useReplica:   useReplica,
		cache:        cache.GetCacheClient(),
		cacheTimeout: time.Duration(cacheLiveDuration) * time.Second,
	}
}

func (r *GameRepository) Create(game Game) error {
	settingsJSON, err := json.Marshal(game.Settings)
	if err != nil {
		logger.Log(logger.ERROR, "[GRM-001] Failed to marshal settings json", fmt.Sprintf("Shortlink: %s, settings: %v, Error: %s", game.Shortlink, game.Settings, err.Error()))
		return err
	}

	query := "INSERT INTO games (shortlink, settings) VALUES (?, ?)"
	_, err = r.db.ExecuteUpdate(query, game.Shortlink, settingsJSON)
	if err == nil {
		cacheErr := r.cache.SetStructValue(CACHE_BY_SHORT_LINK_PREFIX+game.Shortlink, game, r.cacheTimeout)
		if cacheErr != nil {
			logger.Log(logger.ERROR, "[GRC-002] Failed to set value to cache", fmt.Sprintf("Shortlink: %s, settingsJson: %s, Error: %s", game.Shortlink, settingsJSON, cacheErr.Error()))
		}
	}

	if err != nil {
		logger.Log(logger.ERROR, "[GRE-001] Failed to create new game", fmt.Sprintf("Query: %s, Shortlink: %s, settings: %s, Error: %s", query, game.Shortlink, settingsJSON, err.Error()))
	} else {
		logger.Log(logger.DEBUG, "[GRS-001] Successfully create new game", fmt.Sprintf("Query: %s, Shortlink: %s, settings: %s", query, game.Shortlink, settingsJSON))
	}

	return err
}

func (r *GameRepository) CreateOrUpdate(game Game) error {
	settingsJSON, err := json.Marshal(game.Settings)
	if err != nil {
		logger.Log(logger.ERROR, "[GRM-002] Failed to marshal settings json", fmt.Sprintf("Shortlink: %s, settings: %v, Error: %s", game.Shortlink, game.Settings, err.Error()))
		return err
	}

	query := "REPLACE INTO games (shortlink, settings) VALUES (?, ?)"
	_, err = r.db.ExecuteUpdate(query, game.Shortlink, settingsJSON)
	if err == nil {
		cacheErr := r.cache.SetStructValue(CACHE_BY_SHORT_LINK_PREFIX+game.Shortlink, game, r.cacheTimeout)
		if cacheErr != nil {
			logger.Log(logger.ERROR, "[GRC-003] Failed to set value to cache", fmt.Sprintf("Shortlink: %s, settingsJson: %s, Error: %s", game.Shortlink, settingsJSON, cacheErr.Error()))
		}
	}

	if err != nil {
		logger.Log(logger.ERROR, "[GRE-002] Failed to create new game", fmt.Sprintf("Query: %s, Shortlink: %s, settings: %s, Error: %s", query, game.Shortlink, settingsJSON, err.Error()))
	} else {
		logger.Log(logger.DEBUG, "[GRS-002] Successfully create new game", fmt.Sprintf("Query: %s, Shortlink: %s, settings: %s", query, game.Shortlink, settingsJSON))
	}

	return err
}

func (r *GameRepository) Update(game Game) error {
	settingsJSON, err := json.Marshal(game.Settings)
	if err != nil {
		logger.Log(logger.ERROR, "[GRM-003] Failed to marshal settings json", fmt.Sprintf("Shortlink: %s, settings: %v, Error: %s", game.Shortlink, game.Settings, err.Error()))
		return err
	}

	query := "UPDATE games SET shortlink = ?, settings = ? WHERE id = ?"
	_, err = r.db.ExecuteUpdate(query, game.Shortlink, settingsJSON, game.ID)
	if err == nil {
		cacheErr := r.cache.SetStructValue(CACHE_BY_SHORT_LINK_PREFIX+game.Shortlink, game, r.cacheTimeout)
		if cacheErr != nil {
			logger.Log(logger.ERROR, "[GRC-004] Failed to set value to cache", fmt.Sprintf("Shortlink: %s, settingsJson: %s, Error: %s", game.Shortlink, settingsJSON, cacheErr.Error()))
		}
	}

	if err != nil {
		logger.Log(logger.ERROR, "[GRE-003] Failed to create new game", fmt.Sprintf("Query: %s, Shortlink: %s, settings: %s, Error: %s", query, game.Shortlink, settingsJSON, err.Error()))
	} else {
		logger.Log(logger.DEBUG, "[GRS-003] Successfully create new game", fmt.Sprintf("Query: %s, Shortlink: %s, settings: %s", query, game.Shortlink, settingsJSON))
	}

	return err
}

func (r *GameRepository) DeleteByShortLink(shortlink string) error {
	query := "DELETE FROM games WHERE shortlink = ?"
	_, err := r.db.ExecuteUpdate(query, shortlink)
	if err == nil {
		cacheErr := r.cache.DeleteKey(shortlink)
		if cacheErr != nil {
			logger.Log(logger.ERROR, "[GRC-005] Failed to delete key from cache", fmt.Sprintf("Shortlink: %s, Error: %s", shortlink, err.Error()))
		}
	}

	if err != nil {
		logger.Log(logger.ERROR, "[GRE-004] Failed to delete game", fmt.Sprintf("Query: %s, Shortlink: %s, Error: %s", query, shortlink, err.Error()))
	} else {
		logger.Log(logger.DEBUG, "[GRS-004] Successfully delete game", fmt.Sprintf("Query: %s, Shortlink: %s", query, shortlink))
	}

	return err
}

func (r *GameRepository) SelectById(id int) (*Game, error) {

	query := "SELECT id, shortlink, settings FROM games WHERE id = ?"
	var rows *sql.Rows
	var err error
	if r.useReplica {
		rows, err = r.db.ExecuteRead(query, id)
	} else {
		rows, err = r.db.ExecuteReadPrimary(query, id)
	}

	if err != nil {
		if !r.useReplica {
			logger.Log(logger.ERROR, "[GRE-005] Failed get game by id from primary db", fmt.Sprintf("Query: %s, id: %d, Error: %s", query, id, err.Error()))
			return nil, err
		}
		logger.Log(logger.ERROR, "[GRE-006] Failed get game by id from replica db", fmt.Sprintf("Query: %s, id: %d, Error: %s", query, id, err.Error()))
		rows, err = r.db.ExecuteReadPrimary(query, id)
		if err != nil {
			logger.Log(logger.ERROR, "[GRE-007] Failed get game by id from primary db", fmt.Sprintf("Query: %s, id: %d, Error: %s", query, id, err.Error()))
			return nil, err
		}
	}
	defer rows.Close()

	if rows.Next() {
		var game Game
		var settingsJSON []byte
		if err := rows.Scan(&game.ID, &game.Shortlink, &settingsJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(settingsJSON, &game.Settings); err != nil {
			return nil, err
		}
		return &game, nil
	}
	return nil, sql.ErrNoRows
}

func (r *GameRepository) SelectByShortLink(shortlink string) (*Game, error) {
	var game Game

	cacheErr := r.cache.GetStructValue(CACHE_BY_SHORT_LINK_PREFIX+game.Shortlink, game)

	if cacheErr == nil {
		return &game, nil
	}

	if cacheErr != redis.Nil {
		logger.Log(logger.ERROR, "[GRC-001] Failed to get value from cache", fmt.Sprintf("Key: %s, Error: %s", shortlink, cacheErr.Error()))
	}

	query := "SELECT id, shortlink, settings FROM games WHERE shortlink = ?"
	var rows *sql.Rows
	var err error
	if r.useReplica {
		rows, err = r.db.ExecuteRead(query, shortlink)
	} else {
		rows, err = r.db.ExecuteReadPrimary(query, shortlink)
	}
	if err != nil {
		if !r.useReplica {
			logger.Log(logger.ERROR, "[GRE-008] Failed get game by id from primary db", fmt.Sprintf("Query: %s, shortlink: %s, Error: %s", query, shortlink, err.Error()))
			return nil, err
		}
		logger.Log(logger.ERROR, "[GRE-009] Failed get game by id from replica db", fmt.Sprintf("Query: %s, shortlink: %s, Error: %s", query, shortlink, err.Error()))
		rows, err = r.db.ExecuteReadPrimary(query, shortlink)
		if err != nil {
			logger.Log(logger.ERROR, "[GRE-010] Failed get game by id from primary db", fmt.Sprintf("Query: %s, shortlink: %s, Error: %s", query, shortlink, err.Error()))
			return nil, err
		}
	}
	defer rows.Close()

	if rows.Next() {
		var settingsJSON []byte
		if err := rows.Scan(&game.ID, &game.Shortlink, &settingsJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(settingsJSON, &game.Settings); err != nil {
			return nil, err
		}
		return &game, nil
	}
	return nil, sql.ErrNoRows
}

func (r *GameRepository) SelectAll() ([]Game, error) {
	query := "SELECT id, shortlink, settings FROM games"
	var rows *sql.Rows
	var err error
	if r.useReplica {
		rows, err = r.db.ExecuteRead(query)
	} else {
		rows, err = r.db.ExecuteReadPrimary(query)
	}
	if err != nil {
		if !r.useReplica {
			logger.Log(logger.ERROR, "[GRE-011] Failed get game by id from primary db", fmt.Sprintf("Query: %s,, Error: %s", query, err.Error()))
			return nil, err
		}
		logger.Log(logger.ERROR, "[GRE-012] Failed get game by id from replica db", fmt.Sprintf("Query: %s, Error: %s", query, err.Error()))
		rows, err = r.db.ExecuteReadPrimary(query)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var game Game
		var settingsJSON []byte
		if err := rows.Scan(&game.ID, &game.Shortlink, &settingsJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(settingsJSON, &game.Settings); err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

func (r *GameRepository) CreateDefaultGame(Shortlink string) error {
	settings := Settings{
		Deck:            []int{1, 2, 4, 8, 16, 32, 64, 128},
		AllowCustomDeck: true,
	}

	game := Game{
		Shortlink: Shortlink,
		Settings:  settings,
	}

	return r.Create(game)
}
