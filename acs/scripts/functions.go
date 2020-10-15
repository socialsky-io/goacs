package scripts

import "goacs/acs/types"

func (se *ScriptEngine) SetParameter(path string, value string) {
	se.ACSSession.CPE.AddParameter(types.ParameterValueStruct{
		Name:  path,
		Value: value,
		Type:  "",
		Flag: types.Flag{
			Read:  true,
			Write: true,
		},
	})
}
