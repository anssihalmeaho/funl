package std

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/anssihalmeaho/funl/funl"
)

/*
Todo
- options map to new-proxy (timeouts)
- TSL support
*/

func initSTDRPC() (err error) {
	stdModuleName := "stdrpc"
	topFrame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	stdRPCFuncs := []stdFuncInfo{
		{
			Name:   "new-server",
			Getter: getRPCNewServer,
		},
		{
			Name:   "register",
			Getter: getRPCRegister,
		},
		{
			Name:   "unregister",
			Getter: getRPCUnRegister,
		},
		{
			Name:   "new-proxy",
			Getter: getRPCNewExtProxy,
		},
		{
			Name:   "rcall",
			Getter: getRPCRCall,
		},
		{
			Name:   "close",
			Getter: getRPCRClose,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdRPCFuncs)
	return
}

func getRPCNewServer(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string argument", name)
		}

		addr := arguments[0].Data.(string)
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: true,
			},
			{
				Kind: funl.StringValue,
				Data: "",
			},
			funl.Value{Kind: funl.OpaqueValue, Data: NewRServer(addr)},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getRPCUnRegister(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if (l != 1) && (l != 2) {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: assuming opaque argument", name)
		}
		if (l == 2) && arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string argument", name)
		}

		server := arguments[0].Data.(*RServer)
		var err error
		if l == 1 {
			err = server.UnRegisterAll()
		} else {
			err = server.UnRegister(arguments[1].Data.(string))
		}
		var errText string
		if err != nil {
			errText = err.Error()
		}
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: err == nil,
			},
			{
				Kind: funl.StringValue,
				Data: errText,
			},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getRPCRegister(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l != 3 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: assuming opaque argument", name)
		}
		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string argument", name)
		}
		switch arguments[2].Kind {
		case funl.FunctionValue, funl.ExtProcValue:
		default:
			funl.RunTimeError2(frame, "%s: assuming proc/func as argument", name)
		}

		server := arguments[0].Data.(*RServer)
		err := server.Register(arguments[1].Data.(string), arguments[2])
		var errText string
		if err != nil {
			errText = err.Error()
		}
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: err == nil,
			},
			{
				Kind: funl.StringValue,
				Data: errText,
			},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getRPCNewExtProxy(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string argument", name)
		}
		addr := arguments[0].Data.(string)
		retVal = funl.Value{Kind: funl.OpaqueValue, Data: NewExtProxy(addr)}
		return
	}
}

func makeErrorReturnList(frame *funl.Frame, err error) funl.Value {
	values := []funl.Value{
		{
			Kind: funl.BoolValue,
			Data: false,
		},
		{
			Kind: funl.StringValue,
			Data: err.Error(),
		},
		{
			Kind: funl.StringValue,
			Data: "",
		},
	}
	return funl.MakeListOfValues(frame, values)
}

func getRPCRCall(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l < 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: assuming opaque argument", name)
		}
		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string argument", name)
		}
		proxy := arguments[0].Data.(*RProxy)
		rprocName := arguments[1].Data.(string)
		arsgList := funl.MakeListOfValues(frame, arguments[2:])
		return proxy.MakeRemoteCall(frame, rprocName, arsgList)
	}
}

func getRPCRClose(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: assuming opaque argument", name)
		}

		server := arguments[0].Data.(*RServer)
		err := server.Close()
		var errText string
		if err != nil {
			errText = err.Error()
		}
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: err == nil,
			},
			{
				Kind: funl.StringValue,
				Data: errText,
			},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

// Request ...
type Request struct {
	RPName string
	Args   string
}

// Reply ...
type Reply struct {
	ErrorDescr string
	RetValue   string
}

