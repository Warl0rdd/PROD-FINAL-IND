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

func AdScore(cpi, cpc, rel float64) float64 {
	return (2.0/3.0)*(logistic(rel)*cpc+(1.0-logistic(rel))*cpi) + rel/3.0
}

func logistic(x float64) float64 {
	// Кф-нт кривизны - чем меньше, тем плавнее будет переход между влиянием цены за показ и цены за клик
	k := 0.15
	// Значение релевантности, в котором функция будет равна 1
	r0 := 50.0

	return 1 / (1 + math.Exp(-k*(x-r0)))
}
