/** @type {import('tailwindcss').Config} */
import PrimeUI from 'tailwindcss-primeui'

export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
    '../../../lib/vue/src/**/*.{vue,js,ts}',
  ],
  darkMode: ['selector', '[class~="dark"]'],
  plugins: [PrimeUI],
}
