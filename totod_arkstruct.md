# Offene Aufgaben `arkstruct`

Separates Go-Projekt: `/Users/martin/dev/go/arkstruct` (eigenes Git-Repo, installierter Binary-Pfad
`/Users/martin/go/bin/arkstruct`). Generiert aus `apis/business-api/handler` (AST-Analyse) die Datei
`client/packages/shared-ui/src/lib/rpc_client.ts` (arktype-Schemas, RPC-Pfad-Konstanten, `RPC_Client`-Klasse).
Aufruf im physication-Repo: `task -t skripte/client_tasks.yml gen-rpc` (führt intern
`arkstruct generate -i ./handler -o ../../client/packages/shared-ui/src/lib/rpc_client.ts` aus
`apis/business-api` aus).

Diese To-dos sind für eine **neue Session im arkstruct-Projekt** gedacht, nicht im physication-Repo.

---

## T1. Nichtdeterministische RPC-Reihenfolge im generierten Output

**Datei:** `generate/generator.go`, Funktion `get_infos()`, ca. Zeile 225–316.

**Ursache:** `rpc_name_map := map[string]RPC{}` wird während der AST-Traversierung befüllt und am Ende so
ausgelesen:

```go
for _, call := range rpc_name_map {
    if call.name == "" || call.path == "" || call.request.Name == "" || call.response.Name == "" {
        fmt.Printf("Ignoring incomplete RPC definition: %+v\n", call)
        continue
    }
    rpcs = append(rpcs, call)
}
```

Iteration über eine Go-`map` ist nicht deterministisch — jeder Lauf von `arkstruct generate` kann die
RPCs in `rpc_client.ts` in anderer Reihenfolge ausgeben. Das erzeugt bei jedem Regenerieren einen
Diff-Rauschen-Commit im physication-Repo, obwohl sich inhaltlich nichts geändert hat.

**Fix:**

1. Nach dem Befüllen von `rpc_name_map`, vor der Schleife, `rpcs` deterministisch sortieren statt
   direkt aus der Map-Iteration zu befüllen — z. B.:

   ```go
   rpcs := make([]RPC, 0, len(rpc_name_map))
   for _, call := range rpc_name_map {
       if call.name == "" || call.path == "" || call.request.Name == "" || call.response.Name == "" {
           fmt.Printf("Ignoring incomplete RPC definition: %+v\n", call)
           continue
       }
       rpcs = append(rpcs, call)
   }
   sort.Slice(rpcs, func(i, j int) bool { return rpcs[i].name < rpcs[j].name })
   ```

   (`"sort"` importieren.) Sortierung nach `call.name` ist stabil und unabhängig von der
   AST-Traversierungsreihenfolge — bevorzugt gegenüber "first-seen"-Reihenfolge, da diese selbst
   wieder von der (nichtdeterministischen) `go/ast`-Verarbeitung abhängen könnte, falls dort ebenfalls
   Maps im Spiel sind.

2. Prüfen, ob es an anderen Stellen in `generator.go` weitere `map`-Iterationen gibt, die in den
   generierten Output einfließen (z. B. Typ-Definitionen, Enums) — gleiches Muster anwenden.

**Regressionstest:** In `generate/generator_test.go` (bestehende `Test_generate_ts` vergleicht aktuell
nur einen einzelnen Lauf gegen die Golden-Datei `test_data/basic.ts`). Ergänzen:

```go
func Test_generate_ts_ist_deterministisch(t *testing.T) {
    var previous string
    for i := 0; i < 20; i++ {
        result := generate_ts(/* gleiche Eingabe wie Test_generate_ts */)
        if i > 0 && result != previous {
            t.Fatalf("generate_ts ist nicht deterministisch (Lauf %d weicht von Lauf %d ab)", i, i-1)
        }
        previous = result
    }
}
```

`test_data/basic.go` enthält bereits 4 RPC-Definitionen (`grep -c "_Path = " test_data/basic.go`) — das
reicht aus, um Map-Reihenfolge-Nichtdeterminismus zuverlässig zu triggern (bei nur 1 RPC gäbe es nichts
zu sortieren).

**Danach im physication-Repo:** `task -t skripte/client_tasks.yml gen-rpc` erneut ausführen und per
`git diff` bestätigen, dass die Datei jetzt bei wiederholten Läufen ohne inhaltliche Änderung stabil
bleibt.

---

## T4. `RPC_Client.#call` liefert nie den HTTP-Status zurück

**Datei:** `generate/generator.go`, Funktion `generate_ts()`, ca. Zeile 107–128.

