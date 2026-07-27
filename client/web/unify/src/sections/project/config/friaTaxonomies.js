// Fixed regulatory taxonomies for the project section's Fundamental Rights
// Impact Assessment flow (FRIA, EU AI Act Art. 27 — Govern tab, risk
// scenarios). Same pattern as config/kinds.js, config/categories.js and
// config/eventKinds.js: plain FRONTEND CONFIG, not a backend lookup table.
// There is no seed data and nothing here is tenant-specific — every project
// on every tenant sees the same rights, groups, vectors, triggers and
// parties, because the source is EU law and EU AI Act guidance, not
// something an org customises.
//
// CRITICAL — `key` IS A PERSISTED-DATA INTERFACE, NOT A LABEL. The moment a
// risk scenario stores a selection (a right, a vulnerable group, a harm
// vector, a trigger condition, an impacted party), that scenario's data
// holds this `key` string. Renaming a key later silently orphans every
// scenario that selected it — the stored value stops resolving to any entry
// and the UI can no longer show what was actually assessed. Treat every
// `key` below as append-only: safe to add new ones, unsafe to rename or
// reuse an old key for a different meaning. If an entry is ever retired,
// leave its key mapped to a deprecated/legacy label rather than deleting it.
//
// i18n: every `labelKey`/`descriptionKey`/`sublabelKey` resolves against the
// merged human-webapp bundle, in locale/en/human-webapp/fria.yaml (see that
// file's own header for how it's picked up).
//
// COLOURS ARE DELIBERATELY OUT OF SCOPE HERE. The design mockup this was
// transcribed from colour-codes fundamental rights by Charter chapter, but
// config/chartColors.js requires every palette used for data-viz to pass the
// dataviz skill's validator in BOTH light and dark mode. Picking those hexes
// belongs with the UI work that actually renders these taxonomies (chips,
// legends, filters), not with this data-shape config — do not add
// color/hex/tone fields to these entries without running that validation.
//
// ACCURACY NOTE — RIGHTS vs. THE ACTUAL EU CHARTER. The fundamental-rights
// list below was transcribed from the design mockup and then checked against
// the real Charter of Fundamental Rights of the European Union (2000/C
// 364/01, as adapted 2007/C 303/01), which has SIX numbered titles (I
// Dignity, II Freedoms, III Equality, IV Solidarity, V Citizens' Rights, VI
// Justice; a seventh, "General provisions", governs interpretation and has
// no rights of its own). Divergences found and how they were handled:
//   - MISSING WELL-KNOWN RIGHT: Article 20 "Equality before the law" — the
//     opening, most foundational article of Title III — was absent from the
//     mockup, which started the equality group at Art. 21 (non-discrimination).
//     Added back as `equalityBeforeTheLaw`.
//   - WRONG ARTICLE CITATION: the mockup cited "Charter Art. 27–28" for
//     "Workers' Rights & Right to Information", double-counting Art. 28
//     (right of collective bargaining and action), which is already its own
//     separate entry (`collectiveBargainingAction`) with its own correct
//     citation. Fixed to cite Art. 27 alone.
//   - CHAPTER GROUPING vs. THE CHARTER'S SIX TITLES: this taxonomy keeps the
//     mockup's 7th grouping, `dataPrivacy`, which is NOT one of the Charter's
//     six numbered titles — it is a practical FRIA-authoring bucket (the
//     mockup's own UI marks it distinctly, with no roman numeral, unlike the
//     six real titles). It holds one genuine Charter right whose true title
//     is II Freedoms — `protectionOfPersonalData` (Art. 8) — alongside three
//     rights that are not Charter articles at all but GDPR/AI Act rights
//     central to AI-Act FRIA practice (`rightNotSubjectAutomatedDecisions`,
//     `rightToExplanation`, `rightsAccessRectificationErasure`). This is a
//     deliberate, documented choice to keep data-protection/AI-Act rights
//     browsable as one group rather than a factual claim that `dataPrivacy`
//     is a seventh Charter title — see RIGHTS_CHAPTERS below, where its
//     `charterTitle` is `null` for exactly this reason.
//   - Every other citation was checked article-by-article against the
//     Charter text and matches.
// This taxonomy is intentionally a curated shortlist of the Charter articles
// most relevant to AI systems: 40 of the Charter's 50 substantive articles
// (1-50 across Titles I-VI) are represented; the 10 that are not are Art. 9
// (marry & found a family), 13 (arts & sciences), 19 (removal/expulsion/
// extradition), 29 (placement services), 30 (unjustified dismissal), 32
// (child labour), 33 (family & professional life), 36 (services of general
// economic interest), 45 (freedom of movement/residence) and 46 (diplomatic/
// consular protection). That curation itself is unchanged from the mockup —
// only factual errors in what it did include were corrected.

