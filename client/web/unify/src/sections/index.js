// Section registry. Each section is a self-contained module that declares its
// routes (paths prefixed, names namespaced), an optional left-sidebar
// component, and per-section topbar overrides. Adding a section = drop a
// folder here and add it to this list.
import admin from './admin'
import agentic from './agentic'
import chatbot from './chatbot'
import compose from './compose'
import home from './home'
import project from './project'
import taq from './taq'
import workflow from './workflow'

export const sections = [home, agentic, workflow, taq, admin, compose, chatbot, project]

// Flattened route list for the merged router.
export const routes = sections.flatMap(section => section.routes)

export function sectionById(id) {
  return sections.find(section => section.id === id) || null
}

// Server-side locale application bundles to load + deep-merge for the active
// section set. `human-webapp-one` provides the shell chrome (navigation/app),
// the rest add each section's namespaces.
export const localeApplications = [
  'human-webapp-one',
  'human-webapp-home',
  'human-webapp-agentic',
  'human-webapp-workflow',
  'human-webapp-taq',
  'human-webapp-admin',
  'human-webapp-compose',
  'human-webapp-chatbot',
  'human-webapp-project',
]