**Ursache:** Der `#call`-Methodenrumpf ist ein fest einprogrammiertes Template (kein aus Go-Quellcode
abgeleiteter Teil):

```go
ts_code.WriteString("  async #call<TRequest, TResponse>(\n")
ts_code.WriteString("    path: string,\n")
ts_code.WriteString("    args: TRequest,\n")
ts_code.WriteString("  ): Promise<{ value: TResponse; error: null } | { value: null; error: string }> {\n\n")
...
ts_code.WriteString("      if (!result.ok) {\n")
ts_code.WriteString("        console.warn(`Fetch error: ${result.status} ${result.statusText} for ${path}`);\n")
ts_code.WriteString("        if (this.options?.handle_error) this.options.handle_error(result);\n")
ts_code.WriteString("        return {\n")
ts_code.WriteString("          value: null,\n")
ts_code.WriteString("          error: (await result.json())?.message ?? 'Unknown error',\n")
ts_code.WriteString("        };\n")
ts_code.WriteString("      }\n\n")
...
ts_code.WriteString("    } catch (error) {\n")
...
ts_code.WriteString("      return {\n")
ts_code.WriteString("        value: null,\n")
ts_code.WriteString("        error: error instanceof Error ? error.message : \"Unknown error\",\n")
ts_code.WriteString("      };\n")
ts_code.WriteString("    }\n")
```

Der Rückgabetyp (`{ value; error }`) enthält nie den HTTP-Status. Im physication-Repo prüfen mehrere
Aufrufer (`client/apps/plan/src/pages/PlanPage.tsx`, `.../ZustimmungPage.tsx`) `result.status === 404`
— das ist aktuell immer `false`, weil das Feld gar nicht existiert (TypeScript erkennt das nicht, weil
der generierte Typ diesbezüglich keine Struktur vorgibt, die den Fehler auffangen würde — bestätigt via
`npx tsc --noEmit`, 3× TS2339). Betroffene Nutzer sehen die generische Fehlermeldung statt der
404-spezifischen ("... ist nicht mehr verfügbar").

**Fix:**

1. Rückgabetyp um `status` erweitern:

   ```go
   ts_code.WriteString("  ): Promise<{ value: TResponse; error: null; status: number } | { value: null; error: string; status: number | null }> {\n\n")
   ```

2. Im `!result.ok`-Zweig `status: result.status` ergänzen:

   ```go
   ts_code.WriteString("        return {\n")
   ts_code.WriteString("          value: null,\n")
   ts_code.WriteString("          error: (await result.json())?.message ?? 'Unknown error',\n")
   ts_code.WriteString("          status: result.status,\n")
   ts_code.WriteString("        };\n")
   ```

3. Im Erfolgsfall `status: result.status` ergänzen:

   ```go
   ts_code.WriteString("      return {\n")
   ts_code.WriteString("        value: revived as TResponse,\n")
   ts_code.WriteString("        error: null,\n")
   ts_code.WriteString("        status: result.status,\n")
   ts_code.WriteString("      };\n")
   ```

4. Im `catch`-Zweig (Netzwerkfehler, kein HTTP-Response vorhanden) `status: null` ergänzen:

   ```go
   ts_code.WriteString("      return {\n")
   ts_code.WriteString("        value: null,\n")
   ts_code.WriteString("        error: error instanceof Error ? error.message : \"Unknown error\",\n")
   ts_code.WriteString("        status: null,\n")
   ts_code.WriteString("      };\n")
   ```

5. Golden-Fixture `test_data/basic.ts` entsprechend anpassen (die 3 Stellen im generierten `#call` +
   Rückgabetyp), dann `go test ./...` laufen lassen bis `Test_generate_ts` wieder grün ist.

**Danach im physication-Repo:**

1. `task -t skripte/client_tasks.yml gen-rpc` erneut ausführen.
2. `client/apps/plan/src/pages/PlanPage.tsx` und `client/apps/plan/src/pages/ZustimmungPage.tsx`
   prüfen: deren `result.status === 404`-Checks sollten jetzt greifen. Der bestehende e2e-Test
   `client/packages/specs/src/e2e/plan_verlauf_start_page.spec.ts` (Testfall "Ungültiger Eintrag im
   Plan-Verlauf verschwindet nach erneutem Aufruf aus der Liste") war genau daran gescheitert — nach
   dem Fix sollte er grün werden.
3. `npx tsc --project client/tsconfig.json --noEmit` — die 3 vorbestehenden TS2339-Fehler sollten
   verschwinden.
