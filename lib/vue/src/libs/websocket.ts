/* eslint-disable @typescript-eslint/ban-ts-ignore */
/**
 * Default websocket configuration
 */
import { Vue } from 'vue/types/vue'
import { Make } from './url'

export const config = {
  format: 'json',

  // (Boolean) whether to reconnect automatically (false)
  reconnection: true,

  // (Number) number of reconnection attempts before giving up (Infinity),
  reconnectionAttempts: 5,

  // (Number) how long to initially wait before attempting a new (1000)
  reconnectionDelay: 3000,

  connectManually: false,
}

/**
 * Extract websocket endpoint from window props (set via config.js)
 */
export function endpoint (): string {
  // eslint-disable-next-line @typescript-eslint/ban-ts-ignore
  // @ts-ignore
  let { CortezaAPI, CortezaWebsocket, location } = window

  if (!CortezaWebsocket) {
    // Corteza websocket entrypoint not set, use API and append /websocket
    //
    // When CortezaAPI is provided as a path (/api for example); make sure that
    // no fragments/query parameters are provided
    const aux = new URL(Make({ url: `${CortezaAPI}/websocket` }))
    aux.hash = ''
    aux.search = ''

    CortezaWebsocket = aux.toString()
  }

  let proto: string
  if (CortezaWebsocket.startsWith('//')) {
    // No proto in the configured API endpoint, use location
    proto = location.protocol
  } else {
    const sep = '://';
    [proto] = (CortezaWebsocket as string).split(sep, 1)
    CortezaWebsocket = CortezaWebsocket.substring(proto.length + sep.length)
  }

  return `${proto === 'https' ? 'wss' : 'ws'}://${CortezaWebsocket}`
}

/**
 * Binds auth and websocket events so that we can pass current access token
 *  - when ws connection opens
 *  - when auth token in fetched/renewed
 *
 *  @todo get rid of ts-ignore lines
 */
export function init (vue: Vue): void {
  // @ts-ignore
  if (!vue.$socket || !vue.$options) {
    // (web)socket plugin not ready.
    return
  }

  const wsAuth = (accessToken?: string): void => {
    if (accessToken && accessToken.length > 0) {
      // @ts-ignore
      vue.$socket.sendObj({ '@type': 'credentials', '@value': { accessToken } })
    }
  }

  // make sure that we send auth token as soon as we're connected
  // @ts-ignore
  vue.$options.sockets.onopen = (): void => {
    // @ts-ignore
    wsAuth(vue.$auth.accessTokenFn())

    // update connection with new access token
    //
    // If event listener is added before the connection is established
    // we might try to send the message too early.
    //
    // @ts-ignore
    vue.$on('auth-token-processed', ({ accessToken }) => wsAuth(accessToken))
  }

  // @ts-ignore
  vue.$options.sockets.onmessage = (msg): void => vue.$emit('websocket-message', msg)
}
