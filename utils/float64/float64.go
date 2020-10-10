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

//当前功率=电压*电流
func CalculateCurComPower(voltage float64, electricity uint32, bit int32) float64 {
	if electricity > 0 {
		fmt.Println("\nCalculateCurComPower ","voltage",voltage,"electricity",electricity)
		curElectricity := decimal.NewFromFloat(BASE_ELECTRICITY_RATIO).Mul(decimal.NewFromFloat(float64(electricity)))
		curPower := decimal.NewFromFloat(voltage).Mul(curElectricity)
		power, exact := curPower.Float64()
		if exact {
			fmt.Println("power",power,"\n")
			return power
		}
	}
	return 0.00
}

//平均功率=电能(度数)/小时
func CalculateCurAverageComPower(energy uint32, timeSecond uint32, bit int32) float64 {
	if energy == 0 || timeSecond == 0 {
		return 0.00
	}
	fmt.Println("\nCalculateCurAverageComPower ","energy",energy,"timeSecond",timeSecond)
	time := decimal.NewFromFloat(float64(timeSecond)).Div(decimal.NewFromFloat(TIMEBYHOUR))
	curPower := decimal.NewFromFloat(float64(energy)).Div(time)
	power, exact := curPower.Float64()
	if exact {
		fmt.Println("++CalculateCurAverageComPower++",power,"\n")
		return power
	}
	return 0.00
}
