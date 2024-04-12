package shared

// type Item struct {
// 	Name        string
// 	IconUrl     string
// 	MarketPrice string
// 	SteamPrice  string
// }

// type Listing struct {
// 	Name             string
// 	Price            string
// 	CreatedAt        int
// 	UpdatedAt        int
// 	PreviewUrl       string
// 	GoodsId          int
// 	ClassId          string
// 	AssetId          string
// 	TradableCooldown string
// 	// item quality
// 	PaintWear  string
// 	PaintIndex int
// 	PaintSeed  int
// 	Rarity     string
// }

// type Transaction struct {
// }

// WARNING: DO NOT USE THIS IN-MEMORY | {item_name: {tier: []seeds}}
type RarePatternDB map[string]map[string][]int

// {item_name: {seed: tier}}
type RarePatternMap map[string]map[int]string
