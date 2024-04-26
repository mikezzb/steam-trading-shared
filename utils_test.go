package shared

import (
	"testing"

	"github.com/mikezzb/steam-trading-shared/database/model"
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

func TestOthers(t *testing.T) {
	t.Run("GetTimestampNow", func(t *testing.T) {
		t.Log(GetTimestampNow())
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
				t.Errorf("invalid price smaller than 0.01")
			}
		}
	})
}
