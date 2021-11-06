package std

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDHttp(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdhttp"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdHttpFuncs := []stdFuncInfo{
		{
			Name:   "do",
			Getter: getStdHttpDo,
		},
		{
			Name:   "do-with",
			Getter: getStdHttpDoWith,
		},
		{
			Name:   "mux",
			Getter: getStdHttpMux,
		},
		{
			Name:   "shutdown",
			Getter: getStdHttpShutdown,
		},
		{
			Name:   "reg-handler",
			Getter: getStdHttpRegHandler,
		},
		{
			Name:   "listen-and-serve",
			Getter: getStdHttpListenAndServe,
		},
		{
			Name:   "listen-and-serve-tls",
			Getter: getStdHttpListenAndServeTLS,
		},
		{
			Name:   "write-response",
			Getter: getStdHttpWriteResponse,
		},
		{
			Name:   "set-response-header",
			Getter: getStdHttpSetResponseHeader,
		},
		{
			Name:   "add-response-header",
			Getter: getStdHttpAddResponseHeader,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdHttpFuncs, interpreter)
	return
}

type OpaqueHttpMux struct {
	muxv   *http.ServeMux
	server *http.Server
}

func (mux *OpaqueHttpMux) TypeName() string {
	return "http-mux"
}

func (mux *OpaqueHttpMux) Str() string {
	return fmt.Sprintf("%s: %#v, %#v", mux.TypeName(), mux.muxv, mux.server)
}

func (mux *OpaqueHttpMux) Equals(with funl.OpaqueAPI) bool {
	_, ok := with.(*OpaqueHttpMux)
	if !ok {
		return false
	}
	return true
}

func applyForEachKeyVal(frame *funl.Frame, name string, mapVal funl.Value, handler func(keyStr, valStr string)) {
	keyvals := funl.HandleKeyvalsOP(frame, []*funl.Item{&funl.Item{Type: funl.ValueItem, Data: mapVal}})
	kvListIter := funl.NewListIterator(keyvals)
	for {
		nextKV := kvListIter.Next()
		if nextKV == nil {
			break
		}
		kvIter := funl.NewListIterator(*nextKV)
		keyv := *(kvIter.Next())
		valv := *(kvIter.Next())
		if keyv.Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: header key not a string: %v", name, keyv)
		}
		switch valv.Kind {
		case funl.StringValue:
			handler(keyv.Data.(string), valv.Data.(string))
		case funl.ListValue:
			valListIter := funl.NewListIterator(valv)
			for {
				nextHeaVal := valListIter.Next()
				if nextHeaVal == nil {
					break
				}
				if nextHeaVal.Kind != funl.StringValue {
					funl.RunTimeError2(frame, "%s: header value in list not a string: %v", name, nextHeaVal)
				}
				handler(keyv.Data.(string), nextHeaVal.Data.(string))
			}
		default:
			funl.RunTimeError2(frame, "%s: header value not a string or list: %v", name, valv)
		}
	}
}

