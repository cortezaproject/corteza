import type { Auth } from '../plugins/auth'
import { Make } from './url'

export interface RealtimeClientOptions {
  auth: Auth;
  onMessage: (message: MessageEvent<string>) => void;
  onError?: (event: Event) => void;
  onOpen?: (socket: WebSocket) => void;
  onClose?: (event: CloseEvent) => void;
  reconnect?: boolean;
  reconnectAttempts?: number;
  reconnectDelay?: number;
}

export function endpoint(): string {
  let { CortezaAPI, CortezaWebsocket, location } = window

  if (!CortezaWebsocket) {
    const aux = new URL(Make({ url: `${CortezaAPI}/websocket` }))
    aux.hash = ''
    aux.search = ''
    CortezaWebsocket = aux.toString()
  }

  let proto: string
  if (CortezaWebsocket.startsWith('//')) {
    proto = location.protocol
  } else {
    const sep = '://'
    ;[proto] = CortezaWebsocket.split(sep, 1)
    CortezaWebsocket = CortezaWebsocket.substring(proto.length + sep.length)
  }

  return `${proto === 'https' ? 'wss' : 'ws'}://${CortezaWebsocket}`
}

export class RealtimeClient {
  private socket?: WebSocket
  private reconnectTimer?: number
  private reconnectCount = 0
  private closedManually = false
  private readonly reconnect: boolean
  private readonly reconnectAttempts: number
  private readonly reconnectDelay: number
  private readonly authOff?: () => void

  constructor(private readonly options: RealtimeClientOptions) {
    this.reconnect = options.reconnect ?? true
    this.reconnectAttempts = options.reconnectAttempts ?? 5
    this.reconnectDelay = options.reconnectDelay ?? 3000

    this.authOff = this.options.auth.on?.('auth-token-processed', ({ accessToken }) => {
      this.sendCredentials(accessToken)
    })
  }

  connect(): WebSocket {
    this.closedManually = false
    this.socket = new WebSocket(endpoint())

    this.socket.addEventListener('open', () => {
      this.reconnectCount = 0
      this.sendCredentials(this.options.auth.accessTokenFn())
      this.options.onOpen?.(this.socket as WebSocket)
    })

    this.socket.addEventListener('message', this.options.onMessage)
    this.socket.addEventListener('error', event => {
      this.options.onError?.(event)
    })

    this.socket.addEventListener('close', event => {
      this.options.onClose?.(event)
      if (!this.closedManually) {
        this.scheduleReconnect()
      }
    })

    return this.socket
  }

  disconnect(): void {
    this.closedManually = true

    if (this.reconnectTimer) {
      window.clearTimeout(this.reconnectTimer)
      this.reconnectTimer = undefined
    }

    if (this.socket) {
      this.socket.close()
      this.socket = undefined
    }

    this.authOff?.()
  }

  get current(): WebSocket | undefined {
    return this.socket
  }

  private scheduleReconnect(): void {
    if (!this.reconnect) {
      return
    }

    if (this.reconnectCount >= this.reconnectAttempts) {
      return
    }

    this.reconnectCount += 1
    this.reconnectTimer = window.setTimeout(() => {
      this.connect()
    }, this.reconnectDelay)
  }

  private sendCredentials(accessToken?: string): void {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      return
    }

    if (!accessToken) {
      return
    }

    this.socket.send(JSON.stringify({
      '@type': 'credentials',
      '@value': { accessToken },
    }))
  }
}

export function createRealtimeClient(options: RealtimeClientOptions): RealtimeClient {
  return new RealtimeClient(options)
}
