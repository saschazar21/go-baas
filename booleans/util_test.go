package booleans

import "testing"

func TestGenerateRandomId(t *testing.T) {
	t.Run("generate random id", func(t *testing.T) {
		got := generateRandomId()

		if len(got) != 16 {
			t.Errorf("generateRandomId() = %v, want length of 16", got)
		}
	})
}
