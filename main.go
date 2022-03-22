package main

import (
	"github.com/deanishe/awgo"
	"github.com/dennis-tra/alfred-dict.cc-workflow/workflow"
)

func main() {
	wf := aw.New()
	wf.Run(workflow.NewDictcc(wf).Run)
}
