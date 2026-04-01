import { expect } from 'chai'
import { getWeekStartDay } from './locale'

describe(__filename, () => {
  describe('getWeekStartDay', () => {
    it('should return 0 by default for en-US', () => {
      expect(getWeekStartDay('en-US')).to.equal(0)
    })

    // @ts-ignore - checking if environment supports weekInfo/getWeekInfo
    if (new Intl.Locale('hu-HU').weekInfo || (new Intl.Locale('hu-HU') as any).getWeekInfo) {
      it('should return 1 for hu-HU', () => {
        expect(getWeekStartDay('hu-HU')).to.equal(1)
      })
    }
  })
})
