package httpjs

import (
	"io"
	"net/http"
	"os"

	"github.com/goccy/go-json"

	"rogchap.com/v8go"
)

var Clients = map[string]http.Client{}

var ConnStates = map[string]http.ConnState{}

func Handler(Args ...string, ctx *v8go.Context, deferfns *[]func()) *v8go.Value {
	vm := ctx.Isolate()
	switch Args[0] {
	case "http_CanonicalHeaderKey":
		val, _ := v8go.NewValue(vm, http.CanonicalHeaderKey(Args[1]))
		return val

	case "http_DetectContentType":
		val, _ := v8go.NewValue(vm, http.DetectContentType([]byte(Args[1])))
		return val

	case "http_Get", "http_Head", "http_Post", "http_PostForm":
		type result struct {
			Status     string // e.g. "200 OK"
			StatusCode int    // e.g. 200
			Proto      string // e.g. "HTTP/1.0"
			ProtoMajor int    // e.g. 1
			ProtoMinor int    // e.g. 0

			Header map[string][]string
			Body   []byte

			ContentLength    int64
			TransferEncoding []string
			Uncompressed     bool

			Trailer map[string][]string
		}
		var resp *http.Response
		var err error
		switch Args[0] {
		case "http_Get":
			resp, err = http.Get(Args[1])
		case "http_Head":
			resp, err = http.Head(Args[1])
		case "http_Post", "http_PostForm":

			f, _ := os.CreateTemp("", "")
			defer f.Close()
			var ctype string

			switch Args[0] {
			case "PostForm":
				f.Write([]byte(Args[2]))
				ctype = "application/json; charset=UTF-8"
			case "Post":
				f.Write([]byte(Args[3]))
				ctype = Args[2]
			}

			resp, err = http.Post(Args[1], ctype, f)

		}
		defer resp.Body.Close()

		if err != nil {
			errs <- err
			return nil
		}

		body, err := io.ReadAll(resp.Body)
		r := result{Status: resp.Status, StatusCode: resp.StatusCode, Proto: resp.Proto,
			ProtoMajor: resp.ProtoMajor, ProtoMinor: resp.ProtoMinor, Header: resp.Header,
			Body: body, ContentLength: resp.ContentLength, TransferEncoding: resp.TransferEncoding,
			Uncompressed: resp.Uncompressed, Trailer: resp.Trailer}

		barr, _ := json.Marshal(r)

		val, err := v8go.JSONParse(ctx, string(barr))
		if err != nil {
			errs <- err
			return nil
		}

		return val

	}
	return nil
}
