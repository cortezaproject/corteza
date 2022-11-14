import { extractID, genericPermissionUpdater, isFresh, kv, ListResponse, PermissionResource, PermissionRole } from './shared'
import { Attachment } from '../../shared'
import { Compose as ComposeAPI } from '../../api-clients'
import { Namespace, Record, Module, Page } from '../../compose'
import { Values } from '../../compose/types/record'
import { IsCortezaID } from '../../cast'
import { IsOf } from '../../guards'

const emailStyle = `
body { -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; color: #3A393C; font-family: Verdana,Arial,sans-serif; font-size: 14px; height: 100%; margin: 0; padding: 0; width: 100% !important; }
table { margin: 20px auto; background: #ffffff; border-collapse: collapse; max-width: 100%; }
table tr { height: 40px; }
table td { padding-top: 10px; padding-left: 20px; width:100%; max-width:100%; min-width:100%; width:100%; vertical-align: top; }
table tbody { border-top: 3px solid #808080; }
tbody tr:nth-child(even) { background-color: #F3F3F4; }
tbody td:first-child { width: 30%; color: #808080; }
tbody td:nth-child(2) { width: 70%; }
h2, p { padding: 10px 20px; }
p { text-align: justify; line-height: 1.4;}
`

interface ComposeContext {
  ComposeAPI: ComposeAPI;
  $namespace?: Namespace;
  $module?: Module;
  $record?: Record;
}

interface PageListFilter {
  [key: string]: string|number|{[key: string]: string}|undefined;
  namespaceID?: string;
  selfID?: string;
  query?: string;
  handle?: string;
  labels?: {[key: string]: string};
  limit?: number;
  pageCursor?: string;
  sort?: string;
}

interface RecordListFilter {
  [key: string]: string|number|{[key: string]: string}|undefined;
  namespaceID?: string;
  moduleID?: string;
  query?: string;
  labels?: {[key: string]: string};
  limit?: number;
  pageCursor?: string;
  sort?: string;
}

interface ModuleListFilter {
  [key: string]: string|number|{[key: string]: string}|undefined;
  namespaceID?: string;
  query?: string;
  name?: string;
  handle?: string;
  labels?: {[key: string]: string};
  limit?: number;
  pageCursor?: string;
  sort?: string;
}

interface NamespaceListFilter {
  [key: string]: string|number|{[key: string]: string}|undefined;
  query?: string;
  slug?: string;
  labels?: {[key: string]: string};
  limit?: number;
  pageCursor?: string;
  sort?: string;
}

/**
 * Helpers to determine if specific object looks like the type we are interested in.
 * It does not rely on instanceof, because of bundling issues.
 */
function isRecord (o: any) {
  return o && !!o.recordID && o.moduleID && o.namespaceID
}
function isModule (o: any) {
  return o && !!o.moduleID && o.namespaceID
}
function isPage (o: any) {
  return o && !!o.pageID && o.namespaceID
}

/**
 * ComposeHelper provides layer over Compose API and utilities that simplify automation script writing
 *
 * Initiated as Compose object and provides a few handy shortcuts and fallback that will enable you
 * to rapidly develop your automation scripts.
 */
export default class ComposeHelper {
  readonly ComposeAPI: ComposeAPI;
  readonly $namespace?: Namespace;
  readonly $module?: Module;
  readonly $record?: Record;

  /**
   * @param ctx.$namespace - Current namespace
   * @param ctx.$module - Current module
   * @param ctx.$record - Current record
   */
  constructor (ctx: ComposeContext) {
    this.ComposeAPI = ctx.ComposeAPI
    this.$namespace = ctx.$namespace
    this.$module = ctx.$module
    this.$record = ctx.$record
  }

  /**
   * Creates new Page object
   *
   * <p>
   *   Created page is "in-memory" only. To store it, use savePage() method
   * </p>
   *
   * @example
   * // Simple page creation new page on current namespace
   * let myPage = await Compose.makePage({ title: 'My Amazing Page!' })
   *
   * @param values
   * @param ns - defaults to current $namespace
   */
  async makePage (values: Partial<Page> = {}, ns: Namespace|undefined = this.$namespace): Promise<Page> {
    return this.resolveNamespace(ns).then(ns => {
      return new Page({ ...values, namespaceID: ns.namespaceID })
    })
  }

