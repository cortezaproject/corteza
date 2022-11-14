import { describe, it, beforeEach } from 'mocha'
import { expect } from 'chai'
import { StubbedInstance, stubObject } from 'ts-sinon'
import ComposeHelper from './compose'
import { Module, ModuleFieldNumber, ModuleFieldString, Namespace, Page, Record } from '../../compose'
import { Compose as ComposeAPI } from '../../api-clients'
import { kv } from './shared'
import { NoID } from '../../cast'

const pagePayload = Object.freeze({
  title: 'My Amazing Page',
})

const ns = (h: ComposeHelper): { namespaceID: string } => {
  let namespaceID = NoID
  if (h && h.$namespace) {
    namespaceID = h.$namespace.namespaceID
  }

  return { namespaceID }
}

describe(__filename, () => {
  describe('supporting functions', () => {
    describe('resolving module', () => {
      let h: StubbedInstance<ComposeHelper>

      beforeEach(() => {
        h = stubObject<ComposeHelper>(
          new ComposeHelper({ ComposeAPI: new ComposeAPI({}) }),
          [
            'findModuleByHandle',
            'findModuleByName',
            'findModuleByID',
          ],
        )
      })

      it('should return first valid module', async () => {
        const m = new Module({ moduleID: '1', namespaceID: '2' })

        expect(await h.resolveModule(undefined, null, false, 0, '', m)).to.deep.equal(m)

        expect(h.findModuleByHandle.notCalled, 'findModuleByHandle call not expected').true
        expect(h.findModuleByName.notCalled, 'findModuleByName call not expected').true
        expect(h.findModuleByID.notCalled, 'findModuleByID call not expected').true
      })

      it('should resolve ID', async () => {
        const m = new Module({ moduleID: '444', namespaceID: '555' })

        h.findModuleByID.resolves(m)

        expect(await h.resolveModule(m.moduleID)).to.deep.equal(m)

        expect(h.findModuleByID.calledOnceWith(m.moduleID), 'findModuleByID call expected').true
        expect(h.findModuleByHandle.notCalled, 'findModuleByHandle call not expected').true
        expect(h.findModuleByName.notCalled, 'findModuleByName call not expected').true
      })

      it('should resolve handle', async () => {
        const m = new Module({ handle: 'm-handle' })

        h.findModuleByHandle.resolves(m)

        expect(await h.resolveModule(m.handle)).to.deep.equal(m)

        expect(h.findModuleByID.notCalled, 'findModuleByID call not expected').true
        expect(h.findModuleByHandle.calledOnceWith(m.handle), 'findModuleByHandle call expected').true
        expect(h.findModuleByName.notCalled, 'findModuleByName call not expected').true
      })

      it('should resolve name', async () => {
        const m = new Module({ name: 'm-name' })

        h.findModuleByName.resolves(m)
        h.findModuleByHandle.resolves(undefined)

        expect(await h.resolveModule(m.name)).to.deep.equal(m)

        expect(h.findModuleByID.notCalled, 'findModuleByID call not expected').true
        expect(h.findModuleByHandle.calledOnceWith(m.name), 'findModuleByHandle call expected').true
        expect(h.findModuleByName.calledOnceWith(m.name), 'findModuleByName call expected').true
      })

      it('should resolve numeric name', async () => {
        const m = new Module({ name: '555' })

        h.findModuleByID.rejects(Error('compose.repository.ModuleNotFound'))
        h.findModuleByName.resolves(m)
        h.findModuleByHandle.resolves(undefined)

        expect(await h.resolveModule(m.name)).to.deep.equal(m)

        expect(h.findModuleByID.calledOnceWith(m.name), 'findModuleByID call expected').true
        expect(h.findModuleByHandle.calledOnceWith(m.name), 'findModuleByHandle call expected').true
        expect(h.findModuleByName.calledOnceWith(m.name), 'findModuleByName call expected').true
      })
    })

    describe('resolveNamespace', () => {
      let h: StubbedInstance<ComposeHelper>

      beforeEach(() => {
        h = stubObject<ComposeHelper>(
          new ComposeHelper({ ComposeAPI: new ComposeAPI({}) }),
          [
            'findNamespaceBySlug',
            'findNamespaceByName',
            'findNamespaceByID',
          ],
        )
      })

      it('should return first valid namespace', async () => {
        const ns = new Namespace({ namespaceID: '2' })
        expect(await h.resolveNamespace(undefined, null, false, 0, '', ns)).to.deep.equal(ns)
      })

      it('should resolve ID', async () => {
        const ns = new Namespace({ namespaceID: '555' })

        h.findNamespaceByID.resolves(ns)

        expect(await h.resolveNamespace(ns.namespaceID)).to.deep.equal(ns)

        expect(h.findNamespaceByID.calledOnceWith(ns.namespaceID), 'findNamespaceByID call expected').true
      })

      it('should resolve slug', async () => {
        const ns = new Namespace({ slug: 'ns-slug' })

        h.findNamespaceBySlug.resolves(ns)

        expect(await h.resolveNamespace(ns.slug)).to.deep.equal(ns)

        expect(h.findNamespaceBySlug.calledOnceWith(ns.slug), 'findNamespaceBySlug call expected').true
      })

      it('should resolve numeric slug', async () => {
        const ns = new Namespace({ slug: '555' })

        h.findNamespaceByID.rejects(Error('compose.repository.NamespaceNotFound'))
        h.findNamespaceBySlug.resolves(ns)

        expect(await h.resolveNamespace(ns.slug)).to.deep.equal(ns)

        expect(h.findNamespaceByID.calledOnceWith(ns.slug), 'findNamespaceByID call expected').true
        expect(h.findNamespaceBySlug.calledOnceWith(ns.slug), 'findNamespaceBySlug call expected').true
      })
    })
  })

  describe('helpers', () => {
    let h: ComposeHelper
    let composeApiStub: StubbedInstance<ComposeAPI>

    const $namespace = new Namespace({
      namespaceID: '2',
      slug: 'slg',
      name: 'name of the name-space',
    })

    const $module = new Module({
      moduleID: '1',
      fields: [
        new ModuleFieldString({ name: 'str' }),
        new ModuleFieldNumber({ name: 'num' }),
        new ModuleFieldString({ name: 'multi', isMulti: true }),
      ],
    }, $namespace)

    beforeEach(() => {
      composeApiStub = stubObject<ComposeAPI>(new ComposeAPI({}))

      h = new ComposeHelper({
        ComposeAPI: composeApiStub,
        $module,
        $namespace,
      })
    })

    describe('record making', () => {
      it('should make a record using $module', async () => {
        const record = await h.makeRecord({ str: 'foo' })

        expect(record).to.instanceof(Record)
        expect(record.moduleID).equal($module.moduleID)
        expect(record.values.str).to.equal('foo')
      })
    })

    describe('record saving', () => {
      it('should reject invalid record', async () => {
        const tests = [
          { val: null, label: 'Null value' },
          { val: {}, label: 'Empty object' },
        ]
        for (const t of tests) {
          expect(async () => { await h.saveRecord(t.val as unknown as Record) }, t.label).to.throw
        }
      })

      it('should create new', async () => {
        const record = await h.makeRecord({})

        composeApiStub.recordCreate.resolves(kv(record))

        await h.saveRecord(record)

        expect(composeApiStub.recordCreate.calledOnceWith(kv(record))).true
      })

      it('should update existing', async () => {
        const record = new Record($module, { recordID: '222' })

        composeApiStub.recordUpdate.resolves(kv(record))

        await h.saveRecord(record)

        expect(composeApiStub.recordUpdate.calledOnceWith(kv(record))).true
      })
    })

    describe('record deleting', () => {
      it('should validate input', async () => {
        expect(async () => { await h.deleteRecord(null as unknown as Record) }).to.throw
        expect(async () => { await h.deleteRecord(false as unknown as Record) }).to.throw
        expect(async () => { await h.deleteRecord(true as unknown as Record) }).to.throw
        expect(async () => { await h.deleteRecord({} as unknown as Record) }).to.throw
        expect(async () => { await h.deleteRecord($module as unknown as Record) }).to.throw
        expect(async () => { await h.deleteRecord(Promise.resolve($module) as unknown as Record) }).to.throw
      })

      it('should delete existing record', async () => {
        const record = new Record($module, { recordID: '222' })
        composeApiStub.recordDelete.resolves(kv(record))
        await h.deleteRecord(record)
        expect(composeApiStub.recordDelete.calledOnceWith(kv(record))).true
      })

      it('should not delete fresh record', async () => {
        const record = new Record($module)
        composeApiStub.recordDelete.resolves()
        await h.deleteRecord(record)
        expect(composeApiStub.recordDelete.notCalled).true
      })
    })

    describe('record finding', () => {
      it('should find records on $module')
      it('should find records on a different module')
      it('should properly translate filter to ID when numeric')
      it('should properly translate filter to query when string')
      it('should cast retrieved objects to Record')
    })

    describe('record finding by id', () => {
      it('should find record on $module')
      it('should find record on a different module')
      it('should find by ID when given a Record object')
      it('should cast retrieved objects to Record')
    })

    describe('page making', () => {
      it('should make a page', async () => {
        // expect(await h.makeRecord({}, $module)).to.instanceof(Record)
        expect((await h.makePage({ title: 'foo' })).title).to.equal('foo')
      })
    })

    describe('page saving', () => {
      it('should create new', async () => {
        const page = await h.makePage(pagePayload)
        composeApiStub.pageCreate.resolves(pagePayload)

        expect(await h.savePage(page)).instanceOf(Page)
        expect(composeApiStub.pageCreate.calledWith(kv(page))).true
      })

      it('should update existing', async () => {
        const page = await h.makePage({ ...pagePayload, pageID: '123' })
        composeApiStub.pageUpdate.resolves(pagePayload)

        expect(await h.savePage(page)).instanceOf(Page)
        expect(composeApiStub.pageUpdate.calledWith(kv(page))).true
      })
    })

    describe('page deleting', () => {
      it('should delete existing page', async () => {
        const page = await h.makePage({ ...pagePayload, pageID: '123' })
        composeApiStub.pageDelete.resolves(kv(page))

        await h.deletePage(page)
        expect(composeApiStub.pageDelete.calledWith(kv(page))).true
      })

      it('should not delete fresh page', async () => {
        const page = await h.makePage({ ...pagePayload })
        composeApiStub.pageDelete.resolves(kv(page))

        await h.deletePage(page)
        expect(composeApiStub.pageDelete.notCalled).true
      })
    })

    describe('page finding', () => {
      it('should find pages on current namespace', async () => {
        composeApiStub.pageList.resolves({ filter: {}, set: [pagePayload] })

        await h.findPages()

        expect(composeApiStub.pageList.calledWith({ ...ns(h) })).true
      })

      it('should find pages on given namespace', async () => {
        composeApiStub.pageList.resolves({ filter: {}, set: [pagePayload] })

        await h.findPages(undefined, new Namespace({ namespaceID: '22' }))
        expect(composeApiStub.pageList.calledWith({ namespaceID: '22' })).true
      })

      it('should filter (as string) pages on given namespace', async () => {
        composeApiStub.pageList.resolves({ filter: {}, set: [pagePayload] })

        await h.findPages('filter')
        expect(composeApiStub.pageList.calledWith({ ...ns(h), query: 'filter' })).true
      })

      it('should filter (as object) pages on given namespace', async () => {
        composeApiStub.pageList.resolves({ filter: {}, set: [pagePayload] })

        await h.findPages({ query: 'filter', limit: 10 })
        expect(composeApiStub.pageList.calledWith({ ...ns(h), query: 'filter', limit: 10 })).true
      })
    })

    describe('findPageByID', () => {
      it('should find page on current namespace', async () => {
        composeApiStub.pageRead.resolves(pagePayload)

        await h.findPageByID('1000')
        expect(composeApiStub.pageRead.calledWith({ ...ns(h), pageID: '1000' })).true
      })

      it('should find page on given namespace', async () => {
        composeApiStub.pageRead.resolves(pagePayload)

        await h.findPageByID('1000', new Namespace({ namespaceID: '22' }))
        expect(composeApiStub.pageRead.calledWith({ namespaceID: '22', pageID: '1000' })).true
      })

      it('should find page from Object', async () => {
        composeApiStub.pageRead.resolves(pagePayload)
        const page = await h.makePage({ ...pagePayload, pageID: '1001' })

        await h.findPageByID(page, new Namespace({ namespaceID: '22' }))
        expect(composeApiStub.pageRead.calledWith({ namespaceID: '22', pageID: '1001' })).true
      })
    })

    describe('making field name from label ', () => {
      it('should clean it up and present in camelcase', () => {
        const test = 'so.me;me:sa,ge+he_ok ok'
        const exp = 'SoMeMeSaGeHe_okOk'

        expect(h.moduleFieldNameFromLabel(test)).to.eq(exp)
      })
    })

    describe('module making', () => {
      it('should make new', async () => {
        const module = await h.makeModule({ name: 'MyModule' })
        expect(module).to.be.instanceOf(Module)
        expect(module.name).to.equal('MyModule')
      })
    })

    describe('module saving', async () => {
      it('should create new', async () => {
        const module = new Module()

        composeApiStub.moduleCreate.resolves(kv(module))

        await h.saveModule(module)

        expect(composeApiStub.moduleCreate.calledWith(kv(module))).true
      })

      it('should update existing', async () => {
        const module = new Module({ moduleID: '555' })

        composeApiStub.moduleUpdate.resolves(kv(module))

        await h.saveModule(module)

        expect(composeApiStub.moduleUpdate.calledWith(kv(module))).true
      })
    })

    describe('module finding', () => {
      it('should find modules on $namespace')
      it('should find modules on a different namespace')
      it('should properly translate filter to query when string')
      it('should cast retrieved objects to Module')
    })

    describe('module finding by id', () => {
      it('should find module by id on $namespace', async () => {
        const module = new Module({ moduleID: '555' })

        composeApiStub.moduleRead.resolves({ ...module })

        expect(await h.findModuleByID('555')).to.be.instanceOf(Module)

        expect(composeApiStub.moduleRead.calledWith({ moduleID: '555', ...ns(h) })).true
      })
    })

    describe('module finding by handle', () => {
      it('should find module by handle on $namespace', async () => {
        const module = new Module({ moduleID: '555' })

        composeApiStub.moduleList.resolves({ filter: { limit: 1 }, set: [module] })

        expect(await h.findModuleByHandle('some-module')).to.be.instanceOf(Module)

        expect(composeApiStub.moduleList.calledWith({ handle: 'some-module', ...ns(h) })).true
      })
    })

    describe('module finding by name', () => {
      it('should find module by name on $namespace', async () => {
        const module = new Module({ moduleID: '555' })

        composeApiStub.moduleList.resolves({ filter: { limit: 1 }, set: [module] })

        expect(await h.findModuleByName('some-module')).to.be.instanceOf(Module)

        expect(composeApiStub.moduleList.calledWith({ name: 'some-module', ...ns(h) })).true
      })
    })

    describe('namespace making', () => {
      it('should make active namespace', async () => {
        expect((await h.makeNamespace()).enabled).to.be.true
      })

      it('should use slug as name', async () => {
        const ns = await h.makeNamespace({ slug: 'sluggy-slug' })
        expect(ns.slug).to.equal('sluggy-slug')
        expect(ns.name).to.equal('sluggy-slug')
      })
    })

    describe('namespace saving', async () => {
      it('should create new', async () => {
        const ns = new Namespace()
        composeApiStub.namespaceCreate.resolves(kv(ns))
        await h.saveNamespace(ns)
        expect(composeApiStub.namespaceCreate.calledWith(kv(ns))).true
      })

      it('should update existing', async () => {
        const ns = new Namespace({ namespaceID: '555' })
        composeApiStub.namespaceUpdate.resolves(kv(ns))
        await h.saveNamespace(ns)
        expect(composeApiStub.namespaceUpdate.calledWith(kv(ns))).true
      })
    })

    describe('sendEmail', () => {
      it('should write some tests')
    })

    describe('recordToPlainText', () => {
      const m = new Module({
        fields: [
          { name: 'dummy' },
          { name: 'multi', isMulti: true },
          { name: 'empty' },
        ],
      })

      it('should convert a given record to plain text', () => {
        const record = new Record(m, {
          values: [
            { name: 'dummy', value: 'value' },
            { name: 'multi', value: 'v1' },
            { name: 'multi', value: 'v2' },
          ],
        })
        expect(h.recordToPlainText(null, record)).to.eq('dummy:\nvalue\n\nmulti:\nv1, v2\n\nempty:\n/')
      })

      it('should convert white-listed fields of a given record to plain text', () => {
        const record = new Record(m, {
          values: [
            { name: 'dummy', value: 'value' },
            { name: 'multi', value: 'v1' },
            { name: 'multi', value: 'v2' },
          ],
        })
        expect(h.recordToPlainText(['dummy', 'multi'], record)).to.eq('dummy:\nvalue\n\nmulti:\nv1, v2')
      })
    })
  })
})
