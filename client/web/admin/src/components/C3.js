import * as applications from './Application/C3'
import * as role from './Role/C3'
import * as settings from './Settings/C3'

import { default as CAuthclientEditorInfo } from './Authclient/C3'
import { default as CFederationEditorInfo } from './Federation/C3'
import { default as CQueueEditorInfo } from './Queues/C3'
import { default as CTemplateEditorInfo } from './Template/C3'
import { default as CUserEditorInfo } from './User/C3'
import { default as CWorkflowEditorInfo } from './Workflow/C3'
import { default as CResourceListStatusFilter } from './CResourceListStatusFilter.c3'
import { default as CSubmitButton } from './CSubmitButton.c3'

export default {
  ...applications,
  ...role,
  ...settings,

  CAuthclientEditorInfo,
  CFederationEditorInfo,
  CQueueEditorInfo,
  CTemplateEditorInfo,
  CUserEditorInfo,
  CWorkflowEditorInfo,
  CResourceListStatusFilter,
  CSubmitButton,
}
