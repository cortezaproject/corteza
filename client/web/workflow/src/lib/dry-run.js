import { automation } from '@cortezaproject/corteza-js'

export async function encodeInput (initialScope, ComposeAPI, SystemAPI) {
  const ev = { args: {} }

  // Types that can and must be fetched
  const {
    namespace,
    oldNamespace,
    module,
    oldModule,
    page,
    oldPage,
    record,
    oldRecord,
    user,
    oldUser,
    role,
    oldRole,
    application,
    oldApplication,
  } = initialScope

  if (namespace && namespace.value) {
    await ComposeAPI.namespaceList({ slug: namespace.value })
      .then(({ set = [] }) => {
        const [ns] = set
        if (ns) {
          ev.args.namespace = { ...ns, resourceType: 'compose:namespace' }
        }
      })
      .catch(() => {})
  }

  if (oldNamespace && oldNamespace.value) {
    await ComposeAPI.namespaceList({ slug: oldNamespace.value })
      .then(({ set = [] }) => {
        const [ns] = set
        if (ns) {
          ev.args.oldNamespace = { ...ns, resourceType: 'compose:namespace' }
        }
      })
      .catch(() => {})
  }

  if (module && module.value && ev.args.namespace) {
    await ComposeAPI.moduleList({ namespaceID: ev.args.namespace.namespaceID, handle: module.value })
      .then(({ set = [] }) => {
        const [m] = set
        if (m) {
          ev.args.module = { ...m, resourceType: 'compose:module' }
        }
      })
      .catch(() => {})
  }

  if (oldModule && oldModule.value && ev.args.namespace) {
    await ComposeAPI.moduleList({ namespaceID: ev.args.namespace.namespaceID, handle: oldModule.value })
      .then(({ set = [] }) => {
        const [m] = set
        if (m) {
          ev.args.oldModule = { ...m, resourceType: 'compose:module' }
        }
      })
      .catch(() => {})
  }

  if (page && page.value && ev.args.namespace) {
    await ComposeAPI.pageRead({ pageID: page.value, namespaceID: ev.args.namespace.namespaceID })
      .then(p => {
        ev.args.page = { ...p }
      })
  }

  if (oldPage && oldPage.value && ev.args.namespace) {
    await ComposeAPI.pageRead({ pageID: page.value, namespaceID: ev.args.namespace.namespaceID })
      .then(p => {
        ev.args.oldPage = { ...p }
      })
  }

  if (record && record.value && ev.args.module && ev.args.namespace) {
    await ComposeAPI.recordRead({ recordID: record.value, moduleID: ev.args.module.moduleID, namespaceID: ev.args.namespace.namespaceID })
      .then(r => {
        ev.args.record = { ...r, resourceType: 'compose:record' }
      })
  }

  if (oldRecord && oldRecord.value && ev.args.module && ev.args.namespace) {
    await ComposeAPI.recordRead({ recordID: oldRecord.value, moduleID: ev.args.module.moduleID, namespaceID: ev.args.namespace.namespaceID })
      .then(r => {
        ev.args.oldRecord = { ...r, resourceType: 'compose:record' }
      })
  }

  if (user && user.value) {
    await SystemAPI.userRead({ userID: user.value })
      .then(u => {
        ev.args.user = { ...u, resourceType: 'User' }
      })
  }

  if (oldUser && oldUser.value) {
    await SystemAPI.oldUserRead({ userID: oldUser.value })
      .then(u => {
        ev.args.oldUser = { ...u, resourceType: 'User' }
      })
  }

  if (role && role.value) {
    await SystemAPI.roleRead({ roleID: role.value })
      .then(r => {
        ev.args.role = { ...r, resourceType: 'Role' }
      })
  }

  if (oldRole && oldRole.value) {
    await SystemAPI.roleRead({ roleID: oldRole.value })
      .then(r => {
        ev.args.oldRole = { ...r, resourceType: 'Role' }
      })
  }

  if (application && application.value) {
    await SystemAPI.applicationRead({ applicationID: application.value })
      .then(a => {
        ev.args.application = { ...a }
      })
  }

  if (oldApplication && oldApplication.value) {
    await SystemAPI.applicationRead({ applicationID: oldApplication.value })
      .then(a => {
        ev.args.oldApplication = { ...a }
      })
  }

  // Add rest to args
  Object.entries(initialScope).forEach(([key, value]) => {
    if (!ev.args[key]) {
      ev.args[key] = {}
    }
  })

  return automation.Encode(ev.args)
}
