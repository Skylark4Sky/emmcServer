package string

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789~!@#$%^&*()_+")
var chars = []rune("abcdefghijklmnopqrstuvwxyz")
var digitsLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

var digits = []rune("1234567890")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func RandomDigitAndLetters(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(digitsLetters))]
	}
	return string(b)
}

func RandomWord(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(chars))]
	}
	return string(b)
}

func RandEmail() string {
	return RandomWord(6) + "@" + RandomWord(4) + "." + RandomWord(3)
}

func ArgsToJsonData(args []interface{}) zap.Field {
	det := make([]string, 0)
	if len(args) > 0 {
		for _, v := range args {
			det = append(det, fmt.Sprintf("%+v", v))
		}
	}
	zap := zap.Any("detail", det)
	return zap
}

func StringJoin(args []interface{}) string {
	var buffer bytes.Buffer
	if len(args) > 0 {
		for _, v := range args {
			switch v.(type) {
			case uint8, uint16, uint32, uint64, uint:
			case int8, int16, int32, int64, int:
				val := v.(int)
				buffer.WriteString(strconv.Itoa(val))
				break
			case string:
				buffer.WriteString(v.(string))
				break
			}
		}
	}
	return buffer.String()
}

func GetDeviceSN(topic string) string {
	clrString := strings.TrimFunc(topic, func(c rune) bool { return strings.ContainsRune("/", c) })

	stringArray := strings.Split(clrString, "/")

	return stringArray[len(stringArray)-1]
}

//email verify
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//mobile verify
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}
