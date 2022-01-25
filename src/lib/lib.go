package lib

import (
	"fmt"
	"mohazit/lang"
	"strings"
)

func Load() {
	streams["void"] = &DummyStream{}
	lang.Funcs = lang.VFuncMap{
		// user interaction
		"say":     fSay,
		"type-of": fTypeOf,
		// numeric
		"random":         fRandom,
		"limited-random": fLimitedRandom,
		"randi":          fLimitedRandom,
		"atoi":           fAtoi,
		"stringify":      fStringify,
		"inc":            fInc,
		"dec":            fDec,
		"neg":            fNeg,
		// file management
		"file-open":   fFileOpen,
		"file-create": fFileCreate,
		"file-delete": fFileDelete,
		"file-rename": fFileRename,
		"file-list":   fFileList,
		"file-exists": fFileExists,
		"dir":         fFileList,
		"ls":          fFileList,
		"walk":        fWalk,
		"cd":          fWalk,
		// external processes
		"run":   fRun,
		"start": fRun,
		// data streams
		"data-stream": fDataStream,
		"data-read":   fDataRead,
		"data-write":  fDataWrite,
		"data-seek":   fDataSeek,
		"data-close":  fDataClose,
		// http
		"http-get": fHttpGet,
		"http-ok":  fHttpOk,
	}
	lang.Comps = lang.VCompMap{
		"=":  cEquals,
		"==": cEquals,
		"~=": cLike,
		"!=": cNotEquals,
		"<>": cNotEquals,
		">":  cGreater,
		"<":  cLesser,
	}
}

func Cleanup() error {
	unclosedStreams := []string{}
	for streamName, stream := range streams {
		if _, ok := stream.(*DummyStream); !ok {
			unclosedStreams = append(unclosedStreams, streamName)
		}
	}
	if len(unclosedStreams) > 0 {
		return fmt.Errorf("unclosed streams: %s", strings.Join(unclosedStreams, ", "))
	}
	return nil
}
