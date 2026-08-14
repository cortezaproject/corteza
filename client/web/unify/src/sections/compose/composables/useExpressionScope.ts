import { computed, unref, type Ref } from 'vue'
import {
  buildScope,
  useModuleStore,
  type ScopeEntry,
  type ScopeModule,
} from '@planetcrust/human-vue'

type MaybeRef<T> = T | Ref<T>

interface Page {
  moduleID?: string
}

interface Source {
  // The page the expression renders on. `${record...}` reads the page record,
  // so the page's module is what those variables resolve against — on a
  // record-list block that is *not* the module the block queries.
  page?: MaybeRef<Page | null | undefined>
  // Module the filter queries. Its fields are bare identifiers in the QL, not
  // `${...}` variables, and are offered as such.
  queryModule?: MaybeRef<ScopeModule | null | undefined>
  // Overrides the page test, for values authored away from the page they
  // render on — a chart is namespace-level and may land on either kind.
  hasRecord?: MaybeRef<boolean>
}

// Scope and field list for a compose expression input.
export function useExpressionScope(source: Source = {}): {
  scope: Ref<ScopeEntry[]>
  queryFields: Ref<Array<{ name: string; label?: string; kind?: string }>>
  isRecordPage: Ref<boolean>
} {
  const moduleStore = useModuleStore()

  const pageModuleID = computed(() => {
    const id = unref(source.page)?.moduleID
    return id && id !== '0' ? id : ''
  })

  // Read from the page rather than from the loaded module, so the variable
  // list does not flip while that module is still being fetched.
  const isRecordPage = computed(() => unref(source.hasRecord) ?? !!pageModuleID.value)

  const scope = computed(() =>
    buildScope({
      recordModule: pageModuleID.value ? moduleStore.getByID(pageModuleID.value) || null : null,
      hasRecord: isRecordPage.value,
    }),
  )

  const queryFields = computed(() => unref(source.queryModule)?.fields || [])

  return { scope, queryFields, isRecordPage }
}
