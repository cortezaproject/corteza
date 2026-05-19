<template>
  <div>
    <div
      v-for="(item, index) in items"
      :key="index"
      :class="[
        'expression-row border border-surface rounded-border mb-3',
        { 'expression-row--dragging': dragIndex === index, 'expression-row--drop-target': dragOverIndex === index && dragIndex !== index },
      ]"
      @dragover.prevent="onDragOver(index)"
      @dragleave="onDragLeave(index)"
      @drop.prevent="onDrop(index)"
    >
      <div
        class="expression-row__header flex items-start gap-2 px-3 py-2 hover:bg-emphasis"
        draggable="true"
        @click="item._showDetails = !item._showDetails"
        @dragstart="onDragStart(index, $event)"
        @dragend="onDragEnd"
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
  </div>
</template>

<script>
import ExpressionEditor from './ExpressionEditor.vue'
import eventBus from '../lib/eventBus'
import { components } from '@planetcrust/human-vue'

const { CInputDelete } = components

export default {
  components: {
    ExpressionEditor,
    CInputDelete,
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

  data() {
    return {
      dragIndex: -1,
      dragOverIndex: -1,
    }
  },

  methods: {
    emitChange() {
      eventBus.emit('change-detected')
    },

    onDragStart(index, event) {
      this.dragIndex = index
      event.dataTransfer.effectAllowed = 'move'
      // Required for Firefox to actually start the drag.
      try { event.dataTransfer.setData('text/plain', String(index)) } catch {}
      // Use the whole row as the drag image instead of just the handle icon.
      const row = event.currentTarget?.closest?.('.expression-row')
      if (row) {
        const r = row.getBoundingClientRect()
        event.dataTransfer.setDragImage(row, event.clientX - r.left, event.clientY - r.top)
      }
    },

    onDragOver(index) {
      if (this.dragIndex < 0 || this.dragIndex === index) return
      this.dragOverIndex = index
    },

    onDragLeave(index) {
      if (this.dragOverIndex === index) this.dragOverIndex = -1
    },

    onDrop(targetIndex) {
      if (this.dragIndex < 0 || this.dragIndex === targetIndex) return
      const next = [...this.items]
      const moved = next.splice(this.dragIndex, 1)[0]
      next.splice(targetIndex, 0, moved)
      this.$emit('update:items', next)
      this.dragIndex = -1
      this.dragOverIndex = -1
      this.emitChange()
    },

    onDragEnd() {
      this.dragIndex = -1
      this.dragOverIndex = -1
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
  transition: background 0.15s, border-color 0.15s, transform 0.15s;
}

.expression-row--dragging {
  opacity: 0.4;
}

.expression-row--drop-target {
  border-color: var(--p-primary-color);
  background: color-mix(in srgb, var(--p-primary-color) 6%, transparent);
}

.expression-row__header {
  cursor: grab;
  user-select: none;
}

.expression-row__header:active {
  cursor: grabbing;
}
</style>
