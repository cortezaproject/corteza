<template>
  <div v-if="['incl', 'excl'].includes(gatewayKind)">
    <div class="configurator-section">
      <div v-if="outEdges < 2" class="text-muted-color italic">
        {{ $t('steps.gateway.configurator.two-paths') }}
      </div>

      <div v-else class="flex flex-col gap-4">
        <div
          v-for="edge in gatewayEdges"
          :key="edge.id"
          class="flex flex-col gap-1"
        >
          <label class="font-medium text-primary">{{ edge.value }}</label>
          <expression-editor
            v-model="edge.expr"
            show-line-numbers
            :show-popout="false"
            @input="updateEdge(edge.id, $event)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import base from './base.vue'
import ExpressionEditor from '../ExpressionEditor.vue'
import eventBus from '../../lib/eventBus'

export default {
  components: {
    ExpressionEditor,
  },

  extends: base,

  computed: {
    gatewayKind () {
      return this.item.config.ref
    },

    gatewayEdges () {
      const edges = []
      if (['incl', 'excl'].includes(this.gatewayKind)) {
        if (this.outEdges && this.item.node.edges) {
          this.item.node.edges.forEach(({ id, source, target, value = '' }) => {
            if (source.id === this.item.node.id) {
              edges.push({
                id,
                source: source.id,
                target: target.id,
                value,
                expr: (this.edges[id] && this.edges[id].config) ? this.edges[id].config.expr || '' : '',
              })
            }
          })
        }
      }
      return edges
    },
  },

  methods: {
    updateEdge (id, expr) {
      if (this.edges[id] && this.edges[id].config) {
        this.edges[id].config.expr = expr
      }
      eventBus.emit('change-detected')
    },
  },
}
</script>
