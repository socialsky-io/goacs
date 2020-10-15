package scripts

import (
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/vm"
	"goacs/acs"
)

type ScriptEngine struct {
	Env        *env.Env
	ACSSession *acs.ACSSession
	debug      bool
}

func NewScriptEngine(session *acs.ACSSession) ScriptEngine {
	scriptEnv := env.NewEnv()
	_ = scriptEnv.Define("session", session)
	_ = scriptEnv.Define("device", session.CPE)

	se := ScriptEngine{
		ACSSession: session,
		debug:      false,
		Env:        scriptEnv,
	}

	_ = scriptEnv.Define("SetParameter", se.SetParameter)

	return se
}

func (se *ScriptEngine) Execute(script string) (interface{}, error) {
	return vm.Execute(se.Env, nil, script)
}

func (se *ScriptEngine) EnableDebug() {
	se.debug = true
}

func (se *ScriptEngine) DisableDebug() {
	se.debug = false
}
