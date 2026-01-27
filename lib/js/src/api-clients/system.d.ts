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
export default class System {
  protected baseURL?: string
  protected accessTokenFn?: () => (string | undefined)
  protected headers: Headers
  constructor({ baseURL, headers, accessTokenFn }: Ctor);
  setAccessTokenFn(fn: () => string | undefined): System;
  setHeaders(headers?: Headers): System;
  setHeader(name: string, value: string | undefined): System;
  api(): AxiosInstance;
  authImpersonate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  authImpersonateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  authImpersonateEndpoint(): string;
  authClientList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  authClientListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  authClientListEndpoint(): string;
  authClientCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  authClientCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  authClientCreateEndpoint(): string;
  authClientUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  authClientUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  authClientUpdateEndpoint(a: KV): string;
  authClientRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  authClientReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  authClientReadEndpoint(a: KV): string;
  authClientDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  authClientDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  authClientDeleteEndpoint(a: KV): string;
  authClientUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  authClientUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  authClientUndeleteEndpoint(a: KV): string;
  authClientRegenerateSecret(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  authClientRegenerateSecretCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  authClientRegenerateSecretEndpoint(a: KV): string;
  authClientExposeSecret(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  authClientExposeSecretCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  authClientExposeSecretEndpoint(a: KV): string;
  expressionEvaluate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  expressionEvaluateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  expressionEvaluateEndpoint(): string;
  settingsList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  settingsListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  settingsListEndpoint(): string;
  settingsUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  settingsUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  settingsUpdateEndpoint(): string;
  settingsGet(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  settingsGetCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  settingsGetEndpoint(a: KV): string;
  settingsSet(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  settingsSetCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  settingsSetEndpoint(a: KV): string;
  settingsCurrent(extra?: AxiosRequestConfig): Promise<KV>;
  settingsCurrentCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  settingsCurrentEndpoint(): string;
  roleList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleListEndpoint(): string;
  roleCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleCreateEndpoint(): string;
  roleUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleUpdateEndpoint(a: KV): string;
  roleRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleReadEndpoint(a: KV): string;
  roleDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleDeleteEndpoint(a: KV): string;
  roleArchive(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleArchiveCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleArchiveEndpoint(a: KV): string;
  roleUnarchive(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleUnarchiveCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleUnarchiveEndpoint(a: KV): string;
  roleUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleUndeleteEndpoint(a: KV): string;
  roleMove(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleMoveCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleMoveEndpoint(a: KV): string;
  roleMerge(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleMergeCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleMergeEndpoint(a: KV): string;
  roleMemberList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleMemberListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleMemberListEndpoint(a: KV): string;
  roleMemberAdd(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleMemberAddCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleMemberAddEndpoint(a: KV): string;
  roleMemberRemove(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleMemberRemoveCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleMemberRemoveEndpoint(a: KV): string;
  roleTriggerScript(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleTriggerScriptCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleTriggerScriptEndpoint(a: KV): string;
  roleCloneRules(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  roleCloneRulesCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  roleCloneRulesEndpoint(a: KV): string;
  userList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userListEndpoint(): string;
  userCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userCreateEndpoint(): string;
  userUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userUpdateEndpoint(a: KV): string;
  userPartialUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userPartialUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userPartialUpdateEndpoint(a: KV): string;
  userRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userReadEndpoint(a: KV): string;
  userDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userDeleteEndpoint(a: KV): string;
  userSuspend(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userSuspendCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userSuspendEndpoint(a: KV): string;
  userUnsuspend(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userUnsuspendCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userUnsuspendEndpoint(a: KV): string;
  userUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userUndeleteEndpoint(a: KV): string;
  userSetPassword(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userSetPasswordCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userSetPasswordEndpoint(a: KV): string;
  userMembershipList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userMembershipListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userMembershipListEndpoint(a: KV): string;
  userMembershipAdd(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userMembershipAddCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userMembershipAddEndpoint(a: KV): string;
  userMembershipRemove(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userMembershipRemoveCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userMembershipRemoveEndpoint(a: KV): string;
  userTriggerScript(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userTriggerScriptCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userTriggerScriptEndpoint(a: KV): string;
  userSessionsRemove(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userSessionsRemoveCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userSessionsRemoveEndpoint(a: KV): string;
  userListCredentials(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userListCredentialsCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userListCredentialsEndpoint(a: KV): string;
  userDeleteCredentials(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userDeleteCredentialsCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userDeleteCredentialsEndpoint(a: KV): string;
  userProfileAvatar(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userProfileAvatarCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userProfileAvatarEndpoint(a: KV): string;
  userProfileAvatarInitial(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userProfileAvatarInitialCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userProfileAvatarInitialEndpoint(a: KV): string;
  userDeleteAvatar(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userDeleteAvatarCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userDeleteAvatarEndpoint(a: KV): string;
  userExport(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userExportCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userExportEndpoint(a: KV): string;
  userImport(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  userImportCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  userImportEndpoint(): string;
  dalDriverList(extra?: AxiosRequestConfig): Promise<KV>;
  dalDriverListCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalDriverListEndpoint(): string;
  dalSensitivityLevelList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalSensitivityLevelListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalSensitivityLevelListEndpoint(): string;
  dalSensitivityLevelCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalSensitivityLevelCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalSensitivityLevelCreateEndpoint(): string;
  dalSensitivityLevelUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalSensitivityLevelUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalSensitivityLevelUpdateEndpoint(a: KV): string;
  dalSensitivityLevelRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalSensitivityLevelReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalSensitivityLevelReadEndpoint(a: KV): string;
  dalSensitivityLevelDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalSensitivityLevelDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalSensitivityLevelDeleteEndpoint(a: KV): string;
  dalSensitivityLevelUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalSensitivityLevelUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalSensitivityLevelUndeleteEndpoint(a: KV): string;
  dalSchemaAlterationList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalSchemaAlterationListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalSchemaAlterationListEndpoint(): string;
  dalSchemaAlterationRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalSchemaAlterationReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalSchemaAlterationReadEndpoint(a: KV): string;
  dalSchemaAlterationApply(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalSchemaAlterationApplyCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalSchemaAlterationApplyEndpoint(): string;
  dalSchemaAlterationDismiss(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalSchemaAlterationDismissCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalSchemaAlterationDismissEndpoint(): string;
  dalConnectionList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalConnectionListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalConnectionListEndpoint(): string;
  dalConnectionCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalConnectionCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalConnectionCreateEndpoint(): string;
  dalConnectionUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalConnectionUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalConnectionUpdateEndpoint(a: KV): string;
  dalConnectionRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalConnectionReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalConnectionReadEndpoint(a: KV): string;
  dalConnectionDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalConnectionDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalConnectionDeleteEndpoint(a: KV): string;
  dalConnectionUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dalConnectionUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dalConnectionUndeleteEndpoint(a: KV): string;
  applicationList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  applicationListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  applicationListEndpoint(): string;
  applicationCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  applicationCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  applicationCreateEndpoint(): string;
  applicationUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  applicationUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  applicationUpdateEndpoint(a: KV): string;
  applicationUpload(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  applicationUploadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  applicationUploadEndpoint(): string;
  applicationFlagCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  applicationFlagCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  applicationFlagCreateEndpoint(a: KV): string;
  applicationFlagDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  applicationFlagDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  applicationFlagDeleteEndpoint(a: KV): string;
  applicationRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  applicationReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  applicationReadEndpoint(a: KV): string;
  applicationDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  applicationDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  applicationDeleteEndpoint(a: KV): string;
  applicationUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  applicationUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  applicationUndeleteEndpoint(a: KV): string;
  applicationTriggerScript(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  applicationTriggerScriptCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  applicationTriggerScriptEndpoint(a: KV): string;
  applicationReorder(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  applicationReorderCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  applicationReorderEndpoint(): string;
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
  reminderList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reminderListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reminderListEndpoint(): string;
  reminderCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reminderCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reminderCreateEndpoint(): string;
  reminderUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reminderUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reminderUpdateEndpoint(a: KV): string;
  reminderRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reminderReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reminderReadEndpoint(a: KV): string;
  reminderDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reminderDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reminderDeleteEndpoint(a: KV): string;
  reminderDismiss(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reminderDismissCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reminderDismissEndpoint(a: KV): string;
  reminderUndismiss(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reminderUndismissCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reminderUndismissEndpoint(a: KV): string;
  reminderSnooze(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reminderSnoozeCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reminderSnoozeEndpoint(a: KV): string;
  notificationList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  notificationListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  notificationListEndpoint(): string;
  notificationCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  notificationCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  notificationCreateEndpoint(): string;
  notificationUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  notificationUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  notificationUpdateEndpoint(a: KV): string;
  notificationRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  notificationReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  notificationReadEndpoint(a: KV): string;
  notificationDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  notificationDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  notificationDeleteEndpoint(a: KV): string;
  notificationMarkAsRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  notificationMarkAsReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  notificationMarkAsReadEndpoint(a: KV): string;
  notificationMarkAllAsRead(extra?: AxiosRequestConfig): Promise<KV>;
  notificationMarkAllAsReadCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  notificationMarkAllAsReadEndpoint(): string;
  attachmentRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  attachmentReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  attachmentReadEndpoint(a: KV): string;
  attachmentDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  attachmentDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  attachmentDeleteEndpoint(a: KV): string;
  attachmentOriginal(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  attachmentOriginalCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  attachmentOriginalEndpoint(a: KV): string;
  attachmentPreview(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  attachmentPreviewCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  attachmentPreviewEndpoint(a: KV): string;
  templateList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  templateListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  templateListEndpoint(): string;
  templateCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  templateCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  templateCreateEndpoint(): string;
  templateRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  templateReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  templateReadEndpoint(a: KV): string;
  templateUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  templateUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  templateUpdateEndpoint(a: KV): string;
  templateDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  templateDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  templateDeleteEndpoint(a: KV): string;
  templateUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  templateUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  templateUndeleteEndpoint(a: KV): string;
  templateRenderDrivers(extra?: AxiosRequestConfig): Promise<KV>;
  templateRenderDriversCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  templateRenderDriversEndpoint(): string;
  templateRender(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  templateRenderCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  templateRenderEndpoint(a: KV): string;
  reportList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reportListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reportListEndpoint(): string;
  reportCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reportCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reportCreateEndpoint(): string;
  reportUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reportUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reportUpdateEndpoint(a: KV): string;
  reportRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reportReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reportReadEndpoint(a: KV): string;
  reportDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reportDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reportDeleteEndpoint(a: KV): string;
  reportUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reportUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reportUndeleteEndpoint(a: KV): string;
  reportDescribe(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reportDescribeCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reportDescribeEndpoint(): string;
  reportRun(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  reportRunCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  reportRunEndpoint(a: KV): string;
  statsList(extra?: AxiosRequestConfig): Promise<KV>;
  statsListCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  statsListEndpoint(): string;
  automationList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  automationListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  automationListEndpoint(): string;
  automationBundle(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  automationBundleCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  automationBundleEndpoint(a: KV): string;
  automationTriggerScript(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  automationTriggerScriptCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  automationTriggerScriptEndpoint(): string;
  actionlogList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  actionlogListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  actionlogListEndpoint(): string;
  queuesList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  queuesListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  queuesListEndpoint(): string;
  queuesCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  queuesCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  queuesCreateEndpoint(): string;
  queuesRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  queuesReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  queuesReadEndpoint(a: KV): string;
  queuesUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  queuesUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  queuesUpdateEndpoint(a: KV): string;
  queuesDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  queuesDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  queuesDeleteEndpoint(a: KV): string;
  queuesUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  queuesUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  queuesUndeleteEndpoint(a: KV): string;
  apigwRouteList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwRouteListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwRouteListEndpoint(): string;
  apigwRouteCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwRouteCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwRouteCreateEndpoint(): string;
  apigwRouteUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwRouteUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwRouteUpdateEndpoint(a: KV): string;
  apigwRouteRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwRouteReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwRouteReadEndpoint(a: KV): string;
  apigwRouteDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwRouteDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwRouteDeleteEndpoint(a: KV): string;
  apigwRouteUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwRouteUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwRouteUndeleteEndpoint(a: KV): string;
  apigwFilterList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwFilterListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwFilterListEndpoint(): string;
  apigwFilterCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwFilterCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwFilterCreateEndpoint(): string;
  apigwFilterUpdate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwFilterUpdateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwFilterUpdateEndpoint(a: KV): string;
  apigwFilterRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwFilterReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwFilterReadEndpoint(a: KV): string;
  apigwFilterDelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwFilterDeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwFilterDeleteEndpoint(a: KV): string;
  apigwFilterUndelete(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwFilterUndeleteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwFilterUndeleteEndpoint(a: KV): string;
  apigwFilterDefFilter(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwFilterDefFilterCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwFilterDefFilterEndpoint(): string;
  apigwFilterDefProxyAuth(extra?: AxiosRequestConfig): Promise<KV>;
  apigwFilterDefProxyAuthCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwFilterDefProxyAuthEndpoint(): string;
  apigwProfilerAggregation(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwProfilerAggregationCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwProfilerAggregationEndpoint(): string;
  apigwProfilerRoute(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwProfilerRouteCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwProfilerRouteEndpoint(a: KV): string;
  apigwProfilerHit(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwProfilerHitCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwProfilerHitEndpoint(a: KV): string;
  apigwProfilerPurgeAll(extra?: AxiosRequestConfig): Promise<KV>;
  apigwProfilerPurgeAllCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwProfilerPurgeAllEndpoint(): string;
  apigwProfilerPurge(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  apigwProfilerPurgeCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  apigwProfilerPurgeEndpoint(a: KV): string;
  localeListResource(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  localeListResourceCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  localeListResourceEndpoint(): string;
  localeCreateResource(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  localeCreateResourceCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  localeCreateResourceEndpoint(): string;
  localeUpdateResource(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  localeUpdateResourceCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  localeUpdateResourceEndpoint(a: KV): string;
  localeReadResource(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  localeReadResourceCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  localeReadResourceEndpoint(a: KV): string;
  localeDeleteResource(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  localeDeleteResourceCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  localeDeleteResourceEndpoint(a: KV): string;
  localeUndeleteResource(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  localeUndeleteResourceCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  localeUndeleteResourceEndpoint(a: KV): string;
  localeList(extra?: AxiosRequestConfig): Promise<KV>;
  localeListCancellable(extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  localeListEndpoint(): string;
  localeGet(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  localeGetCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  localeGetEndpoint(a: KV): string;
  dataPrivacyConnectionList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dataPrivacyConnectionListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dataPrivacyConnectionListEndpoint(): string;
  dataPrivacyRequestList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dataPrivacyRequestListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dataPrivacyRequestListEndpoint(): string;
  dataPrivacyRequestCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dataPrivacyRequestCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dataPrivacyRequestCreateEndpoint(): string;
  dataPrivacyRequestRead(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dataPrivacyRequestReadCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dataPrivacyRequestReadEndpoint(a: KV): string;
  dataPrivacyRequestUpdateStatus(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dataPrivacyRequestUpdateStatusCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dataPrivacyRequestUpdateStatusEndpoint(a: KV): string;
  dataPrivacyRequestCommentList(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dataPrivacyRequestCommentListCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dataPrivacyRequestCommentListEndpoint(a: KV): string;
  dataPrivacyRequestCommentCreate(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  dataPrivacyRequestCommentCreateCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  dataPrivacyRequestCommentCreateEndpoint(a: KV): string;
  smtpConfigurationCheckerCheck(a: KV, extra?: AxiosRequestConfig): Promise<KV>;
  smtpConfigurationCheckerCheckCancellable(a: KV, extra?: AxiosRequestConfig): {
        response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>;
        cancel: () => void;
    };
  smtpConfigurationCheckerCheckEndpoint(): string;
}
export {}