func getDoHandler(name string, argAdd int, infoHandler func(*funl.Frame, funl.Value) *http.Client) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if (l != (4 + argAdd)) && (l != (3 + argAdd)) {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}

		client := infoHandler(frame, arguments[0])

		if arguments[0+argAdd].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		method := arguments[0+argAdd].Data.(string)
		if arguments[1+argAdd].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		url := arguments[1+argAdd].Data.(string)

		if arguments[2+argAdd].Kind != funl.MapValue {
			funl.RunTimeError2(frame, "%s: requires map value", name)
		}

		var body io.Reader
		if l == (4 + argAdd) {
			bodyData, convOk := arguments[3+argAdd].Data.(*OpaqueByteArray)
			if !convOk {
				funl.RunTimeError2(frame, "%s: last argument not byte array", name)
			}
			body = bytes.NewBuffer(bodyData.data)
		}

		req, err := http.NewRequest(method, url, body)
		if err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: error in creating request (%v)", name, err)}
			return
		}

		// fill header from map
		headerAdder := func(keyStr, valStr string) {
			req.Header.Add(keyStr, valStr)
		}
		applyForEachKeyVal(frame, name, arguments[2+argAdd], headerAdder)

		// do request
		response, err := client.Do(req)
		if err != nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("%s: error in HTTP (%v)", name, err)}
			return
		}
		defer response.Body.Close()

		statusCode := response.StatusCode
		statusText := response.Status

		// lets create map for response header
		headerMapval := funl.HandleMapOP(frame, []*funl.Item{})
		for hk, hv := range response.Header {
			var hvalues []funl.Value
			for _, hit := range hv {
				hvalues = append(hvalues, funl.Value{Kind: funl.StringValue, Data: hit})
			}
			keyv := funl.Value{Kind: funl.StringValue, Data: hk}
			valv := funl.MakeListOfValues(frame, hvalues)
			putArgs := []*funl.Item{
				&funl.Item{Type: funl.ValueItem, Data: headerMapval},
				&funl.Item{Type: funl.ValueItem, Data: keyv},
				&funl.Item{Type: funl.ValueItem, Data: valv},
			}
			headerMapval = funl.HandlePutOP(frame, putArgs)
			if headerMapval.Kind != funl.MapValue {
				funl.RunTimeError2(frame, "%s: failed to put key-value to map", name)
			}
		}

		respKVs := []*funl.Item{
			// status code
			{
				Type: funl.ValueItem,
				Data: funl.Value{
					Kind: funl.StringValue,
					Data: "status-code",
				},
			},
			{
				Type: funl.ValueItem,
				Data: funl.Value{
					Kind: funl.IntValue,
					Data: statusCode,
				},
			},
			// status text
			{
				Type: funl.ValueItem,
				Data: funl.Value{
					Kind: funl.StringValue,
					Data: "status-text",
				},
			},
			{
				Type: funl.ValueItem,
				Data: funl.Value{
					Kind: funl.StringValue,
					Data: statusText,
				},
			},
			// header map
			{
				Type: funl.ValueItem,
				Data: funl.Value{
					Kind: funl.StringValue,
					Data: "header",
				},
			},
			{
				Type: funl.ValueItem,
				Data: headerMapval,
			},
		}

		// get body if there is one
		if body, err := ioutil.ReadAll(response.Body); err == nil {
			bodyAsByteArr := &OpaqueByteArray{data: body}
			respKVs = append(respKVs, &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: "body"}})
			respKVs = append(respKVs, &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.OpaqueValue, Data: bodyAsByteArr}})
		}

		retVal = funl.HandleMapOP(frame, respKVs)

		//fmt.Println(fmt.Sprintf("response : %#v", response))
		return
	}
}

func getStdHttpDo(name string) stdFuncType {
	infoHandler := func(frame *funl.Frame, v funl.Value) *http.Client {
		return &http.Client{}
	}
	return getDoHandler(name, 0, infoHandler)
}

func getStdHttpDoWith(name string) stdFuncType {
	infoHandler := func(frame *funl.Frame, v funl.Value) *http.Client {
		if v.Kind != funl.MapValue {
			funl.RunTimeError2(frame, "%s: requires map value", name)
		}

		client := &http.Client{}

		keyvals := funl.HandleKeyvalsOP(frame, []*funl.Item{&funl.Item{Type: funl.ValueItem, Data: v}})
		kvListIter := funl.NewListIterator(keyvals)
		for {
			nextKV := kvListIter.Next()
			if nextKV == nil {
				break
			}
			kvIter := funl.NewListIterator(*nextKV)
			keyv := *(kvIter.Next())
			valv := *(kvIter.Next())
			if keyv.Kind != funl.StringValue {
				funl.RunTimeError2(frame, "%s: info key not a string: %v", name, keyv)
			}
			switch keyStr := keyv.Data.(string); keyStr {
			case "timeout":
				if valv.Kind != funl.IntValue {
					funl.RunTimeError2(frame, "%s: %s value not int: %v", name, keyStr, keyv)
				}
				client.Timeout = time.Duration(valv.Data.(int))

			case "insecure":
				if valv.Kind != funl.BoolValue {
					funl.RunTimeError2(frame, "%s: %s value not bool: %v", name, keyStr, keyv)
				}
				tr := &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: valv.Data.(bool)},
				}
				client.Transport = tr
			}
		}

		return client
	}
	return getDoHandler(name, 1, infoHandler)
}

func getStdHttpListenAndServeTLS(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if (l != 3) && (l != 4) {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}

		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		hm, ok := arguments[0].Data.(*OpaqueHttpMux)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not mux value", name)
		}

		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value for 2nd arg.", name)
		}
		if arguments[2].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value for 3rd arg.", name)
		}
		certFile := arguments[1].Data.(string)
		keyFile := arguments[2].Data.(string)

		var addr string
		if l > 3 {
			if arguments[3].Kind != funl.StringValue {
				funl.RunTimeError2(frame, "%s: requires string value for 4th arg.", name)
			}
			addr = arguments[3].Data.(string)
		} else {
			addr = ":https"
		}

		hm.server = &http.Server{
			Addr:    addr,
			Handler: hm.muxv,
		}

		err := hm.server.ListenAndServeTLS(certFile, keyFile)
		var errText string
		if err == nil {
			errText = ""
		} else {
			errText = err.Error()
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: errText}
		return
	}
}

