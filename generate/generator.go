package generate

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"

	"github.com/fatih/structtag"
)

type Property struct {
	Name       string // json Name
	Type       string // TS Type
	Validation string // Ark Validation
}

type Schema struct {
	Name       string
	Properties []Property
}

type RPC struct {
	path     string
	request  Schema
	response Schema
}

// Ein freies Schema ist ein DTO
type DTOs []Schema

func (d DTOs) Len() int {
	return len(d)
}

func (d DTOs) Less(i, j int) bool {
	return d[i].Name < d[j].Name
}

func (d DTOs) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type RPCs []RPC

func (r RPCs) Len() int {
	return len(r)
}

func (r RPCs) Less(i, j int) bool {
	return r[i].request.Name < r[j].request.Name
}

func (r RPCs) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func Generate(go_folder_path, target_path string) error {
	folder, err := os.ReadDir(go_folder_path)
	if err != nil {
		return errors.New("Error reading folder: " + err.Error())
	}

	all_dtos := DTOs{}
	all_rpcs := RPCs{}
	for _, file := range folder {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue // no directories, only go files
		}

		file_path := go_folder_path + "/" + file.Name()

		content, err := os.ReadFile(file_path)
		if err != nil {
			return errors.New("Error reading Go file: " + err.Error())
		}

		dtos, rpcs, err := get_infos(string(content))
		if err != nil {
			return errors.New("Error getting RPCs: " + err.Error())
		}

		all_dtos = append(all_dtos, dtos...)
		all_rpcs = append(all_rpcs, rpcs...)
	}

	sort.Sort(all_dtos)
	sort.Sort(all_rpcs)

	ts_code, err := generate_ts(all_dtos, all_rpcs)
	if err != nil {
		return errors.New("Error generating TypeScript code: " + err.Error())
	}

	err = os.WriteFile(target_path, []byte(ts_code), 0644)
	if err != nil {
		return errors.New("Error writing TypeScript file: " + err.Error())
	}

	return nil
}

func generate_ts(dtos DTOs, rpcs RPCs) (string, error) {
	ts_code := &strings.Builder{}
	ts_code.WriteString(`import { type } from "arktype";`)
	ts_code.WriteString("\n\n")

	for _, dto := range dtos {
		write_schema(ts_code, dto)
	}

	for _, rpc := range rpcs {
		write_schema(ts_code, rpc.request)
		write_schema(ts_code, rpc.response)
	}

	// rpc client class
	ts_code.WriteString("export class RPC_Client {\n")
	ts_code.WriteString("  constructor(public base_url: string) {}\n\n")
	ts_code.WriteString("  async #call<TRequest, TResponse>(\n")
	ts_code.WriteString("    path: string,\n")
	ts_code.WriteString("    args: TRequest,\n")
	ts_code.WriteString("  ): Promise<\n")
	ts_code.WriteString("    { value: TResponse; error: null } | { value: null; error: string }\n")
	ts_code.WriteString("  > {\n")
	ts_code.WriteString("    try {\n")
	ts_code.WriteString("      const result = await fetch(new URL(path, this.base_url).href, {\n")
	ts_code.WriteString("        method: \"POST\",\n")
	ts_code.WriteString("        body: JSON.stringify(args),\n")
	ts_code.WriteString("      });\n\n")
	ts_code.WriteString("      if (!result.ok) {\n")
	ts_code.WriteString("        console.error(\n")
	ts_code.WriteString("          `Fetch error: ${result.status} ${result.statusText} for ${path}`,\n")
	ts_code.WriteString("        );\n")
	ts_code.WriteString("        return {\n")
	ts_code.WriteString("          value: null,\n")
	ts_code.WriteString("          error: `Fetch error: ${result.status} ${result.statusText}`,\n")
	ts_code.WriteString("        };\n")
	ts_code.WriteString("      }\n\n")
	ts_code.WriteString("      const data = await result.json();\n\n")
	ts_code.WriteString("      return {\n")
	ts_code.WriteString("        value: data as TResponse,\n")
	ts_code.WriteString("        error: null,\n")
	ts_code.WriteString("      };\n")
	ts_code.WriteString("    } catch (error) {\n")
	ts_code.WriteString("      console.error(`Error during fetch for ${path}:`, error);\n\n")
	ts_code.WriteString("      return {\n")
	ts_code.WriteString("        value: null,\n")
	ts_code.WriteString("        error: error instanceof Error ? error.message : \"Unknown error\",\n")
	ts_code.WriteString("      };\n")
	ts_code.WriteString("    }\n")
	ts_code.WriteString("  }\n\n")

	for idx, rpc := range rpcs {
		ts_code.WriteString(
			"  " +
				strings.ToLower(rpc.request.Name[:strings.LastIndex(rpc.request.Name, "_")]) +
				" = (args: " + rpc.request.Name + ") =>\n")
		ts_code.WriteString(
			"    this.#call<" +
				rpc.request.Name +
				", " +
				rpc.response.Name +
				">(\"" + rpc.path + "\", args);\n")

		if idx < len(rpcs)-1 {
			ts_code.WriteString("\n")
		}
	}

	ts_code.WriteString("}\n")

	return ts_code.String(), nil
}

func write_schema(ts_code *strings.Builder, schema Schema) {
	fmt.Fprintf(ts_code, "export const %s_Schema = type({", schema.Name)

	for idx, prop := range schema.Properties {
		if idx == 0 {
			ts_code.WriteString("\n")
		}
		if strings.HasPrefix(prop.Type, "type:") {
			fmt.Fprintf(ts_code, `  %s: %s,`, prop.Name, strings.TrimLeft(prop.Type, "type:"))
		} else {
			fmt.Fprintf(ts_code, `  %s: "%s",`, prop.Name, prop.Type)
		}
		ts_code.WriteString("\n")
	}
	ts_code.WriteString("});\n")

	fmt.Fprintf(ts_code, "export type %s = typeof %s_Schema.infer;\n\n", schema.Name, schema.Name)
}

