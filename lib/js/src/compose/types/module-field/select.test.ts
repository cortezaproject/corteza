import { expect } from 'chai'
import { ModuleFieldSelect } from './select'

describe('select module field options', () => {
  it('uses value as label when label is empty', () => {
    const f = new ModuleFieldSelect({
      options: {
        options: [
          { value: 'one' },
          { value: 'two', text: 'Two' },
        ],
      },
    } as never)

    expect(f.options.options.map(({ text }) => text)).to.deep.equal(['one', 'Two'])
  })

  it('creates an empty option without a value', () => {
    const f = new ModuleFieldSelect()

    expect(f.createSelectOption()).to.deep.include({ value: '', text: '' })
  })
})
