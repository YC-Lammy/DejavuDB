exports.go = class {
    MethodGet     = "GET";
	MethodHead    = "HEAD";
	MethodPost    = "POST";
	MethodPut     = "PUT";
	MethodPatch   = "PATCH"; // RFC 5789
	MethodDelete  = "DELETE";
	MethodConnect = "CONNECT";
	MethodOptions = "OPTIONS";
	MethodTrace   = "TRACE";


    StatusContinue           = 100; // RFC 7231, 6.2.1
	StatusSwitchingProtocols = 101; // RFC 7231, 6.2.2
	StatusProcessing         = 102; // RFC 2518, 10.1
	StatusEarlyHints         = 103; // RFC 8297

	StatusOK                   = 200; // RFC 7231, 6.3.1
	StatusCreated              = 201; // RFC 7231, 6.3.2
	StatusAccepted             = 202; // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = 203; // RFC 7231, 6.3.4
	StatusNoContent            = 204; // RFC 7231, 6.3.5
	StatusResetContent         = 205; // RFC 7231, 6.3.6
	StatusPartialContent       = 206; // RFC 7233, 4.1
	StatusMultiStatus          = 207; // RFC 4918, 11.1
	StatusAlreadyReported      = 208; // RFC 5842, 7.1
	StatusIMUsed               = 226; // RFC 3229, 10.4.1

	StatusMultipleChoices  = 300; // RFC 7231, 6.4.1
	StatusMovedPermanently = 301; // RFC 7231, 6.4.2
	StatusFound            = 302; // RFC 7231, 6.4.3
	StatusSeeOther         = 303; // RFC 7231, 6.4.4
	StatusNotModified      = 304; // RFC 7232, 4.1
	StatusUseProxy         = 305; // RFC 7231, 6.4.5

	StatusTemporaryRedirect = 307; // RFC 7231, 6.4.7
	StatusPermanentRedirect = 308; // RFC 7538, 3

	StatusBadRequest                   = 400; // RFC 7231, 6.5.1
	StatusUnauthorized                 = 401; // RFC 7235, 3.1
	StatusPaymentRequired              = 402; // RFC 7231, 6.5.2
	StatusForbidden                    = 403; // RFC 7231, 6.5.3
	StatusNotFound                     = 404; // RFC 7231, 6.5.4
	StatusMethodNotAllowed             = 405; // RFC 7231, 6.5.5
	StatusNotAcceptable                = 406; // RFC 7231, 6.5.6
	StatusProxyAuthRequired            = 407; // RFC 7235, 3.2
	StatusRequestTimeout               = 408; // RFC 7231, 6.5.7
	StatusConflict                     = 409; // RFC 7231, 6.5.8
	StatusGone                         = 410; // RFC 7231, 6.5.9
	StatusLengthRequired               = 411; // RFC 7231, 6.5.10
	StatusPreconditionFailed           = 412; // RFC 7232, 4.2
	StatusRequestEntityTooLarge        = 413; // RFC 7231, 6.5.11
	StatusRequestURITooLong            = 414; // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = 415; // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable = 416; // RFC 7233, 4.4
	StatusExpectationFailed            = 417; // RFC 7231, 6.5.14
	StatusTeapot                       = 418; // RFC 7168, 2.3.3
	StatusMisdirectedRequest           = 421; // RFC 7540, 9.1.2
	StatusUnprocessableEntity          = 422; // RFC 4918, 11.2
	StatusLocked                       = 423; // RFC 4918, 11.3
	StatusFailedDependency             = 424; // RFC 4918, 11.4
	StatusTooEarly                     = 425; // RFC 8470, 5.2.
	StatusUpgradeRequired              = 426; // RFC 7231, 6.5.15
	StatusPreconditionRequired         = 428; // RFC 6585, 3
	StatusTooManyRequests              = 429; // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = 431; // RFC 6585, 5
	StatusUnavailableForLegalReasons   = 451; // RFC 7725, 3

	StatusInternalServerError           = 500; // RFC 7231, 6.6.1
	StatusNotImplemented                = 501; // RFC 7231, 6.6.2
	StatusBadGateway                    = 502; // RFC 7231, 6.6.3
	StatusServiceUnavailable            = 503; // RFC 7231, 6.6.4
	StatusGatewayTimeout                = 504; // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = 505; // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = 506; // RFC 2295, 8.1
	StatusInsufficientStorage           = 507; // RFC 4918, 11.5
	StatusLoopDetected                  = 508; // RFC 5842, 7.2
	StatusNotExtended                   = 510; // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511; // RFC 6585, 6


    DefaultMaxHeaderBytes = 1 << 20; // 1 MB

    DefaultMaxIdleConnsPerHost = 2;

    TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT";

    TrailerPrefix = "Trailer:";

    static Error(ResponseWriter, error, code){
        http_Error(ResponseWriter, error, code);
    }

    static Handle(pattern, handler){
        http_Handle(pattern, handler);
    }

    static HandleFunc(pattern, handler){
        http_HandleFunc(pattern, handler);
    }

    static ListenAndServe(addr, handler){
        http_ListenAndServe(addr, handler);
    }

    static ListenAndServeTLS(addr, certFile, keyFile, handler){
        ListenAndServeTLS(addr, certFile, keyFile, handler);
    }

    static MaxBytesReader(ResponseWriter,ReadCloser, int){
        http_MaxBytesReader(ResponseWriter,ReadCloser, int);
    }

    static NotFound(ResponseWriter, Request){
        http_NotFound(ResponseWriter, Request);
    }

    static ParseHTTPVersion(vers){
        http_ParseHTTPVersion(vers);
    }

    static ParseTime(text){
        http_ParseTime(text);
    }

    static ProxyFromEnvironment(Request){
        http_ProxyFromEnvironment(Request);
    }

    static ProxyURL(fixedURL){
        http_ProxyURL(fixedURL);
    }

    static Redirect(ResponseWriter, Request, url , code){
        http_Redirect(ResponseWriter, Request, url , code);
    }

    static Serve(Listener, handler){
        http_Serve(Listener, handler);
    }

    static ServeContent(ResponseWriter,Request, name, modtime, content){
        http_ServeContent(ResponseWriter,Request, name, modtime, content);
    }

    static ServeFile(ResponseWriter,Request, name){
        http_ServeFile(ResponseWriter,Request, name);
    }

    static ServeTLS(Listener, handle, certFile, keyFile){
        http_ServeTLS(Listener, handle, certFile, keyFile);
    }

    static SetCookie(ResponseWriter, cookie){
        http_SetCookie(ResponseWriter, cookie);
    }

    static StatusText(code){
        http_StatusText(code);
    }

}

exports.Get = function(url){
    return call_go_fn("http_Get", url);
}

exports.Head = function(url){
    return call_go_fn("http_Head", url);
}

exports.Post = function(url,contentType, data){
    return call_go_fn("http_Post", url, contentType, data);
}

exports.PostForm = function(url,jsonString){
    return call_go_fn("http_PostForm", url,jsonString);
}

exports.CanonicalHeaderKey= function(s){
    return call_go_fn("http_CanonicalHeaderKey",s);
}

exports.DetectContentType = function(data){
    return call_go_fn("http_DetectContentType",data);
}

exports.Request = class{
    constructor(method,url,body){
        this.method = method
        this.url = url
        this.body = body
    }
}

exports.NewRequest = function(method,url,body){
    return new exports.Request(method,url,body);
}

exports.Do = function(req){
    call_go_fn("http_Client_Do",req.method,req.url, req.body)
}
