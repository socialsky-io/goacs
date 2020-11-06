package types

import (
	"encoding/xml"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Header  Header   `xml:"Header"`
	Body    Body     `xml:"Body"`
}

type Header struct {
	XMLName xml.Name `xml:"Header"`
	ID      string   `xml:"ID"`
}

type Body struct {
	Message XMLMessage `xml:",any"`
}

type XMLMessage struct {
	XMLName xml.Name
}

type Inform struct {
	DeviceId      DeviceId               `xml:"Body>Inform>DeviceId"`
	Events        []Event                `xml:"Body>Inform>Event>EventStruct"`
	ParameterList []ParameterValueStruct `xml:"Body>Inform>ParameterList>ParameterValueStruct"`
}

type Fault struct {
	FaultCode         string `xml:"Body>Fault>faultcode"`
	FaultString       string `xml:"Body>Fault>faultstring"`
	DetailFaultCode   string `xml:"Body>Fault>detail>Fault>FaultCode"`
	DetailFaultString string `xml:"Body>Fault>detail>Fault>FaultString"`
}

type Event struct {
	EventCode  string
	CommandKey string
}

type ParameterValueStruct struct {
	Name  string `db:"name" json:"name"`
	Value string `db:"value" json:"value"`
	Type  string `xml:",attr" db:"type" json:"type"`
	Flag  Flag   `json:"flag" db:"flags"`
}
type ParameterInfo struct {
	Name     string `xml:"Name"`
	Writable string `xml:"Writable"`
}

type DeviceId struct {
	Manufacturer string
	OUI          string
	ProductClass string
	SerialNumber string
}

type GetParameterNamesResponse struct {
	ParameterList []ParameterInfo `xml:"Body>GetParameterNamesResponse>ParameterList>ParameterInfoStruct"`
}

type GetParameterValuesResponse struct {
	ParameterList []ParameterValueStruct `xml:"Body>GetParameterValuesResponse>ParameterList>ParameterValueStruct"`
}

type AddObjectResponseStruct struct {
	InstanceNumber int `xml:"Body>AddObjectResponse>InstanceNumber"`
	Status         int `xml:"Body>AddObjectResponse>Status"`
}

type DeleteObjectResponseStruct struct {
	Status int `xml:"Body>DeleteObjectResponse>Status"`
}

type ACSBool bool

func NewEnvelope() Envelope {
	rand.NewSource(time.Now().UnixNano())
	return Envelope{
		XMLName: xml.Name{},
		Header: Header{
			ID: strconv.Itoa(rand.Int()),
		},
		Body: Body{},
	}
}

func (abool *ACSBool) UnmarshalXMLAttr(attr xml.Attr) (err error) {

	if attr.Value == "0" {
		*abool = false
	}

	*abool = true

	return nil
}

func (abool ACSBool) String() string {
	if abool == true {
		return "1"
	}
	return "0"
}

func GetParametersWithFlag(parametersToFilter []ParameterValueStruct, flag string) []ParameterValueStruct {
	parameters := []ParameterValueStruct{}
	for _, parameter := range parametersToFilter {
		fieldName := parameter.Flag.CharToFieldName(flag)
		flagBool := reflect.ValueOf(parameter.Flag).FieldByName(fieldName).Bool()
		if flagBool == true {
			parameters = append(parameters, parameter)
		}
	}
	return parameters
}

func (envelope *Envelope) Type() string {
	return strings.ToLower(envelope.Body.Message.XMLName.Local)
}

func (envelope *Envelope) InformResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
    <soapenv:Header>
        <cwmp:ID soapenv:mustUnderstand="1">` + envelope.Header.ID + `</cwmp:ID>
    </soapenv:Header>
    <soapenv:Body>
        <cwmp:InformResponse>
            <MaxEnvelopes>1</MaxEnvelopes>
        </cwmp:InformResponse>
    </soapenv:Body>