// ServeHTTP ...
func (server *RServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeErrorResponse := func(err error, rpName string) {
		errVal := funl.Value{Kind: funl.StringValue, Data: ""}
		retvStr := encode(server.TopFrame, server.EncoderVal, errVal)
		reply := Reply{
			ErrorDescr: fmt.Sprintf("Error: %s: %v", rpName, err),
			RetValue:   retvStr,
		}
		dataBytes, err := json.Marshal(reply)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		w.Write(dataBytes)
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(err, "")
		return
	}
	var req Request
	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		writeErrorResponse(err, "")
		return
	}
	receivedListValue := decode(server.TopFrame, server.DecoderVal, req.Args)

	server.ProcsMutex.RLock()
	procVal, found := server.Procs[req.RPName]
	server.ProcsMutex.RUnlock()

	if !found {
		writeErrorResponse(err, req.RPName)
		return
	}

	callArgs := []*funl.Item{
		&funl.Item{
			Type: funl.ValueItem,
			Data: procVal,
		},
	}
	lit := funl.NewListIterator(receivedListValue)
	for {
		nextv := lit.Next()
		if nextv == nil {
			break
		}
		nextItem := &funl.Item{
			Type: funl.ValueItem,
			Data: *nextv,
		}
		callArgs = append(callArgs, nextItem)
	}

	returnedVal, pcallErr := func() (returnValue funl.Value, callErr error) {
		defer func() {
			if r := recover(); r != nil {
				var rtestr string
				if err, isError := r.(error); isError {
					rtestr = err.Error()
				}
				callErr = fmt.Errorf("procedure made RTE: %s", rtestr)
			}
		}()
		return funl.HandleCallOP(server.TopFrame, callArgs), nil
	}()
	if pcallErr != nil {
		writeErrorResponse(pcallErr, req.RPName)
		return
	}

	retvStr := encode(server.TopFrame, server.EncoderVal, returnedVal)
	reply := Reply{
		ErrorDescr: "",
		RetValue:   retvStr,
	}
	dataBytes, err := json.Marshal(reply)
	if err != nil {
		writeErrorResponse(err, req.RPName)
		return
	}
	w.Write(dataBytes)
}

