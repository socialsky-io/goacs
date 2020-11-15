package scripts

import (
	"goacs/acs/types"
	"goacs/repository"
	"goacs/repository/mysql"
)

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

func (se *ScriptEngine) SaveDevice() {
	cpeRepository := mysql.NewCPERepository(repository.GetConnection())
	_ = cpeRepository.BulkInsertOrUpdateParameters(&se.ACSSession.CPE, se.ACSSession.CPE.ParameterValues)
}

func (se *ScriptEngine) AddObject(path string) {

}
