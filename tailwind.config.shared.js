/** @type {import('tailwindcss').Config} */
import PrimeUI from 'tailwindcss-primeui'

export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
    '../../../lib/vue/src/**/*.{vue,js,ts}',
  ],
  darkMode: ['selector', '[class~="dark"]'],
  theme: {
    container: {
      center: true,
      screens: {
        sm: '100%',
        md: '100%',
        lg: '1024px',
        xl: '1600px',
        '2xl': '1900px',
      },
    },
    extend: {
      borderColor: {
        DEFAULT: 'var(--p-content-border-color)',
      },
    },
  },
  plugins: [PrimeUI],
}
