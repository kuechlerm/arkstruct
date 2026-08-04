---
name: implementierung
description: Arbeitet einen Implementierungsplan schrittweise ab. Ändert Code, erstellt Tests und führt sie aus. Immer in eigener Session.
disable-model-invocation: true
---

## Auftrag

Arbeite den vorhandenen Implementierungsplan Schritt für Schritt ab. Ändere Code, erstelle Tests und führe sie aus.

## Kontext laden

1. **Implementierungsplan laden** — Der Plan enthält alle nötigen Informationen.
2. **Glossar laden** — Für korrekte Begriffe im Code.
3. **Code gezielt laden** — Nur die im Plan referenzierten Dateien.
4. **Referenzierte Artefakte bei Bedarf nachschlagen** — ADRs, Storyboard-Abschnitte oder Architekturdoku, die der Plan referenziert, dürfen zum Verständnis nachgeladen werden. Sie ändern nicht den Plan.

## Ablauf

1. **Plan lesen** — Implementierungsplan vollständig lesen und verstehen.
2. **Schritte abarbeiten** — Jeden Schritt einzeln umsetzen:
   - Code ändern
   - Zugehörige Tests erstellen oder anpassen
   - Die zum Schritt gehörenden Tests ausführen und sicherstellen, dass sie bestehen
   - Schritt im Plan abhaken
3. **Bei Problemen** — Triviale Anpassungen (fehlende Imports, Tippfehler im Plan, offensichtliche Kleinigkeiten) dürfen eigenständig gelöst werden. Bei Abweichungen von Architektur, Schnittstellen oder Scope: stoppen und den Entwickler informieren. Nicht eigenmächtig vom Plan abweichen.
4. **Abschluss** — Wenn alle Schritte abgearbeitet sind:
   - Alle Tests ausführen (nicht nur die neuen)
   - Implementierungsplan entfernen
   - Zusammenfassen, was umgesetzt wurde
   - Hinweis auf mögliche Aktualisierung von Produktanleitungen und Produktinformationen
   - Deployment-Bereitschaft bestätigen oder Blocker benennen

## Regeln

- **Eigene Session.** Implementierung findet immer in einer frischen Session statt, getrennt von Diskussion und Planung.
- **Plan ist verbindlich.** Nicht vom Plan abweichen, keine Features hinzufügen, keinen Code „verbessern", der nicht im Plan steht.
- **Tests sind Pflicht.** Kein Schritt ohne Validierung. Akzeptanztests (als Playwright-Tests umgesetzt) sind die bevorzugte Validierungsform. Unit-Tests ergänzen dort, wo es passt — insbesondere bei technischen Schritten ohne nutzersichtbare Änderung.
- **Plan entfernen.** Nach erfolgreicher Implementierung wird der Plan gelöscht — er ist kein Archiv-Artefakt.