package lib

import (
	"bytes"
	"fmt"
	"io"
	"mohazit/lang"
	"os"
)

type Stream interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

var streams = make(map[string]Stream)
var streamsSoFar = 1
var lastStream = ""

func pDataRead(arg *lang.Object) (*lang.Object, error) {
	var amt int
	var streamName string
	if arg.Type != lang.ObjInt {
		return nil, badType.Get("amount must be an integer")
	}
	amt = arg.IntV
	streamName = lastStream
	stream, ok := streams[streamName]
	if !ok {
		return nil, badState.Get("no stream named " + streamName + " is open")
	}

	fmt.Printf("reading %d byte(s) from stream `%s`\n", amt, streamName)

	data := make([]byte, amt)
	_, err := stream.Read(data)
	if err != nil {
		return nil, err
	}
	return lang.NewStr(string(data)), nil
}

func fDataWrite(args []*lang.Object, i lang.InterVar) error {
	var data []byte
	var streamName string
	if len(args) < 1 {
		return moreArgs.Get("need data to write")
	}
	data = []byte(args[0].String())
	if len(args) >= 2 {
		var streamObj = args[1]
		if streamObj.Type != lang.ObjStr {
			return badType.Get("stream name must be a string")
		}
		streamName = streamObj.StrV
	} else {
		if lastStream == "" {
			return badState.Get("could not infer stream name")
		}
		streamName = lastStream
	}
	lastStream = streamName
	stream, ok := streams[streamName]
	if !ok {
		return badState.Get("no stream named " + streamName + " is open")
	}

	fmt.Printf("writing %d byte(s) to stream `%s`\n", len(data), streamName)

	_, err := stream.Write(data)
	return err
}

func fDataSeek(args []*lang.Object, i lang.InterVar) error {
	var pos int
	var streamName string
	if len(args) < 1 {
		return moreArgs.Get("need data to write")
	}
	posObj := args[0]
	if posObj.Type != lang.ObjInt {
		return badType.Get("position must be an integer")
	}
	pos = posObj.IntV
	if len(args) >= 2 {
		var streamObj = args[1]
		if streamObj.Type != lang.ObjStr {
			return badType.Get("stream name must be a string")
		}
		streamName = streamObj.StrV
	} else {
		if lastStream == "" {
			return badState.Get("could not infer stream name")
		}
		streamName = lastStream
	}
	lastStream = streamName
	stream, ok := streams[streamName]
	if !ok {
		return badState.Get("no stream named " + streamName + " is open")
	}

	fmt.Printf("seeking to position %d in stream `%s`\n", pos, streamName)

	_, err := stream.Seek(int64(pos), 0)
	return err
}

func fDataClose(args []*lang.Object, i lang.InterVar) error {
	var streamName string
	if len(args) >= 1 {
		var streamObj = args[0]
		if streamObj.Type != lang.ObjStr {
			return badType.Get("stream name must be a string")
		}
		streamName = streamObj.StrV
	} else {
		if lastStream == "" {
			return badState.Get("could not infer stream name")
		}
		streamName = lastStream
	}
	lastStream = streamName
	stream, ok := streams[streamName]
	if !ok {
		return badState.Get("no stream named " + streamName + " is open")
	}

	fmt.Printf("closing stream `%s`\n", streamName)

	stream.Close()
	delete(streams, streamName)
	return nil
}

func pFileOpen(arg *lang.Object) (*lang.Object, error) {
	var fileName string
	var streamName string
	if arg.Type == lang.ObjNil {
		return nil, moreArgs.Get("need file name")
	}
	if arg.Type != lang.ObjStr {
		return nil, badType.Get("file name must be a string")
	}
	fileName = arg.StrV
	streamName = fmt.Sprintf("filestream%d", streamsSoFar)
	streamsSoFar++

	fmt.Printf("opening file `%s` to stream `%s`\n", fileName, streamName)

	file, err := os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	streams[streamName] = file
	lastStream = streamName
	return lang.NewStr(streamName), nil
}

func pDataStream(args *lang.Object) (*lang.Object, error) {
	streamName := fmt.Sprintf("buffer%d", streamsSoFar)
	streamsSoFar++

	fmt.Printf("opening stream `%s`\n", streamName)

	streams[streamName] = &BufferStream{}
	return lang.NewStr(streamName), nil
}

type DummyStream struct{}

func (s *DummyStream) Read(p []byte) (int, error) {
	return len(p), nil
}
func (s *DummyStream) Write(p []byte) (int, error) {
	return 0, nil
}
func (s *DummyStream) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
func (s *DummyStream) Close() error {
	return nil
}

type BufferStream struct {
	data bytes.Buffer
	read bytes.Reader
}

func (s *BufferStream) Read(p []byte) (int, error) {
	return s.read.Read(p)
}

func (s *BufferStream) Write(p []byte) (int, error) {
	n, err := s.data.Write(p)
	s.read.Seek(int64(n), 0)
	return n, err
}

func (s *BufferStream) Seek(offset int64, whence int) (int64, error) {
	return s.read.Seek(offset, whence)
}

func (s *BufferStream) Close() error {
	s.data.Reset()
	return nil
}

type GenericStream struct {
	data []byte
	pos  int
}

func (s *GenericStream) Read(p []byte) (int, error) {
	var i = 0
	for s.pos < len(s.data) && i < len(p) {
		p[i] = s.data[s.pos]
		i++
		s.pos++
	}
	return i, nil
}

func (s *GenericStream) Write(p []byte) (int, error) {
	var i = 0
	for i < len(p) {
		if s.pos >= len(s.data) {
			s.data = append(s.data, p[i])
		} else {
			s.data[s.pos] = p[i]
		}
		s.pos++
		i++
	}
	return i, nil
}

func (s *GenericStream) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		s.pos = int(offset)
	case 1:
		s.pos += int(offset)
	case 2:
		s.pos = len(s.data) - int(offset)
	default:
		return int64(s.pos), badState.Fail("unknown whence value")
	}
	return int64(s.pos), nil
}

func (s *GenericStream) Close() error {
	s.data = []byte{}
	s.pos = 0
	return nil
}