  /**
   * Creates/updates Page
   *
   * @param page
   */
  async savePage (page: Promise<Page>|Page|Partial<Page>): Promise<Page> {
    return Promise.resolve(page).then(page => {
      if (!isPage(page)) {
        throw Error('expecting Page type')
      }

      if (page.pageID && isFresh(page.pageID)) {
        return this.ComposeAPI.pageCreate(kv(page)).then(page => new Page(page))
      } else {
        return this.ComposeAPI.pageUpdate(kv(page)).then(page => new Page(page))
      }
    })
  }

  /**
   * Deletes a page
   *
   * @example
   * Compose.deletePage(myPage)
   *
   * @param page
   */
  async deletePage (page: Page): Promise<unknown> {
    return Promise.resolve(page).then(page => {
      if (!isPage(page)) {
        throw Error('expecting Page type')
      }

      if (!isFresh(page.pageID)) {
        return this.ComposeAPI.pageDelete(kv(page))
      }
    })
  }

  /**
   * Searches for pages
   *
   * @private
   * @param filter
   * @param ns
   */
  async findPages (filter: undefined|string|PageListFilter = {}, ns: Namespace|undefined = this.$namespace): Promise<ListResponse<Page[], PageListFilter>> {
    if (typeof filter === 'string') {
      filter = { query: filter }
    }

    return this.resolveNamespace(ns).then(ns => {
      const namespaceID = extractID(ns, 'namespaceID')
      return this.ComposeAPI.pageList({ namespaceID, ...filter as object }).then(res => {
        // Casting all we got to to Page
        res.set = (res.set as any[]).map(m => new Page(m))
        return res as unknown as ListResponse<Page[], PageListFilter>
      })
    })
  }

  /**
   * Finds page by ID
   *
   * @example
   * // Explicitly load page and do something with it
   * Compose.finePageByID('2039248239042').then(myPage => {
   *   // do something with myPage
   *   myPage.title = 'My More Amazing Page!'
   *   return myPage
   * }).then(Compose.savePage)
   *
   * @param page - accepts Page, pageID (when string string)
   * @param ns - namespace, defaults to current $namespace
   */
  async findPageByID (page: string|Page, ns: Namespace|undefined = this.$namespace): Promise<Page> {
    return this.resolveNamespace(ns).then((ns) => {
      const pageID = extractID(page, 'pageID')
      const namespaceID = extractID(ns, 'namespaceID')

      return this.ComposeAPI.pageRead({ namespaceID, pageID }).then(page => new Page(page))
    })
  }

  /**
   * Creates new Record object
   *
   * <p>
   *   Created record is "in-memory" only. To store it, use saveRecord() method
   * </p>
   *
   * @example
   * // Simple record creation (new record of current module - $module)
   * let myLead = await Compose.makeRecord()
   * myLead.values.Title = 'My Lead Title'
   *
   * // Create record of type Lead and copy values from another Record
   * // This will copy only values that have the same name in both modules
   * let myLead = await Compose.makeRecord(myContact, 'Lead')
   *
   * // Or use promises:
   * Compose.makeRecord(myContact, 'Lead').then(myLead => {
   *   myLead.values.Title = 'My Lead Title'
   *
   *   // ...
   *
   *   // return record when finished
   *   return myLead
   * }).catch(err => {
   *   // solve the problem
   *   console.error(err)
   * })
   *
   * @param values
   * @param module - defaults to current $module
   */
  async makeRecord (values: Values = {}, module: Module|null = null): Promise<Record> {
    return this.resolveModule(module, this.$module).then(module => {
      const record = new Record(module)

      // Set record values
      record.setValues(values)

      return record
    })
  }

