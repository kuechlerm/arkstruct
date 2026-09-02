package generate

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
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
	name     string
	path     string
	request  Schema
	response Schema
}

// Ein freies Schema ist ein DTO
type (
	DTOs []Schema
	RPCs []RPC
)

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

	ts_code, err := generate_ts(all_dtos, all_rpcs)
	if err != nil {
		return errors.New("Error generating TypeScript code: " + err.Error())
	}

	err = os.WriteFile(target_path, []byte(ts_code), 0o644)
	if err != nil {
		return errors.New("Error writing TypeScript file: " + err.Error())
	}

	return nil
}

func generate_ts(dtos DTOs, rpcs RPCs) (string, error) {
	ts_code := &strings.Builder{}
	ts_code.WriteString(`import { type } from "arktype";`)
	ts_code.WriteString("\n\n")

	geprueft := gepruefte_schemas(dtos, rpcs)

	for _, dto := range dtos {
		write_schema(ts_code, dto, geprueft[dto.Name])
	}

	for _, rpc := range rpcs {
		write_path(ts_code, rpc.name, rpc.path)
		write_schema(ts_code, rpc.request, geprueft[rpc.request.Name])
		write_schema(ts_code, rpc.response, geprueft[rpc.response.Name])
	}

	// rpc client class
	ts_code.WriteString("export class RPC_Client {\n")
	ts_code.WriteString("  constructor(\n")
	ts_code.WriteString("    private base_url: string,\n")
	ts_code.WriteString("    private options?: {\n")
	ts_code.WriteString("      fetch?: (url: string, init: RequestInit) => Promise<Response>;\n")
	ts_code.WriteString("      handle_error?: (response: Response) => void;\n")
	ts_code.WriteString("      log?: { warn(meldung: unknown, ...details: unknown[]): void };\n")
	ts_code.WriteString("    },\n")
	ts_code.WriteString("  ) { }\n\n")

	ts_code.WriteString("  async #call<TRequest, TResponse>(\n")
	ts_code.WriteString("    path: string,\n")
	ts_code.WriteString("    args: TRequest,\n")
	ts_code.WriteString("    schema: (data: unknown) => unknown,\n")
	ts_code.WriteString("  ): Promise<{ value: TResponse; error: null; status: number } | { value: null; error: string; status: number | null; body: unknown }> {\n\n")
	ts_code.WriteString("    const do_fetch = this.options?.fetch ?? globalThis.fetch;\n")
	ts_code.WriteString("    const do_log = this.options?.log ?? console;\n\n")
	ts_code.WriteString("    try {\n")
	ts_code.WriteString("      const result = await do_fetch(new URL(path, this.base_url).href, {\n")
	ts_code.WriteString("        method: \"POST\",\n")
	ts_code.WriteString("        credentials: \"include\",\n")
	ts_code.WriteString("        headers: {\n")
	ts_code.WriteString("          \"Content-Type\": \"application/json\",\n")
	ts_code.WriteString("        },\n")
	ts_code.WriteString("        body: JSON.stringify(args),\n")
	ts_code.WriteString("      });\n\n")
	ts_code.WriteString("      if (!result.ok) {\n")
	ts_code.WriteString("        do_log.warn(`Fetch error: ${result.status} ${result.statusText} for ${path}`);\n")
	ts_code.WriteString("        if (this.options?.handle_error) this.options.handle_error(result);\n")
	// Der Fehlerkoerper folgt dem Fehlermodell, nicht dem Response-Schema, und wird ungeprueft und
	// unveraendert durchgereicht (ADR-0035 §6 im Physication-Repo).
	ts_code.WriteString("        const fehler_body = await result.json().catch(() => null);\n")
	ts_code.WriteString("        return {\n")
	ts_code.WriteString("          value: null,\n")
	ts_code.WriteString("          error: (fehler_body as { message?: string } | null)?.message ?? 'Unknown error',\n")
	ts_code.WriteString("          status: result.status,\n")
	ts_code.WriteString("          body: fehler_body,\n")
	ts_code.WriteString("        };\n")
	ts_code.WriteString("      }\n\n")
	ts_code.WriteString("      const geprueft = schema(await result.json());\n\n")
	// Gemeldet werden Pfad und erwarteter Typ, nie der vorgefundene Wert: `actual`, `summary` und
	// `message` tragen ihn (ADR-0025 §3 im Physication-Repo). Eine Verletzung kommt als unlesbare
	// Antwort zurueck - `status: null` wie bei einer Antwort, die nie zustande kam (ADR-0035 §3).
	ts_code.WriteString("      if (geprueft instanceof type.errors) {\n")
	ts_code.WriteString("        const stellen = geprueft.map((fehler) => `${fehler.path.join(\".\")} erwartet ${fehler.expected}`).join(\"; \");\n")
	ts_code.WriteString("        do_log.warn(`Antwort verletzt Vertrag: ${path} — ${stellen}`);\n\n")
	ts_code.WriteString("        return {\n")
	ts_code.WriteString("          value: null,\n")
	ts_code.WriteString("          error: \"Antwort verletzt den Vertrag\",\n")
	ts_code.WriteString("          status: null,\n")
	ts_code.WriteString("          body: null,\n")
	ts_code.WriteString("        };\n")
	ts_code.WriteString("      }\n\n")
	// Der Rueckgabewert ist das Ergebnis des Schemas, nicht die rohe Antwort - nur so wirken die
	// Datums-Morphs aus map_date_tag.
	ts_code.WriteString("      return {\n")
	ts_code.WriteString("        value: geprueft as TResponse,\n")
	ts_code.WriteString("        error: null,\n")
	ts_code.WriteString("        status: result.status,\n")
	ts_code.WriteString("      };\n")
	ts_code.WriteString("    } catch (error) {\n")
	// Die Aufrufargumente gehen bewusst nicht ins Log: sie tragen bei `nutzer_login` das
	// Klartext-Passwort und bei `push_changes` Gesundheitsdaten, und dieser catch-Zweig laeuft
	// bei jedem Offline-Aufruf. Siehe ADR-0025 §3 im Physication-Repo.
	ts_code.WriteString("      do_log.warn(`RPC_Client Error for ${path}`);\n")
	ts_code.WriteString("      do_log.warn(error);\n\n")
	ts_code.WriteString("      return {\n")
	ts_code.WriteString("        value: null,\n")
	ts_code.WriteString("        error: error instanceof Error ? error.message : \"Unknown error\",\n")
	ts_code.WriteString("        status: null,\n")
	ts_code.WriteString("        body: null,\n")
	ts_code.WriteString("      };\n")
	ts_code.WriteString("    }\n")
	ts_code.WriteString("  }\n\n")

	for idx, rpc := range rpcs {
		trenner_index := strings.LastIndex(rpc.request.Name, "_")
		if trenner_index == -1 {
			continue
		}

		ts_code.WriteString(
			"  " +
				strings.ToLower(rpc.request.Name[:trenner_index]) +
				" = (args: " + rpc.request.Name + ") =>\n")

		ts_code.WriteString(
			"    this.#call<" +
				rpc.request.Name +
				", " +
				rpc.response.Name +
				">(" + rpc.name + "_Path, args, " + rpc.response.Name + "_Schema);\n")

		if idx < len(rpcs)-1 {
			ts_code.WriteString("\n")
		}
	}

	ts_code.WriteString("}\n")

	return ts_code.String(), nil
}

