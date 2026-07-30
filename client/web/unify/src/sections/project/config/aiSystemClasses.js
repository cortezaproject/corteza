// Fixed regulatory catalog for the project section's AI systems — the EU AI
// Act (Regulation (EU) 2024/1689) Art. 6 risk classification an AI system
// carries, and the Annex III high-risk use-case points a high-risk system is
// classified under. Same pattern as config/friaTaxonomies.js, config/kinds.js
// and config/eventKinds.js: plain FRONTEND CONFIG, not a backend lookup
// table. There is no seed data and nothing here is tenant-specific — every
// project on every tenant sees the same risk classes and the same eight Annex
// III points, because the source is EU law, not something an org customises.
//
// CRITICAL — `key` IS A PERSISTED-DATA INTERFACE, NOT A LABEL. The moment an
// AI system stores its classification (a risk class, an Annex III point),
// that system's data holds this `key` string. Renaming a key later silently
// orphans every AI system that selected it — the stored value stops resolving
// to any entry and the UI can no longer show what the system was actually
// classified as, which for a compliance record is worse than showing nothing.
// Treat every `key` below as append-only: safe to add new ones, unsafe to
// rename or reuse an old key for a different meaning. If an entry is ever
// retired, leave its key mapped to a deprecated/legacy label rather than
// deleting it.
//
// i18n: every `labelKey`/`descriptionKey` resolves against the merged
// human-webapp bundle, in locale/en/human-webapp/fria.yaml (see that file's
// own header for how it's picked up). FRIA/AI-Act copy lives in fria.yaml,
// not project.yaml, per the 2026-07-28 ruling.
//
// COLOURS ARE DELIBERATELY OUT OF SCOPE HERE. A four-tier severity ladder
// (prohibited → high → limited → minimal) is exactly the kind of thing a UI
// wants to tint, but config/chartColors.js requires every palette used for
// data-viz to pass the dataviz skill's validator in BOTH light and dark mode.
// Picking those hexes belongs with the UI work that actually renders these
// (chips, badges, filters) — see config/friaScenario.js's severity classes for
// how that was handled there — not with this data-shape config. Do not add
// color/hex/tone fields to these entries without running that validation.
//
// ACCURACY NOTES — WHAT THE ACT ACTUALLY SAYS.
//   - THE FOUR TIERS ARE NOT FOUR LEGAL CATEGORIES. Only "prohibited"
//     (Art. 5) and "high-risk" (Art. 6 + Annexes I/III) are defined as such in
//     the Regulation's text. "Limited risk" and "minimal risk" are the
//     Commission's and the field's shorthand for, respectively, systems that
//     carry only the Art. 50 transparency obligations and systems that carry
//     no AI-Act-specific obligations at all. They are kept here because the
//     four-tier ladder is how practitioners classify systems and how this UI
//     asks the question — not because the Act enumerates them.
//   - ART. 50 IS AN OVERLAY, NOT A TIER. A high-risk system that also talks to
//     people (a chatbot, a generator of synthetic media) owes the Art. 50
//     transparency duties *on top of* its Art. 6 obligations. Classifying a
//     system as `high` therefore does not mean Art. 50 is inapplicable; the
//     `limited` tier means transparency duties are the *only* ones that apply.
//   - GPAI IS A SEPARATE AXIS. General-purpose AI model obligations (Chapter V,
//     Art. 51-56, including systemic-risk models) do not sit anywhere on this
//     ladder — a GPAI model is regulated as a model, independently of how the
//     system built on it classifies under Art. 6. Not modelled here.
//   - `friaRequired` IS A TIER-LEVEL PRECONDITION, NOT THE FULL TEST. Art. 27
//     attaches the FRIA obligation to *deployers of high-risk systems*, and
//     then only to certain ones: bodies governed by public law, private
//     entities providing public services, and deployers of the Annex III
//     point 5(b) creditworthiness and 5(c) life/health-insurance use cases.
//     `friaRequired: true` on `high` says "this tier is the one that can
//     trigger Art. 27"; whether it actually does for a given project is what
//     config/friaDeterminationForm.js's deployer-category questions decide.
//   - The eight Annex III points below are the FINAL adopted text of Annex III
//     (Regulation (EU) 2024/1689), not the 2021 Commission proposal — the two
//     differ (e.g. the proposal's separate crime-analytics point under law
//     enforcement is gone, and emergency-call triage was added under essential
//     services). Each description enumerates that point's own sub-points as
//     adopted.

