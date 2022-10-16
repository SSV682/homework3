package v1

import "github.com/dimfeld/httptreemux"

func RegisterRouters(router *httptreemux.TreeMux) {
	//api := router.NewGroup("/v1")
	registerHeatlthGroup(router)
}
