<template>
  <div>
    <div class="flex items-center gap-2 px-4 py-2 bg-emphasis text-sm font-medium rounded-t-border">
      <div
        v-for="(field, index) in fields"
        :key="index"
        class="flex-1"
      >
        {{ field.label }}
      </div>
    </div>

    <div>
      <div
        v-for="(item, index) in items"
        :key="index"
      >
        <div
          class="flex items-center gap-2 px-2 py-2 cursor-pointer hover:bg-emphasis border-b border-surface"
          @click="item._showDetails = !item._showDetails"
        >
          <div class="p-2 grab cursor-grab">
            <i class="pi pi-bars text-muted-color" />
          </div>

          <div
            v-for="(field, i) in fields"
            :key="i"
            class="flex-1 truncate"
          >
            <div
              v-if="field.key === 'expr'"
              class="flex justify-between items-center"
            >
              <samp class="truncate">{{ item[field.key] }}</samp>

              <Button
                v-if="item._showDetails"
                icon="pi pi-trash"
                severity="danger"
                text
                rounded
                size="small"
                @click.stop="$emit('remove', index)"
              />
            </div>

            <var v-else>
              {{ field.formatter ? field.formatter(item) : item[field.key] }}
            </var>
          </div>
        </div>

        <transition name="fade">
          <div
            v-if="item._showDetails"
            class="px-4 py-3 border-b border-surface bg-emphasis"
          >
            <div class="flex flex-col gap-3">
              <div class="flex flex-col gap-1">
                <label class="font-medium text-primary text-sm">Target</label>
                <InputText
                  v-model="item.target"
                  placeholder="Target"
                  @input="emitChange"
                />
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
                <small
                  v-if="getTypeDescription(item.type)"
                  class="text-muted-color"
                >
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
  </div>
</template>

<script>
import ExpressionEditor from './ExpressionEditor.vue'
import eventBus from '../lib/eventBus'

export default {
  components: {
    ExpressionEditor,
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

  methods: {

    emitChange () {
      eventBus.emit('change-detected')
    },

    getTypeDescription (type) {
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
</style>
