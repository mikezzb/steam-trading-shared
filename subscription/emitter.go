package subscription

import (
	"context"
	"log"
	"strconv"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"github.com/mikezzb/steam-trading-shared/database/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// change streams
	subChangeStream *mongo.ChangeStream
}

func NewNotificationEmitter(subRepo *repository.SubscriptionRepository, itemRepo *repository.ItemRepository, config *NotifierConfig) *NotificationEmitter {
	emitter := &NotificationEmitter{
		subRepo:        subRepo,
		itemRepo:       itemRepo,
		notifer:        NewNotifier(config),
		itemRaritySubs: make(map[string][]*ParsedSubscription),
		itemPrices:     make(map[string]float64),
	}
	emitter.init()
	return emitter
}

func (e *NotificationEmitter) init() {
	// get all subscriptions
	subs, err := e.subRepo.GetAll()
	if err != nil {
		log.Fatalf("NotificationEmitter.Init: %v", err)
		return
	}

	// group subscriptions by item name
	for _, sub := range subs {
		e.addSub(&sub)
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

	// start monitoring
	e.initMonitoring()
}

func (e *NotificationEmitter) addSub(sub *model.Subscription) {
	key := getItemRarityKey(sub.Name, sub.Rarity)
	e.itemRaritySubs[key] = append(e.itemRaritySubs[key], GetParsedSubscription(sub))
}

func (e *NotificationEmitter) initMonitoring() {
	// monitor any update or insert to the subscription collection
	changeStreamOpts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	changeStream, err := e.subRepo.SubCol.Watch(context.Background(), mongo.Pipeline{}, changeStreamOpts)
	if err != nil {
		log.Fatalf("NotificationEmitter.initMonitoring: %v", err)
		return
	}
	e.subChangeStream = changeStream

	go e.monitorSubChanges()
}

func (e *NotificationEmitter) monitorSubChanges() {
	defer e.subChangeStream.Close(context.Background())

	// loop to monitor changes
	for e.subChangeStream.Next(context.Background()) {

		var subChange struct {
			FullDocument  model.Subscription `bson:"fullDocument"`
			OperationType string             `bson:"operationType"`
		}
		if err := e.subChangeStream.Decode(&subChange); err != nil {
			log.Fatalf("NotificationEmitter.monitorSubChanges: %v", err)
			return
		}

		switch subChange.OperationType {
		case "delete":
			// remove the subscription
		case "update":
			// update the subscription
		case "insert":
			// add the subscription
			e.addSub(&subChange.FullDocument)
		default:
			log.Fatalf("NotificationEmitter.monitorSubChanges: unknown operation type %v", subChange.OperationType)
		}

		// update the subscription
		key := getItemRarityKey(subChange.FullDocument.Name, subChange.FullDocument.Rarity)

		//

		e.itemRaritySubs[key] = append(e.itemRaritySubs[key], GetParsedSubscription(&subChange.FullDocument))
	}

	if err := e.subChangeStream.Err(); err != nil {
		log.Fatalf("NotificationEmitter.monitorSubChanges: %v", err)
		return
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
