package http

import (
	"context"
	"crypto/tls"
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

	export.Set("get", vm.ToValue(http_get))
}

func http_get(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {}

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

type http_clientrequest struct {
	this *goja.Object

	context   context.Context
	cancel    context.CancelFunc
	destroyed bool

	request  *http.Request
	response *http.Response
	client   *http.Client

	listeners map[string]goja.Callable

	vm *goja.Runtime
}

func (c *http_clientrequest) Get(key string) goja.Value {
	switch key {
	case "on":
		return c.vm.ToValue(c.addeventListener)
	case "end":
		return c.vm.ToValue(c.end)

	case "destroy":
		return c.vm.ToValue(c.destroy)
	case "destroyed":
		return c.vm.ToValue(c.destroyed)
	}
	return goja.Undefined()
}

func (c *http_clientrequest) Set(key string, val goja.Value) bool

func (c *http_clientrequest) Delete(key string) bool {
	return false
}

func (c *http_clientrequest) Has(key string) bool {
	for _, v := range c.Keys() {
		if v == key {
			return true
		}
	}
	return false
}

func (c *http_clientrequest) Keys() []string {
	return []string{
		"on",
		"end",
		"destroy",
	}
}

func (c *http_clientrequest) addeventListener(arg goja.FunctionCall) goja.Value {
	if len(arg.Arguments) > 1 {
		callback, ok := goja.AssertFunction(arg.Arguments[1])
		if !ok {
			return goja.Undefined()
		}
		switch arg.Arguments[0].String() {
		case "connect":
			d := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			var conn net.Conn
			var err error
			switch c.request.URL.Scheme {
			case "https":
				conn, err = tls.DialWithDialer(d, "tcp", c.request.URL.Hostname()+":https", nil)

			case "http":
				if c.request.URL.Port() == "" {
					conn, err = net.Dial("tcp", c.request.URL.Host+":http")
				} else {
					conn, err = net.Dial("tcp", c.request.URL.Host)
				}

			}
		}
	}
	return goja.Undefined()
}

func (c *http_clientrequest) end(arg goja.FunctionCall) goja.Value {
	return c.this
}

func (c *http_clientrequest) destroy(arg goja.FunctionCall) goja.Value {
	c.cancel()
	c.destroyed = true
	return c.this
}
