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

const A_Path = "/a-brauche-einen-langen-text"

type (
	A_Request struct {
		Msg string `json:"msg" ark:"string > 0"`
	}
	A_Response struct {
		Msg string `json:"msg" ark:"string > 0"`
	}
)

func IgnoreMe() {
	//
}
