/**
 * @cortezaproject/corteza-js-next
 *
 * Core JavaScript/TypeScript library for Corteza
 * Provides API clients, type definitions, utilities, and models
 */

// Core utilities
export { NoID } from './cast'

// Event bus for client-side event handling
export * as eventbus from './eventbus'

// Corredor automation scripting support
export * as corredor from './corredor'

// Validation utilities
export * as validator from './validator/validator'

// Compose module (records, modules, namespaces, pages, charts, etc.)
export * as compose from './compose'

// System module (users, roles, applications, settings, etc.)
export * as system from './system'

// Reporter module (reports, data sources, etc.)
export * as reporter from './reporter'

// Automation module (workflows, triggers, etc.)
export * as automation from './automation'

// Shared utilities and types
export * as shared from './shared'

// API Clients for backend communication
export * as apiClients from './api-clients'

// Formatting utilities (dates, numbers, etc.)
export * as fmt from './formatting'