  /**
   * Saves a record
   *
   * Please note that there is no need to explicitly save (current record) on before/after events,
   * internal systems take care of that.
   *
   * @example
   * // Quick example how to make and save new Lead:
   * let mySavedLead = await Compose.saveRecord(Compose.makeRecord({Title: 'Lead title'}, 'Lead'))
   * if (mySavedLead) {
   *   console.log('Record saved, new ID', mySavedLead.recordID)
   * } else {
   *   // solve the problem
   *   console.error(err)
   * }
   *
   * // Or with promises:
   * Compose.makeRecord({Title: 'Lead title'}, 'Lead')).then(myLead => {
   *   return Compose.saveRecord(myLead)
   * }).then(mySavedLead => {
   *   console.log('Record saved, new ID', mySavedLead.recordID)
   * }).catch(err => {
   *   // solve the problem
   *   console.error(err)
   * })
   *
   * @param record
   */
  async saveRecord (record: Record|Promise<Record>): Promise<Record> {
    return Promise.resolve(record).then(record => {
      if (!isRecord(record)) {
        throw Error('expecting Record type')
      }

      if (isFresh(record.recordID)) {
        return this.ComposeAPI.recordCreate(kv(record)).then(r => new Record(record.module, r))
      } else {
        return this.ComposeAPI.recordUpdate(kv(record)).then(r => new Record(record.module, r))
      }
    })
  }

  /**
   * Deletes a record
   *
   * Please note that there is no need to explicitly delete (current record) on before/after events.
   *
   * @example
   * Compose.deleteRecord(myLead)
   *
   * @param record
   */
  async deleteRecord (record: Record): Promise<unknown> {
    return Promise.resolve(record).then(record => {
      if (!isRecord(record)) {
        throw Error('expecting Record type')
      }

      if (!isFresh(record.recordID)) {
        return this.ComposeAPI.recordDelete(kv(record))
      }
    })
  }

  /**
   * Searches for records of a specific record
   *
   * @example
   * // Find all records (of the current module)
   * Compose.findRecords()
   *
   * // Find Projects where ROI is more than 15%
   * // (assuming we have Project module with netProfit and totalInvestment numeric fields)
   * Compose.findRecords('netProfit / totalInvestment > 0.15', 'Project')
   *
   * // Find Projects where ROI is more than 15%
   * // (assuming we have Project module with netProfit and totalInvestment numeric fields)
   * Compose.findRecords('netProfit / totalInvestment > 0.15', 'Project')
   *
   * // More complex query with sorting:
   * // Returns top 5 Projects with more than 15% ROI in the last year
   * Compose.findRecords({
   *   filter: '(netProfit / totalInvestment > 0.15) AND (YEAR(createdAt) = YEAR(NOW()) - 1)'
   *   sort: 'netProfit / totalInvestment DESC',
   *   limit: 5,
   * }, 'Project')
   *
   * // Accessing returned records
   * Compose.findRecords().then(({ set, filter }) => {
   *    // set: array of records
   *    // filter: object with filter specs
   *
   *    Use internal Array functions
   *    set.forEach(r => {
   *      // r, one of the records each iteration
   *    })
   *
   *    // Or standard for-loop
   *    for (let r of set) {
   *       // r...
   *    }
   * })
   *
   * @param filter - filter object (or filtering conditions when string)
   * @property {string} filter.query - filtering conditions
   * @property {string} filter.sort - sorting rules
   * @property {number} filter.limit - number of max returned records
   * @property {number} filter.pageCursor - hashed string that retrieves a specific page
   * @param [module] - if not set, defaults to $module
   */
  async findRecords (filter: string|RecordListFilter = '', module: Module|undefined = this.$module): Promise<ListResponse<Record[], RecordListFilter>> {
    return this.resolveModule(module).then(module => {
      const { moduleID, namespaceID } = module

      let params = {
        moduleID,
        namespaceID,
      } as { moduleID: string; namespaceID: string; query: string}

      if (typeof filter === 'string') {
        params.query = filter
      } else if (typeof filter === 'object') {
        params = { ...params, ...filter }
      }

      return this.ComposeAPI.recordList(params).then(res => {
        // Casting all we got to to Record
        res.set = (res.set as any[]).map(record => new Record(module, record))
        return res as unknown as ListResponse<Record[], RecordListFilter>
      })
    })
  }

  /**
   * Finds last (created) record in the module
   *
   * @example
   * Compose.findLastRecord('Settings').then(lastSettingRecord => {
   *   // handle lastSettingRecord
   * })
   *
   * @param module
   */
  async findLastRecord (module: Module|undefined = this.$module): Promise<Record> {
    return this.findRecords({ sort: 'createdAt DESC', limit: 1 }, module).then(res => {
      if (!Array.isArray(res.set) || res.set.length === 0) {
        throw new Error('records not found')
      }

      return res.set[0]
    })
  }

