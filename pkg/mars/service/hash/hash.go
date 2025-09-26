package hash

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/yykhomenko/mars/pkg/mars/config"
)

type HashConnector struct {
	config *config.Config
	client *fasthttp.Client
}

type HashResponse struct {
	Value    string `json:"value,omitempty"`
	ErrorID  byte   `json:"error_id,omitempty"`
	ErrorMsg string `json:"error_msg,omitempty"`
}

func NewHashConnector(config *config.Config) *HashConnector {
	return &HashConnector{
		config: config,
		client: getClient(config),
	}
}

func (h *HashConnector) GetHash(msisdn string) (*HashResponse, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(h.config.HashApiAddr + "/hashes/" + msisdn)
	//h.config.Log.Println(string(req.RequestURI()))
	resp := fasthttp.AcquireResponse()
	err := h.client.Do(req, resp)

	fasthttp.ReleaseRequest(req)
	code := resp.StatusCode()

	if err == nil {

		var hr HashResponse
		err1 := json.Unmarshal(resp.Body(), &hr)
		if err1 != nil {
			h.config.Log.Printf("hash: body parse error: %v\n", err1)
		}

		if code >= 200 && code <= 399 {

			h.config.Log.Println("hash: OK", msisdn)
			fasthttp.ReleaseResponse(resp)
			return &hr, nil
		} else {
			h.config.Log.Println("hash: ERROR", msisdn)
			fasthttp.ReleaseResponse(resp)
			return &hr, errors.New("hash: hash connector fail")
		}
	} else {
		h.config.Log.Printf("hash: ERR Connection error: %v\n", err)
		fasthttp.ReleaseResponse(resp)
		return nil, errors.New("hash: hash connector fail")
	}

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
