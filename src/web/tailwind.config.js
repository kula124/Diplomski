module.exports = {
  purge: [],
  darkMode: false, // or 'media' or 'class'
  theme: {
    boxShadow: {
      '2x1t': "0px 0px 50px 0px rgba(0, 128, 128, 0.66)"
    },
    extend: {
      colors: {
        teal: {
          DEFAULT: "teal"
        }
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [
    require('tailwind-scrollbar')
  ],
}
