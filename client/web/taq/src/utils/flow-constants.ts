import type { IconDef } from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import { DEFAULT_ICONS } from '@cortezaproject/corteza-js-next/src/automation/types/icon'

// Default icon fallbacks (re-exported from shared lib)
export const DEFAULT_TRIGGER_ICON: IconDef = DEFAULT_ICONS.TRIGGER
export const DEFAULT_ACTION_ICON: IconDef = DEFAULT_ICONS.ACTION

// Trigger metadata map: eventType -> { icon, i18nKey }
// Icons are IconDef objects with type 'name' for PrimeIcons
// i18nKey maps to builder.nodePicker.nodes.triggers.{i18nKey}.label/description
export const TRIGGER_META: Record<string, { icon: IconDef; i18nKey: string }> = {
  onManual: { icon: { type: 'name', value: 'play' }, i18nKey: 'manual' },
  onInterval: { icon: { type: 'name', value: 'sync' }, i18nKey: 'interval' },
  onTimestamp: { icon: { type: 'name', value: 'clock' }, i18nKey: 'timestamp' },
  afterCreate: { icon: { type: 'name', value: 'plus-circle' }, i18nKey: 'recordCreate' },
  afterUpdate: { icon: { type: 'name', value: 'pencil' }, i18nKey: 'recordUpdate' },
  afterDelete: { icon: { type: 'name', value: 'trash' }, i18nKey: 'recordDelete' },
}

export const NODE_DIMENSIONS = {
  // Default node dimensions
  WIDTH: 260,
  HEIGHT: 100,

  // Spacing
  HORIZONTAL_SEP: 200, // Between sibling branches
  VERTICAL_SEP: 100, // Between parent-child ranks
}

// Helper to get center offset for viewport centering
export function getNodeCenterOffset() {
  const width = NODE_DIMENSIONS.WIDTH
  const height = NODE_DIMENSIONS.HEIGHT

  return {
    x: width / 2,
    y: height / 2,
  }
}
