<template>
  <div>
    <div
      v-for="(item, index) in items"
      :key="index"
      :class="[
        'expression-row border border-surface rounded-border mb-3 overflow-hidden',
        { 'expression-row--dragging': dragIndex === index },
      ]"
      draggable="true"
      @dragstart="onDragStart(index, $event)"
      @dragover.prevent
      @drop.prevent="onDrop(index)"
      @dragend="onDragEnd"
    >
      <div
        class="expression-row__header flex items-center justify-between gap-2 px-3 py-2 cursor-pointer hover:bg-emphasis"
        @click="item._showDetails = !item._showDetails"
      >
        <div class="flex-1 truncate">
          <var>{{ item.target }}</var>
          <samp v-if="item.type" class="text-muted-color ml-1">({{ item.type }})</samp>
        </div>
        <samp class="truncate text-right flex-1">{{ item[valueField] }}</samp>

        <CInputDelete
          class="expression-row__delete"
          icon="pi pi-trash"
          severity="danger"
          text
          size="small"
          :header="item.target || $t('general.label.delete')"
          :message="$t('notification.delete-confirmation')"
          @click.stop
          @confirm="$emit('remove', index)"
        />
      </div>

      <transition name="fade">
        <div v-if="item._showDetails" class="p-3 border-t border-surface bg-emphasis">
          <div class="flex flex-col gap-3">
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary text-sm">Target</label>
              <InputText v-model="item.target" placeholder="Target" @input="emitChange" />
            </div>

            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary text-sm">Type</label>
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
            </div>

            <div class="flex flex-col gap-1">
              <expression-editor
                v-model="item[valueField]"
                show-line-numbers
                @open="$emit('open-editor', index)"
                @input="emitChange"
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
    }
  },

  methods: {
    emitChange() {
      eventBus.emit('change-detected')
    },

    onDragStart(index, event) {
      this.dragIndex = index
      event.dataTransfer.effectAllowed = 'move'
    },

    onDrop(targetIndex) {
      if (this.dragIndex < 0 || this.dragIndex === targetIndex) return
      const next = [...this.items]
      const moved = next.splice(this.dragIndex, 1)[0]
      next.splice(targetIndex, 0, moved)
      this.$emit('update:items', next)
      this.dragIndex = -1
      this.emitChange()
    },

    onDragEnd() {
      this.dragIndex = -1
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
  cursor: grab;
}

.expression-row--dragging {
  opacity: 0.4;
}

.expression-row__delete {
  visibility: hidden;
}

.expression-row__header:hover .expression-row__delete {
  visibility: visible;
}
</style>
