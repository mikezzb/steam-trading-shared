package shared

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	STAT_TRAK_LABEL    = "StatTrak™ "
	STAR_LEBEL         = "★ "
	BUFF_IDS_PATH      = "shared/data/buff/buffids.json"
	RARE_PATTERNS_PATH = "shared/data/items/rare_patterns.json"
)

const (
	MARKET_NAME_BUFF  = "buff"
	MARKET_NAME_STEAM = "steam"
	MARKET_NAME_UU    = "uu"
)

var WEAR_LEVELS = [5]string{"Factory New", "Minimal Wear", "Field-Tested", "Well-Worn", "Battle-Scarred"}

var buffIds = map[string]int{}
var rarePatternMap = RarePatternMap{}
var once sync.Once
var sharedBasePath string

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
	once.Do(func() {
		if err := loadJSON(BUFF_IDS_PATH, &buffIds); err != nil {
			panic(err)
		}
	})
	return buffIds
}

func GetRarePatterns() RarePatternMap {
	once.Do(func() {
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
