package main

import (
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const clusterName = "opa"

func main() {
	proxywasm.SetNewHttpContext(newHTTPContext)
	proxywasm.SetNewRootContext(newRootContext)
}

type httpContext struct {
	proxywasm.DefaultHttpContext
	contextID uint32
}

type rootContext struct {
	proxywasm.DefaultRootContext
}

func newHTTPContext(rootContextID, contextID uint32) proxywasm.HttpContext {
	return &httpContext{contextID: contextID}
}

func newRootContext(contextID uint32) proxywasm.RootContext {
	return &rootContext{}
}

func (ctx *rootContext) OnPluginStart(pluginConfigurationSize int) bool {

	data, err := proxywasm.GetPluginConfiguration(pluginConfigurationSize)
	if err != nil {
		proxywasm.LogCriticalf("unable to get plugin configuration: %v", err)
	}

	proxywasm.LogInfof("plugin configuration: %s", string(data))

	return true
}

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {

	requestID, err := proxywasm.GetHttpRequestHeader("x-request-id")
	if err != nil {
		proxywasm.LogCriticalf("failed to get x-request-id from headers: %v", err)
		return types.ActionContinue
	}

	hs := [][2]string{
		{":authority", "opa:8181"},
		{":path", "/v1/data/envoy/authz"},
		{":method", "POST"},
		{"accept", "*/*"},
		{"user-agent", "proxy-wasm"},
		{"x-request-id", requestID},
	}

	authz, err := proxywasm.GetHttpRequestHeader("authorization")
	if err != nil {
		proxywasm.LogCriticalf("failed to get authz from headers: %v", err)
		return types.ActionContinue
	}

	payload := fmt.Sprintf(`{"input":{"authorization": "%s"}"}`, authz)
	if _, err := proxywasm.DispatchHttpCall(
		clusterName, hs, payload, [][2]string{}, 5000, httpCallResponseCallback,
	); err != nil {
		proxywasm.LogCriticalf("dispatch httpcall failed: %v", err)
		return types.ActionContinue
	}

	return types.ActionPause
}

func httpCallResponseCallback(numHeaders, bodySize, numTrailers int) {

	b, err := proxywasm.GetHttpCallResponseBody(0, bodySize)
	if err != nil {
		proxywasm.LogCriticalf("failed to get response body: %v", err)
		proxywasm.ResumeHttpRequest()
		return
	}

	allowed, err := jsonparser.GetBoolean(b, "result", "allow")
	if err != nil {
		proxywasm.LogCriticalf("failed to parse allowed filed from response body: %v", err)
		proxywasm.ResumeHttpRequest()
		return
	}

	proxywasm.LogInfof("allowed: %v", allowed)

	proxywasm.ResumeHttpRequest()
}
