module GoServer

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jinzhu/gorm v1.9.12
	github.com/lestrrat-go/file-rotatelogs v2.3.0+incompatible
	github.com/lestrrat-go/strftime v1.0.1 // indirect
	github.com/mojocn/base64Captcha v1.3.1
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.5.0
	github.com/ugorji/go v1.1.7 // indirect
	golang.org/x/crypto v0.0.0-20191205180655-e7c4368fe9dd
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	gopkg.in/yaml.v2 v2.2.4
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
