package oval

import (
	"encoding/xml"
	"sync"
)

type GOvalDefinitions struct {
	XMLName     xml.Name     `xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 oval_definitions"`
	Definitions GDefinitions `xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 definitions"`
	Tests       GTests       `xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 tests"`
	Objects     GObjects     `xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 objects"`
	States      GStates      `xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 states"`
}

//
// Definitions
//

type GDefinitions struct {
	XMLName     xml.Name      `xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 definitions"`
	Definitions []GDefinition `xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 definition"`
}

type GDefinition struct {
	XMLName xml.Name `xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 definition"`

	Metadata GDefinitionMeta `xml:"metadata"`
	Criteria GCriteria       `xml:"criteria"`

	ID      string `xml:"id,attr"`
	Version string `xml:"version,attr"`
	Class   string `xml:"class,attr"`

	// Extended struct information, not used by XML parser but used during
	// definition evaluation
	sync.Mutex
	status int
}

type GDefinitionMeta struct {
	XMLName xml.Name `xml:"metadata"`

	Title      string `xml:"title"`
	Decription string `xml:"description"`
}

type GCriteria struct {
	XMLName xml.Name `xml:"criteria"`

	Operator string `xml:"operator,attr"`

	Criterion []GCriterion `xml:"criterion"`
	Criteria  []GCriteria  `xml:"criteria"`

	ExtendDef []GExtendDefinition `xml:"extend_definition"`

	// Extended struct information, not used by XML parser but used during
	// criteria evaluation
	status int
}

type GCriterion struct {
	XMLName xml.Name `xml:"criterion"`

	Test    string `xml:"test_ref,attr"`
	Comment string `xml:"comment,attr"`

	// Extended struct information, not used by XML parser but used during
	// criteria evaluation
	status int
}

type GExtendDefinition struct {
	XMLName xml.Name `xml:"extend_definition"`

	Test    string `xml:"definition_ref,attr"`
	Comment string `xml:"comment,attr"`
}

//
// Tests
//

type GTests struct {
	XMLName xml.Name `xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 tests"`

	RPMInfoTests  []GRPMInfoTest  `xml:"rpminfo_test"`
	DPKGInfoTests []GDPKGInfoTest `xml:"dpkginfo_test"`
	TFC54Tests    []GTFC54Test    `xml:"textfilecontent54_test"`
}

type GTest struct {
	ID     string      `xml:"id,attr"`
	Object GTestObject `xml:"object"`
	State  GTestState  `xml:"state"`

	// Extended struct information, not used by XML parser but used during
	// test evaluation
	status int
	sync.Mutex
}

type GRPMInfoTest struct {
	XMLName xml.Name `xml:"rpminfo_test"`
	GTest
}

type GDPKGInfoTest struct {
	XMLName xml.Name `xml:"dpkginfo_test"`
	GTest
}

type GTFC54Test struct {
	XMLName xml.Name `xml:"textfilecontent54_test"`
	GTest
}

type GTestObject struct {
	XMLName   xml.Name `xml:"object"`
	ObjectRef string   `xml:"object_ref,attr"`
}

type GTestState struct {
	XMLName  xml.Name `xml:"state"`
	StateRef string   `xml:"state_ref,attr"`
}

//
// Objects
//

type GObjects struct {
	XMLName xml.Name `xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 objects"`

	RPMInfoObjects  []GRPMInfoObj  `xml:"rpminfo_object"`
	DPKGInfoObjects []GDPKGInfoObj `xml:"dpkginfo_object"`
	TFC54Objects    []GTFC54Obj    `xml:"textfilecontent54_object"`
}

type GObject struct {
	ID      string `xml:"id,attr"`
	Version string `xml:"version,attr"`
}

type GRPMInfoObj struct {
	XMLName xml.Name `xml:"rpminfo_object"`
	GObject
	Name string `xml:"name"`
}

type GDPKGInfoObj struct {
	XMLName xml.Name `xml:"dpkginfo_object"`
	GObject
	Name string `xml:"name"`
}

type GTFC54Obj struct {
	XMLName xml.Name `xml:"textfilecontent54_object"`
	GObject
	Path     string `xml:"path"`
	Filename string `xml:"filename"`
	Filepath string `xml:"filepath"`
	Pattern  string `xml:"pattern"`
}

//
// States
//

type GStates struct {
	XMLName xml.Name `xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 states"`

	RPMInfoStates  []GRPMInfoState  `xml:"rpminfo_state"`
	TFC54States    []GTFC54State    `xml:"textfilecontent54_state"`
	DPKGInfoStates []GDPKGInfoState `xml:"dpkginfo_state"`
}

type GState struct {
	ID            string        `xml:"id,attr"`
	Version       string        `xml:"version,attr"`
	VersionCheck  GVersionCheck `xml:"version"`
	EVRCheck      GEVRCheck     `xml:"evr"`
	SubExpression string        `xml:"subexpression"`
}

type GTFC54State struct {
	XMLName xml.Name `xml:"textfilecontent54_state"`
	GState
}

type GRPMInfoState struct {
	XMLName xml.Name `xml:"rpminfo_state"`
	GState
}

type GDPKGInfoState struct {
	XMLName xml.Name `xml:"dpkginfo_state"`
	GState
}

type GVersionCheck struct {
	XMLName   xml.Name `xml:"version"`
	Operation string   `xml:"operation,attr"`

	Value string `xml:",chardata"`
}

type GEVRCheck struct {
	XMLName   xml.Name `xml:"evr"`
	Operation string   `xml:"operation,attr"`
	DataType  string   `xml:"datatype,attr"`

	Value string `xml:",chardata"`
}
