// Resource Management governance step (EU AI Act Art. 17(l) — resource
// management, including security-of-supply). The step is custom-rendered (LLM
// catalogue inputs, connection whitelist cards), so unlike Project Summary
// there is no GovernanceForm schema — just the values shape.
//
// Values persist in the generic governance step values map
// (/governance/resource-management); the BE's ProjectConfig.ResourceManagement
// modeling exists but stays unused until it's confirmed.

export const resourceManagementDefaults = () => ({
  ai: {
    // Permitted providers (with alternatives for continuity).
    // [{ id, provider, model }]
    providers: [],
  },
  infra: {
    // Permitted infrastructure providers, each with its own continuity options.
    // [{ id, selfHosting, name, specification, regions, failover, failoverRegion }]
    providers: [],
    // Backup/Restore is a project-wide capability, not per provider.
    backupRestore: false,
  },
  // Whitelist of permitted external connections the Connections step will
  // later instantiate from.
  connections: [],
})

// Merge saved governance values over the defaults (one level deep — the
// nested ai/infra groups must not lose keys the saved copy doesn't carry).
// List entries are copied so editing the working copy never mutates the
// cached governance values.
export function resourceManagementValues(saved = {}) {
  const d = resourceManagementDefaults()
  return {
    ai: { ...d.ai, providers: (saved.ai?.providers || []).map(p => ({ ...p })) },
    infra: {
      ...d.infra,
      backupRestore: saved.infra?.backupRestore ?? d.infra.backupRestore,
      providers: (saved.infra?.providers || []).map(p => ({ ...p })),
    },
    connections: (saved.connections || []).map(c => ({ ...c })),
  }
}
