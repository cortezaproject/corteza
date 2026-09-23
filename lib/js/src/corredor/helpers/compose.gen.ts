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
import { Compose as ComposeAPI } from '../../api-clients'
import { Chart, Module, Namespace, Page, PageLayout } from '../../compose'
import { Attachment } from '../../shared'
import { Args, IdArgs, KV, ListResponse, RefSpec, argsOf, castSet, resolveRef } from './shared'

// Where a handle, a slug or a name arriving in place of an ID is looked up.
const refs: RefSpec = {
  namespace: { idField: 'namespaceID', list: 'namespaceList', by: ['slug'], parents: [] },
  page: { idField: 'pageID', list: 'pageList', by: ['handle'], parents: ['namespaceID'] },
  icon: { idField: 'iconID', list: 'iconList', by: [], parents: [] },
  pageLayout: {
    idField: 'pageLayoutID',
    list: 'pageLayoutList',
    by: ['handle'],
    parents: ['namespaceID', 'pageID'],
  },
  module: {
    idField: 'moduleID',
    list: 'moduleList',
    by: ['handle', 'name'],
    parents: ['namespaceID'],
  },
  record: { idField: 'recordID', list: 'recordList', by: [], parents: ['namespaceID', 'moduleID'] },
  chart: { idField: 'chartID', list: 'chartList', by: ['handle'], parents: ['namespaceID'] },
  attachment: {
    idField: 'attachmentID',
    list: 'attachmentList',
    by: [],
    parents: ['kind', 'namespaceID'],
  },
}

export default class GeneratedComposeHelper {
  readonly ComposeAPI: ComposeAPI

  constructor(ctx: { ComposeAPI: ComposeAPI }) {
    this.ComposeAPI = ctx.ComposeAPI
  }

  // Turns an object, an ID or a handle into the ID the endpoint wants.
  protected resolveRef(entity: string, value: unknown, context: KV = {}): Promise<unknown> {
    return resolveRef(this.ComposeAPI, refs, entity, value, context)
  }

  // Create namespace
  // Mirrors ComposeAPI.namespaceCreate.
  createNamespace(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Namespace> {
    const args = argsOf(a)
    return this.ComposeAPI.namespaceCreate(args, extra).then(r => new Namespace(r))
  }

  // Update namespace
  // Mirrors ComposeAPI.namespaceUpdate.
  async updateNamespace(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Namespace> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.namespaceUpdate(args, extra).then(r => new Namespace(r))
  }

  // Delete namespace
  // Mirrors ComposeAPI.namespaceDelete.
  async deleteNamespace(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.namespaceDelete(args, extra)
  }

  // Upload namespace assets
  namespaceUpload(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.namespaceUpload(args, extra)
  }

  // Clone compose namespace
  async namespaceClone(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.namespaceClone(args, extra)
  }

  // Export compose namespace
  async namespaceExport(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.namespaceExport(args, extra)
  }

  // Initiate namespace import session
  namespaceImportInit(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.namespaceImportInit(args, extra)
  }

  // Run namespace import
  namespaceImportRun(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sessionID')
    return this.ComposeAPI.namespaceImportRun(args, extra)
  }

  // Fire compose:namespace trigger
  async namespaceTriggerScript(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.namespaceTriggerScript(args, extra)
  }

  // List translation
  async namespaceListTranslations(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.namespaceListTranslations(args, extra)
  }

  // Update translation
  async namespaceUpdateTranslations(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.namespaceUpdateTranslations(args, extra)
  }

  // Create page
  // Mirrors ComposeAPI.pageCreate.
  async createPage(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Page> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.pageCreate(args, extra).then(r => new Page(r))
  }

  // Get page all (non-record) pages, hierarchically
  async pageTree(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.pageTree(args, extra)
  }

  // Update page
  // Mirrors ComposeAPI.pageUpdate.
  async updatePage(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Page> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    return this.ComposeAPI.pageUpdate(args, extra).then(r => new Page(r))
  }

  // Reorder pages
  async pageReorder(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.pageReorder(args, extra)
  }

  // Uploads attachment to page
  async pageUpload(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    return this.ComposeAPI.pageUpload(args, extra)
  }

  // Fire compose:page trigger
  async pageTriggerScript(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    return this.ComposeAPI.pageTriggerScript(args, extra)
  }

  // List page translation
  async pageListTranslations(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    return this.ComposeAPI.pageListTranslations(args, extra)
  }

  // Update page translation
  async pageUpdateTranslations(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    return this.ComposeAPI.pageUpdateTranslations(args, extra)
  }

  // Update icon for page
  async pageUpdateIcon(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    return this.ComposeAPI.pageUpdateIcon(args, extra)
  }

  // List icons
  // Mirrors ComposeAPI.iconList.
  findIcons(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.iconList(args, extra)
  }

  // Upload icon
  iconUpload(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.iconUpload(args, extra)
  }

  // Delete icon
  // Mirrors ComposeAPI.iconDelete.
  async deleteIcon(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'iconID')
    args.iconID = await this.resolveRef('icon', args.iconID, args)
    return this.ComposeAPI.iconDelete(args, extra)
  }