// The four Art. 6 risk classes, in DESCENDING severity — prohibited first, so
// a UI that renders this array in order reads worst-to-least like the rest of
// this section (see config/friaScenario.js's severity ladder). `friaRequired`
// marks the single tier that can trigger the EU AI Act Art. 27 FRIA
// obligation; see the accuracy note above for why that is a precondition
// rather than the whole test.
export const AI_ACT_RISK_CLASSES = [
  {
    key: 'prohibited',
    labelKey: 'fria.riskClasses.prohibited.label',
    descriptionKey: 'fria.riskClasses.prohibited.description',
    friaRequired: false,
  },
  {
    key: 'high',
    labelKey: 'fria.riskClasses.high.label',
    descriptionKey: 'fria.riskClasses.high.description',
    friaRequired: true,
  },
  {
    key: 'limited',
    labelKey: 'fria.riskClasses.limited.label',
    descriptionKey: 'fria.riskClasses.limited.description',
    friaRequired: false,
  },
  {
    key: 'minimal',
    labelKey: 'fria.riskClasses.minimal.label',
    descriptionKey: 'fria.riskClasses.minimal.description',
    friaRequired: false,
  },
]

// The eight Annex III high-risk use-case areas, referred to by Art. 6(2). An
// AI system classified `high` is high-risk either because it falls under one
// of these points, or because it is a safety component of (or is itself) a
// product covered by the Annex I harmonisation legislation under Art. 6(1) —
// this catalog covers the Annex III route only.
//
// `point` carries the Annex III numbering separately from the translatable
// label, for the same reason RIGHTS_CHAPTERS carries `charterTitle` in
// config/friaTaxonomies.js: UI can compose "5 · Essential private and public
// services and benefits" without baking the numbering into a string a
// translator can reorder or drop. Keys are `annexIII1`…`annexIII8` and are
// permanent — the Act may amend the *content* of a point (Art. 7 empowers the
// Commission to do exactly that), but the numbering is stable, so a stored
// `annexIII5` keeps resolving even if the point's wording is later updated.
export const ANNEX_III_POINTS = [
  {
    key: 'annexIII1',
    point: 1,
    labelKey: 'fria.annexIII.annexIII1.label',
    descriptionKey: 'fria.annexIII.annexIII1.description',
  },
  {
    key: 'annexIII2',
    point: 2,
    labelKey: 'fria.annexIII.annexIII2.label',
    descriptionKey: 'fria.annexIII.annexIII2.description',
  },
  {
    key: 'annexIII3',
    point: 3,
    labelKey: 'fria.annexIII.annexIII3.label',
    descriptionKey: 'fria.annexIII.annexIII3.description',
  },
  {
    key: 'annexIII4',
    point: 4,
    labelKey: 'fria.annexIII.annexIII4.label',
    descriptionKey: 'fria.annexIII.annexIII4.description',
  },
  {
    key: 'annexIII5',
    point: 5,
    labelKey: 'fria.annexIII.annexIII5.label',
    descriptionKey: 'fria.annexIII.annexIII5.description',
  },
  {
    key: 'annexIII6',
    point: 6,
    labelKey: 'fria.annexIII.annexIII6.label',
    descriptionKey: 'fria.annexIII.annexIII6.description',
  },
  {
    key: 'annexIII7',
    point: 7,
    labelKey: 'fria.annexIII.annexIII7.label',
    descriptionKey: 'fria.annexIII.annexIII7.description',
  },
  {
    key: 'annexIII8',
    point: 8,
    labelKey: 'fria.annexIII.annexIII8.label',
    descriptionKey: 'fria.annexIII.annexIII8.description',
  },
]

// Whether a risk class can trigger the EU AI Act Art. 27 FRIA obligation —
// the one derivation callers should use instead of comparing against the
// literal 'high' string, so the tier→obligation mapping stays in this file.
// Same `.some()` idiom as config/friaDeterminationForm.js's friaRequiredFrom().
// An unknown, empty, or not-yet-classified key answers false: an unclassified
// system has not yet established that it is high-risk.
export function friaRequiredForRiskClass(key) {
  return AI_ACT_RISK_CLASSES.some(riskClass => riskClass.key === key && riskClass.friaRequired)
}
