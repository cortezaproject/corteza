import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig } from 'axios'

export interface LogEntry {
  ts: Date
  level: string
  logger: string
  msg: string
  index: number
  extra: Record<string, unknown>
}

interface LogEntryResponse extends Omit<LogEntry, 'ts'> {
  ts: string
  [_:string]: unknown
}

interface LogFetchParams {
  after?: number
  limit?: number
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
        extra
      }))
    })
}

function api(): AxiosInstance {
  let baseURL = '/console'
  if (!Object.getOwnPropertyNames(window).includes('CortezaAPI')) {
    baseURL = ((window as unknown) as { CortezaAPI: string }).CortezaAPI
  }

  return axios.create({ baseURL })
}