// The 7 groupings rights are filed under. Six mirror the Charter's own
// numbered titles 1:1 (`charterTitle` carries the roman numeral so UI can
// compose "I · Dignity" without baking numbering into the translatable
// label); `dataPrivacy` is the practical 7th bucket described above and
// carries no Charter numeral.
export const RIGHTS_CHAPTERS = [
  { key: 'dignity', charterTitle: 'I', labelKey: 'fria.rightsChapters.dignity' },
  { key: 'freedoms', charterTitle: 'II', labelKey: 'fria.rightsChapters.freedoms' },
  { key: 'equality', charterTitle: 'III', labelKey: 'fria.rightsChapters.equality' },
  { key: 'solidarity', charterTitle: 'IV', labelKey: 'fria.rightsChapters.solidarity' },
  { key: 'citizens', charterTitle: 'V', labelKey: 'fria.rightsChapters.citizens' },
  { key: 'justice', charterTitle: 'VI', labelKey: 'fria.rightsChapters.justice' },
  { key: 'dataPrivacy', charterTitle: null, labelKey: 'fria.rightsChapters.dataPrivacy' },
]

// Fundamental rights a risk scenario may infringe, limit, or put at risk.
// `chapter` is one of RIGHTS_CHAPTERS' keys above; `article` is a citation
// string (Charter / ECHR / GDPR / AI Act as applicable — see the accuracy
// note above for the three dataPrivacy entries that carry no Charter
// article at all). Order within each chapter follows Charter article order.
export const FUNDAMENTAL_RIGHTS = [
  // I — Dignity (Charter Art. 1-5)
  {
    key: 'humanDignity',
    chapter: 'dignity',
    article: 'Charter Art. 1 · ECHR Preamble',
    labelKey: 'fria.rights.humanDignity.label',
    descriptionKey: 'fria.rights.humanDignity.description',
  },
  {
    key: 'rightToLife',
    chapter: 'dignity',
    article: 'Charter Art. 2 · ECHR Art. 2',
    labelKey: 'fria.rights.rightToLife.label',
    descriptionKey: 'fria.rights.rightToLife.description',
  },
  {
    key: 'integrityOfThePerson',
    chapter: 'dignity',
    article: 'Charter Art. 3',
    labelKey: 'fria.rights.integrityOfThePerson.label',
    descriptionKey: 'fria.rights.integrityOfThePerson.description',
  },
  {
    key: 'prohibitionTortureInhumanTreatment',
    chapter: 'dignity',
    article: 'Charter Art. 4 · ECHR Art. 3',
    labelKey: 'fria.rights.prohibitionTortureInhumanTreatment.label',
    descriptionKey: 'fria.rights.prohibitionTortureInhumanTreatment.description',
  },
  {
    key: 'prohibitionSlaveryForcedLabour',
    chapter: 'dignity',
    article: 'Charter Art. 5 · ECHR Art. 4',
    labelKey: 'fria.rights.prohibitionSlaveryForcedLabour.label',
    descriptionKey: 'fria.rights.prohibitionSlaveryForcedLabour.description',
  },

  // II — Freedoms (Charter Art. 6-19; Art. 8 Data Protection and Art. 9
  // Marry & Found a Family are out of scope here — Art. 8 is filed under the
  // dataPrivacy group below per the accuracy note, Art. 9 has no clear AI
  // risk-scenario relevance and was not in the mockup)
  {
    key: 'rightToLibertySecurity',
    chapter: 'freedoms',
    article: 'Charter Art. 6 · ECHR Art. 5',
    labelKey: 'fria.rights.rightToLibertySecurity.label',
    descriptionKey: 'fria.rights.rightToLibertySecurity.description',
  },
  {
    key: 'respectPrivateFamilyLife',
    chapter: 'freedoms',
    article: 'Charter Art. 7 · ECHR Art. 8',
    labelKey: 'fria.rights.respectPrivateFamilyLife.label',
    descriptionKey: 'fria.rights.respectPrivateFamilyLife.description',
  },
  {
    key: 'freedomThoughtConscienceReligion',
    chapter: 'freedoms',
    article: 'Charter Art. 10 · ECHR Art. 9',
    labelKey: 'fria.rights.freedomThoughtConscienceReligion.label',
    descriptionKey: 'fria.rights.freedomThoughtConscienceReligion.description',
  },
  {
    key: 'freedomExpressionInformation',
    chapter: 'freedoms',
    article: 'Charter Art. 11 · ECHR Art. 10',
    labelKey: 'fria.rights.freedomExpressionInformation.label',
    descriptionKey: 'fria.rights.freedomExpressionInformation.description',
  },
  {
    key: 'freedomAssemblyAssociation',
    chapter: 'freedoms',
    article: 'Charter Art. 12 · ECHR Art. 11',
    labelKey: 'fria.rights.freedomAssemblyAssociation.label',
    descriptionKey: 'fria.rights.freedomAssemblyAssociation.description',
  },
  {
    key: 'rightToEducation',
    chapter: 'freedoms',
    article: 'Charter Art. 14',
    labelKey: 'fria.rights.rightToEducation.label',
    descriptionKey: 'fria.rights.rightToEducation.description',
  },
  {
    key: 'freedomChooseOccupation',
    chapter: 'freedoms',
    article: 'Charter Art. 15',
    labelKey: 'fria.rights.freedomChooseOccupation.label',
    descriptionKey: 'fria.rights.freedomChooseOccupation.description',
  },
  {
    key: 'freedomConductBusiness',
    chapter: 'freedoms',
    article: 'Charter Art. 16',
    labelKey: 'fria.rights.freedomConductBusiness.label',
    descriptionKey: 'fria.rights.freedomConductBusiness.description',
  },
  {
    key: 'rightToProperty',
    chapter: 'freedoms',
    article: 'Charter Art. 17 · ECHR Protocol 1 Art. 1',
    labelKey: 'fria.rights.rightToProperty.label',
    descriptionKey: 'fria.rights.rightToProperty.description',
  },
  {
    key: 'rightToAsylum',
    chapter: 'freedoms',
    article: 'Charter Art. 18',
    labelKey: 'fria.rights.rightToAsylum.label',
    descriptionKey: 'fria.rights.rightToAsylum.description',
  },

  // III — Equality (Charter Art. 20-26)
  {
    // Added during the Charter cross-check — see the accuracy note above.
    key: 'equalityBeforeTheLaw',
    chapter: 'equality',
    article: 'Charter Art. 20',
    labelKey: 'fria.rights.equalityBeforeTheLaw.label',
    descriptionKey: 'fria.rights.equalityBeforeTheLaw.description',
  },
  {
    key: 'nonDiscrimination',
    chapter: 'equality',
    article: 'Charter Art. 21 · ECHR Art. 14 · Equality Directives',
    labelKey: 'fria.rights.nonDiscrimination.label',
    descriptionKey: 'fria.rights.nonDiscrimination.description',
  },
  {
    key: 'culturalReligiousLinguisticDiversity',
    chapter: 'equality',
    article: 'Charter Art. 22',
    labelKey: 'fria.rights.culturalReligiousLinguisticDiversity.label',
    descriptionKey: 'fria.rights.culturalReligiousLinguisticDiversity.description',
  },
  {
    key: 'equalityWomenMen',
    chapter: 'equality',
    article: 'Charter Art. 23',
    labelKey: 'fria.rights.equalityWomenMen.label',
    descriptionKey: 'fria.rights.equalityWomenMen.description',
  },
  {
    key: 'rightsOfTheChild',
    chapter: 'equality',
    article: 'Charter Art. 24 · UNCRC · DSA Art. 28',
    labelKey: 'fria.rights.rightsOfTheChild.label',
    descriptionKey: 'fria.rights.rightsOfTheChild.description',
  },
  {
    key: 'rightsOfTheElderly',
    chapter: 'equality',
    article: 'Charter Art. 25',
    labelKey: 'fria.rights.rightsOfTheElderly.label',
    descriptionKey: 'fria.rights.rightsOfTheElderly.description',
  },
  {
    key: 'integrationPersonsDisabilities',
    chapter: 'equality',
    article: 'Charter Art. 26 · CRPD',
    labelKey: 'fria.rights.integrationPersonsDisabilities.label',
    descriptionKey: 'fria.rights.integrationPersonsDisabilities.description',
  },

  // IV — Solidarity (Charter Art. 27-38)
  {
    key: 'workersRightInformation',
    chapter: 'solidarity',
    // Fixed during the Charter cross-check: the mockup cited "Art. 27-28",
    // double-counting Art. 28, which is `collectiveBargainingAction` below.
    article: 'Charter Art. 27',
    labelKey: 'fria.rights.workersRightInformation.label',
    descriptionKey: 'fria.rights.workersRightInformation.description',
  },
  {
    key: 'collectiveBargainingAction',
    chapter: 'solidarity',
    article: 'Charter Art. 28',
    labelKey: 'fria.rights.collectiveBargainingAction.label',
    descriptionKey: 'fria.rights.collectiveBargainingAction.description',
  },
  {
    key: 'fairJustWorkingConditions',
    chapter: 'solidarity',
    article: 'Charter Art. 31',
    labelKey: 'fria.rights.fairJustWorkingConditions.label',
    descriptionKey: 'fria.rights.fairJustWorkingConditions.description',
  },
  {
    key: 'socialSecurityAssistance',
    chapter: 'solidarity',
    article: 'Charter Art. 34',
    labelKey: 'fria.rights.socialSecurityAssistance.label',
    descriptionKey: 'fria.rights.socialSecurityAssistance.description',
  },
  {
    key: 'rightToHealthcare',
    chapter: 'solidarity',
    article: 'Charter Art. 35',
    labelKey: 'fria.rights.rightToHealthcare.label',
    descriptionKey: 'fria.rights.rightToHealthcare.description',
  },
  {
    key: 'environmentalProtection',
    chapter: 'solidarity',
    article: 'Charter Art. 37',
    labelKey: 'fria.rights.environmentalProtection.label',
    descriptionKey: 'fria.rights.environmentalProtection.description',
  },
  {
    key: 'consumerProtection',
    chapter: 'solidarity',
    article: 'Charter Art. 38 · EU Consumer Law',
    labelKey: 'fria.rights.consumerProtection.label',
    descriptionKey: 'fria.rights.consumerProtection.description',
  },

  // V — Citizens' Rights (Charter Art. 39-44)
  {
    key: 'rightToVoteStandElection',
    chapter: 'citizens',
    article: 'Charter Art. 39-40 · ECHR Protocol 1 Art. 3',
    labelKey: 'fria.rights.rightToVoteStandElection.label',
    descriptionKey: 'fria.rights.rightToVoteStandElection.description',
  },
  {
    key: 'rightToGoodAdministration',
    chapter: 'citizens',
    article: 'Charter Art. 41',
    labelKey: 'fria.rights.rightToGoodAdministration.label',
    descriptionKey: 'fria.rights.rightToGoodAdministration.description',
  },
  {
    key: 'accessToDocumentsTransparency',
    chapter: 'citizens',
    article: 'Charter Art. 42 · EU Transparency Reg.',
    labelKey: 'fria.rights.accessToDocumentsTransparency.label',
    descriptionKey: 'fria.rights.accessToDocumentsTransparency.description',
  },
  {
    key: 'rightOfPetitionOmbudsman',
    chapter: 'citizens',
    article: 'Charter Art. 43-44',
    labelKey: 'fria.rights.rightOfPetitionOmbudsman.label',
    descriptionKey: 'fria.rights.rightOfPetitionOmbudsman.description',
  },

  // VI — Justice (Charter Art. 47-50)
  {
    key: 'rightToEffectiveRemedyFairTrial',
    chapter: 'justice',
    article: 'Charter Art. 47 · ECHR Art. 6',
    labelKey: 'fria.rights.rightToEffectiveRemedyFairTrial.label',
    descriptionKey: 'fria.rights.rightToEffectiveRemedyFairTrial.description',
  },
  {
    key: 'presumptionInnocenceRightDefence',
    chapter: 'justice',
    article: 'Charter Art. 48 · ECHR Art. 6(2-3)',
    labelKey: 'fria.rights.presumptionInnocenceRightDefence.label',
    descriptionKey: 'fria.rights.presumptionInnocenceRightDefence.description',
  },
  {
    key: 'legalityProportionalityPenalties',
    chapter: 'justice',
    article: 'Charter Art. 49 · ECHR Art. 7',
    labelKey: 'fria.rights.legalityProportionalityPenalties.label',
    descriptionKey: 'fria.rights.legalityProportionalityPenalties.description',
  },
  {
    key: 'neBisInIdem',
    chapter: 'justice',
    article: 'Charter Art. 50 · ECHR Protocol 7 Art. 4',
    labelKey: 'fria.rights.neBisInIdem.label',
    descriptionKey: 'fria.rights.neBisInIdem.description',
  },

  // Data & Privacy (practical grouping, not a Charter title — see the
  // accuracy note above)
  {
    key: 'protectionOfPersonalData',
    chapter: 'dataPrivacy',
    // True Charter title is II Freedoms (Art. 8) — filed here deliberately,
    // see the accuracy note above.
    article: 'Charter Art. 8 · GDPR · AI Act Art. 10',
    labelKey: 'fria.rights.protectionOfPersonalData.label',
    descriptionKey: 'fria.rights.protectionOfPersonalData.description',
  },
  {
    key: 'rightNotSubjectAutomatedDecisions',
    chapter: 'dataPrivacy',
    // Not a Charter article — GDPR/AI Act right.
    article: 'GDPR Art. 22 · AI Act Art. 14 & 86',
    labelKey: 'fria.rights.rightNotSubjectAutomatedDecisions.label',
    descriptionKey: 'fria.rights.rightNotSubjectAutomatedDecisions.description',
  },
  {
    key: 'rightToExplanation',
    chapter: 'dataPrivacy',
    // Not a Charter article — GDPR/AI Act right.
    article: 'GDPR Art. 13-14 · AI Act Art. 86',
    labelKey: 'fria.rights.rightToExplanation.label',
    descriptionKey: 'fria.rights.rightToExplanation.description',
  },
  {
    key: 'rightsAccessRectificationErasure',
    chapter: 'dataPrivacy',
    // Not a Charter article — GDPR right.
    article: 'GDPR Art. 15-17',
    labelKey: 'fria.rights.rightsAccessRectificationErasure.label',
    descriptionKey: 'fria.rights.rightsAccessRectificationErasure.description',
  },
]

