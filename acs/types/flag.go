package types

import (
	"database/sql/driver"
	"errors"
)

type Flag struct {
	Read         bool `json:"read"`          //R
	Write        bool `json:"write"`         //W
	AddObject    bool `json:"add_object"`    //A
	System       bool `json:"system"`        //X
	PeriodicRead bool `json:"periodic_read"` //P
	Important    bool `json:"important"`     //I
	Send         bool `json:"send"`          //S
}

func FlagFromString(flags string) (Flag, error) {
	var err error = nil

	var flag = Flag{
		Read:         false,
		Write:        false,
		AddObject:    false,
		System:       false,
		PeriodicRead: false,
		Important:    false,
		Send:         false,
	}

	for _, token := range flags {
		switch token {
		case 'R':
			flag.Read = true
		case 'W':
			flag.Write = true
		case 'A':
			flag.AddObject = true
		case 'X':
			flag.System = true
		case 'P':
			flag.PeriodicRead = true
		case 'I':
			flag.Important = true
		case 'S':
			flag.Send = true
		default:
			err = errors.New("Unknown flag " + string(token))
		}
	}

	return flag, err
}

func (flag *Flag) AsString() string {
	stringFlag := ""

	if flag.Read == true {
		stringFlag += "R"
	}
	if flag.Write == true {
		stringFlag += "W"
	}
	if flag.AddObject == true {
		stringFlag += "A"
	}
	if flag.System == true {
		stringFlag += "X"
	}
	if flag.PeriodicRead == true {
		stringFlag += "P"
	}
	if flag.Important == true {
		stringFlag += "I"
	}
	if flag.Send == true {
		stringFlag += "S"
	}

	return stringFlag
}

func (flag *Flag) IsReadable() bool {
	return flag.Read
}

func (flag *Flag) IsWriteable() bool {
	return flag.Write
}

func (flag *Flag) CharToFieldName(char string) string {
	switch char {
	case "W":
		return "Write"
	case "A":
		return "AddObject"
	case "X":
		return "System"
	case "P":
		return "PeriodicRead"
	case "I":
		return "Important"
	case "S":
		return "Send"
	}

	return "Read"
}

func (flag *Flag) Value() (driver.Value, error) {
	return flag.AsString(), nil
}

func (flag *Flag) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []uint8:
		src := src.([]byte)
		*flag, err = FlagFromString(string(src))
	default:
		err = errors.New("Invalid flag")
	}

	return
}
