import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'

export class Prompt {
  public ref = ''
  public sessionID = NoID
  public stateID = NoID
  public createdAt?: Date = undefined
  public payload: any = undefined

  constructor (u?: Partial<Prompt>) {
    this.apply(u)
  }

  apply (u?: Partial<Prompt>): void {
    Apply(this, u, CortezaID, 'sessionID', 'stateID')
    Apply(this, u, String, 'ref')
    Apply(this, u, ISO8601Date, 'createdAt')

    if (u?.payload) {
      this.payload = u.payload
    }
  }
}
