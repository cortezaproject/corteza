# Compose — TODO

Before merging, follow the general checklist in [docs/TODO.md](../../../docs/TODO.md).

## Vue 3 Migration — In Progress

Compose is being ported from Vue 2. The admin CRUD is functional; the public page viewer renders pages with Content blocks.

### Done

- [x] Namespace CRUD (list, create, edit, clone, delete, manage)
- [x] Module CRUD (list, create, edit)
- [x] Page CRUD (list, create, edit)
- [x] Chart CRUD (list, create, edit)
- [x] Record list and create (per module)
- [x] Namespace sidebar with navigation
- [x] All Pinia stores (namespace, module, page, chart, user, page-layout)
- [x] Public page viewer — renders pages with blocks in CSS grid layout
- [x] Page tree navigation in sidebar (visible pages as expandable tree)
- [x] Home page redirect (auto-redirect to first visible page, onboarding fallback)
- [x] Block rendering system with registry pattern
- [x] Content block type

### In Progress / TODO

- [ ] Multi-layout system — visibility expressions, role-based layout selection
- [ ] Record pages — pages with `moduleID`, record loading, RecordEditor block
- [ ] Additional block types — RecordList, Chart, Calendar, File, IFrame, Metric, etc.
- [ ] Block visibility expressions — per-block show/hide based on expressions and roles
- [ ] Expression evaluation engine — shared utility for layout/block/title expressions
- [ ] Page builder — drag-drop block editing, block configurators
- [ ] Magnification modal — full-screen block view
- [ ] Record detail/edit view
- [ ] Record inline editing in list
- [ ] Module field configuration UI
- [ ] Namespace export/import
- [ ] RBAC permission management UI
- [ ] Automation/workflow triggers on compose resources
- [ ] Chart rendering in page builder
- [ ] File/attachment field support
- [ ] Rich text field support
- [ ] Related records / reference fields UI
