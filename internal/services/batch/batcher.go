package batch

type Batcher struct{}

func NewBatcher() *Batcher {
	return &Batcher{}
}

func (b *Batcher) DivideIntoBatches(targetCount, batchSize int) []Batch {
	if batchSize <= 0 {
		batchSize = 100
	}
	if targetCount <= 0 {
		return []Batch{}
	}

	numBatches := (targetCount + batchSize - 1) / batchSize
	batches := make([]Batch, 0, numBatches)

	for i := 0; i < numBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > targetCount {
			end = targetCount
		}
		batches = append(batches, Batch{
			ID:       i + 1,
			Start:    start,
			End:      end,
			Size:     end - start,
			RecordID: start + 1,
		})
	}

	return batches
}

type Batch struct {
	ID       int
	Start    int
	End      int
	Size     int
	RecordID int
}