func getStdHttpListenAndServe(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}

		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		hm, ok := arguments[0].Data.(*OpaqueHttpMux)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not mux value", name)
		}

		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		addr := arguments[1].Data.(string)

		hm.server = &http.Server{
			Addr:    addr,
			Handler: hm.muxv,
		}

		err := hm.server.ListenAndServe()
		var errText string
		if err == nil {
			errText = ""
		} else {
			errText = err.Error()
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: errText}
		return
	}
}

func getStdHttpRegHandler(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 3 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}

		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		hm, ok := arguments[0].Data.(*OpaqueHttpMux)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not mux value", name)
		}

		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		pattern := arguments[1].Data.(string)

		if arguments[2].Kind != funl.FunctionValue {
			funl.RunTimeError2(frame, "%s: requires func/proc value", name)
		}

		handlerWrapper := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					err, _ := r.(error)
					fmt.Printf("\nHTTP handler died: %v\n", err)
				}
			}()

			// lets read query parameters
			queryMapval := funl.HandleMapOP(frame, []*funl.Item{})
			for hk, hv := range r.URL.Query() {
				var hvalues []funl.Value
				for _, hit := range hv {
					hvalues = append(hvalues, funl.Value{Kind: funl.StringValue, Data: hit})
				}
				keyv := funl.Value{Kind: funl.StringValue, Data: hk}
				valv := funl.MakeListOfValues(frame, hvalues)
				putArgs := []*funl.Item{
					&funl.Item{Type: funl.ValueItem, Data: queryMapval},
					&funl.Item{Type: funl.ValueItem, Data: keyv},
					&funl.Item{Type: funl.ValueItem, Data: valv},
				}
				queryMapval = funl.HandlePutOP(frame, putArgs)
				if queryMapval.Kind != funl.MapValue {
					funl.RunTimeError2(frame, "%s: failed to put key-value to map", name)
				}
			}

			// lets create map for request header
			headerMapval := funl.HandleMapOP(frame, []*funl.Item{})
			for hk, hv := range r.Header {
				var hvalues []funl.Value
				for _, hit := range hv {
					hvalues = append(hvalues, funl.Value{Kind: funl.StringValue, Data: hit})
				}
				keyv := funl.Value{Kind: funl.StringValue, Data: hk}
				valv := funl.MakeListOfValues(frame, hvalues)
				putArgs := []*funl.Item{
					&funl.Item{Type: funl.ValueItem, Data: headerMapval},
					&funl.Item{Type: funl.ValueItem, Data: keyv},
					&funl.Item{Type: funl.ValueItem, Data: valv},
				}
				headerMapval = funl.HandlePutOP(frame, putArgs)
				if headerMapval.Kind != funl.MapValue {
					funl.RunTimeError2(frame, "%s: failed to put key-value to map", name)
				}
			}

			reqKVs := []*funl.Item{
				// URI
				{
					Type: funl.ValueItem,
					Data: funl.Value{
						Kind: funl.StringValue,
						Data: "URI",
					},
				},
				{
					Type: funl.ValueItem,
					Data: funl.Value{
						Kind: funl.StringValue,
						Data: r.RequestURI,
					},
				},
				// method
				{
					Type: funl.ValueItem,
					Data: funl.Value{
						Kind: funl.StringValue,
						Data: "method",
					},
				},
				{
					Type: funl.ValueItem,
					Data: funl.Value{
						Kind: funl.StringValue,
						Data: r.Method,
					},
				},
				// header map
				{
					Type: funl.ValueItem,
					Data: funl.Value{
						Kind: funl.StringValue,
						Data: "header",
					},
				},
				{
					Type: funl.ValueItem,
					Data: headerMapval,
				},
				// query map
				{
					Type: funl.ValueItem,
					Data: funl.Value{
						Kind: funl.StringValue,
						Data: "query",
					},
				},
				{
					Type: funl.ValueItem,
					Data: queryMapval,
				},
			}
			// get body if there is one
			if body, err := ioutil.ReadAll(r.Body); err == nil {
				bodyAsByteArr := &OpaqueByteArray{data: body}
				reqKVs = append(reqKVs, &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: "body"}})
				reqKVs = append(reqKVs, &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.OpaqueValue, Data: bodyAsByteArr}})
			}
			reqm := funl.HandleMapOP(frame, reqKVs)

			argsForCall := []*funl.Item{
				&funl.Item{
					Type: funl.ValueItem,
					Data: arguments[2],
				},
				&funl.Item{
					Type: funl.ValueItem,
					Data: funl.Value{
						Kind: funl.OpaqueValue,
						Data: &OpaqueResponseWriter{w: w},
					},
				},
				&funl.Item{
					Type: funl.ValueItem,
					Data: reqm,
				},
			}
			funl.HandleCallOP(frame, argsForCall)
		}

		handlerWrapperWrapper := func(w http.ResponseWriter, r *http.Request) {
			handlerWrapper(w, r)
		}

		hm.muxv.HandleFunc(pattern, handlerWrapperWrapper)
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		return
	}
}

