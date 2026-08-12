// Compose section — mounted under /compose. The legacy app nests everything
// under a redirecting `root`; here we wrap the children in a /compose layout
// route (ComposeHost) that provides the user/record stores + compose overlays.
// Absolute child paths are prefixed with /compose; relative children (the
// namespace.view sub-routes) keep mirroring the legacy paths.
import ComposeHost from './ComposeHost.vue'
import { composeRoutes } from './routes'
import ComposeSidebar from './sidebar/ComposeSidebar.vue'
import { useModuleStore, useNamespaceStore, useRecordStore } from '@planetcrust/human-vue'
import { usePageStore } from '@planetcrust/human-vue'
import { useReminderStore } from './stores/reminder'
import { useRightSidebarStore } from '@planetcrust/human-vue'

const PREFIX = '/compose'

function prefixAbsolute(routes) {
  return routes.map(route => {
    const out = { ...route, meta: { ...(route.meta || {}), section: 'compose' } }
    if (typeof route.path === 'string' && route.path.startsWith('/')) {
      out.path = PREFIX + route.path
    }
    if (typeof route.redirect === 'string' && route.redirect.startsWith('/')) {
      out.redirect = PREFIX + route.redirect
    }
    if (route.children) out.children = prefixAbsolute(route.children)
    return out
  })
}

const legacyRoot = composeRoutes.find(route => route.name === 'root')

// Resolves namespace/page/module/record entities for the AI agent context,
// layered on top of the shell's base route context.
function agentContext(ctx) {
  const params = ctx.routeParams || {}

  if (params.slug) {
    const ns = useNamespaceStore().getByUrlPart?.(params.slug)
    if (ns) ctx.namespace = { namespaceID: String(ns.namespaceID), slug: ns.slug, name: ns.name }
  }
  if (params.pageID) {
    const page = usePageStore().getByID?.(params.pageID)
    if (page) {
      ctx.page = {
        pageID: String(page.pageID),
        handle: page.handle,
        title: page.title,
        moduleID: page.moduleID ? String(page.moduleID) : undefined,
      }
    }
  }
  if (params.moduleID) {
    const mod = useModuleStore().getByID?.(params.moduleID)
    if (mod) ctx.module = { moduleID: String(mod.moduleID), handle: mod.handle, name: mod.name }
  }
  if (params.recordID) {
    const recordStore = useRecordStore()
    const record =
      recordStore.records?.get?.(params.recordID) || recordStore.labelCache?.get?.(params.recordID)
    if (record) {
      ctx.record = {
        recordID: String(record.recordID),
        moduleID: record.moduleID ? String(record.moduleID) : undefined,
      }
      if (record.values) {
        ctx.record.values = Array.isArray(record.values)
          ? record.values.reduce((acc, v) => {
              if (v && v.name) acc[v.name] = v.value
              return acc
            }, {})
          : { ...record.values }
      }
    }
  }

  return ctx
}

// Reminders entry in the topbar profile menu (t passed in from the shell).
function profileItems(t) {
  const reminderStore = useReminderStore()
  const label =
    reminderStore.activeCount > 0
      ? `${t('navigation.userSettings.reminders')} (${reminderStore.activeCount})`
      : t('navigation.userSettings.reminders')
  return [
    {
      label,
      icon: 'pi pi-clock',
      command: () => useRightSidebarStore().toggle('reminders'),
    },
  ]
}

export default {
  id: 'compose',
  routes: [
    {
      // Layout route is intentionally unnamed; the index child carries the
      // `compose` name (Vue Router can't render an unnamed empty-path child
      // when navigating to a named parent).
      path: PREFIX,
      component: ComposeHost,
      meta: { section: 'compose' },
      children: [
        { path: '', name: 'compose', redirect: { name: 'namespace.list' } },
        ...prefixAbsolute(legacyRoot.children),
      ],
    },
  ],
  sidebar: ComposeSidebar,
  sidebarExpandedByDefault: true,
  agentContext,
  profileItems,
}
