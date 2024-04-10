package shared

type Item struct {
	Name        string
	IconUrl     string
	MarketPrice string
	SteamPrice  string
}

type Listing struct {
	Price            string
	CreatedAt        int
	UpdatedAt        int
	PreviewUrl       string
	GoodsId          int
	ClassId          string
	AssetId          string
	TradableCooldown string
	// item quality
	PaintWear  string
	PaintIndex int
	PaintSeed  int
	Rarity     string
}

type Transaction struct {
}
