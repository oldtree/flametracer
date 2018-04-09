package pkg

import (
	"net/http"
	"testing"
	"time"
)

func Test_newNormalRequest(t *testing.T) {
	var urlstr = "http://www.baidu.com"
	var testc = NewClient(nil, nil)
	req, err := testc.NewRequest("GET", urlstr, nil)
	if err != nil {
		t.Fatalf("new ruquest failed  %s \n", err.Error())
		return
	}
	_, err = testc.Do(req)
	if err != nil {
		t.Fatalf("do request failed : %s \n", err.Error())
		return
	}
	return
}

func Test_newTraceRequest(t *testing.T) {
	var urlstr = "https://www.baidu.com"
	var testc = NewClient(nil, nil)
	req, err := testc.NewRequestWithTrace("GET", urlstr, nil)
	if err != nil {
		t.Fatalf("new trace request failed : %s \n", err.Error())
		return
	}
	_, err = testc.Do(req)
	if err != nil {
		t.Fatalf("do request failed : %s \n ", err.Error())
		return
	}
	return
}

func Test_newTraceRequestWithTimeout(t *testing.T) {
	var urlstr = "http://www.baidu.com"
	var testc = NewClient(nil, nil)
	req, cancelfunc, err := testc.NewRequestWithTraceTimeout("GET", urlstr, nil, time.Second*5)
	if err != nil {
		t.Fatalf("new trace request withtimeout failed : %s \n", err.Error())
		return
	}
	_, err = testc.DoWithTimeout(req)
	defer cancelfunc()
	if err != nil {
		t.Fatalf("do request failed : %s \n", err.Error())
		return
	}
	return
}

func Test_newTraceRequestWithTimeout_certain(t *testing.T) {
	var urlstr = "http://www.google.com"
	var testc = NewClient(nil, nil)
	req, cancelfunc, err := testc.NewRequestWithTraceTimeout("GET", urlstr, nil, time.Second*5)
	if err != nil {
		t.Fatalf("new trace request withtimeout failed : %s \n", err.Error())
		return
	}
	_, err = testc.DoWithTimeout(req)
	defer cancelfunc()
	if err != nil {
		return
	}
	t.Fatalf("do request success \n")
	return
}

func Benchmark_NormalRequest(b *testing.B) {
	var urlstr = "http://www.baidu.com"
	var testc = NewClient(nil, nil)
	var req *http.Request
	var err error
	for index := 0; index < b.N; index++ {
		req, err = testc.NewRequest("GET", urlstr, nil)
		if err != nil {
			b.Fatalf("new request failed %s \n", err.Error())
			return
		}
		_, err = testc.Do(req)
		if err != nil {
			b.Fatalf("do request failed : %s \n", err.Error())
			return
		}
	}
	return
}

func Benchmark_TraceRequest(b *testing.B) {
	var urlstr = "http://www.baidu.com"
	var testc = NewClient(nil, nil)
	var req *http.Request
	var err error
	for index := 0; index < b.N; index++ {
		req, err = testc.NewRequestWithTrace("GET", urlstr, nil)
		if err != nil {
			b.Fatalf("new request failed %s \n", err.Error())
			return
		}
		_, err = testc.Do(req)
		if err != nil {
			b.Fatalf("do request failed : %s \n", err.Error())
			return
		}
	}
	return
}

func Benchmark_TraceRequestWithTimeout_normal(b *testing.B) {
	var urlstr = "http://www.baidu.com"
	var testc = NewClient(nil, nil)
	var req *http.Request
	var err error
	for index := 0; index < b.N; index++ {
		req, _, err = testc.NewRequestWithTraceTimeout("GET", urlstr, nil, time.Second*5)
		if err != nil {
			b.Fatalf("new request failed %s \n", err.Error())
			return
		}
		_, err = testc.DoWithTimeout(req)
		if err != nil {
			b.Fatalf("do request failed : %s \n", err.Error())
			return
		}
	}
	return
}

func Benchmark_TraceRequestWithTimeout_failed(b *testing.B) {
	var urlstr = "http://www.google.com"
	var testc = NewClient(nil, nil)
	var req *http.Request
	var err error
	for index := 0; index < b.N; index++ {
		req, _, err = testc.NewRequestWithTraceTimeout("GET", urlstr, nil, time.Second*5)
		if err != nil {
			b.Fatalf("new request failed %s \n", err.Error())
			return
		}
		_, err = testc.DoWithTimeout(req)
		if err == nil {
			b.Fatalf("should error return ")
			return
		}
	}
	return
}