// Vulnerable groups that require heightened analysis under EU AI Act Art. 27
// and Recital 58. `sublabelKey` carries the short qualifier text (age
// threshold, legal basis, framing) shown under the main label.
export const VULNERABLE_GROUPS = [
  {
    key: 'childrenMinors',
    labelKey: 'fria.vulnerableGroups.childrenMinors.label',
    sublabelKey: 'fria.vulnerableGroups.childrenMinors.sublabel',
  },
  {
    key: 'olderAdultsElderly',
    labelKey: 'fria.vulnerableGroups.olderAdultsElderly.label',
    sublabelKey: 'fria.vulnerableGroups.olderAdultsElderly.sublabel',
  },
  {
    key: 'personsWithDisabilities',
    labelKey: 'fria.vulnerableGroups.personsWithDisabilities.label',
    sublabelKey: 'fria.vulnerableGroups.personsWithDisabilities.sublabel',
  },
  {
    key: 'ethnicRacialMinorities',
    labelKey: 'fria.vulnerableGroups.ethnicRacialMinorities.label',
    sublabelKey: 'fria.vulnerableGroups.ethnicRacialMinorities.sublabel',
  },
  {
    key: 'religiousMinorities',
    labelKey: 'fria.vulnerableGroups.religiousMinorities.label',
    sublabelKey: 'fria.vulnerableGroups.religiousMinorities.sublabel',
  },
  {
    key: 'womenGenderMinorities',
    labelKey: 'fria.vulnerableGroups.womenGenderMinorities.label',
    sublabelKey: 'fria.vulnerableGroups.womenGenderMinorities.sublabel',
  },
  {
    key: 'lgbtqiaPersons',
    labelKey: 'fria.vulnerableGroups.lgbtqiaPersons.label',
    sublabelKey: 'fria.vulnerableGroups.lgbtqiaPersons.sublabel',
  },
  {
    key: 'migrantsAsylumSeekers',
    labelKey: 'fria.vulnerableGroups.migrantsAsylumSeekers.label',
    sublabelKey: 'fria.vulnerableGroups.migrantsAsylumSeekers.sublabel',
  },
  {
    key: 'lowIncomeSocioeconomicDisadvantage',
    labelKey: 'fria.vulnerableGroups.lowIncomeSocioeconomicDisadvantage.label',
    sublabelKey: 'fria.vulnerableGroups.lowIncomeSocioeconomicDisadvantage.sublabel',
  },
  {
    key: 'neurodivergentPersons',
    labelKey: 'fria.vulnerableGroups.neurodivergentPersons.label',
    sublabelKey: 'fria.vulnerableGroups.neurodivergentPersons.sublabel',
  },
  {
    key: 'languageMinorities',
    labelKey: 'fria.vulnerableGroups.languageMinorities.label',
    sublabelKey: 'fria.vulnerableGroups.languageMinorities.sublabel',
  },
  {
    key: 'personsWithMentalHealthConditions',
    labelKey: 'fria.vulnerableGroups.personsWithMentalHealthConditions.label',
    sublabelKey: 'fria.vulnerableGroups.personsWithMentalHealthConditions.sublabel',
  },
  {
    key: 'unhousedHomelessPersons',
    labelKey: 'fria.vulnerableGroups.unhousedHomelessPersons.label',
    sublabelKey: 'fria.vulnerableGroups.unhousedHomelessPersons.sublabel',
  },
  {
    key: 'digitallyExcludedLowLiteracy',
    labelKey: 'fria.vulnerableGroups.digitallyExcludedLowLiteracy.label',
    sublabelKey: 'fria.vulnerableGroups.digitallyExcludedLowLiteracy.sublabel',
  },
  {
    key: 'prisonersPersonsInDetention',
    labelKey: 'fria.vulnerableGroups.prisonersPersonsInDetention.label',
    sublabelKey: 'fria.vulnerableGroups.prisonersPersonsInDetention.sublabel',
  },
  {
    key: 'victimsDomesticGenderBasedViolence',
    labelKey: 'fria.vulnerableGroups.victimsDomesticGenderBasedViolence.label',
    sublabelKey: 'fria.vulnerableGroups.victimsDomesticGenderBasedViolence.sublabel',
  },
]

