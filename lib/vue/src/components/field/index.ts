// Dispatcher components
export { default as CFieldEditor } from './CFieldEditor.vue'
export { default as CFieldViewer } from './CFieldViewer.vue'

// Unified registry
export { FIELD_REGISTRY, resolveFieldEditor, resolveFieldViewer } from './registry'

// Editor components
export { default as CFieldStringEditor } from './editors/CFieldStringEditor.vue'
export { default as CFieldNumberEditor } from './editors/CFieldNumberEditor.vue'
export { default as CFieldBoolEditor } from './editors/CFieldBoolEditor.vue'
export { default as CFieldSelectEditor } from './editors/CFieldSelectEditor.vue'
export { default as CFieldDateTimeEditor } from './editors/CFieldDateTimeEditor.vue'

// Viewer components
export { default as CFieldStringViewer } from './viewers/CFieldStringViewer.vue'
export { default as CFieldNumberViewer } from './viewers/CFieldNumberViewer.vue'
export { default as CFieldBoolViewer } from './viewers/CFieldBoolViewer.vue'
export { default as CFieldDateTimeViewer } from './viewers/CFieldDateTimeViewer.vue'
export { default as CFieldEmailViewer } from './viewers/CFieldEmailViewer.vue'
export { default as CFieldUrlViewer } from './viewers/CFieldUrlViewer.vue'
export { default as CFieldSelectViewer } from './viewers/CFieldSelectViewer.vue'

// Utilities
export { trimUrlFragment, trimUrlQuery, trimUrlPath, onlySecureUrl } from './url'
