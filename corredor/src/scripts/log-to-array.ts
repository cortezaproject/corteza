/**
 * List of (pino logger) levels we use
 * when log is received in LogToArray
 */
const levels: { [_: number]: string } = {
  60: 'FATAL',
  50: 'ERROR',
  40: 'WARN ',
  30: 'INFO ',
  20: 'DEBUG',
  10: 'TRACE',
}

/**
 * LogToArray provides stream.Writable compatible class
 * to support gathering emited logs (console.* calls) on script executuin
 */
export class LogToArray {
  protected buf: string[] = []

  public lastTime!: string
  public lastMsg!: string
  public lastLevel!: number

  constructor() {
    // informs pino to assign lastTime, lastMsg, lastLevel
    // (and some other) properties before write() is called
    ;(this as Record<symbol, unknown>)[Symbol.for('pino.metadata')] = true
  }

  /**
   * Satisfying stream.Writable
   */
  write(): void {
    const d = new Date(parseInt(this.lastTime))
    this.buf.push(`${d.toISOString()} ${levels[this.lastLevel]} ${this.lastMsg}`)
  }

  /**
   * Returns collected logs in form of array of strings
   * @returns string[]
   */
  serialize(): string[] {
    // @todo returns a array of strings
    return this.buf
  }
}