func write_path(ts_code *strings.Builder, name string, path string) {
	fmt.Fprintf(ts_code, "export const %s_Path = \"%s\";\n", name, path)
}

func write_schema(ts_code *strings.Builder, schema Schema, geprueft bool) {
	fmt.Fprintf(ts_code, "export const %s_Schema = type({", schema.Name)

	for idx, prop := range schema.Properties {
		if idx == 0 {
			ts_code.WriteString("\n")
		}
		typ := prop.Type
		if geprueft {
			typ = map_date_tag(typ)
		}
		if strings.HasPrefix(typ, "type:") {
			fmt.Fprintf(ts_code, `  %s: %s,`, prop.Name, strings.TrimPrefix(typ, "type:"))
		} else {
			fmt.Fprintf(ts_code, `  %s: "%s",`, prop.Name, typ)
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
				// todo: check / Fehler loggen?
				const_name := const_spec.Names[0].Name

				// Wir suchen nach Konstanten, die mit "_Path" enden
				// todo: check / Fehler loggen?
				literal, ok := const_spec.Values[0].(*ast.BasicLit)
				if !ok || literal.Kind != token.STRING || !strings.HasSuffix(const_name, "_Path") {
					continue
				}

				// muss gleich specName sein, damit Zuordnung stimmt
				trenner_index := strings.LastIndex(const_name, "_")
				if trenner_index == -1 {
					continue
				}
				const_spec_name := const_name[:trenner_index]
				if const_spec_name == "" {
					continue
				}

				// todo: check / Fehler loggen?
				rpc := rpc_name_map[const_spec_name]
				rpc.path = strings.Trim(literal.Value, "\"")
				// todo: check / Fehler loggen?
				rpc_name_map[const_spec_name] = rpc

				continue
			}

			type_spec, ok := spec.(*ast.TypeSpec)
			if !ok || (!strings.HasSuffix(type_spec.Name.Name, "_DTO") && !strings.HasSuffix(type_spec.Name.Name, "_Request") && !strings.HasSuffix(type_spec.Name.Name, "_Response")) {
				continue
			}
			trenner_index := strings.LastIndex(type_spec.Name.Name, "_")
			if trenner_index == -1 {
				continue
			}
			spec_name := type_spec.Name.Name[:trenner_index]
			if spec_name == "" {
				continue
			}

			if _, ok := type_spec.Type.(*ast.StructType); ok {
				if strings.HasSuffix(type_spec.Name.Name, "_DTO") {
					dtos = append(dtos, map_schema(type_spec))
				} else {

					// check, ob Path für diesen Request/Response existiert findet am Ende statt

					// todo: check / Fehler loggen?
					call := rpc_name_map[spec_name]
					call.name = spec_name

					if strings.HasSuffix(type_spec.Name.Name, "_Request") {
						call.request = map_schema(type_spec)
					}

					if strings.HasSuffix(type_spec.Name.Name, "_Response") {
						call.response = map_schema(type_spec)
					}

					// todo: check / Fehler loggen?
					rpc_name_map[spec_name] = call
				}
			}
		}
	}

	for _, call := range rpc_name_map {
		// check, ob path, request und response gesetzt sind
		if call.name == "" || call.path == "" || call.request.Name == "" || call.response.Name == "" {
			fmt.Printf("Ignoring incomplete RPC definition: %+v\n", call)
			continue
		}
		rpcs = append(rpcs, call)
	}
	sort.Slice(rpcs, func(i, j int) bool { return rpcs[i].name < rpcs[j].name })

	return dtos, rpcs, nil
}

