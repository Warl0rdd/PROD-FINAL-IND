package learning

import "solution/internal/adapters/database/postgres"

// Подробнее в algorithm.md

func GenNewR0(oldR0 float64, data []postgres.GetImpressionsForLearningRow) float64 {
	n := 0.1

	var sum float64
	for _, item := range data {
		var trueClick float64
		switch item.ClickedAfter {
		case true:
			trueClick = 1.0
		case false:
			trueClick = 0.0
		}
		sum += item.ModelScore - trueClick
	}

	sum *= 0.15 / float64(len(data))

	return oldR0 - n*sum
}
