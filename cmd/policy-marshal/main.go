package main

import (
	"encoding/json"
	"fmt"
)

type OnOff string

const (
	On  OnOff = "on"
	Off OnOff = "off"
)

func (o OnOff) IsOn() bool { return o == On }

// 프로파일 기능의 활성화 여부.
type State int

const (
	Disabled State = 0
	Enabled  State = 1
)

func (s State) IsEnabled() bool { return s == Enabled }

type ThreatLevel = int

const (
	ThreatLevelLow    ThreatLevel = 1
	ThreatLevelMedium ThreatLevel = 2
	ThreatLevelHigh   ThreatLevel = 3
	ThreatLevelXBA    ThreatLevel = 9
)

// 안티 멀웨어 정책.
type AntiMalwarePolicy struct {
	// 안티멀웨어 설정.
	AntiMalware State `json:"UseAntiMalware,string"`

	// 엔진 업데이트 설정.
	EngineUpdate OnOff `json:"UseAntiMalwareUpdate"`
	// 엔진 업데이트 주기 (hour).
	EngineUpdateInterval int `json:"AntiMalwareUpdateAutoPeriod,string"`

	// 검사할 최대 파일 크기 (MB).
	ScanMaxSize int64 `json:"MaxFileSizeLimit"`
	// 압축 파일 검사 설정.
	ArchiveScan OnOff `json:"ScanCompressFile"`
	// 잠재적 위협 프로그램 탐지 설정.
	SuspiciousApp OnOff `json:"DetectThreatTypeApp"`

	// 자동대응 사용 여부.
	AutoResponse OnOff `json:"AutoResponse"`
	// 자동대응 시 팝업을 생성하는 최소 위험도.
	PopupThreshold ThreatLevel `json:"PopupLevel"`
	// 자동대응 시 프로세스를 종료하는 최소 위험도.
	TerminateThreshold ThreatLevel `json:"KillLevel"`
	// 자동대응 시 파일을 격리하는 최소 위험도.
	QuarantineThreshold ThreatLevel `json:"IsolateLevel"`

	// 안티멀웨어의 예약 검사 목록.
	Reservation AMReservation `json:"Reservation"`

	// 안티멀웨어의 파일 예외 패턴.
	BypassFilePattern AMBypassFilePattern `json:"ExceptFileInfoList"`
	// 안티멀웨어의 위협 예외 패턴.
	BypassThreatPattern AMBypassThreatPattern `json:"ExceptThreatInfoList"`
}

func NewAntiMalwarePolicy() *AntiMalwarePolicy {
	return &AntiMalwarePolicy{
		AntiMalware:          Disabled,
		EngineUpdate:         Off,
		EngineUpdateInterval: 3,
		ScanMaxSize:          100,
		ArchiveScan:          Off,
		SuspiciousApp:        Off,
		AutoResponse:         Off,
		PopupThreshold:       ThreatLevelHigh,
		TerminateThreshold:   ThreatLevelHigh,
		QuarantineThreshold:  ThreatLevelHigh,
	}
}

type PathItem struct {
	PathConfig string `json:"AV_ReservationScanPathConfig"`
}

type PathList struct {
	Items []PathItem `json:"-"`
}

func (a *PathList) UnmarshalJSON(data []byte) error {
	var origin string
	if err := json.Unmarshal(data, &origin); err != nil {
		return err
	}
	var raw []PathItem
	if err := json.Unmarshal([]byte(origin), &raw); err != nil {
		return err
	}
	a.Items = raw
	return nil
}

func (a *PathList) MarshalJSON() ([]byte, error) {
	raw, err := json.Marshal(a.Items)
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(raw))
}

type ScanItem struct {
	Name  string   `json:"Name"`
	Cycle string   `json:"Cycle"`
	Week  string   `json:"DayOfTheWeek"`
	Time  string   `json:"Time"`
	Path  PathList `json:"Path"`
}

type AMReservationScanList struct {
	Items []ScanItem `json:"-"`
}

func (a *AMReservationScanList) UnmarshalJSON(data []byte) error {
	var origin string
	if err := json.Unmarshal(data, &origin); err != nil {
		return err
	}
	var raw []ScanItem
	if err := json.Unmarshal([]byte(origin), &raw); err != nil {
		return err
	}
	a.Items = raw
	return nil
}

func (a *AMReservationScanList) MarshalJSON() ([]byte, error) {
	raw, err := json.Marshal(a.Items)
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(raw))
}

type AMReservation struct {
	ScanList AMReservationScanList `json:"AV_ReservationScanList"`
}

type AMBypassFilePatternItem struct {
	// 패턴 이름.
	Name string `json:"FileInfoName"`
	// 패턴 설명.
	Desc string `json:"FileInfoDesc"`

	// 패턴 타입.
	// 타입에 따라 Path, Hash, CodeSign 중 하나만 값이 존재한다.
	Type string `json:"FileInfoType"`

	// 패턴 타입이 Path 일 때의 값.
	Path string `json:"FileInfoPath"`
	// 패턴 타입이 Hash 일 때의 값.
	Hash string `json:"FileInfoHash"`
	// 패턴 타입이 CodeSign 일 때의 값.
	CodeSign string `json:"FileInfoCodeSign"`
}