// AI-specific mechanisms by which the AI system enables or amplifies a harm
// beyond what an equivalent human process would produce (section "Where
// does AI Enable or Amplify the Harm?" in the mockup).
export const AI_HARM_VECTORS = [
  {
    key: 'scaleAutomation',
    labelKey: 'fria.aiHarmVectors.scaleAutomation.label',
    descriptionKey: 'fria.aiHarmVectors.scaleAutomation.description',
  },
  {
    key: 'opacityInexplicability',
    labelKey: 'fria.aiHarmVectors.opacityInexplicability.label',
    descriptionKey: 'fria.aiHarmVectors.opacityInexplicability.description',
  },
  {
    key: 'historicalBiasEncoding',
    labelKey: 'fria.aiHarmVectors.historicalBiasEncoding.label',
    descriptionKey: 'fria.aiHarmVectors.historicalBiasEncoding.description',
  },
  {
    key: 'misplacedEpistemicAuthority',
    labelKey: 'fria.aiHarmVectors.misplacedEpistemicAuthority.label',
    descriptionKey: 'fria.aiHarmVectors.misplacedEpistemicAuthority.description',
  },
  {
    key: 'proxyVariableUse',
    labelKey: 'fria.aiHarmVectors.proxyVariableUse.label',
    descriptionKey: 'fria.aiHarmVectors.proxyVariableUse.description',
  },
  {
    key: 'feedbackLoopSelfFulfillingProphecy',
    labelKey: 'fria.aiHarmVectors.feedbackLoopSelfFulfillingProphecy.label',
    descriptionKey: 'fria.aiHarmVectors.feedbackLoopSelfFulfillingProphecy.description',
  },
  {
    key: 'automationBiasInOversight',
    labelKey: 'fria.aiHarmVectors.automationBiasInOversight.label',
    descriptionKey: 'fria.aiHarmVectors.automationBiasInOversight.description',
  },
  {
    key: 'dataAggregationReIdentification',
    labelKey: 'fria.aiHarmVectors.dataAggregationReIdentification.label',
    descriptionKey: 'fria.aiHarmVectors.dataAggregationReIdentification.description',
  },
  {
    key: 'systemicEcosystemEffects',
    labelKey: 'fria.aiHarmVectors.systemicEcosystemEffects.label',
    descriptionKey: 'fria.aiHarmVectors.systemicEcosystemEffects.description',
  },
  {
    key: 'lockInDependency',
    labelKey: 'fria.aiHarmVectors.lockInDependency.label',
    descriptionKey: 'fria.aiHarmVectors.lockInDependency.description',
  },
]

