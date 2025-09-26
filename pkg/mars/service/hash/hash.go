package hash

import (
	"time"

	"github.com/valyala/fasthttp"
	"github.com/yykhomenko/mars/pkg/mars/config"
)

type HashConnector struct {
	config *config.Config
	client *fasthttp.Client
}

func NewHashConnector(config *config.Config) *HashConnector {
	return &HashConnector{
		config: config,
		client: getClient(config),
	}
}

func (h *HashConnector) getHash(msisdn string) string {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(h.config.HashApiAddr + "/hash/" + msisdn)
	resp := fasthttp.AcquireResponse()
	err := h.client.Do(req, resp)

	fasthttp.ReleaseRequest(req)
	code := resp.StatusCode()
	fasthttp.ReleaseResponse(resp)

	if err == nil {
		if code >= 200 && code <= 399 {

			h.config.Log.Println("OK", msisdn)
		} else {
			h.config.Log.Println("ERROR", msisdn)
		}
	} else {
		h.config.Log.Printf("ERR Connection error: %v\n", err)
	}

	return string(resp.Body())
}

func getClient(conf *config.Config) *fasthttp.Client {
	client := &fasthttp.Client{
		ReadTimeout:                   500 * time.Millisecond,
		WriteTimeout:                  500 * time.Millisecond,
		MaxIdleConnDuration:           1 * time.Hour,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		Dial: (&fasthttp.TCPDialer{
			//Concurrency:      param.ConnNum,
			Concurrency:      1,
			DNSCacheDuration: 1 * time.Hour,
		}).Dial,
	}
	return client
}