func (a *AMBypassFilePatternItem) Value() string {
	switch a.Type {
	case "Path":
		return a.Path
	case "Hash":
		return a.Hash
	case "CodeSign":
		return a.CodeSign
	default:
		return ""
	}
}

// 안티멀웨어의 파일 예외 패턴.
type AMBypassFilePattern struct {
	Items []AMBypassFilePatternItem `json:"-"`
}

func (a *AMBypassFilePattern) UnmarshalJSON(data []byte) error {
	var origin string
	if err := json.Unmarshal(data, &origin); err != nil {
		return err
	}
	var raw []AMBypassFilePatternItem
	if err := json.Unmarshal([]byte(origin), &raw); err != nil {
		return err
	}
	a.Items = raw
	return nil
}

func (a *AMBypassFilePattern) MarshalJSON() ([]byte, error) {
	raw, err := json.Marshal(a.Items)
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(raw))
}

type AMBypassThreatPatternItem struct {
	// 패턴 이름.
	Name string `json:"ThreatInfoName"`
	// 패턴 설명.
	Desc string `json:"ThreatInfoDesc"`

	// 예외처리 할 위협 이름.
	Pattern string `json:"ThreatInfoThreatName"`
}

// 안티멀웨어의 위협 예외 패턴.
type AMBypassThreatPattern struct {
	Items []AMBypassThreatPatternItem `json:"-"`
}

func (a *AMBypassThreatPattern) UnmarshalJSON(data []byte) error {
	var origin string
	if err := json.Unmarshal(data, &origin); err != nil {
		return err
	}
	var raw []AMBypassThreatPatternItem
	if err := json.Unmarshal([]byte(origin), &raw); err != nil {
		return err
	}
	a.Items = raw
	return nil
}

func (a *AMBypassThreatPattern) MarshalJSON() ([]byte, error) {
	raw, err := json.Marshal(a.Items)
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(raw))
}

var (
	jsonData = string("{\"AdditionalScanFileExtensionList\":\"\",\"AntiMalwareUpdateAutoPeriod\":\"1\",\"AutoResponse\":\"on\",\"DetectThreatTypeApp\":\"on\",\"ExceptFileInfoList\":\"[{\\\"FileInfoName\\\":\\\"exception-1\\\",\\\"FileInfoType\\\":\\\"Path\\\",\\\"FileInfoPath\\\":\\\"/home/minsu/projects/insights-idl\\\",\\\"FileInfoHash\\\":\\\"\\\",\\\"FileInfoDesc\\\":\\\"insights 프로젝트\\\"},{\\\"FileInfoName\\\":\\\"exception-hash-1\\\",\\\"FileInfoType\\\":\\\"Hash\\\",\\\"FileInfoPath\\\":\\\"\\\",\\\"FileInfoHash\\\":\\\"97ff5b3b232b8f24b871284f6684b116085fc27c3f221906ebfc5ae9fed3d7d7\\\",\\\"FileInfoDesc\\\":\\\"테스트 해시 값\\\"}]\",\"ExceptThreatInfoList\":\"[{\\\"ThreatInfoName\\\":\\\"ThreatInfo\\\",\\\"ThreatInfoThreatName\\\":\\\"test-threat\\\",\\\"ThreatInfoDesc\\\":\\\"test information\\\"}]\",\"IsolateLevel\":3,\"KillLevel\":3,\"MaxFileSizeLimit\":100,\"Reservation\":{\"AV_ReservationScanList\":\"[{\\\"Name\\\":\\\"everyDays\\\",\\\"Cycle\\\":\\\"Day\\\",\\\"DayOfTheWeek\\\":\\\"\\\",\\\"Time\\\":\\\"16:44\\\",\\\"Path\\\":\\\"[{\\\\\\\"AV_ReservationScanPathConfig\\\\\\\":\\\\\\\"/home/minsu/projects/virus/locky\\\\\\\"},{\\\\\\\"AV_ReservationScanPathConfig\\\\\\\":\\\\\\\"/home/minsu/projects/virus/jaff\\\\\\\"}]\\\"},{\\\"Name\\\":\\\"everyWeek\\\",\\\"Cycle\\\":\\\"Week\\\",\\\"DayOfTheWeek\\\":\\\"5,0,2\\\",\\\"Time\\\":\\\"13:33\\\",\\\"Path\\\":\\\"[{\\\\\\\"AV_ReservationScanPathConfig\\\\\\\":\\\\\\\"/home/minsu/projects/insights-idl\\\\\\\"}]\\\"}]\"},\"ScanCompressFile\":\"on\",\"UseAntiMalware\":\"1\",\"UseAntiMalwareUpdate\":\"on\",\"UseRecommendedExceptThreatNameList\":\"on\"}")
)

func main() {
	amPolicy := NewAntiMalwarePolicy()
	json.Unmarshal([]byte(jsonData), amPolicy)

	data, _ := json.MarshalIndent(amPolicy, "", "   ")

	fmt.Println(string(data))
}