  /**
   * Finds first (created) record in the module
   *
   * @example
   * Compose.findFirstRecord('Settings').then(firstSettingRecord => {
   *   // handle this firstSettingRecord
   * })
   *
   * @param module
   */
  async findFirstRecord (module: Module|undefined = this.$module): Promise<Record> {
    return this.findRecords({ sort: 'createdAt', limit: 1 }, module).then(res => {
      if (!Array.isArray(res.set) || res.set.length === 0) {
        throw new Error('records not found')
      }

      return res.set[0]
    })
  }

  /**
   * Finds one record by ID
   *
   * @example
   * Compose.findRecordByID("23957823957").then(specificRecord => {
   *   // handle this specificRecord
   * })
   *
   * @param record
   * @param module
   */
  async findRecordByID (record: string|object|Record, module: Module|null = null): Promise<Record> {
    // We're handling module default a bit differently here
    // because we want to allow users to use record's module
    return this.resolveModule(module, (record as Record || {}).module, this.$module).then((module) => {
      const { moduleID, namespaceID } = module
      return this.ComposeAPI.recordRead({
        moduleID,
        namespaceID,
        recordID: extractID(record, 'recordID'),
      }).then(r => new Record(module, r))
    })
  }

  /**
   * Finds a single attachment
   *
   * @param attachment Attachment to find
   * @param ns
   */
  async findAttachmentByID (attachment: string|object|Attachment, ns: Namespace|undefined = this.$namespace): Promise<Attachment> {
    return this.resolveNamespace(ns).then(namespace => {
      const { namespaceID } = namespace
      return this.ComposeAPI.attachmentRead({
        kind: 'original',
        attachmentID: extractID(attachment, 'attachmentID'),
        namespaceID,
      }).then(att => new Attachment(att))
    })
  }

  /**
   * Helper to determine field's name from it's label
   * @param label Field's label
   */
  moduleFieldNameFromLabel (label: string): string {
    return label.split(/[^a-zA-Z0-9_]/g)
      .filter(p => !!p)
      .map(p => `${p[0].toUpperCase()}${p.slice(1)}`)
      .join('')
  }

  /**
   * Creates new Module object
   *
   * @param module
   * @param ns, defaults to current $namespace
   */
  async makeModule (module: Promise<Module>|Module|Partial<Module> = {} as Module, ns: Namespace|undefined = this.$namespace): Promise<Module> {
    return this.resolveNamespace(ns).then((ns) => {
      return new Module({ ...module, namespaceID: ns.namespaceID })
    })
  }

  /**
   * Creates/updates Module
   *
   * @param module
   */
  async saveModule (module: Promise<Module>|Module): Promise<Module> {
    return Promise.resolve(module).then(module => {
      if (!isModule(module)) {
        throw new Error('expecting Module type')
      }

      if (isFresh(module.moduleID)) {
        return this.ComposeAPI.moduleCreate(kv(module)).then(m => new Module(m))
      } else {
        return this.ComposeAPI.moduleUpdate(kv(module)).then(m => new Module(m))
      }
    })
  }

  /**
   * Searches for modules
   *
   * @private
   * @param filter
   * @param ns
   */
  async findModules (filter: string|ModuleListFilter = '', ns: Namespace|undefined = this.$namespace): Promise<ListResponse<Module[], ModuleListFilter>> {
    if (typeof filter === 'string') {
      filter = { query: filter }
    }

    return this.resolveNamespace(ns).then((ns) => {
      const namespaceID = extractID(ns, 'namespaceID')

      return this.ComposeAPI.moduleList({ namespaceID, ...filter as object }).then(res => {
        // Casting all we got to to Module
        res.set = (res.set as any[]).map(m => new Module(m))
        return res as unknown as ListResponse<Module[], ModuleListFilter>
      })
    })
  }

  /**
   * Finds module by ID
   *
   * @example
   * // Explicitly load module and do something with it
   * Compose.findModuleByID('2039248239042').then(myModule => {
   *   // do something with myModule
   *   return Compose.findLastRecord(myModule)
   * }).then((lastRecord) => {})
   *
   * // or
   * Compose.findLastRecord(Compose.findModuleByID('2039248239042')).then(....)
   *
   * // even shorter
   * Compose.findLastRecord('2039248239042').then(....)
   *
   * @param module - accepts Module, moduleID (when string) or Record
   * @param ns - namespace, defaults to current $namespace
   */
  async findModuleByID (module: string|Module|Record, ns: Namespace|undefined = this.$namespace): Promise<Module> {
    return this.resolveNamespace(ns).then((ns) => {
      const moduleID = extractID(module, 'moduleID')
      const namespaceID = extractID(ns, 'namespaceID')

      return this.ComposeAPI.moduleRead({ namespaceID, moduleID }).then(m => new Module(m))
    })
  }

