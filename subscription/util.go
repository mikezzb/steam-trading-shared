package subscription

import (
	"fmt"
	"log"
	"strconv"

	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database/model"
)

func getItemRarityKey(itemName, rarity string) string {
	return itemName + "_" + rarity
}

func GetSubKey(sub *model.Subscription) string {
	return fmt.Sprintf(
		"%s_%s_%s_%s_%s",
		sub.Name,
		sub.Rarity,
		sub.MaxPremium,
		sub.NotiId,
		sub.NotiType,
	)
}

func GetListingMessage(listing *model.Listing, minPrice float64) string {
	return fmt.Sprintf(
		"🌸 NEW LISTING 🌸\nName: %s\nTier: %s\nPrice: %s (Min: %.1f)\nLink: %s",
		listing.Name,
		listing.Rarity,
		listing.Price,
		minPrice,
		shared.GetListingUrl(listing),
	)
}

type ParsedSubscription struct {
	Premium      float64
	PremiumPerc  float64
	Subscription *model.Subscription
}

func GetParsedSubscription(sub *model.Subscription) *ParsedSubscription {
	pSub := &ParsedSubscription{
		Subscription: sub,
		Premium:      -1,
		PremiumPerc:  -1,
	}

	// check if the subscription premium is a percentage
	if sub.MaxPremium[len(sub.MaxPremium)-1] == '%' {
		// convert to float
		perc, err := strconv.ParseFloat(sub.MaxPremium[:len(sub.MaxPremium)-1], 64)
		if err != nil {
			log.Fatalf("GetParsedSubscription: %v", err)
		}
		pSub.PremiumPerc = perc / 100
	} else {
		permium, err := strconv.ParseFloat(sub.MaxPremium, 64)
		if err != nil {
			log.Fatalf("GetParsedSubscription: %v", err)
		}
		pSub.Premium = permium
	}

	return pSub
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
