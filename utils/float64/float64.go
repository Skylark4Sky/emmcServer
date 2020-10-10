package float64

import (
	"github.com/shopspring/decimal"
	"math"
)

const (
	CUR_VOLTAGE float64 = 220.0
	BASE_ELECTRICITY_RATIO float64 = 0.001
)

func CalculateComPower(voltage float64, electricity uint32, bit int) (power float64) {
	power = 0.00
	if electricity > 0 {
		curElectricity := decimal.NewFromFloat(BASE_ELECTRICITY_RATIO).Mul(decimal.NewFromFloat(float64(electricity)))
		curPower := decimal.NewFromFloat(voltage).Mul(curElectricity)
		power,exact := curPower.Float64()

		if exact {
			//保留2位小数
			n10 := math.Pow10(bit)
			power = math.Trunc((power+0.5/n10)*n10) / n10
		}
	}
	return
}