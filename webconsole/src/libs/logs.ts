import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig } from 'axios'

export interface LogEntry {
  ts: Date
  level: string
  logger: string
  msg: string
  index: number
  extra?: Record<string, unknown>
}

interface LogEntryResponse extends Omit<LogEntry, 'ts'> {
  ts: string
  [_:string]: unknown
}

interface LogFetchParams {
  after?: number
  limit?: number
}


let baseURL = '/console'

const storedBaseURL = window.localStorage.getItem('console-api-base-url')
if (storedBaseURL !== null) {
  baseURL = storedBaseURL
  console.warn('using base URL from local store (key: console-api-base-url)', { baseURL })
}

export async function fetchLoggedEvents (params: LogFetchParams = {}): Promise<Array<LogEntry>> {
  const config: AxiosRequestConfig = {
    params
  }

  return api()
    .get<Array<LogEntryResponse>>('/server-log-feed.json', config)
    .then(({ data }) => {
      return data.map(({ ts, index, level, logger, msg, ...extra }) => ({
        ts: new Date(Date.parse(ts)),
        index,
        level,
        logger,
        msg,
        extra: Object.getOwnPropertyNames(extra).length > 0 ? extra : undefined
      }))
    })
}

function api(): AxiosInstance {
  return axios.create({ baseURL })
}
