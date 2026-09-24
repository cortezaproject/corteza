import shared from '../../../tailwind.config.shared.js'

export default {
  ...shared,
  content: ['./src/**/*.{html,vue,ts}', '../../../lib/vue/src/**/*.{vue,js,ts}'],
}