  // List available page layouts
  async pageLayoutListNamespace(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.pageLayoutListNamespace(args, extra)
  }

  // List available page layouts
  // Mirrors ComposeAPI.pageLayoutList.
  async findPageLayouts(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<PageLayout[], KV>> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    return this.ComposeAPI.pageLayoutList(args, extra).then(r => castSet(r, v => new PageLayout(v)))
  }

  // Create page layout
  // Mirrors ComposeAPI.pageLayoutCreate.
  async createPageLayout(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<PageLayout> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    return this.ComposeAPI.pageLayoutCreate(args, extra).then(r => new PageLayout(r))
  }

  // Get page details
  // Mirrors ComposeAPI.pageLayoutRead.
  async findPageLayoutByID(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<PageLayout> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    args.pageLayoutID = await this.resolveRef('pageLayout', args.pageLayoutID, args)
    return this.ComposeAPI.pageLayoutRead(args, extra).then(r => new PageLayout(r))
  }

  // Update page
  // Mirrors ComposeAPI.pageLayoutUpdate.
  async updatePageLayout(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<PageLayout> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    args.pageLayoutID = await this.resolveRef('pageLayout', args.pageLayoutID, args)
    return this.ComposeAPI.pageLayoutUpdate(args, extra).then(r => new PageLayout(r))
  }

  // Reorder page layouts
  async pageLayoutReorder(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    return this.ComposeAPI.pageLayoutReorder(args, extra)
  }

  // Delete page layout
  // Mirrors ComposeAPI.pageLayoutDelete.
  async deletePageLayout(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    args.pageLayoutID = await this.resolveRef('pageLayout', args.pageLayoutID, args)
    return this.ComposeAPI.pageLayoutDelete(args, extra)
  }

  // Undelete soft deleted Delete page layout
  // Mirrors ComposeAPI.pageLayoutUndelete.
  async undeletePageLayout(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    args.pageLayoutID = await this.resolveRef('pageLayout', args.pageLayoutID, args)
    return this.ComposeAPI.pageLayoutUndelete(args, extra)
  }

  // List page layout translation
  async pageLayoutListTranslations(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    args.pageLayoutID = await this.resolveRef('pageLayout', args.pageLayoutID, args)
    return this.ComposeAPI.pageLayoutListTranslations(args, extra)
  }

  // Update page layout translation
  async pageLayoutUpdateTranslations(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.pageID = await this.resolveRef('page', args.pageID, args)
    args.pageLayoutID = await this.resolveRef('pageLayout', args.pageLayoutID, args)
    return this.ComposeAPI.pageLayoutUpdateTranslations(args, extra)
  }

  // Create module
  // Mirrors ComposeAPI.moduleCreate.
  async createModule(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Module> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.moduleCreate(args, extra).then(r => new Module(r))
  }

  // Update module
  // Mirrors ComposeAPI.moduleUpdate.
  async updateModule(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Module> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.moduleUpdate(args, extra).then(r => new Module(r))
  }

  // Delete module
  // Mirrors ComposeAPI.moduleDelete.
  async deleteModule(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.moduleDelete(args, extra)
  }

  // Fire compose:module trigger
  async moduleTriggerScript(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.moduleTriggerScript(args, extra)
  }

  // List moudle translation
  async moduleListTranslations(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.moduleListTranslations(args, extra)
  }

  // Update module translation
  async moduleUpdateTranslations(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.moduleUpdateTranslations(args, extra)
  }

  // Generates report from module records
  async recordReport(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordReport(args, extra)
  }

  // Initiate record import session
  async recordImportInit(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordImportInit(args, extra)
  }

  // Run record import
  async recordImportRun(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordImportRun(args, extra)
  }

  // Get import progress
  async recordImportProgress(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordImportProgress(args, extra)
  }

  // Exports records that match
  async recordExport(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordExport(args, extra)
  }

  // Executes server-side procedure over one or more module records
  async recordExec(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordExec(args, extra)
  }

  // Create record in module section
  // Mirrors ComposeAPI.recordCreate.
  async createRecord(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordCreate(args, extra)
  }

  // Update records in module section
  // Mirrors ComposeAPI.recordUpdate.
  async updateRecord(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    args.recordID = await this.resolveRef('record', args.recordID, args)
    return this.ComposeAPI.recordUpdate(args, extra)
  }

  // Partially update record values
  async recordPatch(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordPatch(args, extra)
  }

  // Delete record row from module section
  async recordBulkDelete(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordBulkDelete(args, extra)
  }

  // Undelete soft-deleted record from module section
  // Mirrors ComposeAPI.recordUndelete.
  async undeleteRecord(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    args.recordID = await this.resolveRef('record', args.recordID, args)
    return this.ComposeAPI.recordUndelete(args, extra)
  }

