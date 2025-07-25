package jobs

type Strategy interface {
	Present(daysAfterCreatedAt int) bool
}

type throughOneStrategy struct {
	preset map[int]struct{}
}

func newThroughOneStrategy(maxValue int) *throughOneStrategy {
	preset := make(map[int]struct{}, maxValue)

	for i := 1; i <= maxValue; i += 2 {
		preset[i] = struct{}{}
	}

	return &throughOneStrategy{
		preset: preset,
	}
}

func (s *throughOneStrategy) Present(i int) bool {
	_, ok := s.preset[i]

	return ok
}
