import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: ["class"],
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "var(--ink)",
        muted: "var(--muted)",
        surface: "var(--surface)",
        canvas: "var(--canvas)",
        accent: "var(--accent)",
        "accent-2": "var(--accent-2)",
        ring: "var(--ring)",
        border: "var(--border)"
      },
      fontFamily: {
        sans: ["var(--font-sans)", "system-ui", "sans-serif"],
        serif: ["var(--font-serif)", "serif"]
      },
      boxShadow: {
        glow: "0 24px 60px -24px rgba(15, 23, 42, 0.35)",
        soft: "0 12px 30px -18px rgba(15, 23, 42, 0.4)"
      },
      keyframes: {
        "float-slow": {
          "0%, 100%": { transform: "translateY(0px)" },
          "50%": { transform: "translateY(-12px)" }
        },
        "fade-up": {
          "0%": { opacity: "0", transform: "translateY(18px)" },
          "100%": { opacity: "1", transform: "translateY(0)" }
        },
        shimmer: {
          "0%": { backgroundPosition: "0% 50%" },
          "100%": { backgroundPosition: "100% 50%" }
        }
      },
      animation: {
        "float-slow": "float-slow 10s ease-in-out infinite",
        "fade-up": "fade-up 0.7s ease-out both",
        shimmer: "shimmer 7s ease infinite"
      }
    }
  },
  plugins: [require("tailwindcss-animate")]
};

export default config;
