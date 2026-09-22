import { expect } from 'chai'
import pino from 'pino'
import { LogToArray } from './log-to-array'

describe('log-to-array', () => {
  it('collects pino lines as ISO timestamp, level and message', () => {
    const buf = new LogToArray()
    const log = pino({ level: 'trace' }, buf as unknown as pino.DestinationStream)

    log.info('hello %s', 'world')
    log.warn('careful')

    const lines = buf.serialize()
    expect(lines).to.have.length(2)
    expect(lines[0]).to.match(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z INFO {2}hello world$/)
    expect(lines[1]).to.match(/ WARN {2}careful$/)

    // the stamp is the log time, not an offset from some other epoch
    const stamped = new Date(lines[0].slice(0, 24)).getTime()
    expect(Math.abs(stamped - Date.now())).to.be.below(5000)
  })
})
