export declare class Prompt {
  ref: string
  sessionID: string
  stateID: string
  createdAt?: Date
  payload: any
  constructor(u?: Partial<Prompt>);
  apply(u?: Partial<Prompt>): void;
}