// MakeRemoteCall ...
func (proxy *RProxy) MakeRemoteCall(frame *funl.Frame, rprocName string, arsgList funl.Value) (retVal funl.Value) {
	req := Request{
		RPName: rprocName,
		Args:   encode(frame, proxy.EncoderVal, arsgList),
	}
	dataBytes, err := json.Marshal(req)
	if err != nil {
		return makeErrorReturnList(frame, err)
	}

	url := fmt.Sprintf("http://%s/rpc", proxy.Addr)
	contentType := "application/json"
	resp, err := proxy.Client.Post(url, contentType, bytes.NewBuffer(dataBytes))
	if err != nil {
		return makeErrorReturnList(frame, err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return makeErrorReturnList(frame, err)
	}
	var reply Reply
	err = json.Unmarshal(respBody, &reply)
	if err != nil {
		return makeErrorReturnList(frame, err)
	}
	if reply.ErrorDescr != "" {
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: false,
			},
			{
				Kind: funl.StringValue,
				Data: reply.ErrorDescr,
			},
			{
				Kind: funl.StringValue,
				Data: "",
			},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
	returnedValue := decode(frame, proxy.DecoderVal, reply.RetValue)

	values := []funl.Value{
		{
			Kind: funl.BoolValue,
			Data: true,
		},
		{
			Kind: funl.StringValue,
			Data: "",
		},
		returnedValue,
	}
	retVal = funl.MakeListOfValues(frame, values)
	return
}

// RProxy ...
type RProxy struct {
	Addr       string
	Client     *http.Client
	DecoderVal funl.Value
	EncoderVal funl.Value
}

// TypeName gives type name
func (proxy *RProxy) TypeName() string {
	return "rproxy"
}

// Str returs value as string
func (proxy *RProxy) Str() string {
	return fmt.Sprintf("rproxy:%s", proxy.Addr)
}

// Equals returns equality
func (proxy *RProxy) Equals(with funl.OpaqueAPI) bool {
	return false
}

// NewExtProxy ...
func NewExtProxy(addr string) *RProxy {
	return &RProxy{
		Addr:       addr,
		Client:     &http.Client{},
		DecoderVal: getDecoder(),
		EncoderVal: getEncoder(),
	}
}

func getDecoder() funl.Value {
	frame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	frame.SetInProcCall(true)

	decItem := &funl.Item{
		Type: funl.ValueItem,
		Data: funl.Value{
			Kind: funl.StringValue,
			Data: "call(proc() import stdser import stdbytes proc(__s) b = call(stdbytes.str-to-bytes __s) _ _ __v = call(stdser.decode b): __v end end)",
		},
	}
	return funl.HandleEvalOP(frame, []*funl.Item{decItem})
}

func getEncoder() funl.Value {
	frame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	frame.SetInProcCall(true)

	encItem := &funl.Item{
		Type: funl.ValueItem,
		Data: funl.Value{
			Kind: funl.StringValue,
			Data: "call(proc() import stdser import stdbytes proc(x) _ _ b = call(stdser.encode x): call(stdbytes.string b) end end)",
		},
	}
	return funl.HandleEvalOP(frame, []*funl.Item{encItem})
}

// RServer ...
type RServer struct {
	Addr            string
	Server          *http.Server
	Procs           map[string]funl.Value
	ProcsMutex      sync.RWMutex
	DecoderVal      funl.Value
	EncoderVal      funl.Value
	TopFrame        *funl.Frame
	IdleConnsClosed chan struct{}
}

func encode(frame *funl.Frame, encoderVal funl.Value, val funl.Value) string {
	encArgs := []*funl.Item{
		&funl.Item{
			Type: funl.ValueItem,
			Data: encoderVal,
		},
		&funl.Item{
			Type: funl.ValueItem,
			Data: val,
		},
	}
	res := funl.HandleCallOP(frame, encArgs)
	return res.Data.(string)
}

func decode(frame *funl.Frame, decoderVal funl.Value, s string) funl.Value {
	decArgs := []*funl.Item{
		&funl.Item{
			Type: funl.ValueItem,
			Data: decoderVal,
		},
		&funl.Item{
			Type: funl.ValueItem,
			Data: funl.Value{Kind: funl.StringValue, Data: s},
		},
	}
	return funl.HandleCallOP(frame, decArgs)
}

// UnRegisterAll ...
func (server *RServer) UnRegisterAll() error {
	server.ProcsMutex.Lock()
	defer server.ProcsMutex.Unlock()

	server.Procs = make(map[string]funl.Value)
	return nil
}

// UnRegister ...
func (server *RServer) UnRegister(procName string) error {
	server.ProcsMutex.Lock()
	defer server.ProcsMutex.Unlock()

	delete(server.Procs, procName)
	return nil
}

// Register ...
func (server *RServer) Register(name string, procVal funl.Value) error {
	server.ProcsMutex.Lock()
	defer server.ProcsMutex.Unlock()

	_, found := server.Procs[name]
	if found {
		return fmt.Errorf("Proc %s already registered", name)
	}
	server.Procs[name] = procVal
	return nil
}

// NewRServer ...
func NewRServer(addr string) *RServer {
	frame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	frame.SetInProcCall(true)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	rserver := &RServer{
		Addr:            addr,
		Server:          server,
		DecoderVal:      getDecoder(),
		EncoderVal:      getEncoder(),
		TopFrame:        frame,
		Procs:           make(map[string]funl.Value),
		IdleConnsClosed: make(chan struct{}),
	}
	mux.Handle("/rpc", rserver)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
		}
		close(rserver.IdleConnsClosed)
	}()
	return rserver
}

// Close ...
func (server *RServer) Close() error {
	err := server.Server.Shutdown(context.Background())
	if err != nil {
		return err
	}
	<-server.IdleConnsClosed
	return nil
}

// TypeName gives type name
func (server *RServer) TypeName() string {
	return "rserver"
}

// Str returs value as string
func (server *RServer) Str() string {
	return fmt.Sprintf("rserver:%s", server.Addr)
}

// Equals returns equality
func (server *RServer) Equals(with funl.OpaqueAPI) bool {
	return false
}
