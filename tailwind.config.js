export default {
	content: ["./views/**/*.html", "./**/*.go"],
	darkMode: "media",
	theme: {
		extend: {
			colors: {
				github: "#181717",
			},
			typography: (theme) => ({
				DEFAULT: {
					css: {
						"--tw-prose-quote-borders": theme("colors.sky[400]"),
						"--tw-prose-hr": theme("colors.neutral[300]"),
						"--tw-prose-invert-quote-borders": theme("colors.sky[600]"),
						"--tw-prose-invert-hr": theme("colors.neutral[700]"),
						".tag a": {
							textDecoration: "none",
						},
						pre: {
							whiteSpace: "pre-wrap",
							wordBreak: "break-word",
						},
						blockquote: {
							fontStyle: "normal",
						},
						"hr ~ ul": {
							listStyle: "none",
							fontSize: "14px",
							paddingLeft: 0,
						},
						"hr ~ ul li": {
							paddingLeft: 0,
						},
					},
				},
			}),
		},
	},
	plugins: [require("@tailwindcss/typography")],
};
