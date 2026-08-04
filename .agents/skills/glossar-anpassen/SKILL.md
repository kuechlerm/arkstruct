---
name: glossar-anpassen
description: Glossar des Projekts anpassen — Begriffe hinzufügen, ändern oder entfernen. Wird aus der Diskussion heraus aufgerufen, wenn sich Begriffe oder Konzepte ändern.
---

## Auftrag

Passe das Glossar des Projekts an. Das Glossar ist das zentrale Werkzeug für einheitliche Begriffe und Konzepte. Es verhindert inhaltliche Konflikte und sorgt für Klarheit — für Entwickler und KI.

## Ablauf

1. **Glossar laden** — Lies die Glossar-Datei unter `docs/glossar.md` des Projekts. Wenn sie nicht existiert, erstelle sie, wenn der erste Begriff hinzugefügt werden soll.
2. **Änderungen verstehen** — Der Entwickler benennt die gewünschte Änderung (neuer Begriff, Umbenennung, Entfernung, Präzisierung). Falls unklar, frage gezielt nach. Wenn der Entwickler einen vagen oder mehrdeutigen Begriff verwendet, einen präzisen kanonischen Begriff vorschlagen, bevor er ins Glossar übernommen wird.
3. **Prüfen gegen Code** — Prüfe per Suche im Code, ob der Begriff bereits verwendet wird und ob die Änderung konsistent ist. Code ist die Quelle der Wahrheit.
4. **Konflikte aufdecken** — Prüfe, ob die Änderung bestehenden Einträgen widerspricht oder Doppelungen erzeugt. Wenn ja, den Widerspruch direkt benennen und den Entwickler zur Entscheidung zwingen, bevor weitergearbeitet wird.
5. **Mit Szenarien prüfen** — Bei neuen oder geänderten Begriffen ein konkretes Szenario durchspielen, das die Abgrenzung zu verwandten Begriffen testet. Besonders bei Begriffen, die Beziehungen zwischen Konzepten beschreiben.
6. **Glossar anpassen** — Führe die Änderung durch. Halte dich an das bestehende Format und die Sortierung der Datei. Wenn mehrere Begriffe in einer Diskussion entstehen, jeden einzeln sofort eintragen, nicht am Ende gesammelt.
7. **Beziehungen gegenprüfen** — Wenn ein Begriff eine `_Beziehungen_`-Angabe bekommt, ändert oder verliert, prüfe den referenzierten Begriff und halte den Gegeneintrag konsistent (anlegen, anpassen oder entfernen).
8. **Zusammenfassung** — Gib kurz an, was geändert wurde und ob Folgeaktionen nötig sind (z.B. Umbenennungen im Code, Anpassung von Storyboard oder Architekturdoku).

## Regeln

- **Keine eigenmächtigen Ergänzungen.** Nur das ändern, was der Entwickler vorgibt oder was sich direkt aus der laufenden Diskussion ergibt.
- **Bestehende Begriffe nicht still umbenennen.** Umbenennungen immer explizit bestätigen lassen, da sie Auswirkungen auf Code und andere Artefakte haben.
- **Format beibehalten.** Struktur, Sortierung und Stil des bestehenden Glossars übernehmen. Einträge folgen diesem Schema:
  ```
  **Begriff**:
  Ein bis zwei Sätze, die beschreiben was der Begriff IST — nicht was er tut.
  _Vermeiden_: Synonym1, Synonym2
  _Beziehungen_: → viele AndererBegriff
  ```
  Die `_Vermeiden_`-Liste enthält bewusst abgelehnte Synonyme oder verwandte Begriffe, die nicht verwendet werden sollen. Die `_Beziehungen_`-Zeile ist optional und beschreibt Kardinalität zu anderen Begriffen (z.B. `→ viele Projekt`, `→ ein Kunde`). Sie wird beidseitig eingetragen: Steht bei „Kunde" eine Beziehung zu „Projekt", bekommt auch „Projekt" den Gegeneintrag zu „Kunde".
- **Keine Implementierungsdetails.** Das Glossar beschreibt Domänenbegriffe, keine technischen Konzepte. Keine Verhaltensbeschreibungen — nur Definitionen. Technische Entscheidungen gehören in ADRs oder Architekturdoku. `_Beziehungen_` beschreibt nur Kardinalität/Struktur (hat, gehört zu) — kein Verhalten (löst aus, erzeugt, verarbeitet).
- **Token-schonend arbeiten.** Nur relevante Code-Stellen laden, nicht den gesamten Code durchsuchen.