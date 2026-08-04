---
name: diskussion
description: Strukturierte Diskussion über Konzepte, Anforderungen und Geschäftslogik. Ergebnis sind Anpassungen an Glossar, Storyboard und Versuche.
disable-model-invocation: true
---

## Auftrag

Führe eine intensive, strukturierte Diskussion über die vom Entwickler eingebrachte Idee, Änderung oder das Problem. Ziel ist ein gemeinsames Verständnis auf konzeptioneller Ebene — Was und Warum klären, nicht das Wie.

## Kontext laden

1. **Immer laden** — Glossar und Storyboard als Basis-Kontext lesen. Bei großem Storyboard nur die für das Thema relevanten Aktivitäten.
2. **Versuche (`docs/versuche/`)** — Nur einbeziehen, wenn der Entwickler es explizit angibt.

Diskussion bleibt konzeptionell: Begriffe schärfen, Verhalten definieren, Anforderungen (funktional und nicht-funktional) klären.

## Ablauf

1. **Thema erfassen** — Der Entwickler bringt den Änderungswunsch ein. Bei Unklarheiten gezielt nachfragen.
2. **Abgleich mit Artefakten** — Geladene Artefakte (Glossar, Storyboard) gegen die eingebrachte Idee prüfen. Widersprüche und Überschneidungen direkt benennen und in die Diskussion einbringen.
3. **Diskutieren** — Intensiv hinterfragen. Für jede Frage eine eigene Empfehlung geben. Entscheidungen dem Entwickler vorlegen. Fakten, die durch Nachschlagen im Projekt klärbar sind, selbst recherchieren statt fragen. Abhängige Entscheidungen schrittweise abarbeiten. Unabhängige Fragen dürfen gebatcht werden (maximal 2–3), um Token zu sparen.
4. **Artefakt-Anpassungen ableiten** — Wenn sich aus der Diskussion Änderungen ergeben, diese benennen:
   - **Glossar** → `glossar-anpassen` aufrufen
   - **Storyboard** → `storyboard-anpassen` aufrufen
   - **Versuche** → Unausgereifte Ideen oder offene Fragen in `docs/versuche/` ablegen
5. **Abschluss** — Zusammenfassen, was entschieden wurde und welche Artefakte angepasst wurden. Wenn sich technische Fragen ergeben haben, diese als offene Punkte für `planen` benennen.

## Regeln

- **Nicht eigenmächtig handeln.** Artefakte erst anpassen, wenn der Entwickler die Richtung bestätigt.
- **Keine technischen Entscheidungen (ADRs, Architektur).** Die Diskussion klärt das Was und Warum. Das Wie folgt in `planen`.
- **Versuche nicht ungefragt laden.** Versuche nur einbeziehen, wenn explizit gewünscht.