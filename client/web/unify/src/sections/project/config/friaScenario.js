// Session-local FRIA risk-scenario data shape + small display helpers shared
// by the scenario list step and the five section-editor steps (see
// components/wizard/steps/Fria*.vue and components/wizard/fria/*). Scenario
// STORAGE itself lives in stores/projects.js governanceByProject, under the
// well-known 'fria-scenarios' step key's values.scenarios array (see that
// file's "FRIA risk scenarios" section) — this file only shapes a fresh
// scenario and derives display-only values; it holds no state of its own.

// A fresh, empty scenario — the shape every FRIA section-editor step reads a
// slice of and patches via stores/projects.js updateFriaScenario. Field
// groups mirror the design mockup's five numbered sections 1:1:
//   1 title/description/severity   (FriaHarmStep)
//   2 triggerTypes/triggerDescription   (FriaTriggerStep)
//   3 impactedParties/vulnerableGroups/vulnerableGroupsNotes  (FriaPartiesStep)
//   4 rights   (FriaRightsStep)
//   5 harmVectors/harmVectorsDescription   (FriaVectorsStep)
export function newFriaScenario() {
  return {
    id: `fria-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`,
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

// Required-field completion, mirroring the design mockup's own progress
// calculation 1:1 (its prog() function): title, description, severity, >=1
// trigger type, >=1 impacted party, >=1 right. Vulnerable groups and harm
// vectors are supporting detail in the mockup's own logic, not required for
// "complete" — kept identical here so the summary list's Complete / In
// Progress status matches what the section-editor steps' own required
// markers (*) promise.
export function friaScenarioCompletion(scenario) {
  const checks = [
    !!scenario?.title?.trim(),
    !!scenario?.description?.trim(),
    !!scenario?.severity,
    (scenario?.triggerTypes?.length || 0) > 0,
    (scenario?.impactedParties?.length || 0) > 0,
    (scenario?.rights?.length || 0) > 0,
  ]
  const filled = checks.filter(Boolean).length
  return { filled, total: checks.length, complete: filled === checks.length }
}
