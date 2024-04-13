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
