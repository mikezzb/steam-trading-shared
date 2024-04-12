package shared

// WARNING: DO NOT USE THIS IN-MEMORY | {item_name: {tier: []seeds}}
type RarePatternDB map[string]map[string][]int

// {item_name: {seed: tier}}
type RarePatternMap map[string]map[int]string
