/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
    theme: {
        container: {
            center: true,
            screens: {},
        },
        screens: {
            sm: "640px",
            md: "768px",
            lg: "1024px",
            xl: "1440px",
            "2xl": "1920px",
        },
        extend: {
            colors: {
                theme: "#FAF5EA",
                themeRed: "#7269FB",
                color333: "#333333",
                color666: "#666666",
                color999: "#999999",
                inputBac: '#F5F5F5',
                default: '#0A2540',
                normal:'#33506D',
                themeBlue:'#7269FB',
            },
            fontFamily: {
                llama: "LlamaFamily, MiSans",
            },
            fontSize: {
                size6: "0.375rem",
                size8: "0.5rem",
                size10: "0.625rem",
                size12: "0.75rem",
                size14: "0.875rem",
                size16: "1rem",
                size18: "1.125rem",
                size20: "1.25rem",
                size22: "1.375rem",
                size26: "1.625rem",
                size28: "1.75rem",
                size32: "2rem",
                size34: "2.125rem",
                size38: "2.375rem",
                size40: "2.5rem",
                size64: "4rem",
            },
        },
    },
    plugins: [
        function ({addComponents}) {
            addComponents({
                ".container": {
                    "@screen md": {
                        maxWidth: "768px",
                    },
                    "@screen lg": {
                        maxWidth: "1000px",
                    },
                    "@screen xl": {
                        maxWidth: "1200px",
                    },
                },
            });
        },
    ],
};
