package generate

import (
	"os"
	"strings"
	"testing"
)

// func TestGetRPCInfosMitBeispiel(t *testing.T) {
// 	go_content, err := os.ReadFile("../beispiel/beispiel_handler.go")
// 	if err != nil {
// 		t.Fatalf("Error reading Go file: %v", err)
// 	}
//
// 	rpcs, err := get_rpcs(string(go_content))
// 	if err != nil {
// 		t.Fatalf("Error getting RPCs: %v", err)
// 	}
//
// 	expected := []RPC{
// 		{
// 			path: "/beispielanlegen",
// 			request: Schema{
// 				Name: "BeispielAnlegen_Request",
// 				Properties: []Property{
// 					{"name", "string", "3 <= string <= 100"},
// 				},
// 			},
// 			response: Schema{
// 				Name: "BeispielAnlegen_Response",
// 				Properties: []Property{
// 					{"id", "string", ""},
// 				},
// 			},
// 		},
// 		{
// 			path: "/beispielaendern",
// 			request: Schema{
// 				Name: "BeispielAendern_Request",
// 				Properties: []Property{
// 					{"name", "string", "3 <= string <= 100"},
// 				},
// 			},
// 			response: Schema{
// 				Name: "BeispielAendern_Response",
// 				Properties: []Property{
// 					{"id", "string", ""},
// 				},
// 			},
// 		},
// 	}
//
// 	if !reflect.DeepEqual(rpcs, expected) {
// 		t.Errorf("slices are not equal: %v vs %v", rpcs, expected)
// 	}
// }

func Test_generate_ts(t *testing.T) {
	go_content, err := os.ReadFile("../test_data/basic.go")
	if err != nil {
		t.Fatalf("Error reading Go file: %v", err)
	}
	dtos, rpcs, err := get_infos(string(go_content))
	if err != nil {
		t.Fatalf("Error getting RPCs: %v", err)
	}

	expected_ts_content, err := os.ReadFile("../test_data/basic.ts")
	if err != nil {
		t.Fatalf("Error reading TS file: %v", err)
	}

	ts_result, err := generate_ts(dtos, rpcs)
	if err != nil {
		t.Fatalf("Error generating TS: %v", err)
	}

	// compare line by line
	expected_ts_lines := strings.Split(string(expected_ts_content), "\n")
	ts_result_lines := strings.Split(ts_result, "\n")

	// for _, line := range ts_result_lines {
	// 	fmt.Println(">" + line)
	// }

	for i, expected_line := range expected_ts_lines {

		result_line := ts_result_lines[i]
		if expected_line != result_line {
			t.Errorf("Line %d mismatch:\nExpected: %s\nGot: %s", i+1, expected_line, result_line)
		}
	}
}

// todo: später implementieren
// func TestMapValidation(t *testing.T) {
// 	tests := []struct {
// 		typ      string
// 		validate string
// 		arktype  string
// 	}{
// 		{"string", "required", "string > 0"},
// 		{"string", "", "string | undefined"},
//
// 		{"number", "required", "number > 0"},
// 		{"number", "", "number | undefined"},
//
// 		{"boolean", "required", "true"},
// 		{"boolean", "", "boolean | undefined"},
//
// 		// todo: date, arrays, ... -> Amplenote
// 	}
//
// 	for _, test := range tests {
// 		result := map_validation(test.typ, test.validate)
// 		if result != test.arktype {
// 			t.Errorf("MapValidation(%q, %q) = %q; want %q", test.typ, test.validate, result, test.arktype)
// 		}
// 	}
// }