  /**
   * Finds module by name
   *
   * @example
   * // Explicitly load module and do something with it
   * Compose.findModuleByName('SomeModule').then(myModule => {
   *   // do something with myModule
   *   return Compose.findLastRecord(myModule)
   * }).then((lastRecord) => {})
   *
   * // or
   * Compose.findLastRecord(Compose.findModuleByName('SomeModule')).then(....)
   *
   * // even shorter
   * Compose.findLastRecord('SomeModule').then(....)
   *
   * @param name - name of the module
   * @param ns - defaults to current $namespace
   */
  async findModuleByName (name: string, ns: string|Namespace|object|undefined = this.$namespace): Promise<Module> {
    return this.resolveNamespace(ns).then((ns) => {
      const namespaceID = extractID(ns, 'namespaceID')
      return this.ComposeAPI.moduleList({ namespaceID, name }).then(res => {
        if (!Array.isArray(res.set) || res.set.length === 0) {
          throw new Error('module not found')
        }

        return new Module(res.set[0])
      })
    })
  }

  /**
   * Finds module by handle
   *
   * @example
   * // Explicitly load module and do something with it
   * Compose.findModuleByHandle('SomeModule').then(myModule => {
   *   // do something with myModule
   *   return Compose.findLastRecord(myModule)
   * }).then((lastRecord) => {})
   *
   * // or
   * Compose.findLastRecord(Compose.findModuleByHandle('SomeModule')).then(....)
   *
   * // even shorter
   * Compose.findLastRecord('SomeModule').then(....)
   *
   * @param handle - handle of the module
   * @param ns - defaults to current $namespace
   */
  async findModuleByHandle (handle: string, ns: string|Namespace|object|undefined = this.$namespace): Promise<Module> {
    return this.resolveNamespace(ns).then((ns) => {
      const namespaceID = extractID(ns, 'namespaceID')
      return this.ComposeAPI.moduleList({ namespaceID, handle }).then(res => {
        if (!Array.isArray(res.set) || res.set.length === 0) {
          throw new Error('module not found')
        }

        return new Module(res.set[0])
      })
    })
  }

  /**
   * Creates new Namespace object
   *
   * @example
   * // Creates enabled (!) namespace with slug & name
   * Compose.saveNamespace(Compose.makeNamespace({
   *   slug: 'my-namespace',
   *   name: 'My Namespace',
   * }))
   *
   * @param namespace
   * @param namespace, defaults to current $namespace
   */
  async makeNamespace (namespace: Promise<Namespace>|Namespace|Partial<Namespace> = {} as Namespace): Promise<Namespace> {
    return new Namespace({
      name: (namespace as Namespace).name || (namespace as Namespace).slug,
      meta: {},
      enabled: true,
      ...namespace,
    })
  }

  /**
   * Creates/updates Namespace
   *
   * @example
   * Compose.saveNamespace(myNamespace)
   *
   * @param namespace
   */
  async saveNamespace (namespace: Promise<Namespace>|Namespace): Promise<Namespace> {
    return Promise.resolve(namespace).then(namespace => {
      if (!(namespace instanceof Namespace)) {
        throw Error('expecting Namespace type')
      }

      if (isFresh(namespace.namespaceID)) {
        return this.ComposeAPI.namespaceCreate(kv(namespace)).then(n => new Namespace(n))
      } else {
        return this.ComposeAPI.namespaceUpdate(kv(namespace)).then(n => new Namespace(n))
      }
    })
  }

  /**
   * Searches for namespaces
   *
   * @private
   * @param filter
   */
  async findNamespaces (filter: string|NamespaceListFilter = ''): Promise<ListResponse<Namespace[], NamespaceListFilter>> {
    if (typeof filter === 'string') {
      filter = { query: filter }
    }

    return this.ComposeAPI.namespaceList({ ...filter }).then(res => {
      // Casting all we got to to Namespace
      res.set = (res.set as any[]).map(m => new Namespace(m))
      return res as unknown as ListResponse<Namespace[], NamespaceListFilter>
    })
  }

