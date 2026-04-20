import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "#f9f9ff",
        primary: {
          DEFAULT: "#0058be",
          container: "#2170e4",
        },
        surface: {
          DEFAULT: "#f9f9ff",
          container: {
            low: "#f0f3ff",
            lowest: "#ffffff",
            high: "#dee8ff",
          },
        },
        on: {
          surface: "#111c2d",
          background: "#111c2d",
        },
        outline: "#727785",
      },
      borderRadius: {
        'md': '12px', // ROUND_EIGHT
      },
      backgroundImage: {
        "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
        "gradient-conic":
          "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
      },
    },
  },
  plugins: [],
};
export default config;
