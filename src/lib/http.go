package lib

import (
	"fmt"
	"mohazit/lang"
	"mohazit/tool"

	"github.com/levigross/grequests"
)

var client = grequests.NewSession(&grequests.RequestOptions{
	UserAgent: fmt.Sprintf("Mohazit/%s%d", tool.Version, tool.Iteration),
})
var resps = make(map[string]*grequests.Response)
var respCount = 0
var lastResp = ""

func fHttpGet(args []*lang.Object) (*lang.Object, error) {
	if len(args) < 1 {
		return lang.NewNil(), moreArgs.Get("need input")
	}
	in := args[0]
	if in.Type != lang.ObjStr {
		return nil, badType.Get("URL must be a string")
	}
	fmt.Printf("sending HTTP request %d: GET %s\n", respCount, in.StrV)
	resp, err := client.Get(in.StrV, nil)
	if err != nil {
		return nil, err
	}
	respName := fmt.Sprintf("response%d", respCount)
	respCount++
	resps[respName] = resp
	lastResp = respName
	streams[respName] = &GenericStream{data: resp.Bytes()}
	lastStream = respName
	return lang.NewStr(respName), nil
}

func fHttpOk(args []*lang.Object) (*lang.Object, error) {
	respName := lastResp
	if len(args) > 1 {
		if args[0].Type == lang.ObjStr {
			respName = args[0].StrV
		} else {
			return lang.NewNil(), badType.Get("response name must be a string")
		}
	} else if respName == "" {
		return lang.NewNil(), badState.Get("could not infer response name")
	}
	resp, ok := resps[respName]
	if !ok {
		return lang.NewNil(), badState.Get("no response named `" + respName + "` exists")
	}
	return lang.NewBool(resp.StatusCode >= 200 && resp.StatusCode < 300), nil
}