  /**
   * Finds namespace by ID
   *
   * @example
   * // Explicitly load namespace and do something with it
   * Compose.findNamespaceByID('2039248239042').then(myNamespace => {
   *   // do something with myNamespace
   *   return Compose.findModules(myNamespace)
   * }).then(modules => {})
   *
   * // even shorter
   * Compose.findModules('2039248239042').then(....)
   *
   * @param ns - accepts Namespace, namespaceID (when string string) or Record
   */
  async findNamespaceByID (ns: string|Namespace|Record|undefined = this.$namespace): Promise<Namespace> {
    const namespaceID = extractID(ns, 'namespaceID')

    return this.ComposeAPI.namespaceRead({ namespaceID }).then(m => new Namespace(m))
  }

  /**
   * Finds namespace by name
   *
   * @example
   * // Explicitly load namespace and do something with it
   * Compose.findNamespaceBySlug('SomeNamespace').then(myNamespace => {
   *   // do something with myNamespace
   *   return Compose.findModules(myNamespace)
   * }).then(modules => {})
   *
   * // even shorter
   * Compose.findModules('SomeNamespace').then(....)
   *
   * @param slug - name of the namespace
   */
  async findNamespaceBySlug (slug: string): Promise<Namespace> {
    return this.ComposeAPI.namespaceList({ slug }).then(res => {
      if (!Array.isArray(res.set) || res.set.length === 0) {
        throw new Error('namespace not found')
      }

      return new Namespace(res.set[0])
    })
  }

  /**
   * Sends a simple email message
   *
   * @example
   * Compose.sendMail('some-address@domain.tld', 'subject...', { html: 'Hello!' })
   *
   * @param to - Recipient(s)
   * @param subject - Mail subject
   * @param body
   * @property {string} body.html - HTML body to be sent
   * @param Any additional addresses we want this to be sent to (carbon-copy)
   */
  async sendMail (to: string|string[], subject: string, { html = '' }: { html?: string } = {}, { cc = [] }: { cc?: string|string[] } = {}): Promise<unknown> {
    if (!to) {
      throw Error('expecting to email address')
    }

    if (!subject) {
      throw Error('expecting subject')
    }

    if (!html) {
      throw Error('expecting HTML body')
    }

    return this.ComposeAPI.notificationEmailSend({
      to: Array.isArray(to) ? to : [to],
      cc: Array.isArray(cc) ? cc : [cc],
      subject,
      content: { html },
    })
  }

  /**
   * Generates HTML with all records fields and sends it to
   *
   * @example
   * // Simplified version, sends current email with generated
   * // subject (<module name> + 'record' +  'update'/'created')
   * Compose.sendRecordToMail('example@domain.tld')
   *
   * // Complex notification with custom subject, header and footer text and custom record
   * Compose.sendRecordToMail(
   *   'asignee@domain.tld',
   *   'New lead assigned to you',
   *   {
   *      header: '<h1>New lead was created and assigned to you</h1>',
   *      footer: 'Review and confirm',
   *      cc: [ 'sales@domain.tld' ],
   *      fields: ['name', 'country', 'amount'],
   *   },
   *   newLead
   * )
   *
   * @param to - Recipient(s)
   * @param subject - Mail subject
   * @param options - Various options for body & email
   * @property {string} options.header - Text (HTML) before the record table
   * @property {string} options.footer - Text (HTML) after the record table
   * @property {string} options.style - Custom CSS styles for the email
   * @param options.fields - List of record fields we want to output
   * @param options.header - Additional mail headers (cc)
   * @param record - record to be converted (or leave for the current $record)
   */
  async sendRecordToMail (
    to: string|string[],
    subject = '',
    {
      header = '',
      footer = '',
      style = emailStyle,
      fields = null,
      ...mailHeader
    }: { header?: string; footer?: string; style?: string; fields?: string[]|null } = {},
    record: Promise<Record>|Record|undefined = this.$record,
  ): Promise<unknown> {
    // Wait for the record if we got a promise

    record = await record

    if (!record) {
      throw Error('record undefined')
    }

    const wb = '<div style="width: 800px; margin: 20px auto;">'
    const wa = '</div>'

    header = `${wb}${header}${wa}`
    footer = `${wb}${footer}${wa}`
    style = `<style type="text/css">${style}</style>`

    const html = style + header + this.recordToHTML(fields, record) + footer

    if (!subject) {
      subject = record.module.name + ' '
      subject += record.updatedAt ? 'record updated' : 'record created'
    }

    return this.sendMail(
      to,
      subject,
      { html },
      { ...mailHeader },
    )
  }

