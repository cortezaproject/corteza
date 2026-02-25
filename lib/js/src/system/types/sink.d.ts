interface KVV {
  [key: string]: string[]
}
export declare class SinkRequest {
  method: string
  path: string
  host: string
  header: KVV
  query: KVV
  postForm: KVV
  username: string
  password: string
  remoteAddress: string
  rawBody: string
  constructor(r?: Partial<SinkRequest>)
}
export declare class SinkResponse {
  status: number
  header: KVV
  body: unknown
  constructor(r?: Partial<SinkResponse>)
}
export {}
