package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MarketPrice struct {
	Price     primitive.Decimal128 `bson:"price" json:"price"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type Item struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Name    string `bson:"name" json:"name"`
	IconUrl string `bson:"iconUrl" json:"iconUrl"`

	// Market prices
	BuffPrice  MarketPrice `bson:"buffPrice" json:"buffPrice"`
	UUPrice    MarketPrice `bson:"uuPrice" json:"uuPrice"`
	IgxePrice  MarketPrice `bson:"igxePrice" json:"igxePrice"`
	SteamPrice MarketPrice `bson:"steamPrice" json:"steamPrice"`
}

type Listing struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`

	Name             string               `bson:"name"`
	Market           string               `bson:"market"`
	Price            primitive.Decimal128 `bson:"price"`
	PreviewUrl       string               `bson:"previewUrl"`
	GoodsId          int                  `bson:"goodsId"`
	ClassId          string               `bson:"classId"`
	AssetId          string               `bson:"assetId"`
	TradableCooldown string               `bson:"tradableCooldown"`

	PaintWear  primitive.Decimal128 `bson:"paintWear"`
	PaintIndex int                  `bson:"paintIndex"`
	PaintSeed  int                  `bson:"paintSeed"`
	Rarity     string               `bson:"rarity"`

	// Market specific ID
	InstanceId string `bson:"instanceId"`
}

// Subscription on the rare patterns of an item
type Subscription struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Name string `bson:"name"`

	// Optional, if not provided, it means subscribe to all rarity
	Rarity string `bson:"rarity"`
	// Optional, can be percentage or absolute value
	MaxPremium string `bson:"maxPremium"`

	// Alarm settings. Example: Telegram, Email
	NotiType string `bson:"notiType"`
	// Example: Telegram chat id, or email address
	NotiId string `bson:"notiId"`

	OwnerId primitive.ObjectID `bson:"ownerId"`
}

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Username string `bson:"username"`
	Password string `bson:"password"`

	SubscriptionIds []primitive.ObjectID `bson:"subscriptionIds"`
	FavItemIds      []primitive.ObjectID `bson:"favItemIds"`
	FavListingIds   []primitive.ObjectID `bson:"favListingIds"`
}

type TransactionMetadata struct {
	Market  string `bson:"market"`
	AssetId string `bson:"assetId"`
}

// Currently same as Listing
type Transaction struct {
	ID       primitive.ObjectID  `bson:"_id,omitempty"`
	Metadata TransactionMetadata `bson:"metadata"`

	Name string `bson:"name"`

	CreatedAt time.Time `bson:"createdAt"`

	Price            primitive.Decimal128 `bson:"price"`
	PreviewUrl       string               `bson:"previewUrl"`
	GoodsId          int                  `bson:"goodsId"`
	ClassId          string               `bson:"classId"`
	TradableCooldown string               `bson:"tradableCooldown"`

	PaintWear  primitive.Decimal128 `bson:"paintWear"`
	PaintIndex int                  `bson:"paintIndex"`
	PaintSeed  int                  `bson:"paintSeed"`

	Rarity string `bson:"rarity"`

	// market specific unique id
	InstanceId string `bson:"instanceId"`
}
