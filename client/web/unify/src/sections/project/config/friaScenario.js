// FRIA risk-scenario data shape + small display helpers shared by the
// scenario list and the single-page scenario editor's five section components
// (see components/wizard/steps/FriaScenariosStep.vue and
// components/wizard/fria/*). This file only shapes a fresh/cloned scenario and
// derives display-only values; it holds no state of its own.
//
// Scenario STORAGE is the backend type ProjectFriaScenario, cached in
// stores/projects.js's "FRIA risk scenarios" section; the shape below is the
// flat form that section maps the API to and from.

// A fresh, empty scenario — the shape the editor seeds a new draft from (see
// FriaScenarioEditor.vue) and every section component reads a slice of and
// patches via its `update` prop. Field groups mirror the design mockup's
// five numbered sections 1:1:
//   1 title/description/severity   (FriaHarmSection)
//   2 triggerTypes/triggerDescription   (FriaTriggerSection)
//   3 impactedParties/vulnerableGroups/vulnerableGroupsNotes  (FriaPartiesSection)
//   4 rights   (FriaRightsSection)
//   5 harmVectors/harmVectorsDescription   (FriaVectorsSection)
export function newFriaScenario() {
  return {
    id: `fria-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`,
    // The AI system this scenario assesses. REQUIRED — a scenario cannot be
    // saved without one: Art. 27 attaches the FRIA to a specific high-risk AI
    // SYSTEM, not to a project, and a project can hold several. It is also what
    // gives detection rules something to scope to. Holds a ProjectAiSystem ID;
    // see the ai-systems Govern step.
    aiSystemID: null,
    title: '',
    description: '',
    severity: null, // 'low' | 'medium' | 'high' | 'critical'
    triggerTypes: [], // config/friaTaxonomies TRIGGER_CONDITIONS keys
    triggerDescription: '',
    impactedParties: [], // config/friaTaxonomies IMPACTED_PARTIES keys
    vulnerableGroups: [], // config/friaTaxonomies VULNERABLE_GROUPS keys
    vulnerableGroupsNotes: '',
    rights: [], // config/friaTaxonomies FUNDAMENTAL_RIGHTS keys
    harmVectors: [], // config/friaTaxonomies AI_HARM_VECTORS keys
    harmVectorsDescription: '',
  }
}

// A standalone copy of a scenario — every array field copied too, so editing
// the clone (the local draft — see components/wizard/fria/
// FriaScenarioEditor.vue) can never mutate the store's version in place
// before an explicit Save commits it. Field list mirrors newFriaScenario()
// above 1:1.
export function cloneFriaScenario(scenario) {
  return {
    ...scenario,
    triggerTypes: [...(scenario.triggerTypes || [])],
    impactedParties: [...(scenario.impactedParties || [])],
    vulnerableGroups: [...(scenario.vulnerableGroups || [])],
    rights: [...(scenario.rights || [])],
    harmVectors: [...(scenario.harmVectors || [])],
  }
}

export const FRIA_SEVERITY_LEVELS = ['low', 'medium', 'high', 'critical']

// Tailwind classes reused 1:1 from the same worst→least warm family already
// used across this section (EventBadge's SEVERITY tints, RiskPips' pip
// colours — see components/dashboard/{EventBadge,RiskPips}.vue): red →
// orange → amber → yellow. FRIA's four-value Low/Medium/High/Critical scale
// doesn't line up with EventBadge's own five-value label set (Informational/
// Minor/Major/Serious/Critical), so the classes are duplicated here rather
// than forcing a label remap through that component — but no new hexes: this
// is Tailwind's own predefined palette, the same families already validated
// and in use for this exact "worst is red, least is yellow" concept, so
// config/chartColors.js's dataviz-validator rule doesn't come into play.
export const FRIA_SEVERITY_BADGE_CLASSES = {
  critical: 'bg-red-100 text-red-700 dark:bg-red-500/15 dark:text-red-300',
  high: 'bg-orange-100 text-orange-700 dark:bg-orange-500/15 dark:text-orange-300',
  medium: 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300',
  low: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-500/15 dark:text-yellow-300',
}
// Same families, used for the severity-card picker's selected border/fill
// (mirrors the mockup's per-severity selection colour).
export const FRIA_SEVERITY_CARD_CLASSES = {
  critical: 'border-red-500 bg-red-50 dark:bg-red-500/10',
  high: 'border-orange-500 bg-orange-50 dark:bg-orange-500/10',
  medium: 'border-amber-500 bg-amber-50 dark:bg-amber-500/10',
  low: 'border-yellow-500 bg-yellow-50 dark:bg-yellow-500/10',
}