</soapenv:Envelope>`
}

func (inform *Inform) IsEvent(event string) bool {
	for idx := range inform.Events {
		if inform.Events[idx].EventCode == event {
			return true
		}
	}

	return false
}

func (inform *Inform) IsBootstrapEvent() bool {
	for idx := range inform.Events {
		if inform.Events[idx].EventCode == "0 BOOTSTRAP" {
			return true
		}
	}

	return false
}

func (inform *Inform) IsBootEvent() bool {
	for idx := range inform.Events {
		if inform.Events[idx].EventCode == "1 BOOT" {
			return true
		}
	}

	return false
}

func (envelope *Envelope) GPNRequest(path string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
    <soapenv:Header>
        <cwmp:ID soapenv:mustUnderstand="1">` + envelope.Header.ID + `</cwmp:ID>
    </soapenv:Header>
    <soapenv:Body>
        <cwmp:GetParameterNames>
			<ParameterPath>` + path + `</ParameterPath>
			<NextLevel>false</NextLevel>
		</cwmp:GetParameterNames>
    </soapenv:Body>
</soapenv:Envelope>`
}

//TODO: zrobić ładniej ;)
//      <cwmp:ID soapenv:mustUnderstand="1">` + envelope.Header.ID + `</cwmp:ID>
func (envelope *Envelope) GPVRequest(info []ParameterInfo) string {
	request := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<soapenv:Header>
		<cwmp:ID soapenv:mustUnderstand="1">` + envelope.Header.ID + `</cwmp:ID>
	</soapenv:Header>
	<soapenv:Body>
		<cwmp:GetParameterValues>
			<ParameterNames soap:arrayType="xsd:string[` + strconv.Itoa(len(info)) + `]">
`
	for _, parameter := range info {
		request += `				<string>` + parameter.Name + `</string>`
		request += "\n"
	}

	request += `			</ParameterNames>
		</cwmp:GetParameterValues>
	</soapenv:Body>
</soapenv:Envelope>`

	return request
}

func (envelope *Envelope) GetRPCMethodsRequest() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
    <soapenv:Header>
        <cwmp:ID soapenv:mustUnderstand="1">` + envelope.Header.ID + `</cwmp:ID>
    </soapenv:Header>
    <soapenv:Body>
        <cwmp:GetRPCMethods>
		</cwmp:GetRPCMethods>
    </soapenv:Body>
</soapenv:Envelope>`
}

func (envelope *Envelope) SetParameterValues(info []ParameterValueStruct) string {
	request := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soapenv:Header>
      <cwmp:ID soapenv:mustUnderstand="1">` + envelope.Header.ID + `</cwmp:ID>
  </soapenv:Header>
  <soapenv:Body>
      <cwmp:SetParameterValues>
			<ParameterList soap:arrayType="cwmp:ParameterValueStruct[` + strconv.Itoa(len(info)) + `]">
`
	for _, parameter := range info {
		request += "<ParameterValueStruct>\n"
		request += `<Name>` + parameter.Name + `</Name>`
		request += "\n"
		request += `<Value>` + parameter.Value + `</Value>`
		request += "\n"
		request += "</ParameterValueStruct>\n"
	}

	request += `</ParameterList>
		</cwmp:SetParameterValues>
  </soapenv:Body>
</soapenv:Envelope>`

	return request
}

func (envelope *Envelope) AddObjectRequest(objectName string, parameterKey string) string {
	request := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soapenv:Header>
      <cwmp:ID soapenv:mustUnderstand="1">` + envelope.Header.ID + `</cwmp:ID>
  </soapenv:Header>
  <soapenv:Body>
      <cwmp:AddObject>
			<ObjectName>` + objectName + `</ObjectName>
			<ParameterKey>` + parameterKey + `</ParameterKey>
		</cwmp:AddObject>
  </soapenv:Body>
</soapenv:Envelope>`

	return request
}

func (envelope *Envelope) DeleteObjectRequest(objectName string, parameterKey string) string {
	request := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soapenv:Header>
      <cwmp:ID soapenv:mustUnderstand="1">` + envelope.Header.ID + `</cwmp:ID>
  </soapenv:Header>
  <soapenv:Body>
      <cwmp:DeleteObject>
			<ObjectName>` + objectName + `</ObjectName>
			<ParameterKey>` + parameterKey + `</ParameterKey>
		</cwmp:DeleteObject>
  </soapenv:Body>
</soapenv:Envelope>`

	return request
}
