/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}', './index.html'],
  theme: {
    extend: {
      fontFamily: {
        rajdhani: ['Rajdhani', 'sans-serif'],
        inter: ['Inter', 'sans-serif'],
      },
      keyframes: {
        'pulse-glow': {
          '0%, 100%': { boxShadow: '0 0 10px rgba(249,115,22,0.3), 0 0 20px rgba(249,115,22,0.1)' },
          '50%': { boxShadow: '0 0 25px rgba(249,115,22,0.6), 0 0 50px rgba(249,115,22,0.2)' },
        },
        'fade-in-up': {
          from: { opacity: '0', transform: 'translateY(20px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        }
      },
      animation: {
        'pulse-glow': 'pulse-glow 2s ease-in-out infinite',
        'fade-in-up': 'fade-in-up 0.4s ease-out forwards',
      }
    },
  },
  plugins: [],
}
