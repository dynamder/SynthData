package batch

import (
	"testing"
)

func TestBatcher_DivideIntoBatches(t *testing.T) {
	batcher := NewBatcher()

	tests := []struct {
		name        string
		targetCount int
		batchSize   int
		want        int
	}{
		{"100 records, batch 100", 100, 100, 1},
		{"100 records, batch 50", 100, 50, 2},
		{"100 records, batch 30", 100, 30, 4},
		{"1 record, batch 100", 1, 100, 1},
		{"0 records", 0, 100, 0},
		{"negative target", -10, 100, 0},
		{"zero batch size", 100, 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batches := batcher.DivideIntoBatches(tt.targetCount, tt.batchSize)
			if len(batches) != tt.want {
				t.Errorf("DivideIntoBatches(%d, %d) = %d batches, want %d", tt.targetCount, tt.batchSize, len(batches), tt.want)
			}
		})
	}
}

func TestBatcher_BatchSizes(t *testing.T) {
	batcher := NewBatcher()
	batches := batcher.DivideIntoBatches(100, 30)

	expectedSizes := []int{30, 30, 30, 10}
	if len(batches) != len(expectedSizes) {
		t.Fatalf("expected %d batches, got %d", len(expectedSizes), len(batches))
	}

	for i, b := range batches {
		if b.Size != expectedSizes[i] {
			t.Errorf("batch %d size = %d, want %d", i, b.Size, expectedSizes[i])
		}
	}
}

func TestBatcher_BatchIDs(t *testing.T) {
	batcher := NewBatcher()
	batches := batcher.DivideIntoBatches(50, 10)

	for i, b := range batches {
		if b.ID != i+1 {
			t.Errorf("batch %d ID = %d, want %d", i, b.ID, i+1)
		}
	}
}
