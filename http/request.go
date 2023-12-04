package http

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/exception"
	"gitee.com/quant1x/gox/logger"
	"io"
	"net/http"
	URL "net/url"
	"strings"
	"time"
)

const (
	GET     = http.MethodGet
	POST    = http.MethodPost
	HEAD    = http.MethodHead
	PUT     = http.MethodPut
	PATCH   = http.MethodPatch // RFC 5789
	DELETE  = http.MethodDelete
	CONNECT = http.MethodConnect
	OPTIONS = http.MethodOptions
	TRACE   = http.MethodTrace

	ContentEncoding = "Content-Encoding"
	ContextType     = "Content-Type"
	LastModified    = "Last-Modified"
	IfModifiedSince = "If-Modified-Since"
	charsetUtf8     = "charset=UTF-8"
	ApplicationJson = "application/json" + ";" + charsetUtf8
	ApplicationForm = "application/x-www-form-urlencoded" + ";" + charsetUtf8
)

var (
	TimeZero = time.Unix(0, 0)
	NotFound = exception.New(http.StatusNotFound, http.StatusText(http.StatusNotFound))
)

// HttpRequest HTTP 请求
func HttpRequest(url string, method string, header ...map[string]any) ([]byte, error) {
	data, lastModified, err := Request(url, method, "", header...)
	_ = lastModified
	return data, err
}

// Get http get请求
func Get(url string, header ...map[string]any) ([]byte, error) {
	data, _, err := Request(url, GET, "", header...)
	return data, err
}

// HttpGet Get请求
func HttpGet(url string) ([]byte, error) {
	return HttpRequest(url, GET)
}

// HttpPost POST 请求
func HttpPost(url string, content string, header ...map[string]any) (data []byte, lastModified time.Time, err error) {
	content = strings.TrimSpace(content)
	length := len(content)
	start := content[0]
	end := content[length-1]
	var requestHeader map[string]any
	if len(header) == 0 {
		requestHeader = make(map[string]any, 0)
	} else {
		requestHeader = header[0]
	}
	if (start == '{' && end == '}') || (start == '[' && end == ']') {
		// 这是json
		requestHeader[ContextType] = ApplicationJson
	} else {
		requestHeader[ContextType] = ApplicationForm
	}
	return Request(url, POST, content, requestHeader)
}

// Request http request, 支持传入header
func Request(url string, method string, content string, header ...map[string]any) (data []byte, lastModified time.Time, err error) {
	u, err := URL.Parse(url)
	if err != nil {
		return nil, TimeZero, err
	}
	reqHeader := make(map[string]string)
	reqHeader["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
	reqHeader["Accept-Encoding"] = "gzip, deflate"
	reqHeader["Accept-Language"] = "zh-CN,zh;q=0.9,en;q=0.8"
	reqHeader["Cache-Control"] = "no-cache"
	reqHeader["Connection"] = "keep-alive"
	reqHeader["Host"] = u.Host
	reqHeader["Pragma"] = "no-cache"
	reqHeader["Upgrade-Insecure-Requests"] = "1"
	reqHeader["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.35"
	if len(header) > 0 {
		mapHeader := header[0]
		for k, v := range mapHeader {
			switch val := v.(type) {
			case time.Time:
				val = val.UTC()
				reqHeader[k] = val.Format(time.RFC1123)
			case float32, float64:
				reqHeader[k] = fmt.Sprintf("%f", val)
			case int8, int16, int32, int64:
				reqHeader[k] = fmt.Sprintf("%d", val)
			case uint8, uint16, uint32, uint64:
				reqHeader[k] = fmt.Sprintf("%d", val)
			case string:
				reqHeader[k] = val
			default:
				reqHeader[k] = fmt.Sprintf("%v", val)
			}
		}
	}

	client := defaultClient()
	var requestBody io.Reader = nil
	if len(content) > 0 {
		requestBody = strings.NewReader(content)
	}
	request, err := http.NewRequest(strings.ToUpper(method), url, requestBody)
	if err != nil {
		return nil, TimeZero, err
	}
	for key, v := range reqHeader {
		request.Header.Add(key, v)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, TimeZero, err
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, TimeZero, NotFound
	}
	lm := response.Header.Get(LastModified)
	if response.StatusCode == http.StatusNotModified && !api.IsEmpty(lm) {
		return nil, TimeZero, nil
	}
	lastModified, err = time.Parse(time.RFC1123, lm)
	defer api.CloseQuietly(response.Body)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, TimeZero, err
	}
	contentEncoding := response.Header.Get(ContentEncoding)
	var reader io.ReadCloser = nil
	if len(contentEncoding) > 0 {
		contentEncoding = strings.ToLower(contentEncoding)
		switch contentEncoding {
		case "gzip":
			reader, err = gzip.NewReader(bytes.NewBuffer(body))
			if err != nil {
				logger.Error(err)
				reader = nil
			}
		case "deflate":
			reader = flate.NewReader(bytes.NewReader(body))
		}
	}
	if reader != nil {
		defer api.CloseQuietly(reader)
		body, err = io.ReadAll(reader)
		if err != nil {
			return nil, TimeZero, err
		}
	}
	return body, lastModified, nil
}
