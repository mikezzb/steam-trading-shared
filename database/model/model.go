package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	Name              string `bson:"name" json:"name"`
	IconUrl           string `bson:"iconUrl" json:"iconUrl"`
	LowestMarketPrice string `bson:"lowestMarketPrice" json:"lowestMarketPrice"`
	LowestMarketName  string `bson:"lowestMarketName" json:"lowestMarketName"`
	SteamPrice        string `bson:"steamPrice" json:"steamPrice"`
	UpdateAt          int    `bson:"updateAt" json:"updateAt"`

	ID primitive.ObjectID `bson:"_id" json:"id"`
}

type Listing struct {
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

	ID primitive.ObjectID `bson:"_id"`
}

type Transaction struct {
	ID primitive.ObjectID `bson:"_id"`
}
