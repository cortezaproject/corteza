import type { compose } from '@planetcrust/human-js'

export interface ComposeUIToast {
  success: (_message: string) => void
  warning: (_message: string) => void
}

export interface ComposeUIContext {
  $namespace?: compose.Namespace
  $module?: compose.Module
  $record?: compose.Record
  pages?: () => compose.Page[]
  toast?: ComposeUIToast
  routePusher?: (_to: object) => void
}

/**
 * The webapp's Compose UI, as a Corredor client script sees it.
 */
export class ComposeUIHelper {
  readonly $namespace?: compose.Namespace
  readonly $module?: compose.Module
  readonly $record?: compose.Record

  protected readonly pages: () => compose.Page[]
  protected readonly toast?: ComposeUIToast
  protected readonly routePusher?: (_to: object) => void

  constructor(ctx: ComposeUIContext) {
    this.$namespace = ctx.$namespace
    this.$module = ctx.$module
    this.$record = ctx.$record
    this.pages = ctx.pages || ((): compose.Page[] => [])
    this.toast = ctx.toast
    this.routePusher = ctx.routePusher
  }

  /**
   * Opens the record page of a record in view mode
   *
   * @example
   * ComposeUI.gotoRecordViewer($record)
   * ComposeUI.gotoRecordViewer()
   */
  gotoRecordViewer(record: compose.Record | undefined = this.$record): void {
    this.gotoRecordPage(record)
  }

  /**
   * Opens the record page of a record in edit mode
   *
   * @example
   * ComposeUI.gotoRecordEditor($record)
   * ComposeUI.gotoRecordEditor()
   */
  gotoRecordEditor(record: compose.Record | undefined = this.$record): void {
    this.gotoRecordPage(record, { edit: '1' })
  }

  /**
   * Shows a success message
   *
   * @example
   * ComposeUI.success('Change was successful')
   */
  success(message: string): void {
    this.toast?.success(message)
  }

  /**
   * Shows a warning message
   *
   * @example
   * ComposeUI.warning('Could not save your changes')
   */
  warning(message: string): void {
    this.toast?.warning(message)
  }

  protected gotoRecordPage(record?: compose.Record, query?: object): void {
    const recordPage = this.getRecordPage(record)
    const pageID = recordPage?.pageID
    const recordID = record?.recordID
    const slug = this.$namespace?.slug

    if (!pageID) {
      throw Error('record page does not exist')
    }

    if (!recordID) {
      throw Error('invalid record')
    }

    if (!slug) {
      throw Error('namespace is unknown')
    }

    this.goto({
      name: 'page.record',
      params: { slug, pageID, recordID },
      ...(query ? { query } : {}),
    })
  }

  /**
   * The record page of a module — the first page bound to it
   */
  protected getRecordPage(
    m: { moduleID?: string } | undefined = this.$module,
  ): compose.Page | undefined {
    if (!m?.moduleID) {
      return undefined
    }

    return this.pages().find(p => p.moduleID === m.moduleID)
  }

  protected goto(to: object): void {
    this.routePusher?.(to)
  }
}
