import { expect } from 'chai'
import moment from 'moment'
import { getDate, setDate, getTime, setTime } from './index'

describe(__filename, () => {
  describe('Can get and set time', () => {
    it('getDate', () => {
      expect(getDate(undefined)).to.equal(undefined)
      expect(getDate('2022-01-01')).to.equal('2022-01-01')
      expect(getDate('2022-01-01T12:00:00Z')).to.equal(moment.utc('2022-01-01T12:00:00Z').local().format('YYYY-MM-DD'))
      expect(getDate('2021-12-31T23:00:00Z')).to.equal(moment.utc('2021-12-31T23:00:00Z').local().format('YYYY-MM-DD'))
    })

    it('setDate', () => {
      // Should undefined if noDate or date is invalid
      expect(setDate(undefined, '2022-01-01T12:00:00Z')).to.equal(undefined)
      expect(setDate('2022-01-02', '2022-01-01T12:00:00Z', true)).to.equal(undefined)

      // If noTime, return just date
      expect(setDate('2022-01-02', '2022-01-01', false, true)).to.equal('2022-01-02')

      expect(setDate('2022-01-02', '2022-01-01T12:00:00Z')).to.equal(moment(`2022-01-02 ${moment.utc('2022-01-01T12:00:00Z').local().format('HH:mm')}`, 'YYYY-MM-DD HH:mm').utc().format())
      expect(setDate('2022-01-02', undefined)).to.equal(moment('2022-01-02 00:00', 'YYYY-MM-DD HH:mm').utc().format())
    })

    it('getTime', () => {
      expect(getTime(undefined)).to.equal(undefined)
      expect(getTime('12:00')).to.equal('12:00')
      expect(getTime('2022-01-01T12:00:00Z')).to.equal(moment.utc('2022-01-01T12:00:00Z').local().format('HH:mm'))
      expect(getTime('2021-12-31T23:00:00Z')).to.equal(moment.utc('2021-12-31T23:00:00Z').local().format('HH:mm'))
    })

    it('setTime', () => {
      // Should undefined if noTime or time is invalid
      expect(setTime(undefined, '2022-01-01T12:00:00Z')).to.equal(undefined)
      expect(setTime('12:00', '2022-01-01T12:00:00Z', false, true)).to.equal(undefined)

      // If noDate, return just time
      expect(setTime('12:00', '15:00', true, false)).to.equal('12:00')

      expect(setTime('12:00', '2022-01-01T12:00:00Z')).to.equal(moment(`${moment.utc('2022-01-01T12:00:00Z').local().format('YYYY-MM-DD')} 12:00`, 'YYYY-MM-DD HH:mm').utc().format())
      expect(setTime('12:00', undefined)).to.equal(moment(`${moment().local().format('YYYY-MM-DD')} 12:00`, 'YYYY-MM-DD HH:mm').utc().format())
    })
  })
})
