---
name: architektur-doc-anpassen
description: Architekturdokumentation aktualisieren. Wird nach ADR-Änderungen (neu, geändert oder abgelöst) aufgerufen.
disable-model-invocation: true
---

## Auftrag

Aktualisiere die Architekturdokumentation (`docs/architektur.md`). Sie bildet den Jetzt-Zustand der Architektur ab und leitet sich aus den ADRs ab.

## Kontext laden

1. **Architekturdokumentation laden** — `docs/architektur.md` lesen. Falls die Datei noch nicht existiert, `ARCHITEKTURDOKUMENTATION-VORLAGE.md` aus dem Skill-Ordner als Grundlage verwenden.
2. **ADRs laden** — Alle aktiven (nicht abgelösten) ADRs aus `docs/adr/` als Grundlage.
3. **Glossar laden** — `docs/glossar.md` für korrekte Begriffe.

## Ablauf

1. **Auslöser klären** — Betroffene ADRs identifizieren (neu, geändert oder abgelöst).
2. **Abgleich** — Architekturdokumentation gegen ADRs prüfen. Veraltete oder widersprüchliche Stellen identifizieren.
3. **Dokumentation anpassen** — Änderungen einarbeiten. Format und Struktur des bestehenden Dokuments beibehalten.
4. **Zusammenfassung** — Kurz angeben, was aktualisiert wurde.

## Regeln

- **Jetzt-Zustand abbilden.** Die Architekturdokumentation beschreibt, wie die Architektur aktuell ist — keine Historie, keine Roadmap.
- **ADRs sind die einzige Quelle.** Inhalte aus ADRs ableiten, nicht frei erfinden.
- **Alle Kapitel beibehalten.** Nicht zutreffende Kapitel mit „Entfällt" kennzeichnen, nicht löschen.
- **Format beibehalten.** Struktur und Stil des bestehenden Dokuments übernehmen.
- **Nicht eigenmächtig erweitern.** Nur die vom Auslöser betroffenen Abschnitte ändern.