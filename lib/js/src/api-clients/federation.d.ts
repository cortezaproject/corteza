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
export default class Federation {
  protected baseURL?: string
  protected accessTokenFn?: () => string | undefined
  protected headers: Headers
  constructor({ baseURL, headers, accessTokenFn }: Ctor)
  setAccessTokenFn(fn: () => string | undefined): Federation
  setHeaders(headers?: Headers): Federation
  setHeader(name: string, value: string | undefined): Federation
  api(): AxiosInstance
  nodeHandshakeInitialize(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  nodeHandshakeInitializeCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  nodeHandshakeInitializeEndpoint(a: KV): string
  nodeSearch(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  nodeSearchCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  nodeSearchEndpoint(): string
  nodeCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  nodeCreateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  nodeCreateEndpoint(): string
  nodeRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  nodeReadCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  nodeReadEndpoint(a: KV): string
  nodeGenerateUri(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  nodeGenerateUriCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  nodeGenerateUriEndpoint(a: KV): string
  nodeUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  nodeUpdateCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  nodeUpdateEndpoint(a: KV): string
  nodeDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  nodeDeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  nodeDeleteEndpoint(a: KV): string
  nodeUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  nodeUndeleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  nodeUndeleteEndpoint(a: KV): string
  nodePair(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  nodePairCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  nodePairEndpoint(a: KV): string
  nodeHandshakeConfirm(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  nodeHandshakeConfirmCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  nodeHandshakeConfirmEndpoint(a: KV): string
  nodeHandshakeComplete(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  nodeHandshakeCompleteCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  nodeHandshakeCompleteEndpoint(a: KV): string
  manageStructureReadExposed(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  manageStructureReadExposedCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  manageStructureReadExposedEndpoint(a: KV): string
  manageStructureCreateExposed(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  manageStructureCreateExposedCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  manageStructureCreateExposedEndpoint(a: KV): string
  manageStructureUpdateExposed(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  manageStructureUpdateExposedCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  manageStructureUpdateExposedEndpoint(a: KV): string
  manageStructureRemoveExposed(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  manageStructureRemoveExposedCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  manageStructureRemoveExposedEndpoint(a: KV): string
  manageStructureReadShared(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  manageStructureReadSharedCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  manageStructureReadSharedEndpoint(a: KV): string
  manageStructureCreateMappings(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  manageStructureCreateMappingsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  manageStructureCreateMappingsEndpoint(a: KV): string
  manageStructureReadMappings(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  manageStructureReadMappingsCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  manageStructureReadMappingsEndpoint(a: KV): string
  manageStructureListAll(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  manageStructureListAllCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  manageStructureListAllEndpoint(a: KV): string
  syncStructureReadExposedInternal(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  syncStructureReadExposedInternalCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  syncStructureReadExposedInternalEndpoint(a: KV): string
  syncStructureReadExposedSocial(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  syncStructureReadExposedSocialCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  syncStructureReadExposedSocialEndpoint(a: KV): string
  syncDataReadExposedAll(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  syncDataReadExposedAllCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  syncDataReadExposedAllEndpoint(a: KV): string
  syncDataReadExposedInternal(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  syncDataReadExposedInternalCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  syncDataReadExposedInternalEndpoint(a: KV): string
  syncDataReadExposedSocial(a: KV, extra?: AxiosRequestConfig): Promise<KV>
  syncDataReadExposedSocialCancellable(
    a: KV,
    extra?: AxiosRequestConfig,
  ): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  }
  syncDataReadExposedSocialEndpoint(a: KV): string
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
}
export {}
