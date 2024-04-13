package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID                primitive.ObjectID `bson:"_id"`
	Name              string             `bson:"name"`
	IconUrl           string             `bson:"icon_url"`
	LowestMarketPrice string             `bson:"lowest_market_price"`
	LowestMarketName  string             `bson:"lowest_market_name"`
	SteamPrice        string             `bson:"steam_price"`
}

type Listing struct {
	ID               primitive.ObjectID `bson:"_id"`
	Name             string             `bson:"name"`
	Price            string             `bson:"price"`
	CreatedAt        int                `bson:"created_at"`
	UpdatedAt        int                `bson:"updated_at"`
	PreviewUrl       string             `bson:"preview_url"`
	GoodsId          int                `bson:"goods_id"`
	ClassId          string             `bson:"class_id"`
	AssetId          string             `bson:"asset_id"`
	TradableCooldown string             `bson:"tradable_cooldown"`

	PaintWear  string `bson:"paint_wear"`
	PaintIndex int    `bson:"paint_index"`
	PaintSeed  int    `bson:"paint_seed"`
	Rarity     string `bson:"rarity"`
}

type Transaction struct {
	ID primitive.ObjectID `bson:"_id"`
}
