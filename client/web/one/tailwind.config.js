/** @type {import('tailwindcss').Config} */
import PrimeUI from 'tailwindcss-primeui'
import sharedConfig from '../../../tailwind.config.shared.js'

export default {
  presets: [sharedConfig],
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
    '../../../lib/vue/src/**/*.{vue,js,ts}',
  ],
  darkMode: ['selector', '[class~="dark"]'],
  theme: {
    extend: {
      borderColor: {
        DEFAULT: 'var(--p-content-border-color)',
      },
    },
  },
  plugins: [PrimeUI],
}
