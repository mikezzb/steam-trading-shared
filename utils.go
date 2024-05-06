package shared

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Get the buff id of an item
func GetItemId(name string) string {
	return strconv.Itoa(GetBuffIds()[name])
}

func DecodeItemFullName(fullName string) (category, skin, exterior string) {
	re := regexp.MustCompile(`(.*?) \| (.*?) \((.*?)\)`)
	matches := re.FindStringSubmatch(fullName)

	if matches == nil {
		return fullName, "", ""
	}

	return matches[1], matches[2], matches[3]
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
	case MARKET_NAME_IGXE:
		return fmt.Sprintf("https://www.igxe.cn/product-%s", listing.InstanceId)
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

func GetUnixNow() int64 {
	return time.Now().Unix()
}

func GetNow() time.Time {
	return time.Now()
}

// consist to json number
func GetUnixFloat() float64 {
	return float64(GetUnixNow())
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
		return item.BuffPrice
	case MARKET_NAME_STEAM:
		return item.SteamPrice
	case MARKET_NAME_IGXE:
		return item.IgxePrice
	case MARKET_NAME_UU:
		return item.UUPrice
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
		if lowestPrice == nil || DecCompareTo(price.Price, lowestPrice.Price) == -1 {
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
		updatedAtTime := price.UpdatedAt
		if updatedAtTime.Add(expireDuration).Before(now) {
			continue
		}

		if lowestPrice == nil || DecCompareTo(price.Price, lowestPrice.Price) == -1 {
			lowestPrice = price
		}
	}

	// if all prices are expired, return the lowest price
	if lowestPrice == nil {
		return GetBestPrice(item)
	}

	return lowestPrice
}

func ConvertToUnix(timestamp string, layout string) (int64, error) {
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		return 0, err
	}

	unixTimestamp := t.Unix()

	return unixTimestamp, nil
}

// Convert a timestamp string in format "2006-01-02T15:04:05" to unix timestamp
func ConvertToUnixTimestamp(timestamp string) (int64, error) {
	layout := "2006-01-02T15:04:05"

	return ConvertToUnix(timestamp, layout)
}

func ConvertChineseDateToUnix(date string) (int64, error) {
	layout := "2006年01月02日"

	return ConvertToUnix(date, layout)
}

func ParseChineseDate(date string) (time.Time, error) {
	layout := "2006年01月02日"

	return time.Parse(layout, date)
}

// Parse a timestamp string in format "2006-01-02T15:04:05" to time.Time
func ParseDateHhmmss(date string) (time.Time, error) {
	layout := "2006-01-02T15:04:05"

	return time.Parse(layout, date)
}

func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// Compare two Decimal128, return -1 if a < b, 0 if a == b, 1 if a > b
func DecCompareTo(d1, d2 primitive.Decimal128) int {
	b1, exp1, err := d1.BigInt()
	if err != nil {
		return 0
	}
	b2, exp2, err := d2.BigInt()
	if err != nil {
		return 0
	}

	sign := b1.Sign()
	if sign != b2.Sign() {
		if b1.Sign() > 0 {
			return 1
		} else {
			return -1
		}
	}

	if exp1 == exp2 {
		return b1.Cmp(b2)
	}

	if sign < 0 {
		if exp1 < exp2 {
			return 1
		}
		return -1
	} else {
		if exp1 < exp2 {
			return -1
		}

		return 1
	}
}

// Compare two numeric strings, return -1 if a < b, 0 if a == b, 1 if a > b
func NumStrCmp(a, b string) int {
	floatA, _ := strconv.ParseFloat(a, 64)
	floatB, _ := strconv.ParseFloat(b, 64)

	if floatA < floatB {
		return -1
	} else if floatA > floatB {
		return 1
	}
	return 0
}

// Get a Decimal128 from a string, ignore error
func GetDecimal128(s string) primitive.Decimal128 {
	d128, err := primitive.ParseDecimal128(s)
	if err != nil {
		return MAX_DECIMAL128
	}
	return d128
}

func GetNowHHMMSS() string {
	return time.Now().Format("2006-01-02T15:04:05")
}

func SameTime(t1, t2 time.Time) bool {
	return t1.Unix() == t2.Unix()
}

func GetTimeBeforeDays(days int) time.Time {
	return time.Now().AddDate(0, 0, -days)
}
