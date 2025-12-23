package main

import (
	"encoding/json"
	"fmt"
)

type Threat map[string]any

type AVName string

func (t Threat) AVName() AVName {
	value, _ := t["AVName"].(string)
	return AVName(value)
}

func (t Threat) SetAVName(value AVName) {
	t["AVName"] = string(value)
}

func (t Threat) Information() map[string]string {
	var value map[string]string
	switch info := t["Information"].(type) {
	case map[string]string:
		return info
	case map[string]any:
		value = make(map[string]string, len(info))
		for key, val := range info {
			value[key], _ = val.(string)
		}
	default:
		value = make(map[string]string)
	}

	t.SetInformation(value)
	return value
}

func (t Threat) SetInformation(value map[string]string) {
	t["Information"] = value
}

func (t Threat) SetInformationName(name string) {
	t.Information()["Name"] = name
}

func (t Threat) InformationName() string {
	return t.Information()["Name"]
}

func (t Threat) SetInformationCategory(category string) {
	t.Information()["Category"] = category
}

func (t Threat) InformationCategory() string {
	return t.Information()["Category"]
}

func main() {
	threat := Threat{}
	threat.SetAVName("Bitdefender")
	threat.SetInformationName("test")
	threat.SetInformationCategory("malware")

	data, _ := json.Marshal(threat)
	fmt.Println("message: ", string(data))

	temp_threat := Threat{}
	_ = json.Unmarshal(data, &temp_threat)

	fmt.Println("Information: ", temp_threat)
	fmt.Println("AVName: ", temp_threat.AVName())
	fmt.Println("Information Name: ", temp_threat.InformationName())
	fmt.Println("Information Category: ", temp_threat.InformationCategory())
}