// Decorative emoji per taxonomy key, transcribed 1:1 from the design mockup's
// trigger/vector cards. Purely presentational (not a colour or a persisted
// value) — kept here rather than in config/friaTaxonomies.js, which is
// explicitly locked/out of scope for this pass and carries no icon fields by
// design (see that file's header comment).
export const TRIGGER_ICONS = {
  dataQualityIssue: '📊',
  modelDrift: '📉',
  misuseRepurposing: '⚠️',
  automationBias: '🤖',
  edgeCaseOodInput: '🔍',
  adversarialInput: '🎯',
  integrationFailure: '🔗',
  bypassedSafeguard: '🛡️',
}
export const AI_HARM_VECTOR_ICONS = {
  scaleAutomation: '⚡',
  opacityInexplicability: '🎭',
  historicalBiasEncoding: '📜',
  misplacedEpistemicAuthority: '🏛️',
  proxyVariableUse: '🔮',
  feedbackLoopSelfFulfillingProphecy: '🔄',
  automationBiasInOversight: '👤',
  dataAggregationReIdentification: '🗃️',
  systemicEcosystemEffects: '🌐',
  lockInDependency: '🔐',
}

// Charter chapter -> PrimeIcon. Stands in for the design mockup's per-chapter
// colour dot: chartColors.js has no validated 7-hue "chapter" categorical
// palette, and this pass is explicitly not the place to invent one (see
// config/chartColors.js's standing dataviz-validator rule) — a non-colour
// affordance keeps the 7 chapters visually distinguishable without picking
// new hexes. A validated chapter palette is still a follow-up worth doing.
export const RIGHTS_CHAPTER_ICONS = {
  dignity: 'pi-heart-fill',
  freedoms: 'pi-unlock',
  equality: 'pi-equals',
  solidarity: 'pi-users',
  citizens: 'pi-building-columns',
  justice: 'pi-verified',
  dataPrivacy: 'pi-lock',
}

// Whether a scenario may be saved at all. The AI system is the one field the
// editor blocks on: everything else can be filled in over time, but a scenario
// that assesses nothing in particular is not a partial FRIA, it is a category
// error. Separate from completion below on purpose — completion measures how
// much of the assessment is written, this measures whether it is addressable.
export function friaScenarioSaveable(scenario) {
  return !!scenario?.aiSystemID
}

// Completion = how many of the editor's FIVE sections carry any content. Both
// vulnerable groups and harm vectors count, so a scenario cannot read
// "complete" having identified no AI harm vector — the very thing Art. 27 asks
// about. One definition, used by both the summary list's Complete / In Progress
// status and the editor's progress bar, so the two can never disagree.
//
// The AI system ref is deliberately NOT a sixth check: a scenario cannot exist
// without one (see friaScenarioSaveable), so counting it would only ever add a
// point every scenario already has.
export function friaScenarioCompletion(scenario) {
  const checks = [
    // 1. Harm narrative + severity
    !!scenario?.title?.trim() || !!scenario?.description?.trim() || !!scenario?.severity,
    // 2. Trigger conditions
    (scenario?.triggerTypes?.length || 0) > 0 || !!scenario?.triggerDescription?.trim(),
    // 3. Impacted parties + vulnerable groups
    (scenario?.impactedParties?.length || 0) > 0 ||
      (scenario?.vulnerableGroups?.length || 0) > 0 ||
      !!scenario?.vulnerableGroupsNotes?.trim(),
    // 4. Fundamental rights
    (scenario?.rights?.length || 0) > 0,
    // 5. AI harm vectors
    (scenario?.harmVectors?.length || 0) > 0 || !!scenario?.harmVectorsDescription?.trim(),
  ]
  const filled = checks.filter(Boolean).length
  return { filled, total: checks.length, complete: filled === checks.length }
}
