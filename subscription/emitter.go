package subscription

import (
	"log"
	"strconv"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"github.com/mikezzb/steam-trading-shared/database/repository"
)

// Event emitter pattern
type NotificationEmitter struct {
	subRepo  *repository.SubscriptionRepository
	itemRepo *repository.ItemRepository

	notifer *Notifier

	// item name + rarity -> subscriptions
	itemRaritySubs map[string][]*ParsedSubscription
	// item name -> min price of item
	itemPrices map[string]float64
}

func NewNotificationEmitter(subRepo *repository.SubscriptionRepository, itemRepo *repository.ItemRepository, config *NotifierConfig) *NotificationEmitter {
	return &NotificationEmitter{
		subRepo:        subRepo,
		itemRepo:       itemRepo,
		notifer:        NewNotifier(config),
		itemRaritySubs: make(map[string][]*ParsedSubscription),
		itemPrices:     make(map[string]float64),
	}
}

func (e *NotificationEmitter) Init() {
	// get all subscriptions
	subs, err := e.subRepo.GetAll()
	if err != nil {
		log.Fatalf("NotificationEmitter.Init: %v", err)
		return
	}

	// group subscriptions by item name
	for _, sub := range subs {
		key := getItemRarityKey(sub.Name, sub.Rarity)
		e.itemRaritySubs[key] = append(e.itemRaritySubs[key], GetParsedSubscription(&sub))
	}

	// get all items
	items, err := e.itemRepo.GetAll()
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
				listingMessage = GetListingMessage(listing)
			}
			e.notifer.Notify(sub.Subscription.NotiType, sub.Subscription.NotiId, listingMessage)
		}
	}
}
