import { describe, it, expect } from 'vitest'
import {
  apiError,
  detailMessage,
  partitionSaveErrors,
  saveWarnings,
  fieldLabeller,
} from './record-errors.js'

// The shapes here are copied from what the dev server actually answers with —
// `POST /compose/namespace/:ns/module/:mod/record/` reports every failure in
// the body of a 200, one detail per issue.
const duplicateValue = {
  kind: 'duplicateValue',
  message: 'The value "{{value}}" already exists in another record',
  meta: { field: 'sku', id: 'parent:0', recordID: 509864056697061377, value: 'SKU-1' },
}

const expressionFailure = {
  kind: 'valueExpression',
  message: 'failed to evaluate formula expression "nonexistent.thing + 1": unknown parameter',
  meta: { field: 'calc', id: 'parent:0' },
}

const required = {
  kind: 'empty',
  message: 'This field is required',
  meta: { field: 'title', id: 'parent:0' },
}

describe('apiError', () => {
  // An upload posts through the raw axios instance, so it never passes
  // stdResolve and reads the body itself. `new Error(data.error)` on the object
  // the server sends is what put '[object Object]' in front of users.
  it('takes the message off the error object', () => {
    const e = apiError({
      error: {
        message: 'not allowed to upload this type of file',
        meta: { resource: 'compose:attachment', type: 'notAllowedToUploadThisType' },
      },
    })

    expect(e).to.be.instanceOf(Error)
    expect(e?.message).to.equal('not allowed to upload this type of file')
    expect(String(e)).to.not.contain('[object Object]')
  })

  it('handles a bare string error too', () => {
    expect(apiError({ error: 'nope' })?.message).to.equal('nope')
  })

  it('falls back when the error object carries no message', () => {
    const e = apiError({ error: { meta: {} } }, 'Upload failed for "a.png"')

    expect(e?.message).to.equal('Upload failed for "a.png"')
    expect(String(e)).to.not.contain('[object Object]')
  })

  it('carries details where the API sends them', () => {
    const e = apiError({ error: { message: '1 issue(s) found', details: [required] } })

    expect((e as Error & { details?: unknown[] })?.details).to.deep.equal([required])
  })

  it('is null for a body that reported no failure', () => {
    expect(apiError({ response: { attachmentID: '1' } })).to.equal(null)
    expect(apiError(undefined)).to.equal(null)
  })
})

describe('detailMessage', () => {
  it('fills a placeholder from the detail meta', () => {
    expect(detailMessage(duplicateValue)).to.equal(
      'The value "SKU-1" already exists in another record',
    )
  })

  it('leaves a placeholder the meta cannot fill alone, rather than blanking it', () => {
    expect(detailMessage({ message: 'no {{nothing}} here', meta: { field: 'a' } })).to.equal(
      'no {{nothing}} here',
    )
  })

  it('is empty for a detail with no message', () => {
    expect(detailMessage({ meta: { field: 'a' } })).to.equal('')
    expect(detailMessage(undefined)).to.equal('')
  })
})

describe('partitionSaveErrors', () => {
  it('places every field issue beside its field when the form shows them all', () => {
    const { fieldErrors, general } = partitionSaveErrors({
      details: [required, expressionFailure],
    })

    expect(fieldErrors).to.deep.equal({
      title: 'This field is required',
      calc: expressionFailure.message,
    })
    expect(general).to.deep.equal([])
  })

  it('names a field the form does not show, rather than writing where nobody looks', () => {
    const { fieldErrors, general } = partitionSaveErrors(
      { details: [required, expressionFailure] },
      { canShow: name => name === 'title', labelOf: name => (name === 'calc' ? 'Calc' : name) },
    )

    expect(fieldErrors).to.deep.equal({ title: 'This field is required' })
    expect(general).to.deep.equal([`Calc: ${expressionFailure.message}`])
  })

  it('treats an issue with no field as general', () => {
    const { fieldErrors, general } = partitionSaveErrors({
      details: [{ kind: 'internal', message: 'something broke', meta: {} }],
    })

    expect(fieldErrors).to.deep.equal({})
    expect(general).to.deep.equal(['something broke'])
  })

  it('keeps both messages when one field carries two issues', () => {
    const { fieldErrors } = partitionSaveErrors({
      details: [required, { ...expressionFailure, meta: { field: 'title' } }],
    })

    expect(fieldErrors.title).to.equal(`This field is required\n${expressionFailure.message}`)
  })

  it('says a repeated complaint once — one issue per offending record', () => {
    const twice = {
      details: [
        { ...duplicateValue, meta: { ...duplicateValue.meta, recordID: 1 } },
        { ...duplicateValue, meta: { ...duplicateValue.meta, recordID: 2 } },
      ],
    }

    expect(partitionSaveErrors(twice).fieldErrors.sku).to.equal(
      'The value "SKU-1" already exists in another record',
    )
    expect(
      partitionSaveErrors(twice, { canShow: () => false, labelOf: () => 'SKU' }).general,
    ).to.deep.equal(['SKU: The value "SKU-1" already exists in another record'])
  })

  it('returns nothing for an error the API sent no details with', () => {
    expect(partitionSaveErrors(new Error('module does not exist'))).to.deep.equal({
      fieldErrors: {},
      general: [],
    })
    expect(partitionSaveErrors(undefined)).to.deep.equal({ fieldErrors: {}, general: [] })
  })
})

describe('saveWarnings', () => {
  it('reads the issues a successful save carried, filled and labelled', () => {
    const saved = {
      recordID: '1',
      valueErrors: {
        message: '1 issue(s) found',
        set: [
          {
            kind: 'duplication_warning',
            message: 'The value "{{value}}" already exists in another record',
            meta: { field: 'code', value: 'AAA', recordID: '2' },
          },
        ],
      },
    }

    expect(saveWarnings(saved, { labelOf: () => 'Code' })).to.deep.equal([
      'Code: The value "AAA" already exists in another record',
    ])
  })

  it('says a repeated warning once', () => {
    const detail = {
      kind: 'duplication_warning',
      message: 'The value "{{value}}" already exists in another record',
      meta: { field: 'code', value: 'AAA' },
    }
    const saved = { valueErrors: { set: [detail, { ...detail, meta: { ...detail.meta } }] } }

    expect(saveWarnings(saved, { labelOf: () => 'Code' })).to.deep.equal([
      'Code: The value "AAA" already exists in another record',
    ])
  })

  it('is empty for a clean save', () => {
    expect(saveWarnings({ valueErrors: { message: '0 issue(s) found' } })).to.deep.equal([])
    expect(saveWarnings(undefined)).to.deep.equal([])
  })
})

describe('fieldLabeller', () => {
  it('prefers the field label and falls back to its name', () => {
    const labelOf = fieldLabeller({ fields: [{ name: 'sku', label: 'SKU' }, { name: 'code' }] })

    expect(labelOf('sku')).to.equal('SKU')
    expect(labelOf('code')).to.equal('code')
    expect(labelOf('gone')).to.equal('gone')
  })

  it('survives a module that never loaded', () => {
    expect(fieldLabeller(null)('sku')).to.equal('sku')
  })
})
