// Module field types — 1:1 with Corteza Compose module field kinds
// (see lib/vue/src/components/field/registry.ts). `id` is the compose kind.
// `Record` references another module (relationship); `User` references a user.
// `labelKey` is an i18n key; components resolve it with $t for display.

// `labelKey` points at Compose's shared field-kind labels (general.fieldKinds.*)
// so the names stay identical to the rest of the app (e.g. Geometry → "Location").
// `icon` is a PrimeVue icon class and `hintKey` an i18n key for the one-line
// description shown in the visual type picker.
export const FIELD_TYPES = [
  {
    id: 'String',
    labelKey: 'general.fieldKinds.String.label',
    icon: 'pi pi-align-left',
    hintKey: 'project.fieldTypeHints.String',
  },
  {
    id: 'Number',
    labelKey: 'general.fieldKinds.Number.label',
    icon: 'pi pi-hashtag',
    hintKey: 'project.fieldTypeHints.Number',
  },
  {
    id: 'Bool',
    labelKey: 'general.fieldKinds.Bool.label',
    icon: 'pi pi-check-square',
    hintKey: 'project.fieldTypeHints.Bool',
  },
  {
    id: 'DateTime',
    labelKey: 'general.fieldKinds.DateTime.label',
    icon: 'pi pi-calendar',
    hintKey: 'project.fieldTypeHints.DateTime',
  },
  {
    id: 'Select',
    labelKey: 'general.fieldKinds.Select.label',
    icon: 'pi pi-list',
    hintKey: 'project.fieldTypeHints.Select',
  },
  {
    id: 'Email',
    labelKey: 'general.fieldKinds.Email.label',
    icon: 'pi pi-envelope',
    hintKey: 'project.fieldTypeHints.Email',
  },
  {
    id: 'Url',
    labelKey: 'general.fieldKinds.Url.label',
    icon: 'pi pi-link',
    hintKey: 'project.fieldTypeHints.Url',
  },
  {
    id: 'User',
    labelKey: 'general.fieldKinds.User.label',
    icon: 'pi pi-user',
    hintKey: 'project.fieldTypeHints.User',
  },
  {
    id: 'Record',
    labelKey: 'general.fieldKinds.Record.label',
    icon: 'pi pi-sitemap',
    hintKey: 'project.fieldTypeHints.Record',
  },
  {
    id: 'File',
    labelKey: 'general.fieldKinds.File.label',
    icon: 'pi pi-file',
    hintKey: 'project.fieldTypeHints.File',
  },
  {
    id: 'Geometry',
    labelKey: 'general.fieldKinds.Geometry.label',
    icon: 'pi pi-map-marker',
    hintKey: 'project.fieldTypeHints.Geometry',
  },
]

// Look up a single field type's display config by id (or null if unknown).
export const fieldType = id => FIELD_TYPES.find(t => t.id === id) || null

export const RECORD_KIND = 'Record'
export const isRecordRef = type => type === RECORD_KIND

// i18n key for a field type's label (or null for an unknown kind, so the
// caller can fall back to the raw id).
export const fieldTypeLabelKey = id => FIELD_TYPES.find(t => t.id === id)?.labelKey || null
