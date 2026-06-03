<template>
  <div class="flex flex-col gap-2">
    <Fieldset
      v-for="section in sections"
      v-show="section.items.length"
      :key="section.key"
      :legend="section.legend"
    >
      <div class="flex flex-col gap-1">
        <div
          v-for="item in section.items"
          :key="item.key"
          class="flex items-center justify-between p-2 rounded-lg hover:bg-emphasis transition-colors cursor-pointer group"
          @click="copyToClipboard(item.value)"
        >
          <span class="text-sm truncate">{{ item.label }}</span>
          <i
            class="pi pi-copy text-xs text-muted-color opacity-0 group-hover:opacity-100 transition-opacity"
          />
        </div>
      </div>
    </Fieldset>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const $toast = inject('$toast')

const props = defineProps({
  partials: {
    type: Array,
    default: () => [],
  },
})

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

const sections = computed(() => [
  {
    key: 'partials',
    legend: t('system.templates.editor.content.toolbox.partials'),
    items: props.partials.map(p => ({
      key: p.templateID,
      label: p.meta?.short || p.handle,
      value: `{{template "${p.handle}" }}`,
    })),
  },
  {
    key: 'snippets',
    legend: t('system.templates.editor.content.toolbox.snippets.label'),
    items: [
      {
        key: 'interpolate',
        label: t('system.templates.editor.content.toolbox.snippets.interpolate'),
        value: '{{.parameter}}',
      },
      {
        key: 'iterator',
        label: t('system.templates.editor.content.toolbox.snippets.iterator'),
        value: '{{range $index, $element := .ListOfItems}}\n\n{{end}}',
      },
      {
        key: 'funcCall',
        label: t('system.templates.editor.content.toolbox.snippets.funcCall'),
        value: '{{funcName param1 param2 paramN}}',
      },
    ],
  },
  {
    key: 'samples',
    legend: t('system.templates.editor.content.toolbox.samples.label'),
    items: [
      {
        key: 'defaultHTML',
        label: t('system.templates.editor.content.toolbox.samples.defaultHTML'),
        value: defaultHTMLSample,
      },
    ],
  },
])

function copyToClipboard(text) {
  navigator.clipboard
    .writeText(text)
    .then(() => {
      $toast.toastSuccess(t('system.templates.editor.content.toolbox.copied'))
    })
    .catch(() => {})
}
</script>
