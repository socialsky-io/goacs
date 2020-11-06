package logic

import (
	"encoding/xml"
	"fmt"
	"goacs/acs"
	acshttp "goacs/acs/http"
	"goacs/acs/methods"
	"goacs/acs/scripts"
	acsxml "goacs/acs/types"
	"goacs/models/tasks"
	"goacs/repository"
	"goacs/repository/mysql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func CPERequestDecision(request *http.Request, w http.ResponseWriter) {
	buffer, err := ioutil.ReadAll(request.Body)
	//session := acs.CreateSession(request)

	//w = acs.AddCookieToResponseWriter(session, w)

	if err != io.EOF && err != nil {
		return
	}

	reqType, envelope := parseEnvelope(buffer)

	var reqRes = acshttp.CPERequest{
		Request:      request,
		Response:     w,
		DBConnection: repository.GetConnection(),
		Session:      acs.GetSessionFromRequest(request),
		Envelope:     &envelope,
		Body:         buffer,
	}

	if reqRes.Session != nil {
		acs.AddCookieToResponseWriter(reqRes.Session, reqRes.Response)
	}

	switch reqType {
	case acsxml.INFORM:
		informDecision := methods.InformDecision{&reqRes}
		informDecision.CpeInformRequestParser()
		informDecision.CpeInformResponse()

	case acsxml.EMPTY:
		log.Println("EMPTY RESPONSE")
		if reqRes.Session.PrevReqType == acsxml.INFORM {
			log.Println("EMPTY ON INFORM", reqRes.Session.IsNew, reqRes.Session.ReadAllParameters)
			if reqRes.Session.IsNew == false && reqRes.Session.ReadAllParameters == true {
				//fmt.Println("GPN REQ")
				//parameterDecisions := methods.ParameterDecisions{&reqRes}
				//parameterDecisions.ParameterNamesRequest(true)
				reqRes.Session.NextJob = acs.JOB_GETPARAMETERNAMES
			}

			ProcessSessionJobs(&reqRes)
			ProcessTasks(&reqRes, reqType)
		}

	case acsxml.GPNResp:
		parameterDecisions := methods.ParameterDecisions{ReqRes: &reqRes}
		parameterDecisions.CpeParameterNamesResponseParser()
		log.Println("GPNResponse next job", reqRes.Session.NextJob)
		//parameterDecisions.GetParameterValuesRequest(reqRes.Session.CPE.ParametersInfo)
		if reqRes.Session.CPE.NewInACS {
			// CPE is new in acs (not exist)
			reqRes.Session.NextJob = acs.JOB_GETPARAMETERVALUES
			//parameterDecisions.GetParameterValuesRequest([]acsxml.ParameterInfo{
			//	{
			//		Name:     reqRes.Session.CPE.Root + ".",
			//		Writable: "0",
			//	},
			//})
		}
	case acsxml.GPVResp:
		parameterDecisions := methods.ParameterDecisions{ReqRes: &reqRes}
		parameterDecisions.GetParameterValuesResponseParser()

	case acsxml.SPVResp:
		paramaterDecisions := methods.ParameterDecisions{ReqRes: &reqRes}
		paramaterDecisions.SetParameterValuesResponse()

	case acsxml.FaultResp:
		var faultresponse acsxml.Fault
		_ = xml.Unmarshal(buffer, &faultresponse)
		reqRes.Session.CPE.Fault = faultresponse
		faultDecision := methods.FaultDecision{ReqRes: &reqRes}
		faultDecision.ResponseDecision()

	default:
		fmt.Println("UNSUPPORTED REQTYPE ", reqType)
	}

}

func ProcessSessionJobs(reqRes *acshttp.CPERequest) {
	log.Println("Processing next job", reqRes.Session.NextJob)
	switch reqRes.Session.NextJob {
	case acs.JOB_SENDPARAMETERS:
		parameterDecisions := methods.ParameterDecisions{ReqRes: reqRes}
		parameterDecisions.SetParameterValuesResponse()
		reqRes.Session.NextJob = acs.JOB_NONE

	case acs.JOB_GETPARAMETERNAMES:
		parameterDecisions := methods.ParameterDecisions{ReqRes: reqRes}
		parameterDecisions.ParameterNamesRequest(true)
		reqRes.Session.NextJob = acs.JOB_GETPARAMETERVALUES

	case acs.JOB_GETPARAMETERVALUES:
		log.Println("JOB_GETPARAMETERVALUES ROOT", reqRes.Session.CPE.Root)
		parameterDecision := methods.ParameterDecisions{ReqRes: reqRes}
		parameterDecision.GetParameterValuesRequest([]acsxml.ParameterInfo{
			{
				Name:     reqRes.Session.CPE.Root + ".",
				Writable: "0",
			},
		})
		reqRes.Session.NextJob = acs.JOB_NONE
	}
}

func ProcessTasks(reqRes *acshttp.CPERequest, event string) {
	tasksRepository := mysql.NewTasksRepository(reqRes.DBConnection)
	cpeTasks := tasksRepository.GetTasksForCPE(reqRes.Session.CPE.UUID)

	if len(cpeTasks) > 0 {
		scriptEngine := scripts.NewScriptEngine(reqRes.Session)
		filteredTasks := tasks.FilterTasksByEvent(event, cpeTasks)
		for _, cpeTask := range filteredTasks {
			if cpeTask.Task == tasks.RunScript {
				_, _ = scriptEngine.Execute(cpeTask.Script)
			}

			if cpeTask.Infinite == false {
				tasksRepository.DoneTask(cpeTask.Id)
			}
		}
	}
}

func parseEnvelope(buffer []byte) (string, acsxml.Envelope) {
	//fmt.Println(string(buffer))
	var envelope acsxml.Envelope
	err := xml.Unmarshal(buffer, &envelope)

	var requestType string = acsxml.EMPTY

	if err == nil {
		switch envelope.Type() {
		case "inform":
			requestType = acsxml.INFORM
		case "getparameternamesresponse":
			requestType = acsxml.GPNResp
		case "getparametervaluesresponse":
			requestType = acsxml.GPVResp
		case "setparametervaluesresponse":
			requestType = acsxml.SPVResp
		case "fault":
			requestType = acsxml.FaultResp
		default:
			fmt.Println("UNSUPPORTED envelope type " + envelope.Type())
			requestType = acsxml.UNKNOWN
		}
	}

	return requestType, envelope
}
