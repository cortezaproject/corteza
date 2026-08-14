// FRIA Determination — the AI Act deployer-category questions that decide
// whether a Fundamental Rights Impact Assessment is required (Art. 27).
//
// The first step of the Govern tab's FRIA flow, not part of project creation
// (see config/pipeline.js's 'fria-determination' entry and
// components/wizard/steps/FriaDeterminationStep.vue). Labels keep the
// project.deployer.* i18n keys in locale/en/human-webapp/project.yaml rather
// than moving to fria.yaml: they are the copy for the backend fields this feeds
// (ProjectDeployerCategories, FriaRequired).
export const DEPLOYER_QUESTIONS = [
  { key: 'publicAuthorityAnnex3', labelKey: 'project.deployer.questions.publicAuthorityAnnex3' },
  {
    key: 'privateEssentialServices',
    labelKey: 'project.deployer.questions.privateEssentialServices',
  },
  { key: 'insuranceBanking', labelKey: 'project.deployer.questions.insuranceBanking' },
]

export function friaDeterminationDefaults() {
  return Object.fromEntries(DEPLOYER_QUESTIONS.map(q => [q.key, false]))
}

// Merge saved governance values over the defaults, same shape-merge idiom as
// config/summaryForm.js's summaryDefaults() / config/resourceManagementForm.js's
// resourceManagementValues().
export function friaDeterminationValues(saved = {}) {
  return { ...friaDeterminationDefaults(), ...saved }
}

// Whether any deployer-category answer is "yes" — the session-local signal
// shown back to the member as the working FRIA-required outcome, mirroring
// (without replacing) the separate backend FriaRequired derivation this same
// data feeds once a backend slice picks it up.
export function friaRequiredFrom(values = {}) {
  return DEPLOYER_QUESTIONS.some(q => !!values[q.key])
}
