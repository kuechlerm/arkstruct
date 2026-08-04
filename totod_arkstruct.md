# Offene Aufgaben `arkstruct`

Separates Go-Projekt: `/Users/martin/dev/go/arkstruct` (eigenes Git-Repo, installierter Binary-Pfad
`/Users/martin/go/bin/arkstruct`). Generiert aus `apis/business-api/handler` (AST-Analyse) die Datei
`client/packages/shared-ui/src/lib/rpc_client.ts` (arktype-Schemas, RPC-Pfad-Konstanten, `RPC_Client`-Klasse).
Aufruf im physication-Repo: `task -t skripte/client_tasks.yml gen-rpc` (führt intern
`arkstruct generate -i ./handler -o ../../client/packages/shared-ui/src/lib/rpc_client.ts` aus
`apis/business-api` aus).

Diese To-dos sind für eine **neue Session im arkstruct-Projekt** gedacht, nicht im physication-Repo.

---
