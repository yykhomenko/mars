package sis

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/yykhomenko/mars/pkg/mars/config"
)

type SisConnector struct {
	config *config.Config
	client *fasthttp.Client
}

type SisErrorResponse struct {
	Value    string `json:"value,omitempty"`
	ErrorID  byte   `json:"error_id,omitempty"`
	ErrorMsg string `json:"error_msg,omitempty"`
}

type SisGoodResponse struct {
	Msisdn       int64     `json:"msisdn"`
	BillingType  int8      `json:"billing_type"`
	LanguageType int8      `json:"language_type"`
	OperatorType int8      `json:"operator_type"`
	ChangeDate   time.Time `json:"change_date"`
}

func NewSisConnector(config *config.Config) *SisConnector {
	return &SisConnector{
		config: config,
		client: getClient(config),
	}
}

func (h *SisConnector) GetSubscriber(msisdn string) (*SisGoodResponse, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(h.config.SisApiAddr + "/subscribers/" + msisdn)
	//h.config.Log.Println(string(req.RequestURI()))
	resp := fasthttp.AcquireResponse()
	err := h.client.Do(req, resp)

	fasthttp.ReleaseRequest(req)
	code := resp.StatusCode()

	if err == nil {

		var sr SisGoodResponse
		err1 := json.Unmarshal(resp.Body(), &sr)
		if err1 != nil {
			h.config.Log.Printf("sis: body parse error: %v\n", err1)
		}

		if code >= 200 && code <= 399 {

			//h.config.Log.Println("sis: OK", msisdn)
			fasthttp.ReleaseResponse(resp)
			return &sr, nil
		} else {
			h.config.Log.Println("sis: ERROR", msisdn)
			fasthttp.ReleaseResponse(resp)
			return &sr, errors.New("sis: sis connector fail")
		}
	} else {
		h.config.Log.Printf("sis: ERR Connection error: %v\n", err)
		fasthttp.ReleaseResponse(resp)
		return nil, errors.New("sis: sis connector fail")
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
			Concurrency:      8,
			DNSCacheDuration: 1 * time.Hour,
		}).Dial,
	}
	return client
}
