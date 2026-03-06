import type { Component } from 'vue'
import { defineAsyncComponent } from 'vue'

const BLOCK_REGISTRY: Record<string, Component> = {
  Chart: defineAsyncComponent(() => import('./Blocks/ChartBlock.vue')),
  Content: defineAsyncComponent(() => import('./Blocks/ContentBlock.vue')),
  Record: defineAsyncComponent(() => import('./Blocks/RecordBlock.vue')),
  RecordList: defineAsyncComponent(() => import('./Blocks/RecordListBlock.vue')),
  Metric: defineAsyncComponent(() => import('./Blocks/MetricBlock.vue')),
  IFrame: defineAsyncComponent(() => import('./Blocks/IFrameBlock.vue')),
  File: defineAsyncComponent(() => import('./Blocks/FileBlock.vue')),
  Automation: defineAsyncComponent(() => import('./Blocks/AutomationBlock.vue')),
  Calendar: defineAsyncComponent(() => import('./Blocks/CalendarBlock.vue')),
  Comment: defineAsyncComponent(() => import('./Blocks/CommentBlock.vue')),
}

export function resolveBlock(kind: string): Component | null {
  return BLOCK_REGISTRY[kind] || null
}