  /**
   * Walks over white listed fields.
   *
   * @param fwl - field white list; if not defined, all fields are used
   * @param record - record to be walked over
   * @param formatter
   *
   * @private
   */
  walkFields (fwl: null|string[]|Record|undefined, record: Record, formatter: (...args: unknown[]) => string): Array<string> {
    if (!formatter) {
      throw new Error('formatter.undefined')
    }

    if (isRecord(fwl)) {
      record = fwl as Record
      fwl = undefined
    }

    if (Array.isArray(fwl) && fwl.length === 0) {
      fwl = null
    }

    return record.module.fields
      .filter(f => !fwl || (fwl as Array<string>).indexOf(f.name) > -1)
      .map(formatter)
  }

  /**
   * Sends a simple record report as HTML
   *
   * @example
   * // generates report for current $record with all fields
   * let report = recordToHTML()
   *
   * // generates report for current $record from a list of fields
   * let report = recordToHTML(['fieldA', 'fieldB', 'fieldC'])
   *
   *
   * @param fwl - field white list (or leave empty/null/false for all fields)
   * @param record - record to be converted (or leave for the current $record)
   */
  recordToHTML (fwl: null|string[]|Record = null, record: Record|undefined = this.$record): string {
    if (!record) {
      throw Error('record undefined')
    }

    const rows = this
      .walkFields(fwl, record, (f): string => {
        const { name, label } = f as { name: string; label: string }
        const v = record.values[name]

        return `<tr><td>${label || name}</td><td>${(Array.isArray(v) ? v : [v]).join(', ') || '&nbsp;'}</td></tr>`
      })
      .join('')

    return `<table width="800" cellspacing="0" cellpadding="0" border="0">${rows}</table>`
  }

  /**
   * Represents a given record as plain text
   *
   * @example
   * // generates report for current $record with all fields
   * let report = recordToPlainText()
   *
   * // generates report for current $record from a list of fields
   * let report = recordToPlainText(['fieldA', 'fieldB', 'fieldC'])
   *
   * @param fwl - field white list (or leave empty/null/false for all fields)
   * @param record - record to be converted (or leave for the current $record)
   */
  recordToPlainText (fwl: null|string[]|Record = null, record: Record|undefined = this.$record): string {
    if (!record) {
      throw Error('record undefined')
    }

    return this
      .walkFields(fwl, record, f => {
        const { name, label } = f as { name: string; label: string }
        const v = record.values[name]
        return `${label || name}:\n${(Array.isArray(v) ? v : [v]).join(', ') || '/'}\n\n`
      })
      .join('')
      .trim()
  }

  /**
   * Scans all given arguments and returns first one that resembles something like a valid module, its name or ID
   *
   * @private
   */
  async resolveModule (...args: unknown[]): Promise<Module> {
    const strResolve = async (module: string): Promise<Module> => {
      return this.findModuleByHandle(module)
        .then(m => {
          if (!m) {
            throw new Error('ModuleNotFound')
          }
          return m
        })
        .catch(() => this.findModuleByName(module))
    }

    for (let module of args) {
      if (!module) {
        continue
      }

      if (typeof module === 'string') {
        if (IsCortezaID(module)) {
          // Looks like an ID
          return this.findModuleByID(module).catch((err = {}) => {
            if (err.message && err.message.indexOf('ModuleNotFound') >= 0) {
              // Not found, let's try if we can find it by slug
              return strResolve(module as string)
            }

            return Promise.reject(err)
          })
        }

        // Assume name
        return strResolve(module)
      }

      if (typeof module !== 'object') {
        continue
      }

      // resolve whatever object we have (maybe it's a promise?)
      // and wait for results
      module = await module

      if (isRecord(module)) {
        const m = module as Record
        return this.resolveModule(m.module, m.moduleID)
      }

      if (IsOf<ListResponse<Module[], ModuleListFilter>>(module, 'set', 'filter')) {
        // We got a result set with modules
        module = module.set
      }

      if (Array.isArray(module)) {
        // We got array of modules
        if (module.length === 0) {
          // Empty array
          continue
        } else {
          // Use first module from the list
          module = module.shift()
        }
      }

      if (!isModule(module)) {
        // not module? is it an object with moduleID & namespaceID?
        if ((module as Module).moduleID === undefined || (module as Module).namespaceID === undefined) {
          break
        }

        return Promise.resolve(new Module(module as Module))
      }

      return Promise.resolve(module as Module)
    }

    return Promise.reject(Error('unexpected input type for module resolver'))
  }