  // Undelete soft-deleted records from module section
  async recordBulkUndelete(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordBulkUndelete(args, extra)
  }

  // Uploads attachment and validates it against record field requirements
  async recordUpload(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordUpload(args, extra)
  }

  // Fire compose:record trigger
  async recordTriggerScript(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    args.recordID = await this.resolveRef('record', args.recordID, args)
    return this.ComposeAPI.recordTriggerScript(args, extra)
  }

  // Fire compose:record trigger
  async recordTriggerScriptOnList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    return this.ComposeAPI.recordTriggerScriptOnList(args, extra)
  }

  // List record revisions
  async recordRevisions(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.moduleID = await this.resolveRef('module', args.moduleID, args)
    args.recordID = await this.resolveRef('record', args.recordID, args)
    return this.ComposeAPI.recordRevisions(args, extra)
  }

  // List records for data privacy
  dataPrivacyRecordList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.dataPrivacyRecordList(args, extra)
  }

  // List modules
  dataPrivacyModuleList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.dataPrivacyModuleList(args, extra)
  }

  // List/read charts
  // Mirrors ComposeAPI.chartList.
  async findCharts(
    a: IdArgs = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<Chart[], KV>> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.chartList(args, extra).then(r => castSet(r, v => new Chart(v)))
  }

  // List/read charts
  // Mirrors ComposeAPI.chartCreate.
  async createChart(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Chart> {
    const args = argsOf(a, 'namespaceID')
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.chartCreate(args, extra).then(r => new Chart(r))
  }

  // Read charts by ID
  // Mirrors ComposeAPI.chartRead.
  async findChartByID(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Chart> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.chartID = await this.resolveRef('chart', args.chartID, args)
    return this.ComposeAPI.chartRead(args, extra).then(r => new Chart(r))
  }

  // Add/update charts
  // Mirrors ComposeAPI.chartUpdate.
  async updateChart(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Chart> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.chartID = await this.resolveRef('chart', args.chartID, args)
    return this.ComposeAPI.chartUpdate(args, extra).then(r => new Chart(r))
  }

  // Delete chart
  // Mirrors ComposeAPI.chartDelete.
  async deleteChart(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.chartID = await this.resolveRef('chart', args.chartID, args)
    return this.ComposeAPI.chartDelete(args, extra)
  }

  // List chart translation
  async chartListTranslations(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.chartID = await this.resolveRef('chart', args.chartID, args)
    return this.ComposeAPI.chartListTranslations(args, extra)
  }

  // Update chart translation
  async chartUpdateTranslations(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.chartID = await this.resolveRef('chart', args.chartID, args)
    return this.ComposeAPI.chartUpdateTranslations(args, extra)
  }

  // Send email from the Compose
  notificationEmailSend(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.notificationEmailSend(args, extra)
  }

  // List, filter all page attachments
  // Mirrors ComposeAPI.attachmentList.
  async findAttachments(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<Attachment[], KV>> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    return this.ComposeAPI.attachmentList(args, extra).then(r => castSet(r, v => new Attachment(v)))
  }

  // Delete attachment
  // Mirrors ComposeAPI.attachmentDelete.
  async deleteAttachment(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.attachmentID = await this.resolveRef('attachment', args.attachmentID, args)
    return this.ComposeAPI.attachmentDelete(args, extra)
  }

  // Serves attached file
  async attachmentOriginal(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.attachmentID = await this.resolveRef('attachment', args.attachmentID, args)
    return this.ComposeAPI.attachmentOriginal(args, extra)
  }

  // Serves preview of an attached file
  async attachmentPreview(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.namespaceID = await this.resolveRef('namespace', args.namespaceID, args)
    args.attachmentID = await this.resolveRef('attachment', args.attachmentID, args)
    return this.ComposeAPI.attachmentPreview(args, extra)
  }

  // Retrieve defined permissions
  permissionsList(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.ComposeAPI.permissionsList(extra)
  }

  // Effective rules for current user
  permissionsEffective(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.permissionsEffective(args, extra)
  }

  // Evaluate rules for given user/role combo
  permissionsTrace(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.permissionsTrace(args, extra)
  }

  // Retrieve role permissions
  permissionsRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    return this.ComposeAPI.permissionsRead(args, extra)
  }

  // Remove all defined role permissions
  permissionsDelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    return this.ComposeAPI.permissionsDelete(args, extra)
  }

  // Update permission settings
  permissionsUpdate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    return this.ComposeAPI.permissionsUpdate(args, extra)
  }

  // List all available automation scripts for compose resources
  automationList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.automationList(args, extra)
  }

  // Serves client scripts bundle
  automationBundle(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.automationBundle(args, extra)
  }

  // Triggers execution of a specific script on a system service level
  automationTriggerScript(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.ComposeAPI.automationTriggerScript(args, extra)
  }
}
