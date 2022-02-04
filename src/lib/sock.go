package lib

import (
	"fmt"
	"mohazit/lang"
	"net"
	"strings"
)

type NetConnStream struct {
	conn net.Conn
}

func (s *NetConnStream) Read(p []byte) (int, error) {
	return s.conn.Read(p)
}

func (s *NetConnStream) Write(p []byte) (int, error) {
	return s.conn.Write(p)
}

func (s *NetConnStream) Seek(offset int64, whence int) (int64, error) {
	return 0, badState.Get("cannot seek in non-buffered socket")
}

func (s *NetConnStream) Close() error {
	return s.conn.Close()
}

var Listeners = make(map[string]net.Listener)

func fSockDial(args []*lang.Object) (*lang.Object, error) {
	var addr string
	var streamName string
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need address")
	}
	addrObj := args[0]
	if addrObj.Type != lang.ObjStr {
		return lang.NewNil(), badType.Get("address must be a string")
	}
	addr = addrObj.Data.String()
	if len(args) != 2 {
		streamName = fmt.Sprintf("socket%d", streamsSoFar)
	} else {
		streamName = args[0].String()
	}
	streamsSoFar++

	fmt.Printf("dialing via socket stream `%s`\n", streamName)

	c, err := net.Dial("tcp", addr)
	if err != nil {
		return lang.NewNil(), err
	}

	streams[streamName] = &NetConnStream{c}
	lastStream = streamName
	return lang.NewStr(streamName), nil
}

func fSockListen(args []*lang.Object) (*lang.Object, error) {
	var addr string
	var sockName string
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need address")
	}
	addrObj := args[0]
	if addrObj.Type != lang.ObjStr {
		return lang.NewNil(), badType.Get("address must be a string")
	}
	addr = addrObj.Data.String()
	if len(args) != 2 {
		sockName = fmt.Sprintf("socket%d", streamsSoFar)
	} else {
		sockName = strings.ToLower(args[1].String())
	}
	streamsSoFar++

	fmt.Printf("listening via socket `%s`\n", sockName)

	c, err := net.Listen("tcp", addr)
	if err != nil {
		return lang.NewNil(), err
	}
	Listeners[sockName] = c
	return lang.NewStr(sockName), nil
}

func fSockAccept(args []*lang.Object) (*lang.Object, error) {
	var sockName string
	if len(args) != 1 {
		return lang.NewNil(), moreArgs.Get("need socket name")
	}
	sockNameObj := args[0]
	if sockNameObj.Type != lang.ObjStr {
		return lang.NewNil(), badArg.Get("socket name must be a string")
	}
	sockName = sockNameObj.Data.String()

	l, ok := Listeners[sockName]
	if !ok {
		return lang.NewNil(), badState.Get("socket does not exist: " + sockName)
	}
	c, err := l.Accept()
	if err != nil {
		return lang.NewNil(), err
	}

	sockName = fmt.Sprintf("socket%d", streamsSoFar)
	streams[sockName] = &NetConnStream{c}
	streamsSoFar++
	lastStream = sockName

	fmt.Printf("receievd connection: socket stream `%s`\n", sockName)

	return lang.NewStr(sockName), nil
}

func fSockAddr(args []*lang.Object) (*lang.Object, error) {
	if len(args) != 4 {
		return lang.NewNil(), badArg.Get("addr must have exactly 4 args")
	}
	b := []byte{}
	for _, a := range args {
		if a.Type != lang.ObjInt {
			return lang.NewNil(), badArg.Get("addr values must be integers")
		}
		if a.Data.Int() > 255 || a.Data.Int() < 0 {
			return lang.NewNil(), badArg.Get("addr values must be 8-bit")
		}
		b = append(b, byte(a.Data.Int()))
	}
	ip := net.IPv4(b[0], b[1], b[2], b[3])
	return lang.NewObject(ip), nil
}
