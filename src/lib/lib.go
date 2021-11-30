package lib

import (
	"fmt"
	"mohazit/lang"
	"strings"
)

func Load() {
	streams["void"] = &DummyStream{}
	lang.Funcs = lang.FuncMap{
		// user interaction
		"say": fSay,
		// file management
		"file-create": fFileCreate,
		"file-delete": fFileDelete,
		"file-rename": fFileRename,
		// data streams
		"data-read":  fDataRead,
		"data-write": fDataWrite,
		"data-seek":  fDataSeek,
		"data-close": fDataClose,
		"file-open":  fFileOpen,
	}
	lang.Comps = lang.CompMap{
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
