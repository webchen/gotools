package util

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/webchen/gotools/base/conf"
	"github.com/webchen/gotools/base/jsontool"
	"github.com/webchen/gotools/help/logs"
	"github.com/webchen/gotools/help/tool/nettool"

	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

// QueryWithZipKin 加上ZIPKIN链接监控
func QueryWithZipKin(method string, url string, jsonMap map[string]interface{}) string {
	tracer := GetTracer()
	if tracer == nil {
		logs.Show("enable zip, but config is empty, doHTTP2 instead...")
		return doHTTP2(method, url, jsonMap)
	}

	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer, zipkinhttp.TagResponseSize(true),
	)

	client, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	client.Timeout = 1 * time.Second
	// 直接走官方默认的，否则会无法上报数据
	//client.Transport = transport
	if err != nil {
		logs.ErrorProcess(err, "unable to create client")
		return ""
	}
	router := mux.NewRouter()
	ts := httptest.NewServer(serverMiddleware(router))
	ts.URL = url
	defer ts.Close()
	req := &http.Request{}
	if method == "GET" {
		req, err = http.NewRequest(method, ts.URL, nil)
		if logs.ErrorProcess(err, "unable to create http GET request") {
			return ""
		}
	}
	if method == "POST" {
		str := jsontool.MarshalToString(jsonMap)
		req, err = http.NewRequest(method, ts.URL, bytes.NewBuffer([]byte(str)))
		if logs.ErrorProcess(err, "unable to create http POST request") {
			return ""
		}
		req.Header.Set("Content-Type", "application/json")
	}
	url1 := strings.Split(url, "//")
	url2 := strings.SplitN(url1[1], "/", 2)
	url3 := strings.Split(url2[1], "?")
	res, err := client.DoWithAppSpan(req, "/"+url3[0])
	if logs.ErrorProcess(err, "unable to do http request") {
		return ""
	}

	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return string(b)
}

// GetTracer 获取tracer对象
func GetTracer() *zipkin.Tracer {
	endPointURL := conf.GetConfig("zipkin.endPoint", "").(string)
	if endPointURL == "" {
		return nil
	}

	reporter := httpreporter.NewReporter(endPointURL)

	// set-up the local endpoint for our service
	endpoint, _ := zipkin.NewEndpoint(conf.GetConfig("zipkin.serviceName", "goZipkin").(string), nettool.GetLocalFirstIPStr())

	// set-up our sampling strategy
	sampler := zipkin.NewModuloSampler(1)

	// initialize the tracer
	tracer, _ := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(sampler),
	)
	return tracer
}
