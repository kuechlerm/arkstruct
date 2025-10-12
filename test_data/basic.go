package test_data

const Eins_Path = "/eins"

type Eins_Request struct {
	RequiredString string `json:"requiredString" validate:"required" ark:"string > 0"`
	OptionalString string `json:"optionalString" ark:"string | undefined"`
	RequiredInt    int    `json:"requiredInt" validate:"required" ark:"number > 0"`
	OptionalInt    int    `json:"optionalInt" ark:"number | undefined"`
	RequiredBool   bool   `json:"requiredBool" validate:"required" ark:"boolean"`
	OptionalBool   bool   `json:"optionalBool" ark:"boolean | undefined"`
}

type Eins_Response struct {
	ResponseString string `json:"responseString" validate:"required" ark:"string > 0"`
}

const Zwei_Path = "/zwei"

type Zwei_Request struct {
	OptionalString string `json:"optionalString" ark:"string | undefined"`
}

type Zwei_Response struct {
	ResponseString string `json:"responseString" validate:"required" ark:"string > 0"`
}

const A_Name_Path = "/a_name"

type (
	A_Name_Request struct {
		Msg string `json:"msg" ark:"string > 0"`
	}
	A_Name_Response struct {
		Msg string `json:"msg" ark:"string > 0"`
	}
)

const Listen_Path = "/listen"

type (
	Ding_DTO struct {
		ID   int    `json:"id" ark:"number"`
		Name string `json:"name" ark:"string > 0"`
	}
	Listen_Request  struct{}
	Listen_Response struct {
		Dinge []Ding_DTO `json:"dinge" ark:"type:Ding_DTO_Schema.array()"`
	}
)

// werden ignoriert, da kein Path
type (
	Request_ohne_Path_Request  struct{}
	Request_ohne_Path_Response struct{}
)

// wird ignoriert, da kein Request
const Path_ohne_Request_Path = "/ohne_request"

type Path_ohne_Request_Response struct{}

func IgnoreMe() {
	//
}
