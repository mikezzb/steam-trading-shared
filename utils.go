package shared

import (
	"strconv"
	"time"
)

var STAT_TRAK_LABEL_LEN = len(STAT_TRAK_LABEL)
var STAR_LABEL_LEN = len(STAR_LEBEL)

// FormatItemName formats the item name with wear and StatTrak™ label
func FormatItemName(name, wear string, isStatTrak bool) (formattedName string) {
	// note: some items do not have wear levels
	formattedName = name
	if isStatTrak {
		formattedName = STAT_TRAK_LABEL + name
	}
	if wear == "" {
		return
	}

	formattedName += " (" + wear + ")"

	// sanity check on name, if it DNE, try adding ★
	buffIds := GetBuffIds()
	if _, ok := buffIds[formattedName]; !ok {
		formattedName = STAR_LEBEL + formattedName
	}
	return
}

// ExtractBaseItemName extracts the base item name from the formatted name
func ExtractBaseItemName(name string) (baseName string) {
	baseName = name
	// remove star label if exists
	if len(baseName) > STAR_LABEL_LEN && baseName[:STAR_LABEL_LEN] == STAR_LEBEL {
		baseName = baseName[STAR_LABEL_LEN:]
	}
	// remove StatTrak label if exists
	if len(baseName) > STAT_TRAK_LABEL_LEN && baseName[:STAT_TRAK_LABEL_LEN] == STAT_TRAK_LABEL {
		baseName = baseName[STAT_TRAK_LABEL_LEN:]
	}
	// remove the wear level in the ending ()
	if baseName[len(baseName)-1] == ')' {
		// find the last (
		i := len(baseName) - 1
		for ; i >= 0; i-- {
			if baseName[i] == '(' {
				break
			}
		}
		// -1 to remove the space before (
		baseName = baseName[:i-1]
	}
	return
}

// Shanghai timezone
func GetTimestampNow() string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	currTime := time.Now().In(loc)
	timestamp := currTime.UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(timestamp, 10)
}
