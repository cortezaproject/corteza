export default {
  data () {
    return {
      c311Dirty: false,
      c311DirtyStorageKey: '',
    }
  },

  mounted () {
    window.addEventListener('beforeunload', this.c311BeforeUnload)
  },

  beforeDestroy () {
    window.removeEventListener('beforeunload', this.c311BeforeUnload)
  },

  beforeRouteLeave (to, from, next) {
    if (!this.c311Dirty) return next()
    const message = this.$t?.('c311:dirty.leave', 'You have unsaved changes. Leave this page?') || 'You have unsaved changes. Leave this page?'
    if (window.confirm(message)) {
      this.c311ClearDirtyDraft()
      if (typeof this.c311ReloadServerState === 'function') this.c311ReloadServerState(to, from)
      return next()
    }
    return next(false)
  },

  methods: {
    c311DraftStorage () {
      if (typeof window === 'undefined') return null
      const storage = this.c311DirtyStorageKey?.startsWith('c311.portal.') ? window.localStorage : window.sessionStorage
      return storage && typeof storage.getItem === 'function' ? storage : null
    },
    c311BeforeUnload (event) {
      if (!this.c311Dirty) return
      event.preventDefault()
      event.returnValue = ''
    },
    c311MarkDirty (value = true) {
      this.c311Dirty = value
    },
    c311SaveDirtyDraft (value) {
      const storage = this.c311DraftStorage()
      if (!this.c311DirtyStorageKey || !storage) return
      try {
        storage.setItem(this.c311DirtyStorageKey, JSON.stringify(sanitizeC311Draft(value)))
      } catch (_error) {}
    },
    c311ReadDirtyDraft () {
      const storage = this.c311DraftStorage()
      if (!this.c311DirtyStorageKey || !storage) return null
      try {
        const value = storage.getItem(this.c311DirtyStorageKey)
        return value ? JSON.parse(value) : null
      } catch (_error) {
        return null
      }
    },
    c311ClearDirtyDraft () {
      const storage = this.c311DraftStorage()
      if (this.c311DirtyStorageKey && storage) {
        try { storage.removeItem(this.c311DirtyStorageKey) } catch (_error) {}
      }
      this.c311Dirty = false
    },
  },
}

export function sanitizeC311Draft (value) {
  if (Array.isArray(value)) return value.map(sanitizeC311Draft)
  if (!value || typeof value !== 'object') return value

  return Object.entries(value).reduce((result, [key, item]) => {
    if (/(?:password|token|secret|credential)/i.test(key)) return result
    result[key] = sanitizeC311Draft(item)
    return result
  }, {})
}
