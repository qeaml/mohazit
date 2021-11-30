package lib

import (
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
var lastStream = ""

func fDataRead(args []*lang.Object, i lang.InterVar) error {
	var amt int
	var target string
	var streamName string
	if len(args) < 2 {
		return moreArgs("need amount of bytes and target variable")
	}
	var amtObj = args[0]
	if amtObj.Type != lang.ObjInt {
		return badType("amount must be an integer")
	}
	amt = amtObj.IntV
	var targetObj = args[1]
	if targetObj.Type != lang.ObjStr {
		return badType("target must be a string")
	}
	target = targetObj.StrV
	if len(args) >= 3 {
		var streamObj = args[1]
		if streamObj.Type != lang.ObjStr {
			return badType("stream name must be a string")
		}
		streamName = streamObj.StrV
	} else {
		if lastStream == "" {
			return badState("could not infer stream name")
		}
		streamName = lastStream
	}
	lastStream = streamName
	stream, ok := streams[streamName]
	if !ok {
		return badState("no stream named " + streamName + " is open")
	}

	fmt.Printf("reading data from stream `%s`\n", streamName)

	data := make([]byte, amt)
	_, err := stream.Read(data)
	if err != nil {
		return err
	}
	obj, err := i.Parse(string(data))
	if err != nil {
		return err
	}
	i.Set(target, obj)
	return nil
}

func fDataWrite(args []*lang.Object, i lang.InterVar) error {
	var data []byte
	var streamName string
	if len(args) < 1 {
		return moreArgs("need data to write")
	}
	data = []byte(args[0].String())
	if len(args) >= 2 {
		var streamObj = args[1]
		if streamObj.Type != lang.ObjStr {
			return badType("stream name must be a string")
		}
		streamName = streamObj.StrV
	} else {
		if lastStream == "" {
			return badState("could not infer stream name")
		}
		streamName = lastStream
	}
	lastStream = streamName
	stream, ok := streams[streamName]
	if !ok {
		return badState("no stream named " + streamName + " is open")
	}

	fmt.Printf("writing %d byte(s) to stream `%s`\n", len(data), streamName)

	_, err := stream.Write(data)
	return err
}

func fDataSeek(args []*lang.Object, i lang.InterVar) error {
	var pos int
	var streamName string
	if len(args) < 1 {
		return moreArgs("need data to write")
	}
	posObj := args[0]
	if posObj.Type != lang.ObjInt {
		return badType("position must be an integer")
	}
	pos = posObj.IntV
	if len(args) >= 2 {
		var streamObj = args[1]
		if streamObj.Type != lang.ObjStr {
			return badType("stream name must be a string")
		}
		streamName = streamObj.StrV
	} else {
		if lastStream == "" {
			return badState("could not infer stream name")
		}
		streamName = lastStream
	}
	lastStream = streamName
	stream, ok := streams[streamName]
	if !ok {
		return badState("no stream named " + streamName + " is open")
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
			return badType("stream name must be a string")
		}
		streamName = streamObj.StrV
	} else {
		if lastStream == "" {
			return badState("could not infer stream name")
		}
		streamName = lastStream
	}
	lastStream = streamName
	stream, ok := streams[streamName]
	if !ok {
		return badState("no stream named " + streamName + " is open")
	}

	fmt.Printf("closing stream `%s`\n", streamName)

	stream.Close()
	delete(streams, streamName)
	return nil
}

func fFileOpen(args []*lang.Object, i lang.InterVar) error {
	var fileName string
	var streamName string
	if len(args) < 1 {
		return moreArgs("need file name")
	}
	fileObj := args[0]
	if fileObj.Type != lang.ObjStr {
		return badType("file name must be a string")
	}
	fileName = fileObj.StrV
	if len(args) >= 2 {
		var streamObj = args[1]
		if streamObj.Type != lang.ObjStr {
			return badType("stream name must be a string")
		}
		streamName = streamObj.StrV
	} else {
		if lastStream == "" {
			return badState("could not infer stream name")
		}
		streamName = lastStream
	}

	fmt.Printf("opening file `%s` to stream `%s`\n", fileName, streamName)

	file, err := os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	streams[streamName] = file
	lastStream = streamName
	return nil
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
