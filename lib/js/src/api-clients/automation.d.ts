import { AxiosInstance, AxiosRequestConfig } from 'axios'
interface KV {
    [header: string]: unknown;
}
interface Headers {
    [header: string]: string;
}
interface Ctor {
    baseURL?: string;
    accessTokenFn?: () => string | undefined;
    headers?: Headers;
}
export default class Automation {
  protected baseURL?: string
  protected accessTokenFn?: () => (string | undefined)
  protected headers: Headers
  constructor({ baseURL, headers, accessTokenFn }: Ctor);
  setAccessTokenFn(fn: () => string | undefined): Automation;
  setHeaders(headers?: Headers): Automation;
  setHeader(name: string, value: string | undefined): Automation;
  api(): AxiosInstance;
  workflowList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  workflowListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  workflowListEndpoint(): string;
  workflowCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  workflowCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  workflowCreateEndpoint(): string;
  workflowUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  workflowUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  workflowUpdateEndpoint(a: KV): string;
  workflowRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  workflowReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  workflowReadEndpoint(a: KV): string;
  workflowDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  workflowDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  workflowDeleteEndpoint(a: KV): string;
  workflowUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  workflowUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  workflowUndeleteEndpoint(a: KV): string;
  workflowTest(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  workflowTestCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  workflowTestEndpoint(a: KV): string;
  workflowExec(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  workflowExecCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  workflowExecEndpoint(a: KV): string;
  triggerList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  triggerListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  triggerListEndpoint(): string;
  triggerCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  triggerCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  triggerCreateEndpoint(): string;
  triggerUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  triggerUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  triggerUpdateEndpoint(a: KV): string;
  triggerRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  triggerReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  triggerReadEndpoint(a: KV): string;
  triggerDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  triggerDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  triggerDeleteEndpoint(a: KV): string;
  triggerUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  triggerUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  triggerUndeleteEndpoint(a: KV): string;
  sessionList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  sessionListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  sessionListEndpoint(): string;
  sessionRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  sessionReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  sessionReadEndpoint(a: KV): string;
  sessionCancel(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  sessionCancelCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  sessionCancelEndpoint(a: KV): string;
  sessionListPrompts(extra?: AxiosRequestConfig): Promise<KV>;
  sessionListPromptsCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  sessionListPromptsEndpoint(): string;
  sessionResumeState(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  sessionResumeStateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  sessionResumeStateEndpoint(a: KV): string;
  functionList(extra?: AxiosRequestConfig): Promise<KV>;
  functionListCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  functionListEndpoint(): string;
  typeList(extra?: AxiosRequestConfig): Promise<KV>;
  typeListCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  typeListEndpoint(): string;
  eventTypesList(extra?: AxiosRequestConfig): Promise<KV>;
  eventTypesListCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  eventTypesListEndpoint(): string;
  permissionsList(extra?: AxiosRequestConfig): Promise<KV>;
  permissionsListCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  permissionsListEndpoint(): string;
  permissionsEffective(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  permissionsEffectiveCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  permissionsEffectiveEndpoint(): string;
  permissionsTrace(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  permissionsTraceCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  permissionsTraceEndpoint(): string;
  permissionsRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  permissionsReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  permissionsReadEndpoint(a: KV): string;
  permissionsDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  permissionsDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  permissionsDeleteEndpoint(a: KV): string;
  permissionsUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  permissionsUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  permissionsUpdateEndpoint(a: KV): string;
}
export {}
