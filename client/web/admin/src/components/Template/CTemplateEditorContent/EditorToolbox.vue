<template>
  <b-card
    class="shadow-sm h-100"
    header-bg-variant="white"
    footer-bg-variant="white"
    no-body
  >
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <b-card-body>
      <span
        v-for="sec in sections"
        :key="sec.key"
      >
        <b-btn
          variant="light"
          class="mb-2"
          block
          @click="openSection(sec.key)"
        >
          {{ $t(sec.key) }}
        </b-btn>
        <b-collapse
          :visible="expandedSections[sec.key]"
          role="tabpanel"
          class="pb-2 px-0"
        >
          <b-list-group flush>
            <b-list-group-item
              v-for="(opt, i) in sec.options"
              :key="opt.label + i"
              class="px-0 text-wrap"
              @click="opt.onClick || (() => {})"
            >
              {{ opt.label }}
              <b-btn
                variant="link"
                class="pr-0 float-right"
                @click="copyToCb(opt.copyValue())"
              >
                <font-awesome-icon
                  v-if="opt.copyValue"
                  :icon="['far', 'copy']"
                />
              </b-btn>
            </b-list-group-item>
          </b-list-group>
        </b-collapse>

      </span>
    </b-card-body>
  </b-card>
</template>

<script>
import copy from 'copy-to-clipboard'

export default {
  i18nOptions: {
    namespaces: 'system.templates',
    keyPrefix: 'editor.content.toolbox',
  },

  props: {
    template: {
      type: Object,
      required: true,
      default: () => ({}),
    },

    partials: {
      type: Array,
      required: false,
      default: () => [],
    },
  },

  data () {
    return {
      expandedSections: {},
    }
  },

  computed: {
    sections () {
      const partials = this.partials.map(p => ({
        label: p.meta.short || p.handle,
        copyValue: () => `{{template "${p.handle}" }}`,
      }))

      const rr = []
      if (partials.length) {
        rr.push({
          key: 'partials',
          options: partials,
        })
      }

      rr.push({
        key: 'snippets.label',
        options: [
          {
            label: this.$t('snippets.interpolate'),
            copyValue: () => `{{.parameter}}`,
          },
          {
            label: this.$t('snippets.iterator'),
            copyValue: () => `{{range $index, $element := .ListOfItems}}\n\n{{end}}`,
          },
          {
            label: this.$t('snippets.funcCall'),
            copyValue: () => `{{funcName param1 param2 paramN}}`,
          },
        ],
      },
      {
        key: 'samples.label',
        options: [
          {
            label: this.$t('samples.defaultHTML'),
            copyValue: () => `<!DOCTYPE html>
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
</html>`,
          },
        ],
      })

      return rr
    },
  },

  methods: {
    openSection (sec) {
      this.$set(this.expandedSections, sec, !this.expandedSections[sec])
    },
    copyToCb: copy,
  },

}
</script>
