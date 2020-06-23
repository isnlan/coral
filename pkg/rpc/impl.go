package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

type methodType struct {
	method    reflect.Method
	argTypes  []reflect.Type
	replyType reflect.Type
}

type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

// rpcImpl represents an RPC implement.
type rpcImpl struct {
	serviceMap sync.Map // map[string]*service
}

// New returns a new RPC.
func New() Rpc {
	return &rpcImpl{}
}

// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
//	- exported method of exported type
//	- two arguments, both of exported type
//	- the second argument is a pointer
//	- one return value, of type error
// It returns an error if the receiver is not an exported type or has
// no suitable methods. It also logs the error using package log.
// The client accesses each method using a string of the form "Type.Method",
// where Type is the receiver's concrete type.
func (rpc *rpcImpl) Register(rcvr interface{}) error {
	return rpc.register(rcvr, "", false)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func (rpc *rpcImpl) RegisterName(name string, rcvr interface{}) error {
	return rpc.register(rcvr, name, true)
}

func (rpc *rpcImpl) Handler(req *Request, baseParam ...interface{}) (interface{}, error) {
	service, mtype, args, err := rpc.readRequest(req, baseParam...)
	if err != nil {
		return nil, err
	}

	ret, err := service.call(mtype, args)
	if err != nil {
		return nil, err
	}

	return ret.Interface(), nil
}

func (rpc *rpcImpl) register(rcvr interface{}, name string, useName bool) (err error) {
	s := new(service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(s.rcvr).Type().Name()
	if useName {
		sname = name
	}
	if sname == "" {
		return errors.New("rpc.Register: no service name for type " + s.typ.String())
	}
	if !isExported(sname) && !useName {
		return errors.New("rpc.Register: type " + sname + " is not exported")
	}
	s.name = sname

	// Install the methods
	s.method, err = suitableMethods(s.typ)
	if err != nil {
		return err
	}

	if len(s.method) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method, err := suitableMethods(reflect.PtrTo(s.typ))
		if len(method) != 0 {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
		}
		return fmt.Errorf("%s, error: %v", str, err)
	}

	for name := range s.method {
		fmt.Printf("rpc.Register functon: %s.%s\n", sname, name)
	}

	if _, dup := rpc.serviceMap.LoadOrStore(sname, s); dup {
		return errors.New("rpc: service already defined: " + sname)
	}
	return nil
}

// suitableMethods returns suitable Rpc methods of typ, it will report
// error using log if reportErr is true.
func suitableArgs(mtype reflect.Type, mname string) ([]reflect.Type, error) {
	argTypes := make([]reflect.Type, 0, mtype.NumIn()-1)
	for i := 1; i < mtype.NumIn(); i++ {
		argType := mtype.In(i)
		if !isExportedOrBuiltinType(argType) {
			return nil, fmt.Errorf("rpc.Register: argument type of method %q is not exported: %q\n", mname, argType)
		}
		argTypes = append(argTypes, argType)
	}
	return argTypes, nil
}

func suitableMethods(typ reflect.Type) (map[string]*methodType, error) {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}

		argTypes, err := suitableArgs(mtype, mname)
		if err != nil {
			return nil, err
		}

		var replyType reflect.Type
		if mtype.NumOut() > 0 {
			replyType = mtype.Out(0)
			if !isExportedOrBuiltinType(replyType) {
				return nil, fmt.Errorf("rpc.Register: return type of method %q is not exported: %q\n", mname, replyType)
			}
		}

		if mtype.NumOut() > 1 {
			lastReplyType := mtype.Out(1)
			if !isExportedOrBuiltinType(lastReplyType) {
				return nil, fmt.Errorf("rpc.Register: return type of method %q is not exported: %q\n", mname, lastReplyType)
			}
			if !lastReplyType.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				return nil, fmt.Errorf("rpc.Register: return type of method %q last reply type not is error type\n", mname)
			}
		}

		if mtype.NumOut() < 1 || mtype.NumOut() > 2 {
			return nil, fmt.Errorf("rpc.Register: method %q has %d output parameters; needs exactly one or two\n", mname, mtype.NumOut())
		}

		methods[mname] = &methodType{method: method, argTypes: argTypes, replyType: replyType}
	}
	return methods, nil
}

func (s *service) call(mtype *methodType, args []reflect.Value) (replyv reflect.Value, err error) {
	function := mtype.method.Func

	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call(args)
	// The return value for the method is an error.
	if len(returnValues) > 0 {
		replyv = returnValues[0]
	}
	if len(returnValues) > 1 {
		v := returnValues[1]
		if !v.IsNil() {
			err = v.Interface().(error)
		}
	}
	return
}

func (rpc *rpcImpl) readRequest(req *Request, defaultParams ...interface{}) (service *service, mtype *methodType, argv []reflect.Value, err error) {
	service, mtype, err = rpc.readRequestServiceMethod(req)
	if err != nil {
		return
	}

	defaultParamsLen := len(defaultParams)

	var lens int
	if req.Params == nil {
		lens = 0
	} else {
		lens = len(req.Params)
	}

	if lens != (mtype.method.Type.NumIn() - defaultParamsLen - 1) {
		err = fmt.Errorf("rpc: params not matched. got %d, need %d", lens, mtype.method.Type.NumIn()-defaultParamsLen-1)
		return
	}

	argv = make([]reflect.Value, len(mtype.argTypes)+1)
	argv[0] = service.rcvr

	for idx, param := range defaultParams {
		argv[idx+1] = reflect.ValueOf(param)
	}

	for i := 0; i < lens; i++ {
		targetType := mtype.argTypes[i+defaultParamsLen] //jump default_params
		var arg reflect.Value
		arg, err = convert(req.Params[i], targetType)
		if err != nil {
			err = fmt.Errorf("rpc: convert param faild. expect %s, found=%v, error: %v\n",
				targetType, string(*req.Params[i]), err)
			return
		}
		argv[i+defaultParamsLen+1] = arg
	}

	return
}

func (rpc *rpcImpl) readRequestServiceMethod(req *Request) (svc *service, mtype *methodType, err error) {
	dot := strings.LastIndex(req.ServiceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc: service/method request ill-formed: " + req.ServiceMethod)
		return
	}
	serviceName := req.ServiceMethod[:dot]
	methodName := req.ServiceMethod[dot+1:]

	// Look up the request.
	svci, ok := rpc.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc: can't find service " + req.ServiceMethod)
		return
	}
	svc = svci.(*service)
	mtype = svc.method[methodName]
	if mtype == nil {
		err = errors.New("rpc: can't find method " + req.ServiceMethod)
	}
	return
}

func convert(msg *json.RawMessage, argType reflect.Type) (argv reflect.Value, err error) {
	// Decode the argument value.
	argIsValue := false // if true, need to indirect before calling.
	if argType.Kind() == reflect.Ptr {
		argv = reflect.New(argType.Elem())
	} else {
		argv = reflect.New(argType)
		argIsValue = true
	}
	// argv guaranteed to be a pointer now.
	if err = json.Unmarshal(*msg, argv.Interface()); err != nil {
		return
	}
	if argIsValue {
		argv = argv.Elem()
	}
	return
}

func isExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}
