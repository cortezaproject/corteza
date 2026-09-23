// This is a generated file.
// Run `pnpm run codegen` in lib/js to rebuild it.
//
// A method here mirrors one endpoint of the REST spec, so a change the server
// makes to that endpoint reaches scripts as soon as this file is regenerated.
// Over the raw API client a method adds two things and nothing else: every ID
// in the endpoint's path also accepts the object it names or its handle, and a
// response is cast into its resource class where one exists. The curated helper
// that extends this class is where a breaking endpoint change is absorbed.

import { AxiosRequestConfig } from 'axios'
import { Automation as AutomationAPI } from '../../api-clients'
import { Function as AutomationFunction, NgAutomation, Workflow } from '../../automation'
import { Args, IdArgs, KV, ListResponse, RefSpec, argsOf, castSet, resolveRef } from './shared'

// Where a handle, a slug or a name arriving in place of an ID is looked up.
const refs: RefSpec = {
  workflow: { idField: 'workflowID', list: 'workflowList', by: [], parents: [] },
  trigger: { idField: 'triggerID', list: 'triggerList', by: [], parents: [] },
  session: { idField: 'sessionID', list: 'sessionList', by: [], parents: [] },
  ngAutomation: { idField: 'automationID', list: 'ngAutomationList', by: [], parents: [] },
}

export default class GeneratedAutomationHelper {
  readonly AutomationAPI: AutomationAPI

  constructor(ctx: { AutomationAPI: AutomationAPI }) {
    this.AutomationAPI = ctx.AutomationAPI
  }

  // Turns an object, an ID or a handle into the ID the endpoint wants.
  protected resolveRef(entity: string, value: unknown, context: KV = {}): Promise<unknown> {
    return resolveRef(this.AutomationAPI, refs, entity, value, context)
  }

  // List workflows
  // Mirrors AutomationAPI.workflowList.
  findWorkflows(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<Workflow[], KV>> {
    const args = argsOf(a)
    return this.AutomationAPI.workflowList(args, extra).then(r => castSet(r, v => new Workflow(v)))
  }

  // Create workflow
  // Mirrors AutomationAPI.workflowCreate.
  createWorkflow(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Workflow> {
    const args = argsOf(a)
    return this.AutomationAPI.workflowCreate(args, extra).then(r => new Workflow(r))
  }

  // Update triger details
  // Mirrors AutomationAPI.workflowUpdate.
  async updateWorkflow(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Workflow> {
    const args = argsOf(a, 'workflowID')
    args.workflowID = await this.resolveRef('workflow', args.workflowID, args)
    return this.AutomationAPI.workflowUpdate(args, extra).then(r => new Workflow(r))
  }

  // Read workflow details
  // Mirrors AutomationAPI.workflowRead.
  async findWorkflowByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Workflow> {
    const args = argsOf(a, 'workflowID')
    args.workflowID = await this.resolveRef('workflow', args.workflowID, args)
    return this.AutomationAPI.workflowRead(args, extra).then(r => new Workflow(r))
  }

  // Remove workflow
  // Mirrors AutomationAPI.workflowDelete.
  async deleteWorkflow(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'workflowID')
    args.workflowID = await this.resolveRef('workflow', args.workflowID, args)
    return this.AutomationAPI.workflowDelete(args, extra)
  }

  // Undelete workflow
  // Mirrors AutomationAPI.workflowUndelete.
  async undeleteWorkflow(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'workflowID')
    args.workflowID = await this.resolveRef('workflow', args.workflowID, args)
    return this.AutomationAPI.workflowUndelete(args, extra)
  }

  // Test workflow details
  async workflowTest(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'workflowID')
    args.workflowID = await this.resolveRef('workflow', args.workflowID, args)
    return this.AutomationAPI.workflowTest(args, extra)
  }

