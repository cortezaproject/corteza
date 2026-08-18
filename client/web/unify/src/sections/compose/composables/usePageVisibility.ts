const NoID = '0'

type Mode = 'view' | 'edit' | 'create'

interface ExpressionVariablesOptions {
  record?: { serialize?: () => Record<string, unknown> } | null
  isRecordPage?: boolean
  mode?: Mode
}

interface VisibilityConfig {
  expression?: string
  roles?: string[]
}

interface Block {
  blockID?: string
  kind?: string
  meta?: {
    tempID?: string
    hidden?: boolean
    visibility?: VisibilityConfig
  }
  options?: {
    tabs?: Array<{ blockID?: string; title?: string }>
  }
  xywh?: number[]
}

interface Layout {
  pageLayoutID: string
  config?: {
    visibility?: VisibilityConfig
  }
}

function fetchBlockID(block: Block): string {
  const bid = block.blockID
  if (bid && bid !== NoID) return bid
  return block.meta?.tempID || ''
}

function getBreakpoint(): string {
  const w = window.innerWidth
  if (w >= 1200) return 'lg'
  if (w >= 996) return 'md'
  if (w >= 768) return 'sm'
  if (w >= 480) return 'xs'
  return 'xxs'
}

// A page no layout matches sends the viewer away, but there is nowhere inside a
// namespace guaranteed not to lead back: `pages` is the landing route, and it
// redirects to the namespace's home page — which may be that very page. So the
// attempt is made once. Arriving a second time, the page stands its ground and
// says why it is empty rather than bouncing forever.
let lastRefusedPageID = ''

function refuseOnce(pageID: string): boolean {
  if (lastRefusedPageID === pageID) return false
  lastRefusedPageID = pageID
  return true
}

// Any page that does resolve a layout ends the streak, so returning to a page
// that once refused is a fresh attempt rather than a silent empty screen.
function clearRefusal(): void {
  lastRefusedPageID = ''
}

export { fetchBlockID, getBreakpoint, refuseOnce, clearRefusal }

export function usePageVisibility(
  $SystemAPI: {
    expressionEvaluate: (_a: {
      variables: Record<string, unknown>
      expressions: Record<string, string>
    }) => Promise<Record<string, boolean>>
  } | null,
  $auth: { user?: { roles?: string[] } } | null,
) {
  function buildExpressionVariables(
    options: ExpressionVariablesOptions = {},
  ): Record<string, unknown> {
    const { record, isRecordPage, mode } = options
    const vars: Record<string, unknown> = {
      user: $auth?.user || {},
      record: record?.serialize ? record.serialize() : {},
      screen: {
        width: window.innerWidth,
        height: window.innerHeight,
        userAgent: navigator.userAgent,
        breakpoint: getBreakpoint(),
      },
    }

    if (isRecordPage && mode) {
      vars.isView = mode === 'view'
      vars.isCreate = mode === 'create'
      vars.isEdit = mode === 'edit'
    }

    return vars
  }

  async function determineLayout(
    layouts: Layout[],
    variables: Record<string, unknown>,
    requestedLayoutID?: string,
  ): Promise<Layout | null> {
    if (!layouts.length) return null

    // Batch evaluate all layouts that have expressions
    const layoutsWithExpressions = layouts.filter(l => l.config?.visibility?.expression)
    let expressionResults: Record<string, boolean> = {}

    if (layoutsWithExpressions.length > 0 && $SystemAPI) {
      const expressions: Record<string, string> = {}
      for (const layout of layoutsWithExpressions) {
        expressions[layout.pageLayoutID] = layout.config!.visibility!.expression!
      }

      try {
        expressionResults = await $SystemAPI.expressionEvaluate({ variables, expressions })
      } catch {
        // On error treat all expressions as false — no layout matches by expression
        for (const id of Object.keys(expressions)) {
          expressionResults[id] = false
        }
      }
    }

    // Find first matching layout
    for (const layout of layouts) {
      if (requestedLayoutID && layout.pageLayoutID !== requestedLayoutID) continue

      const { expression, roles = [] } = layout.config?.visibility || {}

      // Expression must pass if set
      if (expression && !expressionResults[layout.pageLayoutID]) continue

      // If roles are set, user must have at least one
      if (roles.length > 0) {
        const userRoles = $auth?.user?.roles || []
        if (!userRoles.some(roleID => roles.includes(roleID))) continue
      }

      return layout
    }

    // If a specific layout was requested but didn't match, retry without the constraint
    if (requestedLayoutID) {
      return determineLayout(layouts, variables)
    }

    return null
  }

  async function evaluateBlocks(
    blocks: Block[],
    variables: Record<string, unknown>,
  ): Promise<Set<string>> {
    const invisibleIDs = new Set<string>()

    // Batch evaluate all blocks that have visibility expressions
    const blocksWithExpressions = blocks.filter(b => b.meta?.visibility?.expression)
    let expressionResults: Record<string, boolean> = {}

    if (blocksWithExpressions.length > 0 && $SystemAPI) {
      const expressions: Record<string, string> = {}
      for (const block of blocksWithExpressions) {
        const id = fetchBlockID(block)
        if (id) expressions[id] = block.meta!.visibility!.expression!
      }

      try {
        expressionResults = await $SystemAPI.expressionEvaluate({ variables, expressions })
      } catch {
        for (const id of Object.keys(expressions)) {
          expressionResults[id] = false
        }
      }
    }

    // Determine which blocks are invisible
    for (const block of blocks) {
      const id = fetchBlockID(block)
      if (!id) continue

      const { expression, roles = [] } = block.meta?.visibility || {}

      const validExpression = !expression || !!expressionResults[id]
      const userRoles = $auth?.user?.roles || []
      const validRole = !roles.length || userRoles.some(roleID => roles.includes(roleID))

      if (!validExpression || !validRole) {
        invisibleIDs.add(id)
      }
    }

    // Tabs propagation: if all child tabs are invisible, the Tabs block itself becomes invisible
    for (const block of blocks) {
      if (block.kind !== 'Tabs') continue
      const blockID = fetchBlockID(block)
      if (invisibleIDs.has(blockID)) continue

      const tabs = block.options?.tabs || []
      const hasVisibleTab = tabs.some(tab => {
        if (!tab.blockID) return !!tab.title
        const child = blocks.find(b => fetchBlockID(b) === tab.blockID)
        return child ? !invisibleIDs.has(fetchBlockID(child)) : !!tab.title
      })

      if (!hasVisibleTab) {
        invisibleIDs.add(blockID)
      }
    }

    return invisibleIDs
  }

  return { buildExpressionVariables, determineLayout, evaluateBlocks }
}
