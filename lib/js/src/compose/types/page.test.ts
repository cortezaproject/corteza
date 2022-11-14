import { NoID } from '../../cast'
import { expect } from 'chai'
import { Page } from './page'

describe(__filename, () => {
  describe('initialization', () => {
    it('simple', () => {
      const page = new Page({})

      expect(page).to.have.property('pageID').and.to.eq(NoID)
      expect(page).to.have.property('selfID').and.to.eq(NoID)
      expect(page).to.have.property('moduleID').and.to.eq(NoID)
      expect(page).to.have.property('namespaceID').and.to.eq(NoID)

      expect(page).to.have.property('title').and.to.eq('')
      expect(page).to.have.property('handle').and.to.eq('')
      expect(page).to.have.property('description').and.to.eq('')
      expect(page).to.have.property('weight').and.to.eq(0)

      expect(page).to.have.property('labels').and.to.deep.eq({})

      expect(page).to.have.property('visible').and.to.eq(false)

      expect(page).to.have.property('blocks').and.to.deep.eq([])

      expect(page).to.have.property('createdAt').and.to.eq(undefined)
      expect(page).to.have.property('updatedAt').and.to.eq(undefined)
      expect(page).to.have.property('deletedAt').and.to.eq(undefined)

      expect(page).to.have.property('canUpdatePage').and.to.eq(false)
      expect(page).to.have.property('canDeletePage').and.to.eq(false)
      expect(page).to.have.property('canGrant').and.to.eq(false)
    })

    it('complex', () => {
      const dateNow = new Date()

      const page = new Page({
        pageID: '42',
        selfID: '42',
        moduleID: '42',
        namespaceID: '42',

        title: 'title',
        handle: 'handle',
        description: 'description',
        weight: 1,

        labels: {
          foo: 'foo'
        },

        visible: true,

        blocks: [
          { kind: 'RecordList', xywh: [0, 0, 3, 3] },
        ],

        createdAt: dateNow,
        updatedAt: dateNow,
        deletedAt: dateNow,

        canDeletePage: true,
        canUpdatePage: true,
        canGrant: true,
      })

      expect(page).to.have.property('pageID').and.to.eq('42')
      expect(page).to.have.property('selfID').and.to.eq('42')
      expect(page).to.have.property('moduleID').and.to.eq('42')
      expect(page).to.have.property('namespaceID').and.to.eq('42')

      expect(page).to.have.property('title').and.to.eq('title')
      expect(page).to.have.property('handle').and.to.eq('handle')
      expect(page).to.have.property('description').and.to.eq('description')
      expect(page).to.have.property('weight').and.to.eq(1)

      expect(page).to.have.property('labels').and.to.deep.eq({
        foo: 'foo'
      })

      expect(page).to.have.property('visible').and.to.eq(true)

      expect(page).to.have.property('blocks').and.lengthOf(1)

      expect(page).to.have.property('createdAt').and.to.eq(dateNow)
      expect(page).to.have.property('updatedAt').and.to.eq(dateNow)
      expect(page).to.have.property('deletedAt').and.to.eq(dateNow)

      expect(page).to.have.property('canUpdatePage').and.to.eq(true)
      expect(page).to.have.property('canDeletePage').and.to.eq(true)
      expect(page).to.have.property('canGrant').and.to.eq(true)
    })
  }),

  describe('getters', () => {
    const page = new Page({
      pageID: '42',
      selfID: '42',
      moduleID: '42',
      namespaceID: '42',
    })

    const pg = new Page({})

    it('resourceType', () => {
      expect(page.resourceType).to.eq('compose:page')
    })

    it('resourceID', () => {
      expect(page.resourceID).to.eq('compose:page:42')
    })

    it('isRecordPage', () => {
      expect(page.isRecordPage).to.eq(true)
      expect(pg.isRecordPage).to.eq(false)
    })

    it('firstLevel', () => {
      expect(page.firstLevel).to.eq(false)
      expect(pg.firstLevel).to.eq(true)
    })
  })

  describe('methods', () => {
    it('export', () => {
      const page = new Page({
        title: 'title',
        handle: 'handle',
        description: 'description',
        visible: true,
        blocks: [
          { kind: 'RecordList', xywh: [0, 0, 3, 3] }
        ],
      })

      const pageExport = page.export()
      expect(pageExport).to.have.property('title').and.to.eq('title')
      expect(pageExport).to.have.property('handle').and.to.eq('handle')
      expect(pageExport).to.have.property('description').and.to.eq('description')
      expect(pageExport).to.have.property('visible').and.to.eq(true)
      expect(pageExport).to.have.property('blocks').and.lengthOf(1)
    })

    it('validate', () => {
      const page = new Page({
        blocks: [
          { kind: 'RecordList', xywh: [0, 0, 3, 3] }
        ],
      })

      expect(page.validate()).to.deep.eq([])
    })
  })

  describe('blocks', () => {
    it('should make all kinds of block', () => {
      const p = new Page({
        title: 'test page',
        blocks: [
          { kind: 'Automation', xywh: [0, 0, 3, 3] },
          { kind: 'Chart', xywh: [0, 0, 3, 3] },
          { kind: 'Content', xywh: [0, 0, 3, 3] },
          { kind: 'File', xywh: [0, 0, 3, 3] },
          { kind: 'Record', xywh: [0, 0, 3, 3] },
          { kind: 'RecordList', xywh: [0, 0, 3, 3] },
          { kind: 'RecordOrganizer', xywh: [0, 0, 3, 3] },
          { kind: 'SocialFeed', xywh: [0, 0, 3, 3] },
        ],
      })

      expect(p.blocks).lengthOf(8)
    })
  })
})
