package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Name              string `bson:"name" json:"name"`
	IconUrl           string `bson:"iconUrl" json:"iconUrl"`
	LowestMarketPrice string `bson:"lowestMarketPrice" json:"lowestMarketPrice"`
	LowestMarketName  string `bson:"lowestMarketName" json:"lowestMarketName"`
	SteamPrice        string `bson:"steamPrice" json:"steamPrice"`
	UpdateAt          int    `bson:"updateAt" json:"updateAt"`
}

type Listing struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Name             string `bson:"name"`
	Price            string `bson:"price"`
	CreatedAt        int    `bson:"createdAt"`
	UpdatedAt        int    `bson:"updatedAt"`
	PreviewUrl       string `bson:"previewUrl"`
	GoodsId          int    `bson:"goodsId"`
	ClassId          string `bson:"classId"`
	AssetId          string `bson:"assetId"`
	TradableCooldown string `bson:"tradableCooldown"`

	PaintWear  string `bson:"paintWear"`
	PaintIndex int    `bson:"paintIndex"`
	PaintSeed  int    `bson:"paintSeed"`
	Rarity     string `bson:"rarity"`
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
}

// Currently same as Listing
type Transaction struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Name             string `bson:"name"`
	Price            string `bson:"price"`
	CreatedAt        int    `bson:"createdAt"`
	UpdatedAt        int    `bson:"updatedAt"`
	PreviewUrl       string `bson:"previewUrl"`
	GoodsId          int    `bson:"goodsId"`
	ClassId          string `bson:"classId"`
	AssetId          string `bson:"assetId"`
	TradableCooldown string `bson:"tradableCooldown"`

	PaintWear  string `bson:"paintWear"`
	PaintIndex int    `bson:"paintIndex"`
	PaintSeed  int    `bson:"paintSeed"`
	Rarity     string `bson:"rarity"`
}
