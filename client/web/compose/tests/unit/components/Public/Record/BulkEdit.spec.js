/* eslint-disable no-unused-expressions */
/* global jest */
import { expect } from 'chai'
import { shallowMount, createLocalVue } from '@vue/test-utils'
import BootstrapVue from 'bootstrap-vue'
import BulkEdit from 'corteza-webapp-compose/src/components/Public/Record/BulkEdit'

// tests/setup.js stubs corteza-js with empty classes; this spec needs a Record
// that carries values so it can check they are cleared, and a Module that can
// list its fields
jest.mock('@cortezaproject/corteza-js', () => {
  class Record {
    constructor (module, initial = {}) {
      this.recordID = initial.recordID || '0'
      this.values = { ...(initial.values || {}) }
    }
  }

  class Module {
    constructor ({ fields = [] } = {}) {
      this.fields = fields
    }

    systemFields () {
      return []
    }
  }

  return {
    NoID: '0',
    compose: { Namespace: class {}, Module, Record },
    validator: { Validated: class { push () {} filterByMeta () { return [] } } },
  }
}, { virtual: true })

// The field editor pulls in the whole field component tree; not under test here
jest.mock('corteza-webapp-compose/src/components/ModuleFields/Editor', () => ({
  name: 'field-editor',
  render: () => {},
}))

describe('components/Public/Record/BulkEdit', () => {
  let localVue
  let namespace
  let module

  const mount = (props = {}) => shallowMount(BulkEdit, {
    localVue,
    propsData: { namespace, module, ...props },
    mocks: {
      $t: (key) => key,
      $ComposeAPI: {},
      $root: { $emit: () => {}, $on: () => {}, $off: () => {} },
    },
    stubs: ['c-input-select', 'c-input-confirm', 'font-awesome-icon', 'b-modal'],
  })

  beforeEach(() => {
    localVue = createLocalVue()
    localVue.use(BootstrapVue)

    const { compose } = require('@cortezaproject/corteza-js')
    namespace = new compose.Namespace()
    module = new compose.Module({
      fields: [
        { fieldID: '10', name: 'title', kind: 'String', canUpdateRecordValue: true },
        { fieldID: '11', name: 'amount', kind: 'Number', canUpdateRecordValue: true },
      ],
    })
  })

  it('clears the chosen fields and values when the modal is closed', () => {
    const wrapper = mount({ allowAddField: true })

    wrapper.vm.addField('title')
    wrapper.vm.addField('amount')
    wrapper.vm.record.values.title = 'changed'
    wrapper.vm.record.values.amount = 5

    wrapper.vm.onModalHide()

    expect(wrapper.vm.fields).to.deep.equal([])
    expect(wrapper.vm.record.values).to.deep.equal({})
    expect(wrapper.emitted('close')).to.have.length(1)
  })

  it('keeps the fields the parent preselected for inline editing', () => {
    const wrapper = mount({ selectedFields: ['title'], initialRecord: { values: { title: 'given' } } })

    wrapper.vm.addField('amount')
    wrapper.vm.record.values.amount = 5

    wrapper.vm.onModalHide()

    expect(wrapper.vm.fields).to.deep.equal(['title'])
    expect(wrapper.vm.record.values).to.deep.equal({ title: 'given' })
  })

  it('starts clean after a successful bulk update', () => {
    const wrapper = mount({ allowAddField: true })

    wrapper.vm.addField('title')
    wrapper.vm.record.values.title = 'changed'

    wrapper.vm.setDefaultValues()

    expect(wrapper.vm.fields).to.deep.equal([])
    expect(wrapper.vm.record.values).to.deep.equal({})
    expect(wrapper.vm.showModal).to.equal(false)
  })
})
