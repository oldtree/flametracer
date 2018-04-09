package monitor

import (
	"encoding/json"
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

//unit is ms
type TraceInfo struct {
	Hostname               string   `json:"hostname,omitempty"`
	ConnBuildTime          float64  `json:"conn_build_time,omitempty"`
	DnsQueryStartTime      float64  `json:"dns_query_start_time,omitempty"`
	DnsQueryEndTime        float64  `json:"dns_query_end_time,omitempty"`
	IsConcurrently         bool     `json:"is_concurrently,omitempty"`
	ConnectionBuildTime    float64  `json:"connection_build_time,omitempty"`
	DailTime               float64  `json:"dail_time,omitempty"`
	DailInfo               string   `json:"dail_info,omitempty"`
	TlsShakeStart          float64  `json:"tls_shake_start,omitempty"`
	TlsShakeEnd            float64  `json:"tls_shake_end,omitempty"`
	IsReused               bool     `json:"is_reused,omitempty"`
	RemoteAddress          []string `json:"remote_address,omitempty"`
	IsIdle                 bool     `json:"is_idle,omitempty"`
	RequestHeaderWriteTime float64  `json:"request_header_write_time,omitempty"`
}

func (t *TraceInfo) Format2Json() []byte {
	data, _ := json.Marshal(t)
	return data
}

type TraceWorker struct {
	Traceinfo    chan *TraceInfo
	Client       client.Client
	DatabaseName string
}

func (t *TraceWorker) Reportor(info *TraceInfo) {
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  t.DatabaseName,
		Precision: "ms",
	})

	fields := map[string]interface{}{
		"dns.start":      info.DnsQueryStartTime,
		"dns.done":       info.DnsQueryEndTime,
		"conn.build":     info.ConnBuildTime,
		"concurrent":     info.IsConcurrently,
		"conn.connected": info.ConnectionBuildTime,
		"dail.start":     info.DailTime,
		"dail.info":      info.DailInfo,
		"tls.start":      info.TlsShakeStart,
		"tls.done":       info.TlsShakeEnd,
		"reused":         info.IsReused,
		"idle":           info.IsIdle,
		"request.header": info.RequestHeaderWriteTime,
	}
	tags := map[string]string{
		"hostname":      info.Hostname,
		"remoteaddress": info.RemoteAddress[1],
	}
	pt, err := client.NewPoint(
		"traceview",
		tags,
		fields,
		time.Now(),
	)
	bp.AddPoint(pt)
	err = t.Client.Write(bp)
	if err != nil {
		log.Println("post data to influxdb failed : ", err.Error())
		return
	}
}

func (t *TraceWorker) Work() {
	for {
		select {
		case info := <-t.Traceinfo:
			t.Reportor(info)
		}
	}
}
