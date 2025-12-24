/** @type {import('tailwindcss').Config} */
export default {
  content: ["./pkg/view/**/*.templ"],
  theme: {
    extend: {
      animation: {
        'orbit': 'orbit 2s linear infinite',
      },
      keyframes: {
        orbit: {
          '0%': { transform: 'rotate(0deg) translateX(20px) rotate(0deg)' },
          '100%': { transform: 'rotate(360deg) translateX(20px) rotate(-360deg)' },
        }
      }
    }
  },
  plugins: [],
}

