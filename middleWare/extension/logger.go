package middleWare

import (
	. "GoServer/utils/log"
	. "GoServer/utils/time"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLogger() gin.HandlerFunc {

	return func(c *gin.Context) {
		traceID := uuid.New().String()

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		var response interface{}
		if respMap, err := formatResponse(blw.body.Bytes()); err == nil {
			response = respMap
		} else {
			response = blw.body.Bytes()
		}

		var bodyStr string
		if len(bodyBytes) > 0 {
			bodyStr, _ = url.QueryUnescape(string(bodyBytes))
		}

		var cookies = formatCookies(c.Request.Cookies())

		if agentID, ok := c.Get("agent_id"); ok {
			cookies["agent_id"] = strconv.Itoa(int(agentID.(int32)))
			cookies["agent_name"] = c.GetString("agent_name")
		}

		SystemLog("time:",TimeFormat(time.Now()) ," id:" ,traceID," status:",c.Writer.Status(), " method:",c.Request.Method," ip:",c.ClientIP()," cookies:", cookies," latency:",latency)
		SystemLog("agent:",c.Request.UserAgent())
		SystemLog("path:",c.Request.URL.Path ," query:" ,c.Request.URL.RawQuery)
		SystemLog("body:", bodyStr)
		SystemLog("response:", response)
		SystemLog("")

		if len(c.Errors) > 0 {
			SystemLog(c.Errors.ByType(gin.ErrorTypeAny).String())
		}
	}
}

func formatCookies(cookies []*http.Cookie) (cookiesMap map[string]string) {
	cookiesMap = make(map[string]string)
	for _, cookie := range cookies {
		cookiesMap[cookie.Name] = cookie.Value
	}
	return
}

func formatResponse(data []byte) (map[string]interface{}, error) {
	var responseMap = make(map[string]interface{})
	err := json.Unmarshal(data, &responseMap)
	if err != nil {
		return responseMap, err
	}
	return responseMap, nil
}
