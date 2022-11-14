import { expect } from 'chai'
import { Assert, AssertFuzzy, AssertStrict } from './filter'

describe(__filename, () => {
  describe('AssertStrict', () => {
    it('regular query', () => {
      expect(AssertStrict({ k: 'csz' }, 'cs', 'k')).to.be.true
      expect(AssertStrict('csz', 'cs')).to.be.true
    })

    it('query with diacritics', () => {
      expect(AssertStrict({ k: 'čšž' }, 'čš', 'k')).to.be.true
      expect(AssertStrict('čšž', 'čš')).to.be.true
    })

    it('query without diacritics', () => {
      expect(AssertStrict({ k: 'čšž' }, 'cs', 'k')).to.be.true
      expect(AssertStrict('čšž', 'cs')).to.be.true
    })
  })

  describe('AssertFuzzy', () => {
    it('abbreviated query', () => {
      expect(AssertFuzzy({ k: 'ComposeRecordMaker' }, 'cmprecmak', 'k')).to.be.true
      expect(AssertFuzzy('ComposeRecordMaker', 'cmprecmak')).to.be.true
    })

    it('too abbreviated query', () => {
      expect(AssertFuzzy({ k: 'ComposeRecordMaker' }, 'crck', 'k')).to.be.false
      expect(AssertFuzzy('ComposeRecordMaker', 'crck')).to.be.false
    })
  })

  describe('Assert', () => {
    it('abbreviated query', () => {
      expect(Assert({ k: 'ComposeRecordMaker' }, 'cmprecmak', 'k')).to.be.true
      expect(Assert('ComposeRecordMaker', 'cmprecmak')).to.be.true
    })

    it('substring', () => {
      expect(Assert({ k: 'ComposeRecordMaker' }, 'record', 'k')).to.be.true
      expect(Assert('ComposeRecordMaker', 'record')).to.be.true
    })

    it('mismatch', () => {
      expect(Assert({ k: 'ComposeRecordMaker' }, 'recorNOTHERE', 'k')).to.be.false
      expect(Assert('ComposeRecordMaker', 'recorNOTHERE')).to.be.false
    })
  })
})
