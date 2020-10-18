package methods

import (
	"encoding/xml"
	"fmt"
	"goacs/acs/http"
	acsxml "goacs/acs/types"
	"goacs/repository/mysql"
)

type InformDecision struct {
	ReqRes *http.ReqRes
}

func (InformDecision *InformDecision) CpeInformResponse() {
	InformDecision.ReqRes.Session.PrevReqType = acsxml.INFORM
	_, _ = fmt.Fprint(InformDecision.ReqRes.Response, InformDecision.ReqRes.Envelope.InformResponse())
}

func (InformDecision *InformDecision) CpeResponseParser() {
	var inform acsxml.Inform
	_ = xml.Unmarshal(InformDecision.ReqRes.Body, &inform)
	fmt.Println("BOOT", inform.IsBootEvent())
	InformDecision.ReqRes.Session.FillCPESessionFromInform(inform)
	tasksRepository := mysql.NewTasksRepository(InformDecision.ReqRes.DBConnection)
	InformDecision.ReqRes.Session.Tasks = tasksRepository.GetTasksForCPE(InformDecision.ReqRes.Session.CPE.UUID)
	cpeRepository := mysql.NewCPERepository(InformDecision.ReqRes.DBConnection)
	_, cpeExist, _ := cpeRepository.UpdateOrCreate(&InformDecision.ReqRes.Session.CPE)
	InformDecision.ReqRes.Session.ReadAllParameters = cpeExist
	_, _ = cpeRepository.SaveParameters(&InformDecision.ReqRes.Session.CPE)
}
