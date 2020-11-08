package string

import (
	"bytes"
	"math/rand"
	"reflect"
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

func StringJoin(args []interface{}) string {
	var buffer bytes.Buffer
	if len(args) > 0 {
		for _, arg := range args {
			t := reflect.TypeOf(arg)
			val := reflect.ValueOf(arg)
			var scratch [64]byte

			switch t.Kind() {
			case reflect.Bool:
				{
					if val.Bool() {
						buffer.WriteString("1")
					} else {
						buffer.WriteString("0")
					}
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				b := strconv.AppendInt(scratch[:0], val.Int(), 10)
				buffer.Write(b)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				b := strconv.AppendUint(scratch[:0], val.Uint(), 10)
				buffer.Write(b)
			case reflect.String:
				buffer.WriteString(val.String())
			}
		}
	}
	return buffer.String()
}

func GetDeviceSN(topic string, partition string) string {
	// partition --> "/" or ":"
	clrString := strings.TrimFunc(topic, func(c rune) bool { return strings.ContainsRune(partition, c) })

	stringArray := strings.Split(clrString, partition)

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

func ReturnTopicPrefix(topic string) string {
	if strings.HasPrefix(topic, "/point_switch") {
		return "/point_switch"
	}

	if strings.HasPrefix(topic, "/power_run") {
		return "/power_run"
	}

	if strings.HasPrefix(topic, "/point_common") {
		return "/point_common"
	}

	if strings.HasPrefix(topic, "/point_switch_resp") {
		return "/point_switch_resp"
	}

	if strings.HasPrefix(topic, "/device") {
		return "/device"
	}

	if strings.HasPrefix(topic, "/firmware_update") {
		return "/firmware_update"
	}
	return ""
}
