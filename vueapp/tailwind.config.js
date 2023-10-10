/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,ts}"
  ],
  daisyui: {
    themes: ["light", "dark"],
  },
  plugins: [require("daisyui")],
}
