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
		// file management
		"file-create": fFileCreate,
		"file-delete": fFileDelete,
		"file-rename": fFileRename,
		"file-list":   fFileList,
		"dir":         fFileList,
		"ls":          fFileList,
		"walk":        fWalk,
		"cd":          fWalk,
		// external processer
		"run":   fRun,
		"start": fRun,
		"!":     fRun,
		// data streams
		"data-write": fDataWrite,
		"data-seek":  fDataSeek,
		"data-close": fDataClose,
	}
	lang.Comps = lang.VCompMap{
		// general equality
		"equals":       cEquals,
		"eq":           cEquals,
		"is":           cEquals,
		"=":            cEquals,
		"==":           cEquals,
		"not-equals":   cNotEquals,
		"neq":          cNotEquals,
		"is-not":       cNotEquals,
		"isnt":         cNotEquals,
		"!=":           cNotEquals,
		"~=":           cNotEquals,
		"<>":           cNotEquals,
		"like":         cLike,
		"greater":      cGreater,
		"greater-than": cGreater,
		"gt":           cGreater,
		"larger":       cGreater,
		"larger-than":  cGreater,
		">":            cGreater,
		"lesser":       cLesser,
		"lesser-than":  cLesser,
		"lt":           cLesser,
		"smaller":      cLesser,
		"smaller-than": cLesser,
		"<":            cLesser,
		// file management
		"file-exists": cFileExists,
		// http
		"http-ok": cHttpOk,
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
