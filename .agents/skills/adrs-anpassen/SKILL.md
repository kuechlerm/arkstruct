---
name: adrs-anpassen
description: ADRs (Architecture Decision Records) erstellen, aktualisieren oder ablösen. Wird aus der Diskussion heraus aufgerufen, wenn Architekturentscheidungen getroffen oder revidiert werden.
---

## Auftrag

Passe die ADRs des Projekts an. ADRs dokumentieren Architekturentscheidungen mit Kontext, Alternativen und Begründung. Sie sind die Grundlage der Architekturdokumentation.

## Ablauf

1. **Bestehende ADRs sichten** — Lies das ADR-Verzeichnis, um Nummerierung, Format und bestehende Entscheidungen zu kennen.
2. **Änderung verstehen** — Der Entwickler beschreibt die Entscheidung oder Revision. Falls unklar, frage gezielt nach: Was wurde entschieden? Welche Alternativen gab es? Warum diese Wahl?
3. **Prüfen ob ein ADR nötig ist** — Ein ADR ist nur gerechtfertigt, wenn alle drei Kriterien zutreffen:
   - **Schwer umkehrbar** — die Kosten einer späteren Änderung sind bedeutend
   - **Überraschend ohne Kontext** — ein zukünftiger Leser würde sich fragen, warum diese Wahl getroffen wurde
   - **Echtes Trade-off** — es gab reale Alternativen und eine wurde aus konkreten Gründen gewählt

   Typische Kandidaten: Architekturform, Integrationsmuster zwischen Kontexten, Technologiewahl mit Lock-in, Grenzentscheidungen zwischen Modulen, bewusste Abweichungen vom naheliegenden Weg, Constraints die im Code nicht sichtbar sind. Wenn ein Kriterium fehlt, den ADR ablehnen und begründen.
4. **Konflikte prüfen** — Prüfe, ob die neue Entscheidung bestehenden ADRs widerspricht. Bei Widerspruch: betroffene ADRs benennen und klären, ob sie abgelöst oder angepasst werden sollen.
5. **ADR schreiben oder anpassen** — Erstelle einen neuen ADR oder aktualisiere einen bestehenden. Nummerierung: höchste Nummer in `docs/adr/` ermitteln und um eins erhöhen. Dateiname: `NNNN-slug.md`.
6. **Zusammenfassung** — Gib kurz an, was geändert wurde. Weise darauf hin, dass die Architekturdokumentation nachgezogen werden muss (`architektur-doc-anpassen`).

## Regeln

- **Ein ADR pro Entscheidung.** Keine mehreren Entscheidungen in einem ADR vermischen.
- **ADRs nicht löschen.** Überholte Entscheidungen werden als abgelöst markiert und verweisen auf den nachfolgenden ADR.
- **Minimaler Stil.** Ein ADR kann ein einzelner Absatz sein: Kurzer Titel, 1-3 Sätze mit Kontext, Entscheidung und Begründung. Template:
  ```
  # {Kurzer Titel der Entscheidung}

  {1-3 Sätze: Was ist der Kontext, was wurde entschieden, und warum.}
  ```
  Optionale Sektionen (Status, Betrachtete Alternativen, Konsequenzen) nur hinzufügen, wenn sie echten Mehrwert bieten. Bestehenden Stil übernehmen, wenn ADRs bereits vorhanden sind.
- **Architekturdoku-Trigger.** Jede ADR-Änderung erfordert eine Aktualisierung der Architekturdokumentation — darauf explizit hinweisen.