import type { Component } from 'vue'
import { defineAsyncComponent } from 'vue'

const BLOCK_REGISTRY: Record<string, Component> = {
  Content: defineAsyncComponent(() => import('./Blocks/ContentBlock.vue')),
  Record: defineAsyncComponent(() => import('./Blocks/RecordBlock.vue')),
  RecordList: defineAsyncComponent(() => import('./Blocks/RecordListBlock.vue')),
}

export function resolveBlock(kind: string): Component | null {
  return BLOCK_REGISTRY[kind] || null
}
