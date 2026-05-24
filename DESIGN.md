---
version: alpha
name: university-social-design
description: "A warm, sophisticated dark-theme campus social platform built around a terracotta (#C67A6A) brand accent on a deep navy canvas. The system reads premium and approachable: dark cards with warm-tinted neutrals, terracotta brand used deliberately on primary CTAs and key interactives, and a robust campus accent palette for secondary signals. Body type uses DM Sans + Noto Sans SC with generous line-height on a 1.25 modular scale. All colors use oklch for perceptual uniformity."

colors:
  brand-primary: "oklch(0.62 0.12 16)"
  brand-primary-light: "oklch(0.70 0.10 16)"
  brand-primary-dark: "oklch(0.50 0.12 16)"
  brand-primary-hex: "#C67A6A"

  campus-gold: "oklch(0.72 0.10 70)"
  campus-green: "oklch(0.65 0.12 150)"
  campus-blue: "oklch(0.62 0.10 240)"
  campus-purple: "oklch(0.58 0.12 290)"
  campus-cyan: "oklch(0.65 0.08 200)"

  bg-app: "oklch(0.14 0.01 30)"
  bg-card: "oklch(0.17 0.015 30)"
  bg-raised: "oklch(0.20 0.015 30)"
  bg-sidebar: "oklch(0.11 0.01 30)"
  bg-hover: "oklch(0.23 0.015 30)"

  text-primary: "oklch(0.92 0.005 30)"
  text-secondary: "oklch(0.65 0.01 30)"
  text-muted: "oklch(0.45 0.01 30)"

  border-default: "oklch(0.22 0.015 30)"
  border-light: "oklch(0.26 0.015 30)"

  semantic-success: "oklch(0.65 0.12 150)"
  semantic-warning: "oklch(0.72 0.10 70)"
  semantic-danger: "oklch(0.58 0.18 30)"
  semantic-info: "oklch(0.62 0.10 240)"
  semantic-link: "oklch(0.65 0.12 16)"

typography:
  font-family: "'DM Sans', 'Noto Sans SC', 'PingFang SC', sans-serif"
  text-xs: "0.75rem (12px) — 400"
  text-sm: "0.875rem (14px) — 400"
  text-base: "1rem (16px) — 400"
  text-lg: "1.125rem (18px) — 500"
  text-xl: "1.25rem (20px) — 600"
  text-2xl: "1.5rem (24px) — 600"
  text-3xl: "1.875rem (30px) — 700"
  text-4xl: "2.25rem (36px) — 700"
  scale-ratio: "1.25 (major third)"

spacing:
  base-unit: "4px"
  "1": "0.25rem (4px)"
  "2": "0.5rem (8px)"
  "3": "0.75rem (12px)"
  "4": "1rem (16px)"
  "5": "1.25rem (20px)"
  "6": "1.5rem (24px)"
  "8": "2rem (32px)"
  "10": "2.5rem (40px)"
  "12": "3rem (48px)"
  "16": "4rem (64px)"

rounded:
  sm: "0.375rem (6px)"
  md: "0.625rem (10px)"
  lg: "1rem (16px)"
  xl: "1.5rem (24px)"
  full: "9999px"

shadows:
  xs: "0 1px 2px rgba(0,0,0,0.3)"
  sm: "0 1px 3px rgba(0,0,0,0.3)"
  md: "0 4px 12px rgba(0,0,0,0.3)"
  lg: "0 8px 24px rgba(0,0,0,0.4)"
  xl: "0 12px 36px rgba(0,0,0,0.5)"

motion:
  duration-instant: "100ms"
  duration-fast: "150ms"
  duration-normal: "200ms"
  duration-slow: "300ms"
  duration-page: "350ms"
  ease-out: "cubic-bezier(0.22, 1, 0.36, 1)"
  ease-in-out: "cubic-bezier(0.65, 0, 0.35, 1)"
  ease-spring: "cubic-bezier(0.34, 1.56, 0.64, 1)"
  stagger-delay: "60ms"

