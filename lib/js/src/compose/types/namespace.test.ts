import { NoID } from '../../cast'
import { expect } from 'chai'
import { Namespace } from './namespace'

describe(__filename, () => {
  describe('initialization', () => {
    // it('requires name and slug', () => {
    //   expect(() => new Namespace({})).to.throw('field name is empty')
    //   expect(() => new Namespace({ name: 'name' })).to.throw('field slug is empty')
    // })

    it('simple', () => {
      const ns = new Namespace({})

      expect(ns).to.have.property('namespaceID').and.to.eq(NoID)
      expect(ns).to.have.property('name').and.to.eq('')
      expect(ns).to.have.property('slug').and.to.eq('')

      expect(ns).to.have.property('enabled').and.to.eq(false)

      expect(ns).to.have.property('labels').and.to.deep.eq({})
      expect(ns).to.have.property('meta').and.to.deep.eq({})

      expect(ns).to.have.property('createdAt').and.to.eq(undefined)
      expect(ns).to.have.property('updatedAt').and.to.eq(undefined)
      expect(ns).to.have.property('deletedAt').and.to.eq(undefined)

      expect(ns).to.have.property('canCreateChart').and.to.eq(false)
      expect(ns).to.have.property('canCreateModule').and.to.eq(false)
      expect(ns).to.have.property('canCreatePage').and.to.eq(false)
      expect(ns).to.have.property('canDeleteNamespace').and.to.eq(false)
      expect(ns).to.have.property('canUpdateNamespace').and.to.eq(false)
      expect(ns).to.have.property('canManageNamespace').and.to.eq(false)
      expect(ns).to.have.property('canCloneNamespace').and.to.eq(false)
      expect(ns).to.have.property('canGrant').and.to.eq(false)
    })

    it('apply', () => {
      const dateNow = new Date()

      const ns = new Namespace({
        namespaceID: '42',
        name: 'name',
        slug: 'slug',

        enabled: true,

        labels: {
          foo: 'foo'
        },

        meta: {
          subtitle: 'subtitle',
          description: 'description',
          icon: 'icon',
          logo: 'logo',
          logoEnabled: true,
        },

        createdAt: dateNow,
        updatedAt: dateNow,
        deletedAt: dateNow,

        canCreateChart: true,
        canCreateModule: true,
        canCreatePage: true,
        canDeleteNamespace: true,
        canUpdateNamespace: true,
        canManageNamespace: true,
        canCloneNamespace: true,
        canGrant: true,
      })

      expect(ns).to.have.property('namespaceID').and.to.eq('42')
      expect(ns).to.have.property('name').and.to.eq('name')
      expect(ns).to.have.property('slug').and.to.eq('slug')

      expect(ns).to.have.property('enabled').and.to.eq(true)

      expect(ns).to.have.property('labels').and.to.deep.eq({
        foo: 'foo'
      })
      expect(ns).to.have.property('meta').and.to.deep.eq({
        subtitle: 'subtitle',
        description: 'description',
        icon: 'icon',
        logo: 'logo',
        logoEnabled: true,
      })

      expect(ns).to.have.property('createdAt').and.to.eq(dateNow)
      expect(ns).to.have.property('updatedAt').and.to.eq(dateNow)
      expect(ns).to.have.property('deletedAt').and.to.eq(dateNow)

      expect(ns).to.have.property('canCreateChart').and.to.eq(true)
      expect(ns).to.have.property('canCreateModule').and.to.eq(true)
      expect(ns).to.have.property('canCreatePage').and.to.eq(true)
      expect(ns).to.have.property('canDeleteNamespace').and.to.eq(true)
      expect(ns).to.have.property('canUpdateNamespace').and.to.eq(true)
      expect(ns).to.have.property('canManageNamespace').and.to.eq(true)
      expect(ns).to.have.property('canCloneNamespace').and.to.eq(true)
      expect(ns).to.have.property('canGrant').and.to.eq(true)
    })
  })

  describe('getters', () => {
    const ns = new Namespace({
      namespaceID: '42'
    })

    it('resourceType', () => {
      expect(ns.resourceType).to.eq('compose:namespace')
    })

    it('resourceID', () => {
      expect(ns.resourceID).to.eq('compose:namespace:42')
    })
  })
})
