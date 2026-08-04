---
name: storyboard-anpassen
description: Storyboard anpassen — Aktivitäten, Schritte, User Stories und Akzeptanzkriterien hinzufügen, ändern oder entfernen. Wird aus der Diskussion heraus aufgerufen.
---

## Auftrag

Passe das Storyboard des Projekts an. Das Storyboard bildet die Nutzeranforderungen ab und ist direkt mit Tests verknüpft.

## Struktur

### Dateien

Die Storyboard-Dateien liegen in `docs/storyboard/`. Jede Datei beschreibt genau eine Aktivität. Dateinamen in `snake_case.md`. Eine optionale `_order.md` steuert die Anzeigereihenfolge der Aktivitäten mit `[[dateiname]]`-Syntax.

### Hierarchie

Jede Storyboard-Datei ist hierarchisch über Markdown-Überschriften aufgebaut:

- `#` **Aktivität** — Beschreibt einen abgeschlossenen Nutzungsablauf (z.B. „Bestellung aufgeben", „Konto verwalten"). Bildet die Gliederungsebene des Storyboards. Genau eine pro Datei.
- `##` **Schritt** — Unterteilt eine Aktivität in zeitlich oder logisch aufeinanderfolgende Phasen.
- `###` **User Story** — Beschreibt ein konkretes, vom Nutzer erwartetes Verhalten innerhalb eines Schritts. Granularität: testbar durch eine Handvoll Akzeptanzkriterien.
- **Akzeptanzkriterium** — Ein einzelner prüfbarer Fall, implementiert als Playwright-Test. Akzeptanzkriterien werden nicht im Storyboard gepflegt, sondern leben in den Testdateien.

Text unter einer User Story (`###`) beschreibt deren Kontext.

### Verknüpfung mit Tests

Playwright-Tests verweisen über eine Fixture auf die User Story:

```ts
_.user_story('Mit_Email_einloggen');
```

Der Parameter ist der slugifizierte `###`-Titel. Slugifizierung:

1. Leerzeichen → `_`
2. Alles außer `\w` und `äöüÄÖÜß` entfernen

Ergänzende Elemente (Personas, Wireframes etc.) sind möglich.

## Kontext laden

1. **Glossar laden** — `docs/glossar.md` lesen für korrekte Begriffe.
2. **Storyboard laden** — Betroffene Dateien aus `docs/storyboard/` lesen.
3. **Code/Tests gezielt nachladen** — Nur bei Bedarf, um bestehende Akzeptanzkriterien gegen Tests abzugleichen.

## Ablauf

1. **Änderung verstehen** — Der Entwickler beschreibt die gewünschte Anpassung. Falls unklar, gezielt nachfragen.
2. **Auswirkungen prüfen** — Betrifft die Änderung bestehende Akzeptanzkriterien oder Tests? Überschneidungen mit anderen User Stories? Konflikte benennen.
3. **Storyboard anpassen** — Änderungen durchführen. Format und Struktur der bestehenden Dateien beibehalten.
4. **User-Story-Titel umbenannt?** — Wenn ein `###`-Titel geändert wurde: den alten slugifizierten Namen im gesamten Projekt suchen und in den `_.user_story()`-Aufrufen der Testdateien ersetzen. Danach `npm run gen-user-story-headers` ausführen, um die TypeScript-Typen zu aktualisieren.
5. **`_order.md` pflegen** — Bei neuen Aktivitäten oder Änderungen der Ablaufreihenfolge die `_order.md` anpassen.
6. **Zusammenfassung** — Angeben, was geändert wurde. Geänderte Testdateien und Typen-Regenerierung auflisten.

## Regeln

- **Hierarchie einhalten.** User Stories gehören zu Schritten, Schritte zu Aktivitäten. Akzeptanzkriterien gehören in Testdateien und verweisen auf User Stories.
- **Akzeptanzkriterien prüfbar formulieren.** Jedes Kriterium muss als Playwright-Test umsetzbar sein.
- **Keine eigenmächtigen Ergänzungen.** Nur ändern, was der Entwickler vorgibt oder was sich aus der Diskussion ergibt.
- **Format beibehalten.** Struktur und Stil der bestehenden Storyboard-Dateien übernehmen.
- **Eine Aktivität pro Datei.** Jede Storyboard-Datei enthält genau eine `#`-Überschrift.
- **Dateinamen in `snake_case.md`.**