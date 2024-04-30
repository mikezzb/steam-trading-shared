package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MarketPrice struct {
	Price     primitive.Decimal128 `bson:"price" json:"price"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
}

// TODO: make the market prices omitempty

type Item struct {
	ID string `bson:"_id,omitempty" json:"_id"`

	Name    string `bson:"name" json:"name"`
	IconUrl string `bson:"iconUrl" json:"iconUrl"`

	// Market prices
	BuffPrice  *MarketPrice `bson:"buffPrice,omitempty" json:"buffPrice"`
	UUPrice    *MarketPrice `bson:"uuPrice,omitempty" json:"uuPrice"`
	IgxePrice  *MarketPrice `bson:"igxePrice,omitempty" json:"igxePrice"`
	SteamPrice *MarketPrice `bson:"steamPrice,omitempty" json:"steamPrice"`
}

type Listing struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`

	CheckedAt time.Time `bson:"checkedAt" json:"checkedAt"`

	Name             string               `bson:"name" json:"name"`
	Market           string               `bson:"market" json:"market"`
	Price            primitive.Decimal128 `bson:"price" json:"price"`
	PreviewUrl       string               `bson:"previewUrl" json:"previewUrl"`
	GoodsId          int                  `bson:"goodsId" json:"goodsId"`
	ClassId          string               `bson:"classId" json:"classId"`
	AssetId          string               `bson:"assetId" json:"assetId"`
	TradableCooldown string               `bson:"tradableCooldown" json:"tradableCooldown"`

	PaintWear  primitive.Decimal128 `bson:"paintWear" json:"paintWear"`
	PaintIndex int                  `bson:"paintIndex" json:"paintIndex"`
	PaintSeed  int                  `bson:"paintSeed" json:"paintSeed"`
	Rarity     string               `bson:"rarity" json:"rarity"`

	// Market specific ID
	InstanceId string `bson:"instanceId" json:"instanceId"`
}

// Subscription on the rare patterns of an item
type Subscription struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	Name string `bson:"name" json:"name"`

	// Optional, if not provided, it means subscribe to all rarity
	Rarity string `bson:"rarity" json:"rarity"`
	// Optional, can be percentage or absolute value
	MaxPremium string `bson:"maxPremium" json:"maxPremium"`

	// Alarm settings. Example: Telegram, Email
	NotiType string `bson:"notiType" json:"notiType"`
	// Example: Telegram chat id, or email address
	NotiId string `bson:"notiId" json:"notiId"`

	OwnerId primitive.ObjectID `bson:"ownerId" json:"ownerId"`
}

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`

	Email string `bson:"email" json:"email"`

	Role string `bson:"role" json:"role"`

	SubscriptionIds []primitive.ObjectID `bson:"subscriptionIds" json:"subscriptionIds"`
	FavItemIds      []primitive.ObjectID `bson:"favItemIds" json:"favItemIds"`
	FavListingIds   []primitive.ObjectID `bson:"favListingIds" json:"favListingIds"`
}

type TransactionMetadata struct {
	Market  string `bson:"market" json:"market"`
	AssetId string `bson:"assetId" json:"assetId"`
}

// Currently same as Listing
type Transaction struct {
	ID       primitive.ObjectID  `bson:"_id,omitempty" json:"_id"`
	Metadata TransactionMetadata `bson:"metadata" json:"metadata"`

	Name string `bson:"name" json:"name"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`

	Price            primitive.Decimal128 `bson:"price" json:"price"`
	PreviewUrl       string               `bson:"previewUrl" json:"previewUrl"`
	GoodsId          int                  `bson:"goodsId" json:"goodsId"`
	ClassId          string               `bson:"classId" json:"classId"`
	TradableCooldown string               `bson:"tradableCooldown" json:"tradableCooldown"`

	PaintWear  primitive.Decimal128 `bson:"paintWear" json:"paintWear"`
	PaintIndex int                  `bson:"paintIndex" json:"paintIndex"`
	PaintSeed  int                  `bson:"paintSeed" json:"paintSeed"`

	Rarity string `bson:"rarity" json:"rarity"`

	// market specific unique id
	InstanceId string `bson:"instanceId" json:"instanceId"`
}
