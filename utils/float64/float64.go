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

func CmpElectricityException(newElectricity, cacheElectricity, cmpValue float64) bool {
	if (newElectricity != cacheElectricity) && (math.Abs(cacheElectricity-cacheElectricity)) >= cmpValue {
		return true
	}
	return false
}

func CmpPower(val1 float64, val2 float64) int {
	return decimal.NewFromFloat(val1).Cmp(decimal.NewFromFloat(val2))
}

func GetPowerValue(power float64, places int32) string {
	curPower := decimal.NewFromFloat(power)
	return curPower.Round(places).String()
}

//当前功率=电压*电流
func CalculateCurComPower(voltage float64, electricity uint32, bit int32) (power float64) {
	power = 0.00
	if electricity > 0 {
		curElectricity := decimal.NewFromFloat(BASE_ELECTRICITY_RATIO).Mul(decimal.NewFromFloat(float64(electricity)))
		curPower := decimal.NewFromFloat(voltage).Mul(curElectricity)
		power, _ = curPower.Round(bit).Float64()
		return power
	}
	return
}

func CalculateCurComPowerToString(voltage float64, electricity uint32, bit int32) (power string) {
	power = "0"
	if electricity > 0 {
		curElectricity := decimal.NewFromFloat(BASE_ELECTRICITY_RATIO).Mul(decimal.NewFromFloat(float64(electricity)))
		curPower := decimal.NewFromFloat(voltage).Mul(curElectricity)
		power = curPower.Round(bit).String()
	}
	return
}

//平均功率=电能(度数)/小时
func CalculateCurAverageComPower(energy uint32, timeSecond uint32, bit int32) (power float64) {
	power = 0.00
	if energy == 0 || timeSecond == 0 {
		return
	}
	time := decimal.NewFromFloat(float64(timeSecond)).Div(decimal.NewFromFloat(TIMEBYHOUR))
	curPower := decimal.NewFromFloat(float64(energy)).Div(time)
	power, _ = curPower.Round(bit).Float64()
	return power
}

func CalculateCurAverageComPowerToString(energy uint32, timeSecond uint32, bit int32) (power string) {
	power = "0"
	if energy == 0 || timeSecond == 0 {
		return
	}
	time := decimal.NewFromFloat(float64(timeSecond)).Div(decimal.NewFromFloat(TIMEBYHOUR))
	curPower := decimal.NewFromFloat(float64(energy)).Div(time)
	power = curPower.Round(bit).String()
	return
}

//电流=瓦数/电能
func CalculateMaxComElectricity(wattage float64) int64 {
	value := decimal.NewFromFloat(wattage).Div(decimal.NewFromFloat(CUR_VOLTAGE))
	value = value.Mul(decimal.NewFromFloat(float64(1000)))
	electricity, _ := value.Round(0).Float64()
	return int64(electricity)
}
