<template>
  <div>
    <b-card
      v-if="['incl', 'excl'].includes(gatewayKind)"
      class="flex-grow-1 border-bottom border-light rounded-0"
    >
      <b-card-header
        header-tag="header"
        class="p-0 mb-3"
      >
        <h5
          class="mb-0"
        >
          {{ $t('configurator:configuration') }}
        </h5>
      </b-card-header>

      <b-card-body
        class="p-0"
      >
        <var
          v-if="outEdges < 2"
        >
          {{ $t('steps:gateway.configurator.two-paths') }}
        </var>

        <div
          v-else
        >
          <b-form-group
            v-for="edge in gatewayEdges"
            :key="edge.id"
            :label="edge.value"
            label-class="text-primary"
          >
            <expression-editor
              v-model="edge.expr"
              lang="javascript"
              height="100"
              show-line-numbers
              :show-popout="false"
              @input="updateEdge(edge.id, $event)"
            />
          </b-form-group>
        </div>
      </b-card-body>
    </b-card>
  </div>
</template>

<script>
import base from './base'
import ExpressionEditor from '../ExpressionEditor.vue'

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
                expr: this.edges[id].config.expr || '',
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
      this.edges[id].config.expr = expr
      this.$root.$emit('change-detected')
    },
  },
}
</script>
