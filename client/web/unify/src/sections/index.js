// Section registry. Each section is a self-contained module that declares its
// routes (paths prefixed, names namespaced), an optional left-sidebar
// component, and per-section topbar overrides. Adding a section = drop a
// folder here and add it to this list.
import admin from './admin'
import agentic from './agentic'
import chatbot from './chatbot'
import compose from './compose'
import home from './home'
import taq from './taq'
import workflow from './workflow'

export const sections = [home, agentic, workflow, taq, admin, compose, chatbot]

// Flattened route list for the merged router.
export const routes = sections.flatMap(section => section.routes)

export function sectionById(id) {
  return sections.find(section => section.id === id) || null
}
