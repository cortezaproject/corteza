<template>
  <CDraggableList v-model="ordered" handle=".expression-row__header" @update="emitChange">
    <div
      v-for="(item, index) in items"
      :key="index"
      data-drag-item
      class="expression-row border border-surface rounded-border mb-3"
    >
      <!-- The header is the grab target; the details below it hold inputs. -->
      <div
        class="expression-row__header flex items-start gap-2 px-3 py-2 hover:bg-emphasis cursor-grab"
        @click="item._showDetails = !item._showDetails"
      >
        <div class="flex-1 min-w-0 flex flex-col gap-0.5">
          <div class="flex items-center gap-2 min-w-0">
            <div class="font-medium text-primary truncate">
              <span>{{ item.target || $t('general.label.untitled', 'Untitled') }}</span>
            </div>
            <div v-if="item.type" class="ml-auto text-xs text-muted-color truncate shrink-0">
              ({{ item.type }})
            </div>
          </div>
          <div class="text-sm truncate">
            <samp class="text-color">{{ item[valueField] }}</samp>
          </div>
        </div>
      </div>

      <transition name="fade">
        <div v-if="item._showDetails" class="p-3 border-t border-surface bg-emphasis">
          <div class="flex flex-col gap-3">
            <CFormGroup label="Target">
              <InputText v-model="item.target" placeholder="Target" @input="emitChange" />
            </CFormGroup>

            <CFormGroup label="Type">
              <Select
                v-model="item.type"
                :options="types"
                filter
                class="w-full"
                @change="emitChange"
              />
              <small v-if="getTypeDescription(item.type)" class="text-muted-color">
                {{ getTypeDescription(item.type) }}
              </small>
            </CFormGroup>

            <div class="flex flex-col gap-1">
              <expression-editor
                v-model="item[valueField]"
                show-line-numbers
                @open="$emit('open-editor', index)"
                @input="emitChange"
              />
            </div>

            <div class="flex justify-start">
              <CInputDelete
                :label="$t('steps.expressions.configurator.delete-expression', 'Delete expression')"
                icon="pi pi-trash"
                severity="danger"
                text
                size="small"
                :header="item.target || $t('general.label.delete')"
                :message="$t('notification.delete-confirmation')"
                @confirm="$emit('remove', index)"
              />
            </div>
          </div>
        </div>
      </transition>
    </div>
  </CDraggableList>
</template>

<script>
import ExpressionEditor from './ExpressionEditor.vue'
import eventBus from '../lib/eventBus'
import { components } from '@planetcrust/human-vue'

const { CInputDelete, CDraggableList } = components

export default {
  components: {
    ExpressionEditor,
    CInputDelete,
    CDraggableList,
  },

  props: {
    valueField: {
      type: String,
      required: true,
    },

    items: {
      type: Array,
      required: true,
    },

    fields: {
      type: Array,
      required: true,
    },

    types: {
      type: Array,
      required: true,
    },
  },

  emits: ['update:items', 'remove', 'open-editor'],

  computed: {
    ordered: {
      get() {
        return this.items
      },
      set(next) {
        this.$emit('update:items', next)
      },
    },
  },

  methods: {
    emitChange() {
      eventBus.emit('change-detected')
    },

    getTypeDescription(type) {
      const typeDescriptions = {
        ID: 'Make sure to provide the ID in double quotes if you\'re using a literal value. Example "123"',
      }

      return typeDescriptions[type]
    },
  },
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.expression-row {
  transition:
    background 0.15s,
    border-color 0.15s,
    transform 0.15s;
}

.expression-row__header {
  cursor: grab;
  user-select: none;
}

.expression-row__header:active {
  cursor: grabbing;
}
</style>
