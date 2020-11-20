package scripts

import (
	"fmt"
	"goacs/acs/types"
	"goacs/repository"
	"goacs/repository/mysql"
)

func (se *ScriptEngine) SetParameter(path string, value string) {
	se.ReqRes.Session.CPE.AddParameter(types.ParameterValueStruct{
		Name:  path,
		Value: value,
		Type:  "",
		Flag: types.Flag{
			Read:  true,
			Write: true,
		},
	})
}

func (se *ScriptEngine) SaveDevice() {
	cpeRepository := mysql.NewCPERepository(repository.GetConnection())
	_ = cpeRepository.BulkInsertOrUpdateParameters(&se.ReqRes.Session.CPE, se.ReqRes.Session.CPE.ParameterValues)
}

func (se *ScriptEngine) AddObject(path string) {
	reqBody := se.ReqRes.Envelope.AddObjectRequest(path, "")
	_, _ = fmt.Fprint(se.ReqRes.Response, reqBody)
}