components:
  card-post:
    backgroundColor: "{colors.bg-card}"
    textColor: "{colors.text-primary}"
    rounded: "{rounded.md}"
    padding: "{spacing.6}"
    border: "1px solid {colors.border-default}"
    hover-transform: "translateY(-2px)"
    hover-shadow: "0 8px 30px rgba(198,122,106,0.15)"
  card-post-pressed:
    transform: "scale(0.98)"
  card-product:
    backgroundColor: "{colors.bg-card}"
    rounded: "{rounded.md}"
    overflow: "hidden"
    hover-transform: "translateY(-2px)"
    hover-shadow: "0 4px 16px rgba(0,0,0,0.3)"
  button-primary:
    backgroundColor: "{colors.brand-primary}"
    textColor: "{colors.text-primary}"
    rounded: "{rounded.md}"
    padding: "{spacing.4} {spacing.8}"
    hover-shadow: "0 4px 20px rgba(198,122,106,0.35)"
    glow-animation: "glow-pulse 3s ease-in-out infinite"
    active-transform: "translateY(2px)"
  button-secondary:
    backgroundColor: "transparent"
    textColor: "{colors.text-primary}"
    border: "1px solid {colors.border-default}"
    rounded: "{rounded.md}"
    padding: "{spacing.4} {spacing.8}"
  input:
    backgroundColor: "{colors.bg-app}"
    textColor: "{colors.text-primary}"
    rounded: "{rounded.md}"
    padding: "{spacing.5} {spacing.6}"
    border: "1px solid {colors.border-default}"
    focus-border: "{colors.brand-primary}"
    placeholder-color: "{colors.text-muted}"
  sidebar:
    width: "240px"
    backgroundColor: "{colors.bg-sidebar}"
    noise-overlay: "true"
  topbar:
    height: "56px"
    background: "rgba(26, 30, 46, 0.85)"
    backdrop-filter: "blur(12px) saturate(1.8)"
  stagger-entrance:
    initial-opacity: "0"
    initial-transform: "translateY(12px)"
    duration: "0.4s"
    easing: "{motion.ease-out}"
    per-item-delay: "60ms"

---

## Overview

University Social's PC design uses a **dark, warm theme** — the product of a campus social app that needs to feel both premium and approachable. The brand color (warm terracotta, `{colors.brand-primary}`) is the system's most recognizable signal: it appears on primary buttons, active states, link text, and branded accents.

The canvas is a deep navy (`{colors.bg-app}`), with content cards sitting on a slightly lighter surface (`{colors.bg-card}`) using subtle borders (`{colors.border-default}`) and layered shadows (`{shadows.md}`). Interactive elements use the terracotta brand color selectively — a CTA button, a link, a hover glow — so the brand stays special.

Typography uses **DM Sans** (Latin) and **Noto Sans SC** (CJK) with a 1.25 modular scale. Body text at `{typography.text-base}` (16px, 1.55 line-height) is the comfortable reading size for social content. A full `oklch` color system provides perceptual uniformity across all surfaces.

**Key Characteristics:**
- **Warm terracotta brand** — the brand lives in the warm earth tone, not a gradient.
- **Dark-on-dark card system** — cards (`{colors.bg-card}`) on a deeper page canvas (`{colors.bg-app}`) with warm-tinted neutrals.
- **Three-tier surface ladder**: bg-app (page) → bg-card (cards) → bg-raised (hover/modals).
- **Conservative use of brand color**: primary buttons, link text, active indicators, glow effects.
- **1.25 modular type scale**: consistent `$text-xs` through `$text-4xl` with DM Sans.
- **Hairline borders + layered shadows**: 1px `{colors.border-default}` borders + `{shadows.md}` elevation.
- **Stagger entrance system**: 60ms per-item delay for list entrance animations.
- **Noise texture overlay**: subtle SVG fractal noise on sidebar for depth.

## Colors

### Brand & Campus
- **Brand Primary** (`{colors.brand-primary}`): Warm terracotta — primary buttons, link text, active indicators, brand glow.
- **Brand Light** (`{colors.brand-primary-light}`): Lighter terracotta — hover states, secondary brand elements.
- **Brand Dark** (`{colors.brand-primary-dark}`): Deeper terracotta — pressed states.
- **Campus Gold** (`{colors.campus-gold}`): Points, achievements, top-3 rankings.
- **Campus Green** (`{colors.campus-green}`): Success states, verified badges, join confirmations.
- **Campus Blue** (`{colors.campus-blue}`): Informational badges, links.
- **Campus Purple** (`{colors.campus-purple}`): Premium/featured badges.
- **Campus Cyan** (`{colors.campus-cyan}`): Alternative accent for tags.

