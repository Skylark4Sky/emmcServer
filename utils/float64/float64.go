package float64

import (
	"github.com/shopspring/decimal"
	"math"
)

const (
	CUR_VOLTAGE            float64 = 240.0 //市电 电压 240
	BASE_ELECTRICITY_RATIO float64 = 0.001
	TIMEBYHOUR             float64 = 3600.0
)

//功率=电压*电流
func CalculateComPowerToCustomer(voltage float64, electricity uint32, bit int) (power float64) {
	power = 0.00
	if electricity > 0 {
		curElectricity := decimal.NewFromFloat(BASE_ELECTRICITY_RATIO).Mul(decimal.NewFromFloat(float64(electricity)))
		curPower := decimal.NewFromFloat(voltage).Mul(curElectricity)
		var exact bool
		power, exact = curPower.Float64()
		if exact {
			//保留n位小数
			n10 := math.Pow10(bit)
			power = math.Trunc((power+0.5/n10)*n10) / n10
		}
	}
	return
}

//功率=电能(度数)/小时
func CalculateComPowerToPartner(energy uint32, timeSecond uint32, bit int) (power float64) {
	power = 0.00
	if energy == 0 || timeSecond == 0 {
		return
	}

	time := decimal.NewFromFloat(float64(timeSecond)).Div(decimal.NewFromFloat(TIMEBYHOUR))
	curPower := decimal.NewFromFloat(float64(energy)).Div(time)
	var exact bool
	power, exact = curPower.Float64()
	if exact {
		//保留n位小数
		n10 := math.Pow10(bit)
		power = math.Trunc((power+0.5/n10)*n10) / n10
	}
	return
}
