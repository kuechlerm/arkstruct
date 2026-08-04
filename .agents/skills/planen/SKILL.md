---
name: planen
description: Trifft technische Entscheidungen (ADRs) und erstellt konkrete Implementierungspläne aus Glossar, Storyboard, Architekturdokumentation und Code.
disable-model-invocation: true
---

## Auftrag

Triff technische Entscheidungen und erstelle einen konkreten, schrittweisen Implementierungsplan. Der Plan muss so klar sein, dass er in einer frischen Session ohne weiteren Kontext abgearbeitet werden kann.

## Kontext laden

1. **Immer laden** — Glossar, Storyboard (relevante Teile) und Architekturdokumentation.
2. **Gezielt nachladen** — Code-Bereiche, die von der Änderung betroffen sind.

Code ist die Quelle der Wahrheit. Den Plan gegen den tatsächlichen Code-Stand erstellen, nicht gegen Annahmen.

## Ablauf

1. **Anforderungen erfassen** — Welche Storyboard-Elemente oder Diskussionsergebnisse sollen umgesetzt werden? Falls unklar, gezielt nachfragen.
2. **Technische Entscheidungen treffen** — Wie werden die Anforderungen umgesetzt? Bei Entscheidungen, die einer Erklärung bedürfen:
   - **ADRs** → `adrs-anpassen` aufrufen
   - **Architekturdoku** → `architektur-doc-anpassen` aufrufen (direkt oder aus ADRs abgeleitet)
   ADRs werden finalisiert, bevor der Plan erstellt wird. Bei großen Entscheidungen kann die Session nach den ADRs enden.
3. **Betroffenen Code analysieren** — Relevante Dateien und Strukturen identifizieren. Abhängigkeiten und Seiteneffekte erkennen. Prefactoring-Möglichkeiten identifizieren: Kann der bestehende Code so umstrukturiert werden, dass die eigentliche Änderung einfacher wird?
4. **Plan erstellen** — Vertikale Schritte: Jeder Schritt soll nach Möglichkeit einen durchgängigen Pfad durch die betroffenen Schichten abdecken und für sich validierbar sein. Akzeptanzkriterien aus dem Storyboard als bevorzugte Validierung. Wenn Prefactoring nötig ist, wird es als eigener, vorangestellter Schritt aufgenommen. Jeder Schritt benennt:
   - Was geändert wird (Datei, Funktion, Komponente)
   - Was das erwartete Ergebnis ist
   - Welche Tests erstellt oder angepasst werden (Akzeptanzkriterien aus dem Storyboard bevorzugen)
5. **Review mit Entwickler** — Plan vorlegen und bestätigen lassen, bevor er geschrieben wird.
6. **Plan schreiben** — Als Implementierungsplan-Datei unter `docs/implementierungsplaene/{slug}.md` ablegen. Template:

```markdown
# Implementierungsplan: {Titel}

## Kontext
{1–2 Sätze: Was wird umgesetzt und warum}

## Referenzen
{Betroffene Storyboard-Aktivitäten, ADRs}

## Schritte
### 1. {Schritt-Titel}
- **Änderung:** {Was wird geändert}
- **Ergebnis:** {Erwartetes Ergebnis / Akzeptanzkriterium}
- **Validierung:** {Welcher Test}
- [ ] Erledigt
```

## Regeln

- **Selbstständig abarbeitbar.** Der Plan muss alle nötigen Informationen enthalten — Dateinamen, Funktionsnamen, Glossar-Begriffe. Eine frische Session mit ~140k Token Kontextbudget muss ihn ohne Rückfragen umsetzen können.
- **Vertikale Schritte.** Jeder Schritt deckt einen durchgängigen Pfad durch die betroffenen Schichten ab und ist für sich validierbar. Reine Refactoring-Schritte (Prefactoring) sind die Ausnahme.
- **Plan aufteilen bei großer Komplexität.** Wenn der Plan zu groß für eine einzelne Session wird, in mehrere Teil-Pläne aufteilen. Jeder Teil-Plan ist für sich in einer Session abarbeitbar und benennt seine Abhängigkeiten zu anderen Teil-Plänen.
- **Tests einplanen.** Akzeptanzkriterien aus dem Storyboard in Playwright-Tests übersetzen. Unit-Tests wo sinnvoll.
- **Konzeptionelle Lücken melden.** Wenn beim Planen auffällt, dass Glossar-Begriffe fehlen, Anforderungen unklar sind oder User Stories widersprüchlich sind: den Entwickler informieren und klären, ob die Lücke direkt gelöst werden kann oder ein Rücksprung in die `diskussion` nötig ist.
- **Kein Code schreiben.** Der Plan beschreibt, was zu tun ist — die Umsetzung erfolgt durch `implementierung`.