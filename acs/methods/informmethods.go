package methods

import (
	"encoding/xml"
	"fmt"
	"goacs/acs"
	"goacs/acs/http"
	acsxml "goacs/acs/types"
	"goacs/lib"
	"goacs/repository/mysql"
	"log"
)

type InformDecision struct {
	ReqRes *http.CPERequest
}

func (InformDecision *InformDecision) CpeInformResponse() {
	InformDecision.ReqRes.Session.PrevReqType = acsxml.INFORM
	acs.AddCookieToResponseWriter(InformDecision.ReqRes.Session, InformDecision.ReqRes.Response)

	_, _ = fmt.Fprint(InformDecision.ReqRes.Response, InformDecision.ReqRes.Envelope.InformResponse())
}

func (InformDecision *InformDecision) CpeInformRequestParser() {
	env := new(lib.Env)

	var inform acsxml.Inform
	_ = xml.Unmarshal(InformDecision.ReqRes.Body, &inform)
	InformDecision.ReqRes.Session = acs.GetOrCreateSession(inform.DeviceId.SerialNumber)
	log.Println("SESSION FROM INFORM", InformDecision.ReqRes.Session.IsNew, InformDecision.ReqRes.Session.ReadAllParameters)
	if env.Get("DEBUG", "false") == "true" {
		InformDecision.ReqRes.Session.IsBoot = true
	}

	InformDecision.ReqRes.Session.FillCPESessionFromInform(inform)
	tasksRepository := mysql.NewTasksRepository(InformDecision.ReqRes.DBConnection)
	InformDecision.ReqRes.Session.Tasks = tasksRepository.GetTasksForCPE(InformDecision.ReqRes.Session.CPE.UUID)
	cpeRepository := mysql.NewCPERepository(InformDecision.ReqRes.DBConnection)
	_, cpeExist, _ := cpeRepository.UpdateOrCreate(&InformDecision.ReqRes.Session.CPE)
	InformDecision.ReqRes.Session.ReadAllParameters = !cpeExist
	_, _ = cpeRepository.SaveParameters(&InformDecision.ReqRes.Session.CPE)
}
