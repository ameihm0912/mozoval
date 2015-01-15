package oval

import (
	"encoding/xml"
	"sync"
)

type GOvalDefinitions struct {
	XMLName		xml.Name	`xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 oval_definitions"`
	Definitions	GDefinitions	`xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 definitions"`
	Tests		GTests		`xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 tests"`
	Objects		GObjects	`xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 objects"`
	States		GStates		`xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 states"`
}

//
// Definitions
//

type GDefinitions struct {
	XMLName		xml.Name	`xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 definitions"`
	Definitions	[]GDefinition	`xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 definition"`
}

type GDefinition struct {
	XMLName		xml.Name	`xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 definition"`

	Metadata	GDefinitionMeta `xml:"metadata"`
	Criteria	GCriteria	`xml:"criteria"`

	ID		string		`xml:"id,attr"`
	Version		string		`xml:"version,attr"`
	Class		string		`xml:"class,attr"`

	// Extended struct information, not used by XML parser but used during
	// definition evaluation
	sync.Mutex
}

type GDefinitionMeta struct {
	XMLName		xml.Name	`xml:"metadata"`

	Title		string		`xml:"title"`
	Decription	string		`xml:"description"`
}

type GCriteria struct {
	XMLName		xml.Name	`xml:"criteria"`

	Operator	string		`xml:"operator,attr"`

	Criterion	[]GCriterion	`xml:"criterion"`
	Criteria	[]GCriteria	`xml:"criteria"`

	ExtendDef	[]GExtendDefinition	`xml:"extend_definition"`

	// Extended struct information, not used by XML parser but used during
	// criteria evaluation
	status		int
}

type GCriterion struct {
	XMLName		xml.Name	`xml:"criterion"`

	Test		string		`xml:"test_ref,attr"`
	Comment		string		`xml:"comment,attr"`

	// Extended struct information, not used by XML parser but used during
	// criteria evaluation
	status		int
}

type GExtendDefinition struct {
	XMLName		xml.Name	`xml:"extend_definition"`

	Test		string		`xml:"definition_ref,attr"`
	Comment		string		`xml:"comment,attr"`
}

//
// Tests
//

type GTests struct {
	XMLName		xml.Name	`xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 tests"`

	RPMInfoTests	[]GRPMInfoTest	`xml:"rpminfo_test"`
}

type GTest struct {
	ID		string		`xml:"id,attr"`
	Object		GTestObject	`xml:"object"`
	State		GTestState	`xml:"state"`
}	

type GRPMInfoTest struct {
	XMLName		xml.Name	`xml:"rpminfo_test"`
	GTest
}

type GTestObject struct {
	XMLName		xml.Name	`xml:"object"`
	ObjectRef	string		`xml:"object_ref,attr"`
}

type GTestState struct {
	XMLName		xml.Name	`xml:"state"`
	StateRef	string		`xml:"state_ref,attr"`
}

//
// Objects
//

type GObjects struct {
	XMLName		xml.Name	`xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 objects"`

	RPMInfoObjects	[]GRPMInfoObj	`xml:"rpminfo_object"`
}

type GObject struct {
	ID		string		`xml:"id,attr"`
	Version		string		`xml:"version,attr"`
}

type GRPMInfoObj struct {
	XMLName		xml.Name	`xml:"rpminfo_object"`
	GObject
	Name		string		`xml:"name"`
}

//
// States
//

type GStates struct {
	XMLName		xml.Name	`xml:"http://oval.mitre.org/XMLSchema/oval-definitions-5 states"`

	RPMInfoStates	[]GRPMInfoState	`xml:"rpminfo_state"`
}

type GState struct {
	ID		string		`xml:"id,attr"`
	Version		string		`xml:"version,attr"`
	VersionCheck	GVersionCheck	`xml:"version"`
	EVRCheck	GEVRCheck	`xml:"evr"`
}

type GRPMInfoState struct {
	XMLName		xml.Name	`xml:"rpminfo_state"`
	GState
}

type GVersionCheck struct {
	XMLName		xml.Name	`xml:"version"`
	Operation	string		`xml:"operation,attr"`

	Value		string		`xml:",chardata"`
}

type GEVRCheck struct {
	XMLName		xml.Name	`xml:"evr"`
	Operation	string		`xml:"operation,attr"`
	DataType	string		`xml:"datatype,attr"`

	Value		string		`xml:",chardata"`
}
