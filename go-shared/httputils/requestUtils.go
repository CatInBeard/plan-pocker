package httputils

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type RequestInfo struct {
	Headers    http.Header
	URL        *url.URL
	Params     url.Values
	Cookies    []*http.Cookie
	RemoteAddr string
}

func CollectRequestInfo(r *http.Request) RequestInfo {
	info := RequestInfo{
		Headers:    r.Header,
		URL:        r.URL,
		Params:     r.URL.Query(),
		Cookies:    nil,
		RemoteAddr: r.RemoteAddr,
	}

	cookies := r.Cookies()
	for _, cookie := range cookies {
		info.Cookies = append(info.Cookies, cookie)
	}

	return info
}

func CollectRequestInfoString(r *http.Request) string {
	info := CollectRequestInfo(r)
	return info.String()
}

func (ri RequestInfo) String() string {
	var sb strings.Builder

	sb.WriteString("Headers:\n")
	for key, values := range ri.Headers {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", ")))
	}

	sb.WriteString(fmt.Sprintf("URL: %s\n", ri.URL.String()))
	sb.WriteString("Params:\n")
	for key, values := range ri.Params {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", ")))
	}

	sb.WriteString("Cookies:\n")
	for _, cookie := range ri.Cookies {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", cookie.Name, cookie.Value))
	}

	sb.WriteString(fmt.Sprintf("Remote Address: %s\n", ri.RemoteAddr))

	return sb.String()
}
