import { ref, type InjectionKey, type Ref, inject, provide } from 'vue'

export interface PermissionDialogOptions {
  /** The RBAC resource string, e.g. 'corteza::system:role/12345' or 'corteza::system:role/*' */
  resource: string
  /** Display name for the resource (e.g. role name) */
  title?: string
  /** Target name for specific-resource i18n labels */
  target?: string
  /** If true, use 'all-specific' i18n pattern instead of 'all/specific' */
  allSpecific?: boolean
}

export interface PermissionsContext {
  visible: Ref<boolean>
  options: Ref<PermissionDialogOptions | null>
  open: (_opts: PermissionDialogOptions) => void
  close: () => void
}

export const PermissionsKey: InjectionKey<PermissionsContext> = Symbol('permissions')

/**
 * Provides the permission-dialog state at the app level.
 * Call once in your root App.vue or layout.
 */
export function providePermissions(): PermissionsContext {
  const visible = ref(false)
  const options = ref<PermissionDialogOptions | null>(null)

  function open(opts: PermissionDialogOptions) {
    options.value = opts
    visible.value = true
  }

  function close() {
    visible.value = false
    options.value = null
  }

  const ctx: PermissionsContext = { visible, options, open, close }
  provide(PermissionsKey, ctx)
  return ctx
}

/**
 * Injects the permission-dialog context.
 * Use in any component that needs to open the permission dialog.
 */
export function usePermissions(): PermissionsContext {
  const ctx = inject(PermissionsKey)
  if (!ctx) {
    throw new Error('usePermissions() requires providePermissions() to be called in a parent component')
  }
  return ctx
}
