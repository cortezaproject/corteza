# VueFlow Implementation Reference

## Installation

```bash
pnpm add @vue-flow/core @vue-flow/background @vue-flow/controls
```

## Required Styles

```js
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css' // optional
import '@vue-flow/controls/dist/style.css'
```

## Basic Setup

```vue
<script setup>
import { VueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
</script>

<template>
  <VueFlow :nodes="nodes" :edges="edges">
    <Background />
    <Controls />
    <!-- Custom nodes via slots -->
    <template #node-custom="props">
      <CustomNode v-bind="props" />
    </template>
  </VueFlow>
</template>
```

## Node Structure

```js
{
  id: 'unique-id',
  type: 'custom',           // Maps to #node-custom slot
  position: { x: 0, y: 0 }, // Required
  data: { /* custom data */ }
}
```

## Edge Structure

```js
{
  id: 'edge-id',
  source: 'node-1',
  target: 'node-2',
  sourceHandle: 'yes',      // For multiple handles
  type: 'smoothstep'        // bezier | step | smoothstep | straight
}
```

## Custom Node Component

```vue
<script setup>
import { Handle, Position } from '@vue-flow/core'
defineProps(['id', 'data', 'selected'])
</script>

<template>
  <div>
    <Handle type="target" :position="Position.Top" />
    <!-- Node content -->
    <Handle type="source" :position="Position.Bottom" />
  </div>
</template>
```

## Multiple Handles (for branches)

```vue
<Handle type="source" :position="Position.Bottom" id="yes" />
<Handle type="source" :position="Position.Bottom" id="no" style="left: 75%" />
```

## Key Composables

```js
import { useVueFlow } from '@vue-flow/core'
const { addNodes, addEdges, removeNodes, findNode, fitView, onNodeClick } = useVueFlow()
```

## Events

```vue
<VueFlow
  @node-click="(event, node) => {}"
  @edge-click="(event, edge) => {}"
  @pane-click="() => {}"
/>
```

## Dagre Auto-Layout

```js
import dagre from 'dagre'

function layoutNodes(nodes, edges) {
  const g = new dagre.graphlib.Graph()
  g.setDefaultEdgeLabel(() => ({}))
  g.setGraph({ rankdir: 'TB', nodesep: 80, ranksep: 100 })

  nodes.forEach(n => g.setNode(n.id, { width: 280, height: 80 }))
  edges.forEach(e => g.setEdge(e.source, e.target))

  dagre.layout(g)

  return nodes.map(n => ({
    ...n,
    position: { x: g.node(n.id).x - 140, y: g.node(n.id).y - 40 },
  }))
}
```
