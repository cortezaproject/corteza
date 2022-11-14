import { expect } from 'chai'
import { PageBlockCalendar } from './page-block'

describe(__filename, () => {
  describe('check namespace casting', () => {
    it('simple assignment', () => {
      const cal = new PageBlockCalendar({
        title: 'My Calendar',
        options: {
          defaultView: 'month',
          feeds: [
            { endField: 'EndDateTime', moduleID: '69055788747849745', startField: 'ActivityDate', titleField: 'Subject' },
            { endField: null, moduleID: '69055789049839633', startField: 'ActivityDate', titleField: 'Subject' },
          ],
          header: { views: ['timeGridWeek', 'dayGridMonth', 'timeGridDay', 'dayGridMonth', 'month'] },
        },
        style: { variants: { bodyBg: 'white', border: 'dark', headerBg: 'white', headerText: 'dark' } },
        kind: 'Calendar',
        xywh: [0, 0, 6, 14],
      },
      )

      expect(cal.getHeader()).to.have.property('center').equal('title')
      expect(cal.getHeader()).to.have.property('right').equal('dayGridMonth,timeGridWeek,timeGridDay')
    })
  })
})