type OpaqueResponseWriter struct {
	w http.ResponseWriter
}

func (orw *OpaqueResponseWriter) TypeName() string {
	return "http-response-writer"
}

func (orw *OpaqueResponseWriter) Str() string {
	return fmt.Sprintf("%s: %#v", orw.TypeName(), orw.w)
}

func (orw *OpaqueResponseWriter) Equals(with funl.OpaqueAPI) bool {
	_, ok := with.(*OpaqueResponseWriter)
	if !ok {
		return false
	}
	return true
}

func getStdHttpAddResponseHeader(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}

		// 1st arg is response writer
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		respWriter, ok := arguments[0].Data.(*OpaqueResponseWriter)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not response writer value", name)
		}

		// 2nd arg is map containing key values
		if arguments[1].Kind != funl.MapValue {
			funl.RunTimeError2(frame, "%s: requires map value", name)
		}

		// fill header from map
		respHeader := respWriter.w.Header()
		headerSetter := func(keyStr, valStr string) {
			respHeader.Add(keyStr, valStr)
		}
		applyForEachKeyVal(frame, name, arguments[1], headerSetter)
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		return
	}
}

func getStdHttpSetResponseHeader(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}

		// 1st arg is response writer
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		respWriter, ok := arguments[0].Data.(*OpaqueResponseWriter)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not response writer value", name)
		}

		// 2nd arg is map containing key values
		if arguments[1].Kind != funl.MapValue {
			funl.RunTimeError2(frame, "%s: requires map value", name)
		}

		// fill header from map
		respHeader := respWriter.w.Header()
		headerSetter := func(keyStr, valStr string) {
			respHeader.Set(keyStr, valStr)
		}
		applyForEachKeyVal(frame, name, arguments[1], headerSetter)
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		return
	}
}

func getStdHttpWriteResponse(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 3 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}

		// 1st arg is response writer
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		respWriter, ok := arguments[0].Data.(*OpaqueResponseWriter)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not response writer value", name)
		}

		// 2nd arg is status code
		if arguments[1].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value", name)
		}
		statusCode := arguments[1].Data.(int)

		// 3rd arg is body data
		if arguments[2].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		bodyByteArr, ok := arguments[2].Data.(*OpaqueByteArray)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not bytearray value", name)
		}

		bodyData := bodyByteArr.data
		retVal = funl.Value{Kind: funl.StringValue, Data: ""}
		bodyLen := len(bodyData)
		if (statusCode != http.StatusOK) || (bodyLen == 0) {
			respWriter.w.WriteHeader(statusCode)
		}
		if bodyLen > 0 {
			n, err := respWriter.w.Write(bodyData)
			if err != nil {
				retVal = funl.Value{Kind: funl.StringValue, Data: fmt.Sprintf("error in writing body (%d bytes written): %v", n, err.Error())}
			}
		}
		return
	}
}

func getStdHttpShutdown(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}

		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		hm, ok := arguments[0].Data.(*OpaqueHttpMux)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not mux value", name)
		}

		if hm.server == nil {
			retVal = funl.Value{Kind: funl.StringValue, Data: "no http server found"}
			return
		}

		err := hm.server.Shutdown(context.Background())
		var errText string
		if err == nil {
			errText = ""
		} else {
			errText = err.Error()
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: errText}
		return
	}
}

func getStdHttpMux(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		muxv := OpaqueHttpMux{muxv: http.NewServeMux()}
		retVal = funl.Value{Kind: funl.OpaqueValue, Data: &muxv}
		return
	}
}
