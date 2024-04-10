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

var WEAR_LEVELS = [5]string{"Factory New", "Minimal Wear", "Field-Tested", "Well-Worn", "Battle-Scarred"}

var BuffIds = map[string]int{}
var once sync.Once
var sharedBasePath string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	sharedBasePath = filepath.Dir(filepath.Dir(filename))
}

func loadJSON(path string) {
	// construct absolute file path
	filePath := filepath.Join(sharedBasePath, path)

	// load json file
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// unmarshal json
	var data map[string]int
	if err := json.Unmarshal(jsonData, &data); err != nil {
		panic(err)
	}
	BuffIds = data
}

// GetBuffIds returns the map of item name to buff id
func GetBuffIds() map[string]int {
	once.Do(func() {
		loadJSON(BUFF_IDS_PATH)
	})
	return BuffIds
}
