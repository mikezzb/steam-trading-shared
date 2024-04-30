package shared

import (
	"testing"

	"github.com/mikezzb/steam-trading-shared/database/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUtils_Naming(t *testing.T) {
	t.Run("ExtractBaseItemName", func(t *testing.T) {
		testPairs := []struct {
			itemName string
			expected string
		}{
			{"★ Bayonet | Doppler (Factory New)", "Bayonet | Doppler"},
			{"StatTrak™ AK-47 | Redline (Field-Tested)", "AK-47 | Redline"},
			{"Five-SeveN | Case Hardened (Battle-Scarred)", "Five-SeveN | Case Hardened"},
		}

		for _, pair := range testPairs {
			actual := ExtractBaseItemName(pair.itemName)
			if actual != pair.expected {
				t.Errorf("Expected %s, got %s", pair.expected, actual)
			}
		}
	})

	t.Run("Name decode", func(t *testing.T) {

		tests := []struct {
			name         string
			input        string
			wantName     string
			wantSkin     string
			wantExterior string
		}{
			{
				name:         "Standard item decode",
				input:        "AK-47 | Redline (Field-Tested)",
				wantName:     "AK-47",
				wantSkin:     "Redline",
				wantExterior: "Field-Tested",
			},
			{
				name:         "Item with no skin and exterior",
				input:        "Glock-18",
				wantName:     "Glock-18",
				wantSkin:     "",
				wantExterior: "",
			},
			{
				name:         "Item with unusual characters",
				input:        "M4A1-S | Chantico's Fire (Well-Worn)",
				wantName:     "M4A1-S",
				wantSkin:     "Chantico's Fire",
				wantExterior: "Well-Worn",
			},
			{
				name:         "Knife",
				input:        "★ Bayonet | Doppler (Factory New)",
				wantName:     "★ Bayonet",
				wantSkin:     "Doppler",
				wantExterior: "Factory New",
			},
			{
				name:         "Incorrect format missing pipe",
				input:        "USP-S Orion (Factory New)",
				wantName:     "USP-S Orion (Factory New)",
				wantSkin:     "",
				wantExterior: "",
			},
			{
				name:         "Empty string input",
				input:        "",
				wantName:     "",
				wantSkin:     "",
				wantExterior: "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				gotName, gotSkin, gotExterior := DecodeItemFullName(tt.input)
				if gotName != tt.wantName || gotSkin != tt.wantSkin || gotExterior != tt.wantExterior {
					t.Errorf("DecodeItemFullName(%q) = %q, %q, %q; want %q, %q, %q",
						tt.input, gotName, gotSkin, gotExterior, tt.wantName, tt.wantSkin, tt.wantExterior)
				}
			})
		}
	})
}

func TestUtils_Tier(t *testing.T) {
	t.Run("GetListingTier", func(t *testing.T) {
		testPairs := []struct {
			listing  model.Listing
			expected string
		}{
			{
				model.Listing{Name: "★ Flip Knife | Marble Fade (Factory New)", PaintSeed: 872},
				"Tricolor",
			},
			{
				model.Listing{Name: "★ Karambit | Doppler (Factory New)", PaintSeed: 741},
				"Good Phase 2",
			},
			{
				model.Listing{Name: "★ Bayonet | Marble Fade (Factory New)", PaintSeed: 727},
				"FFI",
			},
		}

		for _, pair := range testPairs {
			listing := pair.listing
			actual := GetTier(listing.Name, listing.PaintSeed)
			if actual != pair.expected {
				t.Errorf("Expected %s, got %s", pair.expected, actual)
			}
		}
	})
}

func TestRandSleep(t *testing.T) {
	t.Run("RandSleep", func(t *testing.T) {
		RandomSleep(5, 9)
	})
}

func TestGetListingUrl(t *testing.T) {
	t.Run("GetListingUrlBuff", func(t *testing.T) {
		listing := model.Listing{
			Name:       "★ Skeleton Knife | Stained (Minimal Wear)",
			Market:     "buff",
			ClassId:    "3608173279",
			InstanceId: "188530139",
			AssetId:    "36356925895",
		}
		url := GetListingUrl(&listing)
		if url != "https://buff.163.com/goods/776793?appid=730&assetid=36356925895&classid=3608173279&instanceid=188530139" {
			t.Errorf("Expected %s, got %s", "https://buff.163.com/goods/776793?appid=730&assetid=36356925895&classid=3608173279&instanceid=188530139", url)
		}
	})
}

func TestStringPrice(t *testing.T) {
	t.Run("StringPrice", func(t *testing.T) {
		testVals := []string{
			"-1",
			"0.000000000000000000000000534",
			"0.01",
			"0.1",
			"1",
			"10",
			"900",
			"23",
			"952340",
			"100000",
			"1000000",
		}

		for _, val := range testVals {
			if MAX_DECIMAL128.String() < val {
				t.Errorf("invalid price smaller than MAX_DECIMAL128: %s", val)
			}
		}
	})
}

func TestDecCompareTo(t *testing.T) {
	testCases := []struct {
		name     string
		a, b     primitive.Decimal128
		expected int
	}{
		{"a < b", primitive.NewDecimal128(123, 2), primitive.NewDecimal128(456, 2), -1},
		{"a > b", primitive.NewDecimal128(789, 2), primitive.NewDecimal128(123, 2), 1},
		{"a == b", primitive.NewDecimal128(123, 2), primitive.NewDecimal128(123, 2), 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := DecCompareTo(tc.a, tc.b)
			if actual != tc.expected {
				t.Errorf("Expected %d, but got %d", tc.expected, actual)
			}
		})
	}
}