  // Executes workflow on a specific step (must be orphan step and connected to 'onManual' trigger)
  async workflowExec(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'workflowID')
    args.workflowID = await this.resolveRef('workflow', args.workflowID, args)
    return this.AutomationAPI.workflowExec(args, extra)
  }

  // List triggers
  // Mirrors AutomationAPI.triggerList.
  findTriggers(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.AutomationAPI.triggerList(args, extra)
  }

  // Create trigger
  // Mirrors AutomationAPI.triggerCreate.
  createTrigger(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.AutomationAPI.triggerCreate(args, extra)
  }

  // Update trigger details
  // Mirrors AutomationAPI.triggerUpdate.
  async updateTrigger(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'triggerID')
    args.triggerID = await this.resolveRef('trigger', args.triggerID, args)
    return this.AutomationAPI.triggerUpdate(args, extra)
  }

  // Read trigger details
  // Mirrors AutomationAPI.triggerRead.
  async findTriggerByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'triggerID')
    args.triggerID = await this.resolveRef('trigger', args.triggerID, args)
    return this.AutomationAPI.triggerRead(args, extra)
  }

  // Remove trigger
  // Mirrors AutomationAPI.triggerDelete.
  async deleteTrigger(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'triggerID')
    args.triggerID = await this.resolveRef('trigger', args.triggerID, args)
    return this.AutomationAPI.triggerDelete(args, extra)
  }

  // Undelete trigger
  // Mirrors AutomationAPI.triggerUndelete.
  async undeleteTrigger(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'triggerID')
    args.triggerID = await this.resolveRef('trigger', args.triggerID, args)
    return this.AutomationAPI.triggerUndelete(args, extra)
  }

  // List sessions
  // Mirrors AutomationAPI.sessionList.
  findSessions(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.AutomationAPI.sessionList(args, extra)
  }

  // Read session details
  // Mirrors AutomationAPI.sessionRead.
  async findSessionByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sessionID')
    args.sessionID = await this.resolveRef('session', args.sessionID, args)
    return this.AutomationAPI.sessionRead(args, extra)
  }

  // Cancel session
  async sessionCancel(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sessionID')
    args.sessionID = await this.resolveRef('session', args.sessionID, args)
    return this.AutomationAPI.sessionCancel(args, extra)
  }

  // Returns pending prompts from all sessions
  sessionListPrompts(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.AutomationAPI.sessionListPrompts(extra)
  }

  // Resume session
  async sessionResumeState(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.sessionID = await this.resolveRef('session', args.sessionID, args)
    return this.AutomationAPI.sessionResumeState(args, extra)
  }

  // Available workflow functions
  // Mirrors AutomationAPI.functionList.
  findFunctions(extra: AxiosRequestConfig = {}): Promise<ListResponse<AutomationFunction[], KV>> {
    return this.AutomationAPI.functionList(extra).then(r =>
      castSet(r, v => new AutomationFunction(v)),
    )
  }

  // Available workflow types
  typeList(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.AutomationAPI.typeList(extra)
  }

  // Available workflow types
  eventTypesList(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.AutomationAPI.eventTypesList(extra)
  }

  // Retrieve defined permissions
  permissionsList(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.AutomationAPI.permissionsList(extra)
  }

  // Effective rules for current user
  permissionsEffective(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.AutomationAPI.permissionsEffective(args, extra)
  }

  // Evaluate rules for given user/role combo
  permissionsTrace(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.AutomationAPI.permissionsTrace(args, extra)
  }

  // Retrieve role permissions
  permissionsRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    return this.AutomationAPI.permissionsRead(args, extra)
  }

  // Remove all defined role permissions
  permissionsDelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    return this.AutomationAPI.permissionsDelete(args, extra)
  }

  // Update permission settings
  permissionsUpdate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    return this.AutomationAPI.permissionsUpdate(args, extra)
  }

  // List functions
  constructLibraryFunctions(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.AutomationAPI.constructLibraryFunctions(extra)
  }

  // List triggers
  constructLibraryTriggers(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.AutomationAPI.constructLibraryTriggers(extra)
  }

  // List automations
  // Mirrors AutomationAPI.ngAutomationList.
  findTAQs(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<NgAutomation[], KV>> {
    const args = argsOf(a)
    return this.AutomationAPI.ngAutomationList(args, extra).then(r =>
      castSet(r, v => new NgAutomation(v)),
    )
  }

  // Create automation
  // Mirrors AutomationAPI.ngAutomationCreate.
  createTAQ(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<NgAutomation> {
    const args = argsOf(a)
    return this.AutomationAPI.ngAutomationCreate(args, extra).then(r => new NgAutomation(r))
  }

  // Update triger details
  // Mirrors AutomationAPI.ngAutomationUpdate.
  async updateTAQ(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<NgAutomation> {
    const args = argsOf(a, 'automationID')
    args.automationID = await this.resolveRef('ngAutomation', args.automationID, args)
    return this.AutomationAPI.ngAutomationUpdate(args, extra).then(r => new NgAutomation(r))
  }

  // Read automation details
  // Mirrors AutomationAPI.ngAutomationRead.
  async findTAQByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<NgAutomation> {
    const args = argsOf(a, 'automationID')
    args.automationID = await this.resolveRef('ngAutomation', args.automationID, args)
    return this.AutomationAPI.ngAutomationRead(args, extra).then(r => new NgAutomation(r))
  }

  // Remove automation
  // Mirrors AutomationAPI.ngAutomationDelete.
  async deleteTAQ(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'automationID')
    args.automationID = await this.resolveRef('ngAutomation', args.automationID, args)
    return this.AutomationAPI.ngAutomationDelete(args, extra)
  }

  // Undelete automation
  // Mirrors AutomationAPI.ngAutomationUndelete.
  async undeleteTAQ(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'automationID')
    args.automationID = await this.resolveRef('ngAutomation', args.automationID, args)
    return this.AutomationAPI.ngAutomationUndelete(args, extra)
  }

  // Test automation details
  async ngAutomationTest(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'automationID')
    args.automationID = await this.resolveRef('ngAutomation', args.automationID, args)
    return this.AutomationAPI.ngAutomationTest(args, extra)
  }

  // Executes automation on a specific step (must be orphan step and connected to 'onManual' trigger)
  async ngAutomationExec(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'automationID')
    args.automationID = await this.resolveRef('ngAutomation', args.automationID, args)
    return this.AutomationAPI.ngAutomationExec(args, extra)
  }

  // Get automation executions
  async ngAutomationExecutions(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'automationID')
    args.automationID = await this.resolveRef('ngAutomation', args.automationID, args)
    return this.AutomationAPI.ngAutomationExecutions(args, extra)
  }

  // Get automation execution trace
  async ngAutomationExecutionTrace(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.automationID = await this.resolveRef('ngAutomation', args.automationID, args)
    return this.AutomationAPI.ngAutomationExecutionTrace(args, extra)
  }

  // List all ng automation execution sessions
  ngAutomationAllExecutions(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.AutomationAPI.ngAutomationAllExecutions(args, extra)
  }
}
