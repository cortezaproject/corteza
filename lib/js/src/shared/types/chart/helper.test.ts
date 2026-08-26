import { expect } from 'chai'

import { getColorschemeColors } from './helper'

const custom = [{ id: 'custom-1', name: 'Brand', colors: ['#FF00AA', '#00FFD5'] }]

describe('getColorschemeColors', () => {
  it('resolves a built-in scheme by its family-qualified name', () => {
    expect(getColorschemeColors('tableau.Tableau10')).to.have.lengthOf(10)
  })

  it('resolves a custom scheme against the supplied list', () => {
    expect(getColorschemeColors('custom-1', custom)).to.deep.equal(['#FF00AA', '#00FFD5'])
  })

  // Every branch below used to hand echarts an undefined or empty palette,
  // which draws the legend and none of the series — a blank panel that reads
  // as "no data" rather than as a colour scheme that has gone missing.
  it('falls back to the default palette when no scheme is named', () => {
    expect(getColorschemeColors()).to.have.lengthOf(13)
  })

  it('falls back when a custom id is not in the list', () => {
    expect(getColorschemeColors('custom-gone', custom)).to.deep.equal(getColorschemeColors())
  })

  it('falls back when a custom id is named but no list is supplied', () => {
    expect(getColorschemeColors('custom-1')).to.deep.equal(getColorschemeColors())
  })

  it('falls back on a built-in name that does not exist', () => {
    expect(getColorschemeColors('tableau.NoSuchScheme')).to.deep.equal(getColorschemeColors())
  })

  it('hands back a fresh default array, so a caller cannot poison the next one', () => {
    getColorschemeColors()[0] = '#000000'
    expect(getColorschemeColors()[0]).to.equal('#37A2DA')
  })
})