  /**
   * Scans all given arguments and returns first one that resembles something like a valid namespace, its slug or ID
   *
   * @private
   */
  async resolveNamespace (...args: unknown[]): Promise<Namespace> {
    for (let ns of args) {
      if (!ns) {
        continue
      }

      if (typeof ns === 'string') {
        if (IsCortezaID(ns)) {
          // Looks like an ID
          return this.findNamespaceByID(ns).catch((err = {}) => {
            if (err.message && err.message.indexOf('NamespaceNotFound') >= 0) {
              // Not found, let's try if we can find it by slug
              return this.findNamespaceBySlug(ns as string)
            }

            return Promise.reject(err)
          })
        }

        // Assume namespace slug
        return this.findNamespaceBySlug(ns)
      }

      if (typeof ns !== 'object') {
        continue
      }

      // resolve whatever object we have (maybe it's a promise?)
      // and wait for results
      ns = await ns

      if (isRecord(ns)) {
        const n = ns as Record
        return this.resolveNamespace(n.namespaceID)
      }

      if ((ns as ListResponse<Namespace[], NamespaceListFilter>).set && (ns as ListResponse<Namespace[], NamespaceListFilter>).filter) {
        // We got a result set with modules
        ns = (ns as ListResponse<Namespace[], NamespaceListFilter>).set
      }

      if (Array.isArray(ns)) {
        // We got array of modules
        if (ns.length === 0) {
          // Empty array
          continue
        } else {
          // Use first module from the list
          ns = ns.shift()
        }
      }

      if (!(ns instanceof Namespace)) {
        // not Namespace? is it an object with namespaceID?
        if ((ns as Namespace).namespaceID === undefined) {
          break
        }

        return Promise.resolve(new Namespace(ns as Namespace))
      }

      return Promise.resolve(ns)
    }

    return Promise.reject(Error('unexpected input type for namespace resolver'))
  }

  /**
   * Allows access for the given role for the given Compose resource
   *
   * @example
   * // Allows users with `someRole` to access the newly created namespace
   * await Compose.allow({
   *    role: someRole,
   *    resource: newNamespace,
   *    operation: 'read',
   * })
   */
  async allow (...pr: { role: PermissionRole; resource: PermissionResource; operation: string }[]) {
    const rr = pr.map(p => ({
      role: p.role,
      resource: p.resource,
      operation: p.operation,
      access: 'allow',
    }))
    return genericPermissionUpdater(this.ComposeAPI, rr)
  }

  /**
   * Denies access for the given role for the given Compose resource
   *
   * @example
   * // Denies users with `someRole` from accessing the newly created namespace
   * await Compose.deny({
   *    role: someRole,
   *    resource: newNamespace,
   *    operation: 'read',
   * })
   */
  async deny (...pr: { role: PermissionRole; resource: PermissionResource; operation: string }[]) {
    const rr = pr.map(p => ({
      role: p.role,
      resource: p.resource,
      operation: p.operation,
      access: 'deny',
    }))
    return genericPermissionUpdater(this.ComposeAPI, rr)
  }

  /**
   * Inherits access for the given role for the given Compose resource
   *
   * @example
   * // Uses inherited permissions for the `sameRole` for the newly created namespace
   * await Compose.inherit({
   *    role: someRole,
   *    resource: newNamespace,
   *    operation: 'read',
   * })
   */
  async inherit (...pr: { role: PermissionRole; resource: PermissionResource; operation: string }[]) {
    const rr = pr.map(p => ({
      role: p.role,
      resource: p.resource,
      operation: p.operation,
      access: 'inherit',
    }))
    return genericPermissionUpdater(this.ComposeAPI, rr)
  }
}
