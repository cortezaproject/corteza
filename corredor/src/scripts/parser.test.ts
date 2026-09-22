import { describe, it } from 'mocha'
import * as chai from 'chai'
import SourceParser from './parser'
import chaiAsPromised from 'chai-as-promised'
import { Script } from '../types'
import { Trigger } from './trigger'

chai.use(chaiAsPromised)

const expect = chai.expect

function p(source: string): Promise<Partial<Script>> {
  return new SourceParser(source).parse()
}

const errExpDefExp = 'expecting default export'
const errExpObjAsDefExp = 'expecting object as default export'
const errTrigDefMis = 'triggers or iterator definition missing'
const errNoTrigDef = 'triggers not defined'
const errExecFunMis = 'exec function missing'

describe(__filename, () => {
  describe('when parsing empty script', () => {
    it('should throw an error about expecting default export', async () => {
      await expect(p('')).to.be.rejectedWith(errExpDefExp)
    })
  })

  describe('when parsing invalid export', () => {
    it('should throw an error about expecting default export', () => {
      expect(p('invalid')).to.be.rejectedWith(errExpDefExp)
      expect(p('export const foo = {}')).to.be.rejectedWith(errExpDefExp)
    })

    it('should throw an error about expecting object', () => {
      expect(p('export default foo')).to.be.rejectedWith(errExpObjAsDefExp)
      expect(p('export default false')).to.be.rejectedWith(errExpObjAsDefExp)
      expect(p('export default "string"')).to.be.rejectedWith(errExpObjAsDefExp)
    })
  })

  describe('when parsing invalid script', () => {
    it('should throw an error about missing trigger definition', () => {
      expect(p('export default {}')).to.be.rejectedWith(errTrigDefMis)
    })

    it('should throw an error about triggers not defined 1', () => {
      expect(p('export default { triggers: [] }')).to.be.rejectedWith(errNoTrigDef)
    })

    it('should throw an error about triggers not defined 2', () => {
      expect(p('export default { triggers () {} }')).to.be.rejectedWith(errNoTrigDef)
    })

    it('should throw an error about triggers not defined 3', () => {
      expect(p('export default { * triggers () {} }')).to.be.rejectedWith(errNoTrigDef)
    })

    it('should throw an error about triggers not defined 4', () => {
      expect(p("export default { triggers ({ on }) { return on('manual') } }")).to.be.rejectedWith(
        errExecFunMis,
      )
    })
  })

  describe('when parsing valid script', () => {
    it('return parsed script', () => {
      const script = p(`
      export default {
        label: 'Label',
        description: 'Desc',
        triggers ({ on }) { return on('manual') },
        exec () {}
      }
      `)

      expect(script).to.eventually.have.property('label').equal('Label')

      expect(script).to.eventually.have.property('description').equal('Desc')

      expect(script)
        .to.eventually.have.property('triggers')
        .deep.equal([new Trigger({ eventTypes: ['onManual'] })])
    })
  })

  describe('when parsing script with trigger generator function', () => {
    it('return parsed script', () => {
      const script = p(`
      export default {
        * triggers ({ on }) { yield on('manual') },
        exec () {}
      }
      `)

      expect(script)
        .to.eventually.have.property('triggers')
        .deep.equal([new Trigger({ eventTypes: ['onManual'] })])
    })
  })
})
