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

func getItemPaintSeedKey(itemName string, paintSeed int) string {
	return fmt.Sprintf("%s_%d", itemName, paintSeed)
}

func GetSubKey(sub *model.Subscription) string {
	return sub.ID.Hex()
}

func GetListingMessage(listing *model.Listing, minPrice float64) string {
	return fmt.Sprintf(
		"🌸 NEW LISTING 🌸\nName: %s\nTier: %s (#%d)\nPrice: %s (Min: %.1f)\nLink: %s",
		listing.Name,
		listing.Rarity,
		listing.PaintSeed,
		listing.Price,
		minPrice,
		shared.GetListingUrl(listing),
	)
}

type ParsedSubscription struct {
	Premium      float64
	PremiumPerc  float64
	Subscription model.Subscription
}

func GetParsedSubscription(sub *model.Subscription) *ParsedSubscription {
	pSub := &ParsedSubscription{
		Subscription: *sub,
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
