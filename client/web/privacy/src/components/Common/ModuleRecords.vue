<template>
  <div
    class="overflow-auto"
  >
    <b-card
      v-for="m in modules"
      :key="m.moduleID"
      header-class="d-flex justify-content-between bg-white border-bottom"
      class="border shadow-sm mb-3"
    >
      <template #header>
        <b-form-group
          :label="$t('module')"
          label-class="text-primary"
          class="mb-0"
        >
          {{ m.module }}
        </b-form-group>
        <b-form-group
          :label="$t('namespace')"
          label-class="text-primary"
          class="mb-0"
        >
          {{ m.namespace }}
        </b-form-group>
      </template>

      <h6
        v-if="!m.records.length"
        class="text-center"
      >
        {{ $t('no-records') }}
      </h6>

      <div
        v-for="(r, ri) in m.records"
        :key="r.recordID"
        class="mb-0"
      >
        <b-form-group
          label="RecordID"
          label-class="text-primary"
        >
          {{ r.recordID }}
        </b-form-group>

        <b-row>
          <b-col
            v-for="value in r.values"
            :key="value.name"
            cols="12"
            md="6"
            lg="6"
          >
            <b-form-group
              :label="value.name"
              label-class="text-primary"
            >
              <slot
                :namespace="{ namespaceID: m.namespaceID, name: m.namespace }"
                :module="{ moduleID: m.moduleID, name: m.module }"
                :record="r"
                :value="value"
              />
            </b-form-group>
          </b-col>
        </b-row>

        <hr
          v-if="ri < m.records.length - 1"
        >
      </div>
    </b-card>
  </div>
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'module-records',
  },

  props: {
    modules: {
      type: Array,
      required: true,
    },
  },
}
</script>