// Technical/operational failure modes that trigger a harm scenario (section
// "Trigger" in the mockup) — what condition or event causes the harm, as
// distinct from the AI-specific vectors above that explain why AI makes it
// worse once triggered.
export const TRIGGER_CONDITIONS = [
  {
    key: 'dataQualityIssue',
    labelKey: 'fria.triggerConditions.dataQualityIssue.label',
    descriptionKey: 'fria.triggerConditions.dataQualityIssue.description',
  },
  {
    key: 'modelDrift',
    labelKey: 'fria.triggerConditions.modelDrift.label',
    descriptionKey: 'fria.triggerConditions.modelDrift.description',
  },
  {
    key: 'misuseRepurposing',
    labelKey: 'fria.triggerConditions.misuseRepurposing.label',
    descriptionKey: 'fria.triggerConditions.misuseRepurposing.description',
  },
  {
    key: 'automationBias',
    labelKey: 'fria.triggerConditions.automationBias.label',
    descriptionKey: 'fria.triggerConditions.automationBias.description',
  },
  {
    key: 'edgeCaseOodInput',
    labelKey: 'fria.triggerConditions.edgeCaseOodInput.label',
    descriptionKey: 'fria.triggerConditions.edgeCaseOodInput.description',
  },
  {
    key: 'adversarialInput',
    labelKey: 'fria.triggerConditions.adversarialInput.label',
    descriptionKey: 'fria.triggerConditions.adversarialInput.description',
  },
  {
    key: 'integrationFailure',
    labelKey: 'fria.triggerConditions.integrationFailure.label',
    descriptionKey: 'fria.triggerConditions.integrationFailure.description',
  },
  {
    key: 'bypassedSafeguard',
    labelKey: 'fria.triggerConditions.bypassedSafeguard.label',
    descriptionKey: 'fria.triggerConditions.bypassedSafeguard.description',
  },
]

