package shared

import (
	"github.com/mikezzb/steam-trading-shared/database/model"
	"testing"
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
		}

		for _, pair := range testPairs {
			actual := GetListingTier(pair.listing)
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
