package shared

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	STAT_TRAK_LABEL    = "StatTrak™ "
	STAR_LEBEL         = "★ "
	BUFF_IDS_PATH      = "shared/data/buff/buffids.json"
	IGXE_IDS_PATH      = "shared/data/igxe/igxeids.json"
	RARE_PATTERNS_PATH = "shared/data/items/rare_patterns.json"
)

const (
	MARKET_NAME_BUFF  = "buff"
	MARKET_NAME_STEAM = "steam"
	MARKET_NAME_UU    = "uu"
	MARKET_NAME_IGXE  = "igxe"
)

const (
	BUFF_ITEM_PREVIEW_BASE_URL = "https://buff.163.com/goods"
	BUFF_CS_APPID              = "730"
)

const (
	SECRET_TELEGRAM_TOKEN = "telegramToken"
)

const (
	MAX_NUM_STR = "99999999"
)

var MAX_DECIMAL128, _ = primitive.ParseDecimal128(MAX_NUM_STR)

// configs
const (
	FRESH_PRICE_DURATION = 60 * time.Minute
)

var WEAR_LEVELS = []string{"Factory New", "Minimal Wear", "Field-Tested", "Well-Worn", "Battle-Scarred"}

var ITEM_MARKET_NAMES = []string{MARKET_NAME_BUFF, MARKET_NAME_STEAM, MARKET_NAME_UU, MARKET_NAME_IGXE}

var ITEM_FIXED_VAL_FILTER_KEYS = []string{"name", "category", "skin", "exterior"}

var buffIds = map[string]int{}
var igxeIds = map[string]int{}
var rarePatternMap = RarePatternMap{}
var sharedBasePath string

// syncs
var buffIdOnce sync.Once
var igxeIdOnce sync.Once
var rarePatternOnce sync.Once

func init() {
	_, filename, _, _ := runtime.Caller(0)
	sharedBasePath = filepath.Dir(filepath.Dir(filename))
}

func loadJSON(path string, data interface{}) error {
	// construct absolute file path
	filePath := filepath.Join(sharedBasePath, path)

	// load json file
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// unmarshal json
	if err := json.Unmarshal(jsonData, data); err != nil {
		return err
	}
	return nil
}

// GetBuffIds returns the map of item name to buff id
func GetBuffIds() map[string]int {
	buffIdOnce.Do(func() {
		if err := loadJSON(BUFF_IDS_PATH, &buffIds); err != nil {
			panic(err)
		}
	})
	return buffIds
}

// GetIgxeIds returns the map of item name to igxe id
func GetIgxeIds() map[string]int {
	igxeIdOnce.Do(func() {
		if err := loadJSON(IGXE_IDS_PATH, &igxeIds); err != nil {
			panic(err)
		}
	})
	return igxeIds
}

func GetRarePatterns() RarePatternMap {
	rarePatternOnce.Do(func() {
		rarePatternDb := RarePatternDB{}
		if err := loadJSON(RARE_PATTERNS_PATH, &rarePatternDb); err != nil {
			panic(err)
		}
		// format rarePatternDb to rarePatternMap
		rarePatternMap = RarePatternMap{}
		for itemName, tiers := range rarePatternDb {
			rarePatternMap[itemName] = map[int]string{}
			for tier, seeds := range tiers {
				for _, seed := range seeds {
					rarePatternMap[itemName][seed] = tier
				}
			}
		}
	})
	return rarePatternMap
}