// Who is directly or indirectly affected by the scenario (section "Impacted
// Parties" in the mockup). Plain label-only entries — unlike the taxonomies
// above, the mockup carries no supporting description text for these.
export const IMPACTED_PARTIES = [
  { key: 'jobSeekersApplicants', labelKey: 'fria.impactedParties.jobSeekersApplicants' },
  { key: 'employeesWorkers', labelKey: 'fria.impactedParties.employeesWorkers' },
  { key: 'customersServiceUsers', labelKey: 'fria.impactedParties.customersServiceUsers' },
  { key: 'citizensGeneralPublic', labelKey: 'fria.impactedParties.citizensGeneralPublic' },
  {
    key: 'patientsHealthcareRecipients',
    labelKey: 'fria.impactedParties.patientsHealthcareRecipients',
  },
  { key: 'studentsLearners', labelKey: 'fria.impactedParties.studentsLearners' },
  { key: 'tenantsHousingApplicants', labelKey: 'fria.impactedParties.tenantsHousingApplicants' },
  { key: 'suspectsAccusedPersons', labelKey: 'fria.impactedParties.suspectsAccusedPersons' },
  { key: 'societyLabourMarket', labelKey: 'fria.impactedParties.societyLabourMarket' },
  { key: 'bystandersThirdParties', labelKey: 'fria.impactedParties.bystandersThirdParties' },
  {
    key: 'competitorsOtherBusinesses',
    labelKey: 'fria.impactedParties.competitorsOtherBusinesses',
  },
]
