package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	nHttp "net/http"
	"time"
)

type Req struct {
	Method  ReqMethod
	Url     string `validate:"required,url"`
	Body    []byte
	Headers map[string]string
	Timeout int
}

type Resp struct {
	Body          []byte
	StatusCode    int
	Header        nHttp.Header
	ContentLength int64
}

type ReqMethod int

const (
	Get = iota + 1
	Post
	Head
	Put
	Patch
)

func (method ReqMethod) String() string {
	switch method {
	case Get:
		return "GET"
	case Post:
		return "POST"
	case Head:
		return "HEAD"
	case Put:
		return "PUT"
	case Patch:
		return "PATCH"
	default:
		return "GET"
	}
}

func (r Req) Do() (Resp, error) {
	rp := Resp{}
	validate := validator.New()
	err := validate.Struct(&r)
	if err != nil {
		return rp, err
	}
	
	if r.Timeout == 0 {
		r.Timeout = 5
	}
	client := &nHttp.Client{
		Timeout: time.Duration(r.Timeout) * time.Second,
	}
	
	req, err := nHttp.NewRequest(r.Method.String(), r.Url, bytes.NewBuffer(r.Body))
	if err != nil {
		return rp, err
	}
	
	req.Header.Set("user-agent", "startops")
	
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return rp, err
	}
	
	rp.Header = resp.Header
	rp.StatusCode = resp.StatusCode
	rp.ContentLength = resp.ContentLength
	
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rp, err
	}
	
	rp.Body = respByte
	defer resp.Body.Close()
	
	return rp, nil
}

// 查看 HTTP 站点证书

func GetHttpCertificate(url string) {
	tr := &nHttp.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &nHttp.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()
	
	for _, cert := range resp.TLS.PeerCertificates {
		fmt.Println("cert DNSNames: ", cert.DNSNames)
		fmt.Println("cert NotAfter: ", cert.NotAfter)
		fmt.Println("cert NotBefore: ", cert.NotBefore)
		fmt.Println("cert Subject: ", cert.Subject)
		fmt.Println("cert Issuer: ", cert.Issuer)
		expDate := int(cert.NotAfter.Sub(time.Now()).Hours() / 24)
		if expDate <= 30 {
			return
		}
	}
	return
}
