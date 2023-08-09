<template>
  <v-tour
    v-if="tour !== null"
    :name="tour.name"
    :steps="tour.steps"
    :options="tour"
    :callbacks="this.callbacks"
    @onStop="onStop"
  >
    <template slot-scope="tour">
      <transition name="fade">
        <template v-for="(step, index) of tour.steps">
          <v-step
            v-if="tour.currentStep === index"
            :key="index"
            class="shadow p-3 text-light text-center"
            :step="step"
            :previous-step="tour.previousStep"
            :next-step="tour.nextStep"
            :stop="tour.stop"
            :is-first="tour.isFirst"
            :is-last="tour.isLast"
            :labels="tour.labels"
            @onStop="onStop"
          >
            <div slot="header">
              <h5
                v-if="(step.header || {}).title"
                v-html="$t(step.header.title)"
              />
            </div>

            <div slot="content">
              <div class="p-1 mb-2">
                <div v-html="$t(step.content)" />
              </div>
            </div>

            <div
              slot="actions"
            >
              <b-button
                data-test-id="button-stop-tour"
                variant="primary"
                @click="onStop"
              >
                <template v-if="tour.isLast && !(callbacks || {}).onNextRedirect">
                  {{ $t('buttons.end') }}
                </template>
                <template v-else>
                  {{ $t('buttons.skip') }}
                </template>
              </b-button>
              <b-button
                data-test-id="button-previous"
                v-if="tour.isFirst && (callbacks || {}).onPrevRedirect"
                variant="primary"
                :href="callbacks.onPrevRedirect"
              >
                {{ $t('buttons.previous') }}
              </b-button>
              <b-button
                data-test-id="button-previous"
                v-else-if="!tour.isFirst"
                variant="primary"
                @click="tour.previousStep"
              >
                {{ $t('buttons.previous') }}
              </b-button>
              <b-button
                data-test-id="button-next"
                v-if="tour.isLast && (callbacks || {}).onNextRedirect"
                variant="primary"
                :href="callbacks.onNextRedirect"
              >
                {{ $t('buttons.next') }}
              </b-button>
              <b-button
                data-test-id="button-next"
                v-else-if="!tour.isLast"
                variant="primary"
                @click="tour.nextStep"
              >
                {{ $t('buttons.next') }}
              </b-button>
            </div>
          </v-step>
        </template>
      </transition>
    </template>
  </v-tour>
</template>
<script>
export default {
  name: 'TourComponent',
  i18nOptions: {
    namespaces: 'onboarding-tour',
  },

  props: {
    name: String,

    callbacks: {
      type: Object,
    },

    steps: {
      default: () => [],
      type: Array,
    },
  },

  computed: {
    labels () {
      return this.tour.labels
    },

    tour () {
      return  {
        name: this.name,
        steps: this.steps.map(step => {
          return {
            name: step,
            target: `[data-v-onboarding="${step}"]`,
            header: {
              title: this.$t(`steps.${step}.title`),
            },
            content: this.$t(`steps.${step}.content`),
          }
        }),
      }
    },
  },

  methods: {
    onStart () {
      if (JSON.parse(localStorage.getItem('corteza.tour')) && this.tour.steps.length) {
        this.$tours[this.tour.name].start()
      }
    },

    onStop () {
      localStorage.setItem('corteza.tour', JSON.stringify(false))
      this.$tours[this.tour.name].stop()
    },

    onStartClick () {
      localStorage.setItem('corteza.tour', JSON.stringify(true))
      this.onStart()
    },
  },
}
</script>

<style lang="scss">
.v-step {
  background: #5C5C5C;
  max-width: 320px;
  z-index: 10000;
  border-radius: 1.1rem;
}
.v-step--sticky {
  position: fixed;
  top: 50%;
  left: 50%;
  -webkit-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
}
.v-step--sticky .v-step__arrow {
  display: none;
}
.v-step__arrow,
.v-step__arrow:before {
  position: absolute;
  width: 10px;
  height: 10px;
  background: inherit;
}
.v-step__arrow {
  visibility: hidden;
}
.v-step__arrow--dark:before {
  background: #5C5C5C;
}
.v-step__arrow:before {
  visibility: visible;
  content: "";
  -webkit-transform: rotate(45deg);
  transform: rotate(45deg);
  margin-left: -5px;
}
.v-step[data-popper-placement^="top"] > .v-step__arrow {
  bottom: -5px;
}
.v-step[data-popper-placement^="bottom"] > .v-step__arrow {
  top: -5px;
}
.v-step[data-popper-placement^="right"] > .v-step__arrow {
  left: -5px;
}
.v-step[data-popper-placement^="left"] > .v-step__arrow {
  right: -5px;
}
</style>