func map_schema(typeSpec *ast.TypeSpec) Schema {
	properties := []Property{}

	for _, field := range typeSpec.Type.(*ast.StructType).Fields.List {
		// if typeSpec.Name.Name == "Ding_DTO" {
		// 	fmt.Printf("Processing field: %+v\n", field)
		// }

		// todo: fix?
		// if field.Names == nil {
		// 	// This is an embedded field!
		// 	// field.Type will be the embedded type (e.g., *ast.Ident or *ast.SelectorExpr)
		// 	fmt.Printf("Ignoring embedded field in type %s\n", typeSpec.Name.Name)
		// 	switch t := field.Type.(type) {
		// 	case *ast.Ident:
		// 		println("Embedded:", t.Name)
		// 	case *ast.SelectorExpr:
		// 		// For imported embedded types
		// 		println("Embedded:", t.X.(*ast.Ident).Name+"."+t.Sel.Name)
		// 	}
		// 	continue
		// }

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
					field_type = tag.Value()
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

		// todo: check / Fehler loggen?
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

// gepruefte_schemas nennt die Schemas, die zur Laufzeit wirklich durchlaufen werden: jede
// *_Response und alles, was von ihr aus erreichbar ist. Nur dort wird aus `Date` eine Morph.
//
// Ein Request-Schema bekommt sie nicht. Es wird nie ausgefuehrt - und eine Morph in einer
// unsortierten Union (`type.or(...)`, wie im Payload von Push_Changes) macht die Union fuer
// arktype unbestimmt und bricht schon beim Laden des Moduls.
func gepruefte_schemas(dtos DTOs, rpcs RPCs) map[string]bool {
	nach_name := map[string]Schema{}
	for _, dto := range dtos {
		nach_name[dto.Name] = dto
	}
	for _, rpc := range rpcs {
		nach_name[rpc.response.Name] = rpc.response
	}

	geprueft := map[string]bool{}
	var lauf func(string)
	lauf = func(name string) {
		if geprueft[name] {
			return
		}
		geprueft[name] = true
		for _, prop := range nach_name[name].Properties {
			for _, referenz := range referenzierte_schemas(prop.Type) {
				lauf(referenz)
			}
		}
	}
	for _, rpc := range rpcs {
		lauf(rpc.response.Name)
	}

	return geprueft
}

// referenzierte_schemas liest aus einem `type:`-Tag die Namen der Schemas heraus, auf die er sich
// beruft - egal ob `Ding_DTO_Schema.array()` oder `type.or(A_Schema, B_Schema)`.
func referenzierte_schemas(ark_tag string) []string {
	return schema_referenz.FindAllString(ark_tag, -1)
}

var schema_referenz = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_]*(?:_DTO|_Response)`)

// map_date_tag bildet die drei Datums-Tags auf die arktype-Morph ab, die den ISO-String beim
// Pruefen in ein Date ueberfuehrt. Der Go-Tag bleibt lesbar `Date`, der abgeleitete
// TypeScript-Typ bleibt `Date` - die Umwandlung passiert genau an den deklarierten Stellen und
// sonst nirgends. Jeder andere Tag geht unveraendert durch.
func map_date_tag(ark_tag string) string {
	switch ark_tag {
	case "Date":
		return "string.date.iso.parse"
	case "Date?":
		return "string.date.iso.parse?"
	case "Date | null":
		return "string.date.iso.parse | null"
	default:
		return ark_tag
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
