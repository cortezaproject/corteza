<template>
  <div class="rounded-lg border border-dashed border-surface p-3 bg-emphasis">
    <InputText
      v-model="draft.name"
      placeholder="Group title"
      class="w-full mb-2"
      size="small"
      fluid
      @keydown.enter="submit"
    />
    <Textarea
      v-model="draft.description"
      rows="2"
      placeholder="Description (optional)"
      class="w-full mb-2"
      size="small"
      fluid
    />
    <div class="flex items-center justify-between">
      <span class="text-xs text-muted-color">
        Link roles and resources after creating.
      </span>
      <Button
        icon="pi pi-plus"
        label="Add group"
        severity="secondary"
        size="small"
        :disabled="!draft.name.trim()"
        @click="submit"
      />
    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'

const emit = defineEmits(['create'])

const draft = reactive({ name: '', description: '' })

const submit = () => {
  if (!draft.name.trim()) return
  emit('create', {
    name: draft.name.trim(),
    description: draft.description.trim(),
    roleIds: [],
    resourceIds: [],
  })
  draft.name = ''
  draft.description = ''
}
</script>
