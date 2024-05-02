package subscription

import (
	"log"
	"strconv"

	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database/model"
	"github.com/mikezzb/steam-trading-shared/database/repository"
)

// Event emitter pattern
type NotificationEmitter struct {
	notifer *Notifier

	// item name + rarity -> subscription key -> subscription (facilates delete & update operations)
	itemRaritySubs map[string]map[string]*ParsedSubscription
	// item name + paint seed -> subscription key -> subscription
	itemPaintSeedSubs map[string]map[string]*ParsedSubscription
	// item name -> min price of item
	itemPrices map[string]float64
}

func NewNotificationEmitter(config *NotifierConfig) *NotificationEmitter {
	emitter := &NotificationEmitter{
		notifer:           NewNotifier(config),
		itemRaritySubs:    make(map[string]map[string]*ParsedSubscription),
		itemPaintSeedSubs: make(map[string]map[string]*ParsedSubscription),
		itemPrices:        make(map[string]float64),
	}
	return emitter
}

func (e *NotificationEmitter) Init(repos repository.RepoFactory) {
	subRepo := repos.GetSubscriptionRepository()
	itemRepo := repos.GetItemRepository()
	// get all subscriptions
	subs, err := subRepo.GetAll()
	if err != nil {
		log.Fatalf("NotificationEmitter.Init: %v", err)
		return
	}

	// group subscriptions by item name
	for _, sub := range subs {
		e.addSub(&sub)
	}

	// get all items
	items, err := itemRepo.GetAll()
	if err != nil {
		log.Fatalf("NotificationEmitter.Init: %v", err)
		return
	}

	// group items by name
	for _, item := range items {
		// get the lowest market price
		bestPrice := shared.GetFreshBestPrice(&item, shared.FRESH_PRICE_DURATION)
		priceFloat, _ := strconv.ParseFloat(bestPrice.Price.String(), 64)
		e.itemPrices[item.Name] = priceFloat
	}
}

func (e *NotificationEmitter) EmitListing(listing *model.Listing) {
	key := getItemRarityKey(listing.Name, listing.Rarity)
	// find all subscriptions for this item & rarity
	subs := e.itemRaritySubs[key]
	var listingMessage string
	for _, sub := range subs {
		// check if price exceeds the subscription config
		if e.IsPriceMatch(listing.Price.String(), sub) {
			// notify user
			if listingMessage == "" {
				listingMessage = GetListingMessage(listing, e.itemPrices[listing.Name])
			}
			e.notifer.Notify(sub.Subscription.NotiType, sub.Subscription.NotiId, listingMessage)
		}
	}

	// find all subscriptions for this item & paint seed
	key = getItemPaintSeedKey(listing.Name, listing.PaintSeed)
	subs = e.itemPaintSeedSubs[key]
	for _, sub := range subs {
		if e.IsPriceMatch(listing.Price.String(), sub) {
			if listingMessage == "" {
				listingMessage = GetListingMessage(listing, e.itemPrices[listing.Name])
			}
			e.notifer.Notify(sub.Subscription.NotiType, sub.Subscription.NotiId, listingMessage)
		}
	}
}

func (e *NotificationEmitter) EmitListings(listings []model.Listing) {
	for _, listing := range listings {
		e.EmitListing(&listing)
	}
}

// Dangerous to use pointer:
// when I create a parsed sub, the sub pointer is from the & of a range result, which got overwritten, so the pointer points to the same sub always

func (e *NotificationEmitter) addSub(sub *model.Subscription) {
	subKey := GetSubKey(sub)
	parsedSub := GetParsedSubscription(sub)
	// add rarities
	for _, rarity := range sub.Rarities {
		key := getItemRarityKey(sub.Name, rarity)
		// if no sub on this item rarity yet, create a map
		if _, ok := e.itemRaritySubs[key]; !ok {
			e.itemRaritySubs[key] = make(map[string]*ParsedSubscription)
		}
		// add sub to the maps
		e.itemRaritySubs[key][subKey] = parsedSub
	}
	// add paint seeds
	for _, paintSeed := range sub.PaintSeeds {
		key := getItemPaintSeedKey(sub.Name, paintSeed)
		// if no sub on this item rarity yet, create a map
		if _, ok := e.itemPaintSeedSubs[key]; !ok {
			e.itemPaintSeedSubs[key] = make(map[string]*ParsedSubscription)
		}
		// add sub to the maps
		e.itemPaintSeedSubs[key][subKey] = parsedSub
	}
}

func (e *NotificationEmitter) DelSub(sub *model.Subscription) {
	for _, rarity := range sub.Rarities {
		key := getItemRarityKey(sub.Name, rarity)
		subKey := GetSubKey(sub)
		delete(e.itemRaritySubs[key], subKey)
	}
	for _, paintSeed := range sub.PaintSeeds {
		key := getItemPaintSeedKey(sub.Name, paintSeed)
		subKey := GetSubKey(sub)
		delete(e.itemPaintSeedSubs[key], subKey)
	}
}

func (e *NotificationEmitter) UpdateSub(sub *model.Subscription) {
	e.DelSub(sub)
	e.addSub(sub)
}

func (e *NotificationEmitter) SubChangeStreamHandler(data interface{}, operationType string) {
	sub, _ := data.(*model.Subscription)
	// Find the sub in the map
	switch operationType {
	case "insert":
		e.addSub(sub)
	case "delete":
		e.DelSub(sub)
	case "update":
		e.UpdateSub(sub)
	default:
		log.Fatalf("NotificationEmitter.EmitSub: invalid operation type")
	}
}

func (e *NotificationEmitter) ListingChangeStreamHandler(data interface{}, operationType string) {
	listing, _ := data.(*model.Listing)
	switch operationType {
	case "insert":
		e.EmitListing(listing)
	case "delete":
		// do nothing
	case "update":
		e.EmitListing(listing)
	default:
		log.Fatalf("NotificationEmitter.EmitListing: invalid operation type")
	}
}
func (e *NotificationEmitter) IsPriceMatch(price string, sub *ParsedSubscription) bool {
	minPrice, ok := e.itemPrices[sub.Subscription.Name]
	priceFloat, _ := strconv.ParseFloat(price, 64)

	// if price is less than current min price, update item price
	if priceFloat < minPrice || !ok {
		e.itemPrices[sub.Subscription.Name] = priceFloat
		return true
	}

	var maxPriceMatch float64
	if sub.PremiumPerc != -1 {
		maxPriceMatch = minPrice * (1 + sub.PremiumPerc)
	} else {
		maxPriceMatch = minPrice + sub.Premium
	}

	return priceFloat <= maxPriceMatch
}
