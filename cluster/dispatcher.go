package cluster

import (
	"github.com/sniperHW/sanguo/cluster/addr"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/kendynet"
)

type MsgHandler func(addr.LogicAddr, proto.Message)

var mtxHandler sync.Mutex
var handlers map[string]MsgHandler = map[string]MsgHandler{}
var onPeerDisconnected func(addr.LogicAddr, error)

func RegisterPeerDisconnected(cb func(addr.LogicAddr, error)) {
	defer mtxHandler.Unlock()
	mtxHandler.Lock()
	onPeerDisconnected = cb
}

func Register(msg proto.Message, handler MsgHandler) {
	defer mtxHandler.Unlock()
	mtxHandler.Lock()

	msgName := reflect.TypeOf(msg).String()
	if nil == handler {
		//记录日志
		Errorf("Register %s failed: handler is nil\n", msgName)
		return
	}
	_, ok := handlers[msgName]
	if ok {
		//记录日志
		Errorf("Register %s failed: duplicate handler\n", msgName)
		return
	}

	handlers[msgName] = handler
}

func pcall(handler MsgHandler, from addr.LogicAddr, name string, msg proto.Message) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 65535)
			l := runtime.Stack(buf, false)
			Errorf("error on Dispatch:%s\nstack:%v,%s\n", name, r, buf[:l])
		}
	}()
	handler(from, msg)
}

func dispatch(from addr.LogicAddr, name string, msg proto.Message) {
	if nil != msg {
		mtxHandler.Lock()
		handler, ok := handlers[name]
		mtxHandler.Unlock()
		if ok {
			pcall(handler, from, name, msg)
		} else {
			//记录日志
			Errorf("unkonw msg:%s\n", name)
		}
	}
}

func dispatchPeerDisconnected(peer addr.LogicAddr, err error) {
	var h func(addr.LogicAddr, error)
	mtxHandler.Lock()
	h = onPeerDisconnected
	mtxHandler.Unlock()
	if nil != h {
		h(peer, err)
	}
}

func dispatch_(from addr.LogicAddr, session kendynet.StreamSession, name string, msg proto.Message, server bool) {
	if nil != msg {
		switch msg.(type) {
		case *Heartbeat:
			if server {
				heartbeat := msg.(*Heartbeat)
				heartbeat_resp := &Heartbeat{}
				heartbeat_resp.Timestamp1 = proto.Int64(time.Now().UnixNano())
				heartbeat_resp.Timestamp2 = proto.Int64(heartbeat.GetTimestamp1())
				session.Send(heartbeat_resp)
			}
			break
		case *AddForginServicesH2S:
			onAddForginServicesH2S(msg.(*AddForginServicesH2S))
			break
		case *RemForginServicesH2S:
			onRemForginServicesH2S(msg.(*RemForginServicesH2S))
			break
		case *AddForginServicesH2H:
			onAddForginServicesH2H(msg.(*AddForginServicesH2H))
			break
		case *RemForginServicesH2H:
			onRemForginServicesH2H(msg.(*RemForginServicesH2H))
			break
		default:
			dispatch(from, name, msg)
			break
		}
	}
}

func dispatchServer(from addr.LogicAddr, session kendynet.StreamSession, name string, msg proto.Message) {
	dispatch_(from, session, name, msg, true)
}

func dispatchClient(from addr.LogicAddr, session kendynet.StreamSession, name string, msg proto.Message) {
	dispatch_(from, session, name, msg, false)
}
