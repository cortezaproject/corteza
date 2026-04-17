import type { App, InjectionKey, Plugin } from 'vue'

export type EventHandler<T = unknown> = (_payload?: T) => void

export class EventBus {
  private listeners = new Map<string, Set<EventHandler>>()

  on<T = unknown>(event: string, handler: EventHandler<T>): () => void {
    const set = this.listeners.get(event) ?? new Set<EventHandler>()
    set.add(handler as EventHandler)
    this.listeners.set(event, set)

    return () => this.off(event, handler)
  }

  off<T = unknown>(event: string, handler: EventHandler<T>): void {
    const set = this.listeners.get(event)
    if (!set) {
      return
    }

    set.delete(handler as EventHandler)
    if (set.size === 0) {
      this.listeners.delete(event)
    }
  }

  emit<T = unknown>(event: string, payload?: T): void {
    const set = this.listeners.get(event)
    if (!set) {
      return
    }

    set.forEach(handler => handler(payload))
  }
}

export const EventBusKey: InjectionKey<EventBus> = Symbol('CortezaEventBus')

export const EventBusPlugin: Plugin = {
  install(app: App) {
    const eventBus = new EventBus()

    app.config.globalProperties.$eventBus = eventBus
    app.provide(EventBusKey, eventBus)
    app.provide('$eventBus', eventBus)
  },
}
