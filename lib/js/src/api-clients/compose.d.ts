import { AxiosInstance, AxiosRequestConfig } from 'axios'
interface KV {
  [header: string]: unknown
}
interface Headers {
  [header: string]: string
}
interface Ctor {
  baseURL?: string
  accessTokenFn?: () => string | undefined
  headers?: Headers
}
export default class Compose {
  protected baseURL?: string
  protected accessTokenFn?: () => string | undefined
  protected headers: Headers
  constructor({ baseURL, headers, accessTokenFn }: Ctor)
  setAccessTokenFn(fn: () => string | undefined): Compose
  setHeaders(headers?: Headers): Compose
  setHeader(name: string, value: string | undefined): Compose
  api(): AxiosInstance
  namespaceList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceListEndpoint(): string
  namespaceCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceCreateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceCreateEndpoint(): string
  namespaceRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceReadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceReadEndpoint(a: KV): string
  namespaceUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceUpdateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceUpdateEndpoint(a: KV): string
  namespaceDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceDeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceDeleteEndpoint(a: KV): string
  namespaceUpload(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceUploadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceUploadEndpoint(): string
  namespaceClone(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceCloneCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceCloneEndpoint(a: KV): string
  namespaceExport(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceExportCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceExportEndpoint(a: KV): string
  namespaceImportInit(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceImportInitCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceImportInitEndpoint(): string
  namespaceImportRun(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceImportRunCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceImportRunEndpoint(a: KV): string
  namespaceTriggerScript(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceTriggerScriptCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceTriggerScriptEndpoint(a: KV): string
  namespaceListTranslations(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceListTranslationsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceListTranslationsEndpoint(a: KV): string
  namespaceUpdateTranslations(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  namespaceUpdateTranslationsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  namespaceUpdateTranslationsEndpoint(a: KV): string
  pageList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageListEndpoint(a: KV): string
  pageCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageCreateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageCreateEndpoint(a: KV): string
  pageRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageReadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageReadEndpoint(a: KV): string
  pageTree(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageTreeCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageTreeEndpoint(a: KV): string
  pageUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageUpdateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageUpdateEndpoint(a: KV): string
  pageReorder(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageReorderCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageReorderEndpoint(a: KV): string
  pageDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageDeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageDeleteEndpoint(a: KV): string
  pageUpload(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageUploadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageUploadEndpoint(a: KV): string
  pageTriggerScript(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageTriggerScriptCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageTriggerScriptEndpoint(a: KV): string
  pageListTranslations(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageListTranslationsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageListTranslationsEndpoint(a: KV): string
  pageUpdateTranslations(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageUpdateTranslationsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageUpdateTranslationsEndpoint(a: KV): string
  pageUpdateIcon(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageUpdateIconCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageUpdateIconEndpoint(a: KV): string
  iconList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  iconListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  iconListEndpoint(): string
  iconUpload(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  iconUploadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  iconUploadEndpoint(): string
  iconDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  iconDeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  iconDeleteEndpoint(a: KV): string
  pageLayoutListNamespace(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageLayoutListNamespaceCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageLayoutListNamespaceEndpoint(a: KV): string
  pageLayoutList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageLayoutListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageLayoutListEndpoint(a: KV): string
  pageLayoutCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageLayoutCreateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageLayoutCreateEndpoint(a: KV): string
  pageLayoutRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageLayoutReadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageLayoutReadEndpoint(a: KV): string
  pageLayoutUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageLayoutUpdateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageLayoutUpdateEndpoint(a: KV): string
  pageLayoutReorder(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageLayoutReorderCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageLayoutReorderEndpoint(a: KV): string
  pageLayoutDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageLayoutDeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageLayoutDeleteEndpoint(a: KV): string
  pageLayoutUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageLayoutUndeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageLayoutUndeleteEndpoint(a: KV): string
  pageLayoutListTranslations(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageLayoutListTranslationsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageLayoutListTranslationsEndpoint(a: KV): string
  pageLayoutUpdateTranslations(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  pageLayoutUpdateTranslationsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  pageLayoutUpdateTranslationsEndpoint(a: KV): string
  moduleList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  moduleListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  moduleListEndpoint(a: KV): string
  moduleCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  moduleCreateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  moduleCreateEndpoint(a: KV): string
  moduleRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  moduleReadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  moduleReadEndpoint(a: KV): string
  moduleUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  moduleUpdateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  moduleUpdateEndpoint(a: KV): string
  moduleDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  moduleDeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  moduleDeleteEndpoint(a: KV): string
  moduleTriggerScript(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  moduleTriggerScriptCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  moduleTriggerScriptEndpoint(a: KV): string
  moduleListTranslations(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  moduleListTranslationsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  moduleListTranslationsEndpoint(a: KV): string
  moduleUpdateTranslations(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  moduleUpdateTranslationsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  moduleUpdateTranslationsEndpoint(a: KV): string
  recordReport(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordReportCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordReportEndpoint(a: KV): string
  recordList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordListEndpoint(a: KV): string
  recordImportInit(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordImportInitCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordImportInitEndpoint(a: KV): string
  recordImportRun(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordImportRunCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordImportRunEndpoint(a: KV): string
  recordImportProgress(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordImportProgressCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordImportProgressEndpoint(a: KV): string
  recordExport(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordExportCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordExportEndpoint(a: KV): string
  recordExec(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordExecCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordExecEndpoint(a: KV): string
  recordCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordCreateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordCreateEndpoint(a: KV): string
  recordRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordReadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordReadEndpoint(a: KV): string
  recordUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordUpdateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordUpdateEndpoint(a: KV): string
  recordPatch(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordPatchCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordPatchEndpoint(a: KV): string
  recordBulkDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordBulkDeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordBulkDeleteEndpoint(a: KV): string
  recordDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordDeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordDeleteEndpoint(a: KV): string
  recordUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordUndeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordUndeleteEndpoint(a: KV): string
  recordBulkUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordBulkUndeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordBulkUndeleteEndpoint(a: KV): string
  recordUpload(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordUploadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordUploadEndpoint(a: KV): string
  recordTriggerScript(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordTriggerScriptCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordTriggerScriptEndpoint(a: KV): string
  recordTriggerScriptOnList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordTriggerScriptOnListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordTriggerScriptOnListEndpoint(a: KV): string
  recordRevisions(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  recordRevisionsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  recordRevisionsEndpoint(a: KV): string
  dataPrivacyRecordList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  dataPrivacyRecordListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  dataPrivacyRecordListEndpoint(): string
  dataPrivacyModuleList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  dataPrivacyModuleListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  dataPrivacyModuleListEndpoint(): string
  chartList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  chartListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  chartListEndpoint(a: KV): string
  chartCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  chartCreateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  chartCreateEndpoint(a: KV): string
  chartRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  chartReadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  chartReadEndpoint(a: KV): string
  chartUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  chartUpdateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  chartUpdateEndpoint(a: KV): string
  chartDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  chartDeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  chartDeleteEndpoint(a: KV): string
  chartListTranslations(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  chartListTranslationsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  chartListTranslationsEndpoint(a: KV): string
  chartUpdateTranslations(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  chartUpdateTranslationsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  chartUpdateTranslationsEndpoint(a: KV): string
  notificationEmailSend(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  notificationEmailSendCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  notificationEmailSendEndpoint(): string
  attachmentList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  attachmentListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  attachmentListEndpoint(a: KV): string
  attachmentRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  attachmentReadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  attachmentReadEndpoint(a: KV): string
  attachmentDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  attachmentDeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  attachmentDeleteEndpoint(a: KV): string
  attachmentOriginal(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  attachmentOriginalCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  attachmentOriginalEndpoint(a: KV): string
  attachmentPreview(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  attachmentPreviewCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  attachmentPreviewEndpoint(a: KV): string
  permissionsList(extra?: AxiosRequestConfig): Promise<KV>
  permissionsListCancellable(extra?: AxiosRequestConfig): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  permissionsListEndpoint(): string
  permissionsEffective(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  permissionsEffectiveCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  permissionsEffectiveEndpoint(): string
  permissionsTrace(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  permissionsTraceCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  permissionsTraceEndpoint(): string
  permissionsRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  permissionsReadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  permissionsReadEndpoint(a: KV): string
  permissionsDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  permissionsDeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  permissionsDeleteEndpoint(a: KV): string
  permissionsUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  permissionsUpdateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  permissionsUpdateEndpoint(a: KV): string
  automationList(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  automationListCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  automationListEndpoint(): string
  automationBundle(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  automationBundleCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  automationBundleEndpoint(a: KV): string
  automationTriggerScript(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  automationTriggerScriptCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  automationTriggerScriptEndpoint(): string
}
export {}
