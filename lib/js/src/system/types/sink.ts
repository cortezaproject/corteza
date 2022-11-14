import { merge } from 'lodash'

interface KVV {
  [key: string]: string[];
}

export class SinkRequest {
  method = ''
  path = ''
  host = ''
  header: KVV = {}
  query: KVV = {}
  postForm: KVV = {}
  username = ''
  password = ''
  remoteAddress = ''
  rawBody = ''

  constructor (r: Partial<SinkRequest> = {}) {
    merge(this, r)
  }
}

export class SinkResponse {
  status = 200
  header: KVV = {}
  body: unknown

  constructor (r: Partial<SinkResponse> = {}) {
    merge(this, r)
  }
}
