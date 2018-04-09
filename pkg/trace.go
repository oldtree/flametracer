package pkg

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

type MonitorCallback func(start time.Time)

type Client struct {
	Callback MonitorCallback
	CClient  *http.Client
}

func DefaultTrace() *httptrace.ClientTrace {
	var timeReq = time.Now()
	trace := &httptrace.ClientTrace{
		GetConn: func(host string) {
			log.Printf("start conn : [%s] timestamp : [%s] \n", host, time.Now().Sub(timeReq).String())
		},
		GotConn: func(info httptrace.GotConnInfo) {
			log.Printf("remote address : [%s] reused : [%t] isIdle :[%t] timestamp : [%s] \n", info.Conn.RemoteAddr().String(),
				info.Reused,
				info.WasIdle,
				time.Now().Sub(timeReq).String())
		},
		PutIdleConn: func(err error) {
			if err != nil {
				log.Printf("put connection back failed : [%s] \n timestamp [%s] ", err.Error(), time.Now().Sub(timeReq).String())
			} else {
				log.Printf("put connection back success [%s] \n", time.Now().Sub(timeReq).String())
			}
		},
		DNSStart: func(dnsinfo httptrace.DNSStartInfo) {
			log.Printf("start query dns : [%s] timestamp [%s] \n", dnsinfo.Host, time.Now().Sub(timeReq).String())
		},
		DNSDone: func(dnsinfo httptrace.DNSDoneInfo) {
			if dnsinfo.Err != nil {
				log.Printf("query dns error : [%s] timestamp : [%s] \n", dnsinfo.Err, time.Now().Sub(timeReq).String())
			}
			log.Printf("query dns done : [%s]  isConcurrently : [%t] timestamp : [%s]\n", dnsinfo.Addrs, dnsinfo.Coalesced, time.Now().Sub(timeReq).String())
		},
		ConnectStart: func(network, addr string) {
			log.Printf("build connect : [%s] address : [%s] timestamp [%s] \n", network, addr, time.Now().Sub(timeReq).String())
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				log.Printf("dail : [%s] address : [%s] failed [%s] timestamp [%s] \n", network, addr, err.Error(),
					time.Now().Sub(timeReq).String())
				return
			}
			log.Printf("dail : [%s] address : [%s] success timestamp [%s] \n", network, addr,
				time.Now().Sub(timeReq).String())
			return
		},
		TLSHandshakeStart: func() {
			log.Printf("tls hand shake start : [%s] \n", time.Now().Sub(timeReq).String())
		},
		TLSHandshakeDone: func(tlss tls.ConnectionState, err error) {
			if err != nil {
				log.Printf("tls shake done with failed : [%s] timestamp [%s] \n", err.Error(), time.Now().Sub(timeReq).String())
				return
			}
			log.Printf("tls hand shake done with success : reused ? [%t] timestamp [%s] \n", tlss.DidResume,
				time.Now().Sub(timeReq).String())
			return
		},
		WroteHeaders: func() {
			log.Printf("request header is write finished timestamp : [%s] \n", time.Now().Sub(timeReq).String())
			return
		},
		WroteRequest: func(w httptrace.WroteRequestInfo) {
			if w.Err != nil {
				log.Printf("request info failed [%s] timestamp [%s] \n", w.Err.Error(), time.Now().Sub(timeReq).String())
			}
			log.Printf("request info [%s] \n", time.Now().Sub(timeReq).String())
		},
	}
	return trace
}

func NewClient(client *http.Client, callback MonitorCallback) *Client {
	if client == nil {
		return &Client{
			Callback: callback,
			CClient:  http.DefaultClient,
		}
	}
	return &Client{
		Callback: callback,
		CClient:  client,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.Callback != nil {
		log.Println("statistic not enable")
	}
	resp, err := c.CClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) DoWithTimeout(req *http.Request) (*http.Response, error) {
	if c.Callback != nil {
		log.Println("statistics is not enable")
	}
	var resp *http.Response
	var err error
	var respChan = make(chan *http.Response, 1)
	go func() {
		resp, err = c.CClient.Do(req)
		respChan <- resp
	}()
	if err != nil {
		log.Println("do request failed : ", err.Error())
		return nil, err
	}
	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case value := <-respChan:
		if value != nil {
			return value, nil
		}
		return nil, err
	}
	return nil, err
}

func (c *Client) NewRequestWithTrace(method string, url string, body io.Reader) (*http.Request, error) {
	sourceReq, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("build request failed : [%s] \n", err.Error())
		return nil, err
	}
	sourceReq = sourceReq.WithContext(httptrace.WithClientTrace(sourceReq.Context(), DefaultTrace()))
	if sourceReq == nil {
		return nil, errors.New("new request failed")
	}
	return sourceReq, nil
}

func (c *Client) NewRequestWithTraceTimeout(method string, url string, body io.Reader, timeDur time.Duration) (*http.Request, context.CancelFunc, error) {
	sourceReq, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("build request failed : [%s] \n", err.Error())
		return nil, nil, err
	}
	ctx, cancelfunc := context.WithTimeout(sourceReq.Context(), timeDur)

	sourceReq = sourceReq.WithContext(httptrace.WithClientTrace(ctx, DefaultTrace()))
	if sourceReq == nil {
		return nil, nil, errors.New("new request failed ")
	}
	return sourceReq, cancelfunc, nil
}

func (c *Client) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}
