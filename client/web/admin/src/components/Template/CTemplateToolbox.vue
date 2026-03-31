<template>
  <div class="flex flex-col gap-2">
    <!-- Partials section -->
    <Panel
      v-if="partials.length"
      :header="$t('system.templates.editor.content.toolbox.partials')"
      toggleable
      :collapsed="false"
    >
      <div class="flex flex-col gap-1">
        <div
          v-for="p in partials"
          :key="p.templateID"
          class="flex items-center justify-between p-2 rounded-lg hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors cursor-pointer group"
          @click="copyToClipboard(`{{template &quot;${p.handle}&quot; }}`)"
        >
          <span class="text-sm font-medium truncate">{{ p.meta?.short || p.handle }}</span>
          <i class="pi pi-copy text-xs text-muted-color opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>
        <small v-if="!partials.length" class="text-muted-color text-center p-2">
          {{ $t('system.templates.editor.content.toolbox.noPartials') }}
        </small>
      </div>
    </Panel>

    <!-- Snippets section -->
    <Panel
      :header="$t('system.templates.editor.content.toolbox.snippets.label')"
      toggleable
      :collapsed="false"
    >
      <div class="flex flex-col gap-1">
        <div
          v-for="snippet in snippets"
          :key="snippet.label"
          class="flex items-center justify-between p-2 rounded-lg hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors cursor-pointer group"
          @click="copyToClipboard(snippet.value)"
        >
          <span class="text-sm">{{ snippet.label }}</span>
          <i class="pi pi-copy text-xs text-muted-color opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>
      </div>
    </Panel>

    <!-- Samples section -->
    <Panel
      :header="$t('system.templates.editor.content.toolbox.samples.label')"
      toggleable
      :collapsed="true"
    >
      <div class="flex flex-col gap-1">
        <div
          class="flex items-center justify-between p-2 rounded-lg hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors cursor-pointer group"
          @click="copyToClipboard(defaultHTMLSample)"
        >
          <span class="text-sm">{{ $t('system.templates.editor.content.toolbox.samples.defaultHTML') }}</span>
          <i class="pi pi-copy text-xs text-muted-color opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>
      </div>
    </Panel>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const $toast = inject('$toast')

defineProps({
  partials: {
    type: Array,
    default: () => [],
  },
})

const snippets = computed(() => [
  {
    label: t('system.templates.editor.content.toolbox.snippets.interpolate'),
    value: '{{.parameter}}',
  },
  {
    label: t('system.templates.editor.content.toolbox.snippets.iterator'),
    value: '{{range $index, $element := .ListOfItems}}\n\n{{end}}',
  },
  {
    label: t('system.templates.editor.content.toolbox.snippets.funcCall'),
    value: '{{funcName param1 param2 paramN}}',
  },
])

const defaultHTMLSample = `<!DOCTYPE html>
<html>
<head>
  <meta charset='utf-8'>
  <meta http-equiv='X-UA-Compatible' content='IE=edge'>
  <title>Title</title>
  <meta name='viewport' content='width=device-width, initial-scale=1'>
</head>
<body>
  <h1>Hello, world!</h1>
</body>
</html>`

function copyToClipboard(text) {
  navigator.clipboard
    .writeText(text)
    .then(() => {
      $toast.toastSuccess(t('system.templates.editor.content.toolbox.copied'))
    })
    .catch(() => {})
}
</script>
