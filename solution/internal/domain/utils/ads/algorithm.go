package ads

import (
	"math"
)

/*
CPI - Цена за показ
CPC - Цена за клик
rel - Релевантность (ml_score)
Первая часть функции - кф-нт прибыли, зависящий тем больше от CPC, чем выше релевантность и тем больше от CPI, чем она меньше
Вторая - релевантность
Итого оценка складывается из 2/3 кф-нта прибыли и 1/3 релевантности
*/

func AdScore(cpi, cpc, rel, r0 float64) float64 {
	return Logistic(rel, r0)*(cpc+cpi) + (1.0-Logistic(rel, r0))*cpi
}

func Logistic(x, r0 float64) float64 {
	// Кф-нт кривизны - чем меньше, тем плавнее будет переход между влиянием цены за показ и цены за клик
	k := 25.0

	return 1 / (1 + math.Exp(-k*(x-r0)))
}

// Нормализация ml_score из абсолютных значений в пределы от 0 до 1

func Normalization(rel float64) float64 {
	return rel / 2147483647.0
}
