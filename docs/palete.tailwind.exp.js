// tailwind.config.js
module.exports = {
    content: [
      "./index.html",
      "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
      extend: {
        colors: {
          primary: "#2C7BE5",       // Azul océano
          accent: "#FF6B6B",        // Coral suave
          background: "#F7E9D7",    // Arena clara
          dark: "#2E2E2E",          // Carbón
          light: "#FFFFFF",         // Blanco puro
        },
        fontFamily: {
          sans: ["Poppins", "sans-serif"],
        },
      },
    },
    plugins: [],
  }