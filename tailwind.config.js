module.exports = {
  content: [
    './app/**/*.{ts,tsx,jsx,js}',
    './data/**/*.md',
    './node_modules/highlight.js/styles/night-owl.css',
  ],
  darkMode: 'media',
  theme: {
    extend: {
      colors: {
        github: '#181717',
        linkedin: '#0077b5',
        twitter: '#1da1f2',
        polywork: '#66ba8c',
        brandBlue: {
          50: '#E5F3FF',
          100: '#CCE7FF',
          200: '#99CFFF',
          300: '#66B8FF',
          400: '#33A0FF',
          500: '#0088FF',
          600: '#006DCC',
          700: '#005299',
          800: '#003666',
          900: '#001B33',
        },
      },
      typography: (theme) => ({
        DEFAULT: {
          css: {
            '--tw-prose-pre-bg': theme('colors.gray[800]'),
            '--tw-prose-pre-code': theme('colors.gray[800]'),
            '--tw-prose-invert-pre-bg': theme('colors.gray[800]'),
            '--tw-prose-invert-pre-code': theme('colors.gray[800]'),
            '.tag a': {
              textDecoration: 'none',
            },
            pre: {
              padding: 0,
              whiteSpace: 'pre-wrap',
              wordBreak: 'break-word',
            },
            blockquote: {
              fontStyle: 'normal',
            },
            hr: {
              marginBottom: '20px',
              marginTop: '20px',
            },
            'hr ~ ul': {
              listStyle: 'none',
              fontSize: '14px',
              paddingLeft: 0,
            },
            'hr ~ ul li': {
              paddingLeft: 0,
            },
          },
        },
      }),
    },
  },
  plugins: [require('@tailwindcss/typography')],
}
