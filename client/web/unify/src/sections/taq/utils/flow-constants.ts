import type { IconDef } from '@planetcrust/human-js/src/automation/types/icon'
import { DEFAULT_ICONS } from '@planetcrust/human-js/src/automation/types/icon'

// Default icon fallbacks (re-exported from shared lib)
export const DEFAULT_TRIGGER_ICON: IconDef = DEFAULT_ICONS.TRIGGER
export const DEFAULT_ACTION_ICON: IconDef = DEFAULT_ICONS.ACTION

// Trigger metadata map: eventType -> { icon, i18nKey }
// Icons are IconDef objects with type 'name' for PrimeIcons
// Event to Icon Mapping
export const EVENT_ICONS: Record<string, IconDef> = {
  // General events
  onManual: { type: 'name', value: 'play' },
  onInterval: { type: 'name', value: 'sync' },
  onTimestamp: { type: 'name', value: 'clock' },

  // Record lifecycle events
  beforeCreate: { type: 'name', value: 'plus-circle' },
  afterCreate: { type: 'name', value: 'plus-circle' },
  beforeUpdate: { type: 'name', value: 'pencil' },
  afterUpdate: { type: 'name', value: 'pencil' },
  beforeDelete: { type: 'name', value: 'trash' },
  afterDelete: { type: 'name', value: 'trash' },
  beforeSuspend: { type: 'name', value: 'ban' },
  afterSuspend: { type: 'name', value: 'ban' },
}

// Helper to resolve trigger metadata taking resource type into context
export function getTriggerMeta(eventType: string, _resourceType?: string): { icon: IconDef } {
  // Icon maps based on the eventType
  return {
    icon: EVENT_ICONS[eventType] || DEFAULT_TRIGGER_ICON,
  }
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