func get_infos(file_content string) (DTOs, RPCs, error) {
	dtos := DTOs{}
	rpcs := RPCs{}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", file_content, parser.AllErrors)
	if err != nil {
		return dtos, rpcs, errors.New("Error parsing Go file: " + err.Error())
	}

	rpc_name_map := map[string]RPC{}

	for _, decl := range node.Decls {
		// fmt.Println(decl)
		gen_decl, ok := decl.(*ast.GenDecl)

		if !ok || (gen_decl.Tok != token.TYPE && gen_decl.Tok != token.CONST) {
			// wir interessieren uns nur für Typ- und Konstantendeklarationen mit passenden Endungen
			continue
		}

		for _, spec := range gen_decl.Specs {
			const_spec, ok := spec.(*ast.ValueSpec)

			if ok {
				// todo: check für Names[0]?
				const_name := const_spec.Names[0].Name

				// Wir suchen nach Konstanten, die mit "_Path" enden
				literal, ok := const_spec.Values[0].(*ast.BasicLit)
				if !ok || literal.Kind != token.STRING || !strings.HasSuffix(const_name, "_Path") {
					continue
				}

				// muss gleich specName sein, damit Zuordnung stimmt
				const_spec_name := const_name[:strings.LastIndex(const_name, "_")]
				if const_spec_name == "" {
					continue
				}

				rpc := rpc_name_map[const_spec_name]
				rpc.path = strings.Trim(literal.Value, "\"")
				rpc_name_map[const_spec_name] = rpc

				continue
			}

			type_spec, ok := spec.(*ast.TypeSpec)
			if !ok || (!strings.HasSuffix(type_spec.Name.Name, "_DTO") && !strings.HasSuffix(type_spec.Name.Name, "_Request") && !strings.HasSuffix(type_spec.Name.Name, "_Response")) {
				continue
			}
			spec_name := type_spec.Name.Name[:strings.LastIndex(type_spec.Name.Name, "_")]
			if spec_name == "" {
				continue
			}

			if _, ok := type_spec.Type.(*ast.StructType); ok {
				if strings.HasSuffix(type_spec.Name.Name, "_DTO") {
					dtos = append(dtos, map_schema(type_spec))
				} else {

					call := rpc_name_map[spec_name]

					if strings.HasSuffix(type_spec.Name.Name, "_Request") {
						call.request = map_schema(type_spec)
					}

					if strings.HasSuffix(type_spec.Name.Name, "_Response") {
						call.response = map_schema(type_spec)
					}

					rpc_name_map[spec_name] = call
				}
			}
		}
	}

	for _, call := range rpc_name_map {
		rpcs = append(rpcs, call)
	}

	return dtos, rpcs, nil
}

func map_schema(typeSpec *ast.TypeSpec) Schema {
	properties := []Property{}

	for _, field := range typeSpec.Type.(*ast.StructType).Fields.List {

		// ##### Type
		field_type := ""
		switch ft := field.Type.(type) {
		case *ast.Ident:
			field_type = go_type_to_ark_type(ft.Name)
		default:
			field_type = "any"
		}

		// ##### Tags
		json_property_name := ""
		if field.Tag != nil {

			tags, err := structtag.Parse(strings.Trim(field.Tag.Value, "`"))
			if err != nil {
				fmt.Printf("Error parsing tags for field %s: %v\n", field.Names[0].Name, err)
				continue
			}

			// fmt.Printf("Processing field: %s || Tags: %s \n", field.Names, tags.Tags())

			for _, tag := range tags.Tags() {
				if tag.Key == "json" {
					// ist der erste Tag-Wert
					json_property_name = tag.Name
				}

				if tag.Key == "ark" {
					// fmt.Printf("Ark tag found: %s\n", tag.Name)
					// hier wird der Ark-Type gesetzt
					field_type = tag.Name
				}

				// if tag.Key == "validate" {
				// 	// fmt.Printf("Validation tag found: %s\n", tag.Name)
				// 	// fmt.Printf("Validation OPTIONS tag found: %s\n", tag.Options)
				// 	fieldType = map_validation(fieldType, tag.Name)
				// 	validate_tag_used = true
				// }
			}

			// if !validate_tag_used {
			// 	fieldType = map_validation(fieldType, "")
			// }
		}

		name := field.Names[0].Name
		if json_property_name != "" {
			name = json_property_name // wenn json-Name vorhanden, dann diesen verwenden
		}
		properties = append(properties, Property{
			Name:       name, // json name
			Type:       field_type,
			Validation: "TODO", // TODO: hier müsste die Validation aus den Struct-Tags geholt werden
		})

	}

	return Schema{
		Name:       typeSpec.Name.Name,
		Properties: properties,
	}
}

// Converts Go type to ArkType type
func go_type_to_ark_type(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64", "float32", "float64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "number"
	case "bool":
		return "boolean"
	default:
		return "any" // fallback
	}
}

// todo: später implementieren
// func map_validation(ts_typ, validation string) string {
// 	if ts_typ == "string" || ts_typ == "number" {
// 		switch validation {
// 		case "required":
// 			return ts_typ + " > 0"
// 		case "":
// 			return ts_typ + " | undefined"
// 		}
// 	}
//
// 	if ts_typ == "boolean" {
// 		switch validation {
// 		case "required":
// 			return "true"
// 		case "":
// 			return ts_typ + " | undefined"
// 		}
// 	}
//
// 	return "todo"
// }
