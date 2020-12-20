package scripts

import (
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/vm"
	acshttp "goacs/acs/http"
	"log"
)

type ScriptEngine struct {
	Env    *env.Env
	ReqRes *acshttp.CPERequest
	debug  bool
}

func NewScriptEngine(reqRes *acshttp.CPERequest) ScriptEngine {
	scriptEnv := env.NewEnv()
	_ = scriptEnv.Define("session", reqRes.Session)
	_ = scriptEnv.Define("device", reqRes.Session.CPE)

	se := ScriptEngine{
		ReqRes: reqRes,
		debug:  false,
		Env:    scriptEnv,
	}

	_ = scriptEnv.Define("SetParameter", se.SetParameter)
	_ = scriptEnv.Define("SendParameters", se.SendParameters)
	_ = scriptEnv.Define("SaveDevice", se.SaveDevice)
	_ = scriptEnv.Define("AddObject", se.AddObject)

	return se
}

func (se *ScriptEngine) Execute(script string) (interface{}, error) {
	log.Println("Script execution", script)
	return vm.Execute(se.Env, nil, script)
}

func (se *ScriptEngine) EnableDebug() {
	se.debug = true
}

func (se *ScriptEngine) DisableDebug() {
	se.debug = false
}
