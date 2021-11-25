package http

import (
	"net"
	"net/http"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

var _ require.ModuleLoader = HttpModuleLoader

func HttpModuleLoader(vm *goja.Runtime, module *goja.Object) {
	export := module.Get("export").ToObject(vm)

	export.DefineDataProperty("Agent", vm.ToValue(
		func(c *goja.ConstructorCall, vm *goja.Runtime) *goja.Object {
			agent := http_Agent{
				Transport: *http.DefaultTransport.(*http.Transport),
			}
			if len(c.Arguments) > 0 {
				obj := c.Arguments[0].ToObject(vm)

				agent.Set("keepAlive", obj.Get("keepAlive"))
				agent.Set("keepAliveMsecs", obj.Get("keepAliveMsecs"))
				agent.Set("maxSockets", obj.Get("maxSockets"))
				agent.Set("maxTotalSockets", obj.Get("maxTotalSockets"))
				agent.Set("maxFreeSockets", obj.Get("maxFreeSockets"))
				agent.Set("scheduling", obj.Get("scheduling"))
				agent.Set("timeout", obj.Get("timeout"))
			}
			return nil
		}), 1, 1, 1)
}

var _ goja.DynamicObject = (*http_Agent)(nil)

type http_Agent struct {
	Client    *http.Client
	Transport http.Transport
	vm        *goja.Runtime

	keepAlive,
	keepAliveMsecs,
	maxSockets,
	maxTotalSockets,
	maxFreeSockets,
	scheduling,
	timeout goja.Value
}

func (a *http_Agent) Get(key string) goja.Value {
	switch key {
	case "keepAlive":
		return a.keepAlive
	case "keepAliveMsecs":
		return a.keepAliveMsecs
	case "maxSockets":
		return a.maxSockets
	case "maxTotalSockets":
		return a.maxTotalSockets
	case "maxFreeSockets":
		return a.maxFreeSockets
	case "scheduling":
		return a.scheduling
	case "timeout":
		return a.timeout

	case "createConnection":
		return a.vm.ToValue(func(arg *goja.FunctionCall) goja.Value {
			return goja.Undefined()
		})

	case "keepSocketAlive":
		return a.vm.ToValue(func(arg *goja.FunctionCall) goja.Value {
			return goja.Undefined()
		})

	case "reuseSocket":
		return a.vm.ToValue(func(arg *goja.FunctionCall) goja.Value {
			return goja.Undefined()
		})

	case "destroy":
		return a.vm.ToValue(func(arg *goja.FunctionCall) goja.Value {
			if a.Client != nil {
				a.Client.CloseIdleConnections()
			}
			return goja.Undefined()
		})
	case "getName":
	case "freeSockets", "requests", "sockets":
		return goja.Null()
	}
	return goja.Undefined()
}

func (agent *http_Agent) Set(key string, val goja.Value) bool {

	switch key {
	case "keepAlive":
		agent.keepAlive = val
		if v, ok := agent.keepAlive.Export().(bool); ok {
			agent.Transport.DisableKeepAlives = !v
		}
	case "keepAliveMsecs":
		agent.keepAliveMsecs = val
		if v, ok := agent.keepAliveMsecs.Export().(int64); ok {
			agent.Transport.IdleConnTimeout = time.Duration(v) * time.Millisecond
		}
	case "maxSockets":
		agent.maxSockets = val
		if v, ok := agent.maxSockets.Export().(int64); ok {
			agent.Transport.MaxConnsPerHost = int(v)
		}
	case "maxTotalSockets":
		agent.maxTotalSockets = val
		if v, ok := agent.maxTotalSockets.Export().(int64); ok {
			agent.Transport.MaxIdleConns = int(v)
		}
	case "maxFreeSockets":
		agent.maxFreeSockets = val
		if v, ok := agent.maxFreeSockets.Export().(int64); ok {
			agent.Transport.MaxIdleConnsPerHost = int(v)
		}
	case "scheduling":
		agent.scheduling = val
	case "timeout":
		agent.timeout = val
		if v, ok := agent.timeout.Export().(int64); ok {
			agent.Transport.DialContext = (&net.Dialer{
				Timeout:   time.Duration(v) * time.Millisecond,
				KeepAlive: 30 * time.Second,
			}).DialContext
		}
	default:
		return false
	}
	return true
}

func (a *http_Agent) Has(key string) bool {
	return false
}

func (a *http_Agent) Delete(key string) bool {
	return false
}

func (a *http_Agent) Keys() []string {
	return []string{
		"keepAlive",
		"keepAliveMsecs",
		"maxSockets",
		"maxTotalSockets",
		"maxFreeSockets",
		"scheduling",
		"timeout",
		"createConnection"}
}
