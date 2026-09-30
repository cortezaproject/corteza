import { expect } from 'chai'
import { DalConnection } from './dalConnection'

describe('system DalConnection', () => {
  const base = {
    connectionID: '515955079844331521',
    type: 'corteza::system:dal-connection',
    canManageDalConfig: true,
  }

  it('offers the DSN connection type when the stored DAL type is empty', () => {
    const c = new DalConnection({
      ...base,
      config: {
        privacy: { sensitivityLevelID: '0' },
        dal: { type: '', params: undefined, modelIdent: '', modelIdentCheck: undefined },
      },
    })

    expect(c.config.dal?.type).to.equal('corteza::dal:connection:dsn')
    expect(c.config.dal?.params).to.deep.equal({ dsn: '' })
  })

  it('keeps a stored DAL type and DSN', () => {
    const c = new DalConnection({
      ...base,
      config: {
        privacy: { sensitivityLevelID: '0' },
        dal: {
          type: 'corteza::dal:connection:rest',
          params: { dsn: 'postgres://h/db' },
          modelIdent: 'records',
        },
      },
    })

    expect(c.config.dal?.type).to.equal('corteza::dal:connection:rest')
    expect(c.config.dal?.params).to.deep.equal({ dsn: 'postgres://h/db' })
    expect(c.config.dal?.modelIdent).to.equal('records')
  })
})
