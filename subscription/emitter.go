package subscription

import (
	"log"
	"strconv"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"github.com/mikezzb/steam-trading-shared/database/repository"
)

// Event emitter pattern
type NotificationEmitter struct {
	notifer *Notifier

	// item name + rarity -> subscription key -> subscription (facilates delete & update operations)
	itemRaritySubs map[string]map[string]*ParsedSubscription
	// item name -> min price of item
	itemPrices map[string]float64
}

func NewNotificationEmitter(config *NotifierConfig) *NotificationEmitter {
	emitter := &NotificationEmitter{
		notifer:        NewNotifier(config),
		itemRaritySubs: make(map[string]map[string]*ParsedSubscription),
		itemPrices:     make(map[string]float64),
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
		priceFloat, _ := strconv.ParseFloat(item.LowestMarketPrice, 64)
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
		if e.IsPriceMatch(listing.Price, sub) {
			// notify user
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
	key := getItemRarityKey(sub.Name, sub.Rarity)
	subKey := GetSubKey(sub)

	// add sub to the map
	if _, ok := e.itemRaritySubs[key]; !ok {
		e.itemRaritySubs[key] = make(map[string]*ParsedSubscription)
	}
	e.itemRaritySubs[key][subKey] = GetParsedSubscription(sub)
}

func (e *NotificationEmitter) SubChangeStreamHandler(data interface{}, operationType string) {
	sub, _ := data.(*model.Subscription)
	// Find the sub in the map
	switch operationType {
	case "insert":
		e.addSub(sub)
	case "delete":
		// find the sub by key
		key := getItemRarityKey(sub.Name, sub.Rarity)
		subKey := GetSubKey(sub)
		delete(e.itemRaritySubs[key], subKey)
	case "update":
		// find the sub by key
		key := getItemRarityKey(sub.Name, sub.Rarity)
		subKey := GetSubKey(sub)
		e.itemRaritySubs[key][subKey] = GetParsedSubscription(sub)
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
	minPrice := e.itemPrices[sub.Subscription.Name]
	priceFloat, _ := strconv.ParseFloat(price, 64)

	// if price is less than current min price, update item price
	if priceFloat < minPrice {
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