### Surface
- **App BG** (`{colors.bg-app}`): Deep navy — default page background.
- **Card BG** (`{colors.bg-card}`): Slightly lighter — cards, dialogs, drawers.
- **Raised** (`{colors.bg-raised}`): Hover highlights, elevated surfaces.
- **Sidebar** (`{colors.bg-sidebar}`): Darkest — navigation sidebar.
- **Hover** (`{colors.bg-hover}`): Row/card hover background.
- **Border Default** (`{colors.border-default}`): Card borders, dividers.
- **Border Light** (`{colors.border-light}`): Lighter borders for secondary dividers.

### Text
- **Primary** (`{colors.text-primary}` #E8E6E3): Headlines, body text — warm off-white.
- **Secondary** (`{colors.text-secondary}` #9B97A0): Secondary text, metadata.
- **Muted** (`{colors.text-muted}` #6B6872): Placeholder text, captions, timestamps.

### Semantic
- **Success Green**: Points earned, verified badges, success toasts.
- **Warning Gold**: Warning states, pending badges.
- **Danger Red**: Errors, destructive actions, unread badges.
- **Info Blue**: Informational badges, help text.
- **Link**: Brand terracotta — clickable text.

## Typography

### Font Family

The system uses **DM Sans** for Latin text and **Noto Sans SC** for CJK, loaded through `@fontsource` packages. Fallback stack: `'DM Sans', 'Noto Sans SC', 'PingFang SC', sans-serif`.

### Modular Scale (1.25 ratio)

| Token | Size | Weight | Use |
|---|---|---|---|
| `$text-xs` | 0.75rem (12px) | 400 | Captions, timestamps, metadata |
| `$text-sm` | 0.875rem (14px) | 400 | Card details, descriptions |
| `$text-base` | 1rem (16px) | 400 | Body text — default |
| `$text-lg` | 1.125rem (18px) | 500 | Card titles, lead text |
| `$text-xl` | 1.25rem (20px) | 600 | Section titles, page headers |
| `$text-2xl` | 1.5rem (24px) | 600 | Page titles, modal headers |
| `$text-3xl` | 1.875rem (30px) | 700 | Hero headlines |
| `$text-4xl` | 2.25rem (36px) | 700 | Display — hero emphasis |

### Principles

- Body text at `1rem` (16px) is the comfortable default for social content reading.
- Card titles use `$text-lg` (18px) weight 500 — bold enough to anchor a card.
- Captions at `$text-xs` (12px) with muted color fit metadata rows without overwhelming.
- Buttons use Element Plus default with terracotta brand background.
- H1 page headers use `22px` (between $text-xl and $text-2xl) weight 700 for visual hierarchy.

## Layout

### Spacing System

- **Base unit**: 4px (`$space-1`).
- **Tokens**: `$space-1` 4px · `$space-2` 8px · `$space-3` 12px · `$space-4` 16px · `$space-5` 20px · `$space-6` 24px · `$space-8` 32px · `$space-10` 40px · `$space-12` 48px · `$space-16` 64px.
- Card interior padding: `$space-6` (24px).
- Input padding: `$space-5` vertical, `$space-6` horizontal.
- Page max-width: 640-900px depending on content type (narrow for social, wider for products/grid).

### Page Structure

- Content area: max-width 640-780px for social pages, up to 900px for grid-based pages.
- Page header: flex row with h1 (left) and action button (right), 20px bottom margin.
- List cards stack vertically with 6-14px gap between items (vertical lists) or `auto-fill` CSS grid for product/found cards.
- Section headers use 24px bottom margin with h2 at 17px weight 600.
- Sidebar: 240px fixed width with deeper background (`{colors.bg-sidebar}`) and noise texture.

### Whitespace Philosophy

The dark canvas provides immersive contrast. Cards use `{colors.bg-card}` to lift content above the page. Generous interior padding (`$space-6`) creates clear content zones within cards. Grid layouts use 14px gap between items for visual breathing room.

## Elevation & Depth

| Level | Treatment | Use |
|---|---|---|
| 0 (page) | `{colors.bg-app}` background | Default page surface |
| 1 (card) | `{colors.bg-card}` bg, 1px `{colors.border-default}` border, `{shadows.md}` | Post cards, default cards |
| 1-hover | TranslateY(-2px), `{shadows.lg}` + terracotta glow | Card hover state |
| 2 (raised) | `{colors.bg-raised}` bg, `{shadows.lg}` | Modals, drawers, dropdowns |
| 3 (popover) | `{colors.bg-card}` bg, `{shadows.xl}` | Select dropdowns, notifications |
| glow (accent) | `box-shadow: 0 4px 20px rgba(198,122,106,0.35)` | Branded button hover |

## Shapes

| Token | Value | Use |
|---|---|---|
| `$radius-sm` | 0.375rem (6px) | Badges, small tags, user rows |
| `$radius-md` | 0.625rem (10px) | Cards, buttons |
| `$radius-lg` | 1rem (16px) | Login card, profile cards |
| `$radius-xl` | 1.5rem (24px) | Large containers |
| `$radius-full` | 9999px | Avatars, FAB, tags |

- Card content areas use `$radius-md` (10px) for consistent rounded corners.
- Avatars are always circle (`$radius-full`).
- Inputs use `$radius-md` matching card radius.

## Components

### Post Card (`card-post`)
The dominant surface in the app — every user post renders as a card.
- Background `{colors.bg-card}`, 1px `{colors.border-default}` border.
- Interior padding `{spacing.6}`, `{rounded.md}` 10px corners.
- Hover: `translateY(-2px)`, `box-shadow: 0 8px 30px rgba(198,122,106,0.15)`.
- User avatar at 36px, name in $text-base, timestamp in $text-xs.
- Action row (like/comment/share) with 32px spacing between icon buttons.
- Content images with consistent aspect ratio and object-fit.

### Product Card (`card-product`)
Grid-based card for marketplace items.
- Background `{colors.bg-card}`, `{rounded.md}`, overflow hidden.
- Image area: 160px height, object-fit cover.
- Card body: 12px padding with title, price, category, condition.
- Hover: translateY(-2px), `{shadows.lg}`, border highlights terracotta.
- Price in `$text-lg` weight 700 with brand color.

### Search Result Row
Compact list items for search results.
- Flex row with avatar, name, meta, and status tag.
- 10px padding, `$radius-sm` corners, hover: `{colors.bg-hover}`.
- 6px gap between items.

### Buttons

**`button-primary`** — Brand CTA (e.g., "发布动态", "进入校园社").
- Background `{colors.brand-primary}`, text `{colors.text-primary}`, rounded `{rounded.md}`.
- Hover: `box-shadow: 0 4px 20px rgba(198,122,106,0.35)`, lighter background.
- Active: `translateY(2px)`.
- Has `glow-pulse` animation (3s infinite).

**`button-secondary`** — Ghost/outline CTA.
- Transparent background, 1px `{colors.border-default}` border, rounded `{rounded.md}`.
- Text `{colors.text-primary}`.

### FAB (Floating Action Button)
- Background `{colors.bg-card}` with terracotta gradient tint, 1px `{colors.border-default}` border.
- Full round shape, 48×48px.
- Appears after scrolling 300px, transitions opacity.
- Click scrolls to top smoothly.

### Sidebar
- Fixed 240px width, `{colors.bg-sidebar}` background.
- SVG fractal noise overlay at 0.03 opacity for subtle texture.
- Navigation items with icons, hover highlights.
- Collapsible on mobile.

### Top Bar
- Fixed 56px height, glass effect (`backdrop-filter: blur(12px) saturate(1.8)`).
- Semi-transparent background: `rgba(26, 30, 46, 0.85)`.
- Search bar, notifications icon, user menu.

### Input Fields
- Background `{colors.bg-app}`, 1px `{colors.border-default}` border, `{rounded.md}` corners.
- Interior padding `$space-5` vertical / `$space-6` horizontal.
- Focus: 1px `{colors.brand-primary}` border.
- Placeholder: `{colors.text-muted}`.

### Notifications
- Flex row with icon circle, avatar, body content, unread dot.
- Icon circles: 32px with colored background per type (like=red, comment=blue, follow=green).
- Unread: subtle brand-tinted background + 8px terracotta dot.
- Hover: `{colors.bg-hover}` background.
- Stagger entrance on list.

## Page Transitions

- Route changes use `<Transition name="page" mode="out-in">`:
  - Enter: `opacity 0→1` + `translateY(12px→0)` over 0.25s ease-out.
  - Leave: `opacity 1→0` + `translateY(0→-8px)` over 0.15s ease-in.
- List items use `stagger-item` class for cascading entrance:
  - 0.4s animation per item with 60ms delay between siblings.
  - Animates `opacity` and `translateY`.
- Cards have hover transition: `transform 0.25s ease-out`, `box-shadow 0.25s ease`.
- `prefers-reduced-motion` respected: all animations zeroed when enabled.
- Buttons use `active: translateY(2px)` for tactile feedback.

## Element Plus Overrides

Full dark theme applied via CSS custom properties in `global.scss`:

| Element | Key Overrides |
|---|---|
| Card | bg `{colors.bg-card}`, border `{colors.border-default}` |
| Dialog | bg `{colors.bg-card}`, title text `{colors.text-primary}` |
| Table | bg `{colors.bg-card}`, header `{colors.bg-app}`, hover `{colors.bg-raised}` |
| Input | wrapper bg `{colors.bg-app}`, focus border brand |
| Select/Dropdown | bg `{colors.bg-card}`, hover `{colors.bg-raised}` |
| Tag | bg `{colors.bg-raised}`, border `{colors.border-light}` |
| Pagination | bg transparent, active color brand |
| Tabs | inactive `{colors.text-secondary}`, active brand |
| Badge | danger background for unread count |
| Empty | description text `{colors.text-secondary}` |
| Notification | bg `{colors.bg-card}`, border `{colors.border-default}` |
| Drawer | bg `{colors.bg-card}` |

## Do's and Don'ts

### Do

- Reserve the terracotta brand for CTAs, link text, and active states — not decorative fills.
- Use `{colors.bg-card}` for all card backgrounds against `{colors.bg-app}` page.
- Apply `$radius-md` (10px) corner radius to cards as the default.
- Use hairline borders (`{colors.border-default}`) on cards and dividers.
- Use `{shadows.md}` for card elevation, `{shadows.lg}` for hover.
- Use stagger-item class for list entrance animations (60ms per item).
- Use `transform` and `opacity` for all motion — never layout properties.
- Apply noise texture overlay to sidebar for visual depth.
- Use glass effect `backdrop-filter` for overlay elements (topbar).

### Don't

- Don't apply the brand color to passive elements (plain backgrounds, inactive tabs).
- Don't use pure white `#ffffff` for text — always use `{colors.text-primary}`.
- Don't use pure black `#000000` for backgrounds — always use `{colors.bg-app}`.
- Don't add shadows to nested elements inside cards.
- Don't animate layout properties (width, height, margin, padding) — use only `transform` and `opacity`.
- Don't use `transition: all` — always specify properties.
- Don't use flat rgb/hex colors where perceptual uniformity matters — prefer oklch.

## Responsive Behavior

### Breakpoints

| Name | Width | Key Changes |
|---|---|---|
| Mobile | <768px | Sidebar collapses to hamburger, single-column layouts |
| Desktop | 768px+ | Sidebar visible, grid layouts available |

### Touch Targets
- All interactive elements ≥ 44px tap height.
- Buttons minimum 112px width.
- FAB at 48×48px with adequate click area.

### Collapsing Strategy
- Sidebar collapses behind hamburger on mobile; static on desktop.
- Right panels (notifications) collapse on mobile.
- Feed always single-column; grid (2+ columns) above 768px.

## Iteration Guide

1. Reference any component by its `components:` token name.
2. When styling a new surface, decide first which level of the surface ladder it belongs to.
3. Default body text to `$text-base` at 16px weight 400.
4. Cards inherit `{components.card-post}` structure and vary by content type.
5. Apply brand color only to high-affordance interactive elements.
6. Use the stagger delay system (60ms per sibling) for list entrance animations.
7. Prefer `transform` and `opacity` for all animations.
8. All new surfaces should be wrapped in `.page-enter` / `.stagger-item` for consistency.

---

*Generated from `awesome-design-md` schema (9 sections). Inspired by Linear, Vercel, and Stripe design systems. See `variables.scss` for SCSS token implementation. Current as of 2026-05-23.*
