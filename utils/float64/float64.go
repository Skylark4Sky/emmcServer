package float64

import (
	"fmt"
	"github.com/shopspring/decimal"
)

const (
	CUR_VOLTAGE            float64 = 240.0 //市电 电压 240
	BASE_ELECTRICITY_RATIO float64 = 0.001
	TIMEBYHOUR             float64 = 3600.0
)

func CmpPower(val1 float64,val2 float64) int {
	return decimal.NewFromFloat(val1).Cmp(decimal.NewFromFloat(val2))
}

func GetFloat64(val decimal.Decimal)  float64 {
	value, exact := val.Float64()
	if exact {
		return value
	}
	return 0.0
}

//当前功率=电压*电流
func CalculateCurComPower(voltage float64, electricity uint32, bit int32) float64 {
	if electricity > 0 {
		curElectricity := decimal.NewFromFloat(BASE_ELECTRICITY_RATIO).Mul(decimal.NewFromFloat(float64(electricity)))
		curPower := decimal.NewFromFloat(voltage).Mul(curElectricity).RoundBank(bit)

		power, exact := curPower.Float64()
		if exact {
			fmt.Println("------------->",power)
			return power
			//保留n位小数
			//n10 := math.Pow10(bit)
			//power = math.Trunc((power+0.5/n10)*n10) / n10
		}
	}
	return 0.0
}

//平均功率=电能(度数)/小时
func CalculateCurAverageComPower(energy uint32, timeSecond uint32, bit int32) float64 {
	if energy == 0 || timeSecond == 0 {
		return 0.00
	}

	time := decimal.NewFromFloat(float64(timeSecond)).DivRound(decimal.NewFromFloat(TIMEBYHOUR),2)
	curPower := decimal.NewFromFloat(float64(energy)).DivRound(time,2).RoundBank(bit)
	power, exact := curPower.Float64()
	if exact {
		fmt.Println("------------->",power)
		return power
		//保留n位小数
		//n10 := math.Pow10(bit)
		//power = math.Trunc((power+0.5/n10)*n10) / n10
	}
	return 0.00
}
