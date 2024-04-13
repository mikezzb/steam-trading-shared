package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID                primitive.ObjectID `bson:"_id"`
	Name              string             `bson:"name"`
	IconUrl           string             `bson:"iconUrl"`
	LowestMarketPrice string             `bson:"lowestMarketPrice"`
	LowestMarketName  string             `bson:"lowestMarketName"`
	SteamPrice        string             `bson:"steamPrice"`
	UpdateAt          int                `bson:"updateAt"`
}

type Listing struct {
	ID               primitive.ObjectID `bson:"_id"`
	Name             string             `bson:"name"`
	Price            string             `bson:"price"`
	CreatedAt        int                `bson:"createdAt"`
	UpdatedAt        int                `bson:"updatedAt"`
	PreviewUrl       string             `bson:"previewUrl"`
	GoodsId          int                `bson:"goodsId"`
	ClassId          string             `bson:"classId"`
	AssetId          string             `bson:"assetId"`
	TradableCooldown string             `bson:"tradableCooldown"`

	PaintWear  string `bson:"paintWear"`
	PaintIndex int    `bson:"paintIndex"`
	PaintSeed  int    `bson:"paintSeed"`
	Rarity     string `bson:"rarity"`
}

type Transaction struct {
	ID primitive.ObjectID `bson:"_id"`
}
