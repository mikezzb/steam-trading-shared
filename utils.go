package shared

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mikezzb/steam-trading-shared/database/model"
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

func GetListingUrl(listing *model.Listing) string {
	switch listing.Market {
	default:
		// default as buff
		buffId := GetBuffIds()[listing.Name]
		params := url.Values{}
		params.Add("appid", BUFF_CS_APPID)
		params.Add("classid", listing.ClassId)
		params.Add("instanceid", listing.InstanceId)
		params.Add("assetid", listing.AssetId)

		return fmt.Sprintf("%s/%d?%s", BUFF_ITEM_PREVIEW_BASE_URL, buffId, params.Encode())

	}
}

// Shanghai timezone
func GetTimestampNow() string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	currTime := time.Now().In(loc)
	timestamp := currTime.UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(timestamp, 10)
}

func GetUnixNow() int64 {
	return time.Now().Unix()
}

func GetTier(name string, paintSeed int) string {
	baseName := ExtractBaseItemName(name)
	tiers, ok := GetRarePatterns()[baseName]
	if !ok {
		return ""
	}

	// find rarity definition
	tier, ok := tiers[paintSeed]
	if !ok {
		return ""
	}
	return tier
}

func PrintCookies(cookies []*http.Cookie, label string) {
	for _, cookie := range cookies {
		log.Printf("[%s] Cookie: %v\n", label, cookie)
	}
}

func RandomFloat(min, max int) float32 {
	return float32(min) + rand.Float32()*(float32(max)-float32(min))
}

func GetRandomSleepDuration(min, max int) time.Duration {
	randVal := RandomFloat(min, max)
	return time.Duration(randVal) * time.Second
}

func RandomSleep(min, max int) {
	randVal := GetRandomSleepDuration(min, max)
	log.Printf("Sleeping for %v\n", randVal)
	time.Sleep(randVal)
}

func GetMarketPrice(item *model.Item, marketName string) *model.MarketPrice {
	switch marketName {
	case MARKET_NAME_BUFF:
		return &item.BuffPrice
	case MARKET_NAME_STEAM:
		return &item.SteamPrice
	case MARKET_NAME_IGXE:
		return &item.IgxePrice
	case MARKET_NAME_UU:
		return &item.UUPrice
	}
	return nil
}

// @return best price, market name
func GetBestPrice(item *model.Item) *model.MarketPrice {
	if item == nil {
		return nil
	}

	var lowestPrice *model.MarketPrice = nil

	for _, marketName := range ITEM_MARKET_NAMES {
		price := GetMarketPrice(item, marketName)
		if price == nil {
			continue
		}
		if lowestPrice == nil || price.Price < lowestPrice.Price {
			lowestPrice = price
		}
	}

	return lowestPrice
}

// @return fresh best price, market name
func GetFreshBestPrice(item *model.Item, expireDuration time.Duration) *model.MarketPrice {
	if item == nil {
		return nil
	}

	var lowestPrice *model.MarketPrice = nil
	now := time.Now()

	for _, marketName := range ITEM_MARKET_NAMES {
		price := GetMarketPrice(item, marketName)
		if price == nil {
			continue
		}

		// check if price is expired
		updatedAtTime := time.Unix(price.UpdatedAt, 0)
		if updatedAtTime.Add(expireDuration).Before(now) {
			continue
		}

		if lowestPrice == nil || price.Price < lowestPrice.Price {
			lowestPrice = price
		}
	}

	// if all prices are expired, return the lowest price
	if lowestPrice == nil {
		return GetBestPrice(item)
	}

	return lowestPrice
}
