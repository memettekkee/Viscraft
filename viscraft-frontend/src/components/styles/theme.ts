import { createSystem, defaultConfig, defineConfig } from "@chakra-ui/react";

/**
 * Viscraft "Cartographer's Atlas" design system theme.
 *
 * Color palette:
 *  - ink (#16140F): Shell background, primary dark
 *  - parchment (#FAF6EC): Content surfaces, cards
 *  - amber (#C9762C): Accents, CTA, borders
 *  - moss (#3E5C4E): Success states
 *  - oxblood (#8B2E2E): Error states
 *  - warmgray (#6B6555): Secondary text, muted elements
 *
 * Typography:
 *  - Fraunces (serif): Display text — headings, modal headers, project names
 *  - Inter: Body/UI text — labels, buttons, navigation
 *  - JetBrains Mono: Utility/mono — prompts, timestamps, error codes
 *
 * Surfaces:
 *  - Flat (no drop shadows), thin 1px amber-tinted borders on parchment cards
 *  - Ink background for outer shell, parchment for content areas
 *
 * Validates: Requirements 14.1, 14.2, 14.5
 */

const config = defineConfig({
  theme: {
    tokens: {
      colors: {
        ink: { value: "#16140F" },
        parchment: { value: "#FAF6EC" },
        amber: { value: "#C9762C" },
        moss: { value: "#3E5C4E" },
        oxblood: { value: "#8B2E2E" },
        warmgray: { value: "#6B6555" },
      },
      fonts: {
        display: { value: "'Fraunces', serif" },
        body: { value: "'Inter', sans-serif" },
        mono: { value: "'JetBrains Mono', monospace" },
      },
    },
    semanticTokens: {
      colors: {
        "shell.bg": { value: "{colors.ink}" },
        "surface.bg": { value: "{colors.parchment}" },
        "accent.default": { value: "{colors.amber}" },
        "success.default": { value: "{colors.moss}" },
        "error.default": { value: "{colors.oxblood}" },
        "text.muted": { value: "{colors.warmgray}" },
        "border.accent": { value: "{colors.amber}" },
      },
    },
    recipes: {
      button: {
        className: "viscraft-button",
        base: {
          fontFamily: "body",
          fontWeight: "medium",
          borderRadius: "sm",
          cursor: "pointer",
          boxShadow: "none",
          _focusVisible: {
            outline: "2px solid",
            outlineColor: "{colors.amber}",
            outlineOffset: "2px",
          },
        },
        variants: {
          variant: {
            solid: {
              bg: "{colors.amber}",
              color: "white",
              _hover: { opacity: 0.9 },
            },
            outline: {
              borderWidth: "1px",
              borderColor: "{colors.amber}",
              color: "{colors.amber}",
              bg: "transparent",
              _hover: { bg: "{colors.parchment}" },
            },
            ghost: {
              color: "{colors.amber}",
              bg: "transparent",
              _hover: { bg: "{colors.parchment}" },
            },
            subtle: {
              bg: "{colors.parchment}",
              color: "{colors.ink}",
              _hover: { opacity: 0.85 },
            },
            surface: {
              bg: "{colors.parchment}",
              color: "{colors.ink}",
              borderWidth: "1px",
              borderColor: "{colors.amber}",
              boxShadow: "none",
              _hover: { opacity: 0.9 },
            },
            plain: {
              color: "{colors.amber}",
              bg: "transparent",
            },
          },
        },
      },
      input: {
        className: "viscraft-input",
        base: {
          fontFamily: "body",
          bg: "{colors.parchment}",
          borderWidth: "1px",
          borderColor: "{colors.amber}",
          borderRadius: "sm",
          color: "{colors.ink}",
          boxShadow: "none",
          _placeholder: { color: "{colors.warmgray}" },
          _focusVisible: {
            outline: "2px solid",
            outlineColor: "{colors.amber}",
            outlineOffset: "0px",
            borderColor: "{colors.amber}",
          },
        },
      },
    },
    slotRecipes: {
      dialog: {
        className: "viscraft-dialog",
        slots: [
          "trigger",
          "backdrop",
          "positioner",
          "content",
          "title",
          "description",
          "closeTrigger",
          "header",
          "body",
          "footer",
        ],
        base: {
          content: {
            bg: "{colors.parchment}",
            borderWidth: "1px",
            borderColor: "{colors.amber}",
            borderRadius: "md",
            boxShadow: "none",
            fontFamily: "body",
          },
          header: {
            fontFamily: "display",
            color: "{colors.ink}",
          },
          title: {
            fontFamily: "display",
            color: "{colors.ink}",
          },
          body: {
            color: "{colors.ink}",
          },
          footer: {
            color: "{colors.ink}",
          },
          backdrop: {
            bg: "blackAlpha.700",
          },
        },
      },
    },
  },
  globalCss: {
    "html, body": {
      bg: "{colors.ink}",
      color: "{colors.parchment}",
      fontFamily: "body",
      margin: 0,
      padding: 0,
    },
    "*": {
      boxSizing: "border-box",
    },
  },
});

export const system = createSystem(defaultConfig, config);
