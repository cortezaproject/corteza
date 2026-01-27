import { ListResponse, PermissionResource, PermissionRole } from './shared'
import { Attachment } from '../../shared'
import { Compose as ComposeAPI } from '../../api-clients'
import { Namespace, Record, Module, Page } from '../../compose'
import { Values } from '../../compose/types/record'
interface ComposeContext {
    ComposeAPI: ComposeAPI;
    $namespace?: Namespace;
    $module?: Module;
    $record?: Record;
}
interface PageListFilter {
    [key: string]: string | number | {
        [key: string]: string;
    } | undefined;
    namespaceID?: string;
    selfID?: string;
    query?: string;
    handle?: string;
    labels?: {
        [key: string]: string;
    };
    limit?: number;
    pageCursor?: string;
    sort?: string;
}
interface RecordListFilter {
    [key: string]: string | number | {
        [key: string]: string;
    } | undefined;
    namespaceID?: string;
    moduleID?: string;
    query?: string;
    labels?: {
        [key: string]: string;
    };
    limit?: number;
    pageCursor?: string;
    sort?: string;
}
interface ModuleListFilter {
    [key: string]: string | number | {
        [key: string]: string;
    } | undefined;
    namespaceID?: string;
    query?: string;
    name?: string;
    handle?: string;
    labels?: {
        [key: string]: string;
    };
    limit?: number;
    pageCursor?: string;
    sort?: string;
}
interface NamespaceListFilter {
    [key: string]: string | number | {
        [key: string]: string;
    } | undefined;
    query?: string;
    slug?: string;
    labels?: {
        [key: string]: string;
    };
    limit?: number;
    pageCursor?: string;
    sort?: string;
}
/**
 * ComposeHelper provides layer over Compose API and utilities that simplify automation script writing
 *
 * Initiated as Compose object and provides a few handy shortcuts and fallback that will enable you
 * to rapidly develop your automation scripts.
 */
export default class ComposeHelper {
  readonly ComposeAPI: ComposeAPI
  readonly $namespace?: Namespace
  readonly $module?: Module
  readonly $record?: Record
  /**
     * @param ctx.$namespace - Current namespace
     * @param ctx.$module - Current module
     * @param ctx.$record - Current record
     */
  constructor(ctx: ComposeContext);
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
  makePage(values?: Partial<Page>, ns?: Namespace | undefined): Promise<Page>;
  /**
     * Creates/updates Page
     *
     * @param page
     */
  savePage(page: Promise<Page> | Page | Partial<Page>): Promise<Page>;
  /**
     * Deletes a page
     *
     * @example
     * Compose.deletePage(myPage)
     *
     * @param page
     */
  deletePage(page: Page): Promise<unknown>;
  /**
     * Searches for pages
     *
     * @private
     * @param filter
     * @param ns
     */
  findPages(filter?: undefined | string | PageListFilter, ns?: Namespace | undefined): Promise<ListResponse<Page[], PageListFilter>>;
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
  findPageByID(page: string | Page, ns?: Namespace | undefined): Promise<Page>;
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
  makeRecord(values?: Values, module?: Module | null): Promise<Record>;
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
  saveRecord(record: Record | Promise<Record>): Promise<Record>;
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
  deleteRecord(record: Record): Promise<unknown>;
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
  findRecords(filter?: string | RecordListFilter, module?: Module | undefined): Promise<ListResponse<Record[], RecordListFilter>>;
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
  findLastRecord(module?: Module | undefined): Promise<Record>;
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
  findFirstRecord(module?: Module | undefined): Promise<Record>;
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
  findRecordByID(record: string | object | Record, module?: Module | null): Promise<Record>;
  /**
     * Finds a single attachment
     *
     * @param attachment Attachment to find
     * @param ns
     */
  findAttachmentByID(attachment: string | object | Attachment, ns?: Namespace | undefined): Promise<Attachment>;
  /**
     * Helper to determine field's name from it's label
     * @param label Field's label
     */
  moduleFieldNameFromLabel(label: string): string;
  /**
     * Creates new Module object
     *
     * @param module
     * @param ns, defaults to current $namespace
     */
  makeModule(module?: Promise<Module> | Module | Partial<Module>, ns?: Namespace | undefined): Promise<Module>;
  /**
     * Creates/updates Module
     *
     * @param module
     */
  saveModule(module: Promise<Module> | Module): Promise<Module>;
  /**
     * Searches for modules
     *
     * @private
     * @param filter
     * @param ns
     */
  findModules(filter?: string | ModuleListFilter, ns?: Namespace | undefined): Promise<ListResponse<Module[], ModuleListFilter>>;
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
  findModuleByID(module: string | Module | Record, ns?: Namespace | undefined): Promise<Module>;
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
  findModuleByName(name: string, ns?: string | Namespace | object | undefined): Promise<Module>;
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
  findModuleByHandle(handle: string, ns?: string | Namespace | object | undefined): Promise<Module>;
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
  makeNamespace(namespace?: Promise<Namespace> | Namespace | Partial<Namespace>): Promise<Namespace>;
  /**
     * Creates/updates Namespace
     *
     * @example
     * Compose.saveNamespace(myNamespace)
     *
     * @param namespace
     */
  saveNamespace(namespace: Promise<Namespace> | Namespace): Promise<Namespace>;
  /**
     * Searches for namespaces
     *
     * @private
     * @param filter
     */
  findNamespaces(filter?: string | NamespaceListFilter): Promise<ListResponse<Namespace[], NamespaceListFilter>>;
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
  findNamespaceByID(ns?: string | Namespace | Record | undefined): Promise<Namespace>;
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
  findNamespaceBySlug(slug: string): Promise<Namespace>;
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
  sendMail(to: string | string[], subject: string, { html }?: {
        html?: string;
    }, { cc }?: {
        cc?: string | string[];
    }): Promise<unknown>;
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
  sendRecordToMail(to: string | string[], subject?: string, { header, footer, style, fields, ...mailHeader }?: {
        header?: string;
        footer?: string;
        style?: string;
        fields?: string[] | null;
    }, record?: Promise<Record> | Record | undefined): Promise<unknown>;
    /**
     * Walks over white listed fields.
     *
     * @param fwl - field white list; if not defined, all fields are used
     * @param record - record to be walked over
     * @param formatter
     *
     * @private
     */
  walkFields(fwl: null | string[] | Record | undefined, record: Record, formatter: (...args: unknown[]) => string): Array<string>;
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
  recordToHTML(fwl?: null | string[] | Record, record?: Record | undefined): string;
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
  recordToPlainText(fwl?: null | string[] | Record, record?: Record | undefined): string;
  /**
     * Scans all given arguments and returns first one that resembles something like a valid module, its name or ID
     *
     * @private
     */
  resolveModule(...args: unknown[]): Promise<Module>;
  /**
     * Scans all given arguments and returns first one that resembles something like a valid namespace, its slug or ID
     *
     * @private
     */
  resolveNamespace(...args: unknown[]): Promise<Namespace>;
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
  allow(...pr: {
        role: PermissionRole;
        resource: PermissionResource;
        operation: string;
    }[]): Promise<void>;
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
  deny(...pr: {
        role: PermissionRole;
        resource: PermissionResource;
        operation: string;
    }[]): Promise<void>;
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
  inherit(...pr: {
        role: PermissionRole;
        resource: PermissionResource;
        operation: string;
    }[]): Promise<void>;
}
export {}
