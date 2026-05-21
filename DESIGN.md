---
version: alpha
name: university-social-design
description: "A clean, modern campus social platform built around #667eea (brand-primary) and #764ba2 (brand-secondary) gradient accents on a white/gray canvas. The system reads approachable but polished: rounded cards with subtle shadows, purple-blue brand gradient used deliberately on actions and highlights, and a warm neutral surface ladder for content hierarchy. Body type uses the system sans with generous line-height for readability. Feed cards sit on white (#ffffff) with 12-16px hairline borders above the page canvas (#f0f2f5), and the brand gradient appears on primary CTAs, active tab indicators, and key interactive elements."

colors:
  brand-primary: "#667eea"
  brand-secondary: "#764ba2"
  brand-gradient: "linear-gradient(135deg, #667eea, #764ba2)"
  brand-gradient-alt: "linear-gradient(135deg, #1a1a2e, #16213e, #0f3460)"
  brand-primary-hover: "#8b9cf6"
  brand-primary-dark: "#4c63d2"
  on-primary: "#ffffff"

  ink: "#303133"
  ink-muted: "#606266"
  ink-subtle: "#909399"
  ink-tertiary: "#c0c4cc"

  canvas: "#f0f2f5"
  surface: "#ffffff"
  surface-1: "#fafafa"
  surface-2: "#f5f5f5"
  hairline: "#e4e7ed"
  hairline-light: "#ebeef5"

  semantic-success: "#43e97b"
  semantic-warning: "#f5a623"
  semantic-error: "#e74c3c"
  semantic-info: "#4facfe"

typography:
  display-xl:
    fontSize: 64rpx
    fontWeight: 700
    lineHeight: 1.15
  display-lg:
    fontSize: 48rpx
    fontWeight: 600
    lineHeight: 1.20
  display:
    fontSize: 40rpx
    fontWeight: 600
    lineHeight: 1.25
  headline:
    fontSize: 34rpx
    fontWeight: 600
    lineHeight: 1.30
  card-title:
    fontSize: 30rpx
    fontWeight: 600
    lineHeight: 1.35
  body-lg:
    fontSize: 30rpx
    fontWeight: 400
    lineHeight: 1.50
  body:
    fontSize: 28rpx
    fontWeight: 400
    lineHeight: 1.55
  body-sm:
    fontSize: 26rpx
    fontWeight: 400
    lineHeight: 1.50
  caption:
    fontSize: 22rpx
    fontWeight: 400
    lineHeight: 1.40
  button:
    fontSize: 28rpx
    fontWeight: 500
    lineHeight: 1.20
  eyebrow:
    fontSize: 22rpx
    fontWeight: 500
    lineHeight: 1.30
    letterSpacing: 0.5px

rounded:
  xs: 8rpx
  sm: 12rpx
  md: 16rpx
  lg: 24rpx
  xl: 32rpx
  pill: 9999rpx
  full: 50%

spacing:
  xxs: 4rpx
  xs: 8rpx
  sm: 16rpx
  md: 24rpx
  lg: 32rpx
  xl: 48rpx
  xxl: 64rpx
  section: 96rpx

components:
  card-post:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    typography: "{typography.body}"
    rounded: "{rounded.md}"
    padding: 24rpx
    border: "1px solid {colors.hairline}"
    shadow: "0 2rpx 16rpx rgba(0,0,0,0.04)"
  card-post-pressed:
    backgroundColor: "{colors.surface}"
    shadow: "0 1rpx 8rpx rgba(0,0,0,0.06)"
    transform: "scale(0.98)"
  card-user:
    backgroundColor: "{colors.surface}"
    rounded: "{rounded.lg}"
    padding: 32rpx
    shadow: "0 4rpx 20rpx rgba(0,0,0,0.06)"
  button-primary:
    backgroundColor: "{colors.brand-primary}"
    textColor: "{colors.on-primary}"
    typography: "{typography.button}"
    rounded: "{rounded.md}"
    padding: "16rpx 32rpx"
  button-primary-hover:
    backgroundColor: "{colors.brand-primary-hover}"
  button-secondary:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    border: "1px solid {colors.hairline}"
    rounded: "{rounded.md}"
    padding: "16rpx 32rpx"
  tab-active:
    color: "{colors.brand-primary}"
    border-bottom: "3px solid {colors.brand-primary}"
  input:
    backgroundColor: "{colors.surface-2}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
    padding: "20rpx 24rpx"
    border: "1px solid {colors.hairline}"
  input-focus:
    border: "1px solid {colors.brand-primary}"
  avatar:
    rounded: "{rounded.full}"
    size: 80rpx
  fab:
    backgroundColor: "{colors.brand-gradient}"
    textColor: "{colors.on-primary}"
    rounded: "{rounded.full}"

motion:
  default-easing: "cubic-bezier(0.4, 0, 0.2, 1)"
  micro: "150ms"
  fast: "200ms"
  normal: "300ms"
  slow: "500ms"
  stagger-delay: "50ms"
  entrance: "{motion.normal} {motion.default-easing}"
  press: "{motion.micro} {motion.default-easing}"
  hover: "{motion.fast} {motion.default-easing}"

---

## Overview

University Social's design leans into **approachable sophistication** — the product of a campus social app that needs to feel both trustworthy and energetic. The brand gradient (purple-blue, #667eea → #764ba2) is the system's most recognizable signal: it appears on the tab bar indicator, primary buttons, and branding elements. 

The canvas is light and airy (`{colors.canvas}` #f0f2f5), with content cards lifted onto `{colors.surface}` (#ffffff) using subtle shadows (`{components.card-post.shadow}`) and hairline borders (`{components.card-post.border}`). Interactive elements use the brand gradient selectively — a CTA button here, a "like" active state there — so the gradient stays special.

Typography relies on the system sans font stack, with `{typography.body}` (28rpx, 1.55 line-height) as the comfortable reading size for social content. Card titles and headlines step up to 500-600 weight for hierarchy. Captions and metadata use `{typography.caption}` weight 400 with a muted ink.

**Key Characteristics:**
- **Purple-blue gradient signature** — the brand lives in the gradient, not flat color.
- **Light-on-light card system** — white cards (`{colors.surface}`) on a light gray page (`{colors.canvas}`).
- **Three-tier surface ladder**: canvas (page) → surface (cards) → surface-2 (inputs).
- **Conservative use of gradient**: tab active indicator, primary CTA, empty-state illustrations.
- **Generous body rhythm**: 1.55 line-height and ~24rpx paragraph gap for social reading.
- **Hairline borders + subtle shadows**: 1px `{colors.hairline}` rules + 2rpx-4rpx elevation.

## Colors

### Brand & Gradient
- **Brand Primary** (`{colors.brand-primary}` #667eea): Flat variant of the gradient start — used when gradient isn't practical (text links, loading indicators).
- **Brand Secondary** (`{colors.brand-secondary}` #764ba2): Gradient end — rarely appears in isolation.
- **Brand Gradient** (`{colors.brand-gradient}`): The signature 135° purple-blue gradient — tab active indicator, primary buttons, app-wide branding.
- **Brand Alt Gradient** (`{colors.brand-gradient-alt}`): Deep dark blue gradient for hero/splash sections.
- **Brand Hover** (`{colors.brand-primary-hover}` #8b9cf6): Lighter primary for hovered buttons.
- **On Primary** (`{colors.on-primary}` #ffffff): Text and icons on brand backgrounds.

### Surface
- **Canvas** (`{colors.canvas}` #f0f2f5): Default page background — light warm gray.
- **Surface** (`{colors.surface}` #ffffff): Card backgrounds, sheet backgrounds, modal content.
- **Surface 1** (`{colors.surface-1}` #fafafa): Subtle elevation for nested cards.
- **Surface 2** (`{colors.surface-2}` #f5f5f5): Input backgrounds, tag/chip backgrounds.
- **Hairline** (`{colors.hairline}` #e4e7ed): 1px borders on cards, dividers.
- **Hairline Light** (`{colors.hairline-light}` #ebeef5): Lighter 1px border for table rows, inner dividers.

### Text
- **Ink** (`{colors.ink}` #303133): Primary body text, headlines — dark gray, not pure black.
- **Ink Muted** (`{colors.ink-muted}` #606266): Secondary text, section subtitles.
- **Ink Subtle** (`{colors.ink-subtle}` #909399): Tertiary text, timestamps, footer disclaimers.
- **Ink Tertiary** (`{colors.ink-tertiary}` #c0c4cc): Placeholder text, disabled state text.

### Semantic
- **Success Green** (`{colors.semantic-success}` #43e97b): Positive actions, verified badges, success toasts.
- **Warning Orange** (`{colors.semantic-warning}` #f5a623): Warnings, pending states.
- **Error Red** (`{colors.semantic-error}` #e74c3c): Errors, destructive actions.
- **Info Blue** (`{colors.semantic-info}` #4facfe): Informational badges, help text.

## Typography

### Font Family

The system uses the default system sans-serif stack across all platforms: `-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif`.

### Hierarchy

| Token | Size (rpx) | Weight | Line Height | Use |
|---|---|---|---|---|
| `{typography.display-xl}` | 64 | 700 | 1.15 | Splash/section hero |
| `{typography.display-lg}` | 48 | 600 | 1.20 | Page titles |
| `{typography.display}` | 40 | 600 | 1.25 | Modal headers |
| `{typography.headline}` | 34 | 600 | 1.30 | Section titles |
| `{typography.card-title}` | 30 | 600 | 1.35 | Card title |
| `{typography.body-lg}` | 30 | 400 | 1.50 | Lead text |
| `{typography.body}` | 28 | 400 | 1.55 | Default body |
| `{typography.body-sm}` | 26 | 400 | 1.50 | Card details |
| `{typography.caption}` | 22 | 400 | 1.40 | Timestamps, meta |
| `{typography.button}` | 28 | 500 | 1.20 | Button labels |
| `{typography.eyebrow}` | 22 | 500 | 1.30 | Section labels |

### Principles

- Body text at 28rpx (1.55 line-height) is the comfortable default for social content reading.
- Card titles are 30rpx weight 600 — bold enough to anchor a card, not heavy enough to compete with headlines.
- Captions at 22rpx with 1.40 line-height fit metadata rows without overwhelming.
- Buttons use 28rpx weight 500 — slightly denser than body for clear affordance.

## Layout

### Spacing System

- **Base unit**: 8rpx.
- **Tokens (front matter)**: `{spacing.xxs}` 4rpx · `{spacing.xs}` 8rpx · `{spacing.sm}` 16rpx · `{spacing.md}` 24rpx · `{spacing.lg}` 32rpx · `{spacing.xl}` 48rpx · `{spacing.xxl}` 64rpx · `{spacing.section}` 96rpx.
- Card interior padding: 24rpx (default post cards); 32rpx (profile cards).
- Input padding: 20rpx vertical, 24rpx horizontal.

### Page Structure

- Feed cards stack vertically with 16-24rpx gap between cards.
- Section headers have 24rpx padding bottom, 16rpx padding top.
- Tab bars are sticky at the top with a border rule separator.
- Profile page sections are separated by 32rpx gaps.

### Whitespace Philosophy

The `{colors.canvas}` background provides breathing room between cards and sections. Within cards, generous 24rpx interior padding creates clear content zones. Section separation uses `{spacing.lg}` (32rpx) between logical groups, with `{spacing.section}` (96rpx) for major page partitions.

## Elevation & Depth

| Level | Treatment | Use |
|---|---|---|
| 0 (page) | `{colors.canvas}` background | Default page surface |
| 1 (card) | `{colors.surface}` bg, 1px `{colors.hairline}` border, 0 2rpx 16rpx shadow | Post cards, default cards |
| 2 (raised) | `{colors.surface}` bg, 0 4rpx 20rpx shadow | Profile card, modals |
| 3 (popover) | `{colors.surface}` bg, 0 8rpx 40rpx shadow | Dropdown, action sheet |
| glow (accent) | `shadow-glow` 0 4rpx 24rpx rgba(102,126,234,0.3) | Branded elements |

## Shapes

| Token | Value | Use |
|---|---|---|
| `{rounded.xs}` | 8rpx | Badges, small tags, inputs |
| `{rounded.sm}` | 12rpx | Secondary inputs |
| `{rounded.md}` | 16rpx | Cards, buttons |
| `{rounded.lg}` | 24rpx | Profile cards, hero blocks |
| `{rounded.xl}` | 32rpx | Search bars, large containers |
| `{rounded.pill}` | 9999rpx | FAB, tags, pills |
| `{rounded.full}` | 50% | Avatars, icons |

- Card content images use `{rounded.sm}` (12rpx) or `{rounded.md}` (16rpx) corners.
- Avatars are always `{rounded.full}` (circle).

## Components

### Post Card (`card-post`)
The dominant surface in the app — every user post renders as a card.
- White background (`{colors.surface}`), 1px hairline border (`{colors.hairline}`), subtle shadow.
- Interior padding 24rpx, `{rounded.md}` 16rpx corners.
- Content image area with `{rounded.sm}` 12rpx corners (or 14rpx for feed).
- User avatar at 72rpx, name in body, timestamp in caption.
- Action row (like/comment/share) with 32rpx spacing between icons.
- Pressed state: same background, shadow reduces to 1px/8rpx, card scales to 0.98.

### User Card (`card-user`)
Profile header card on user pages.
- 32rpx padding, `{rounded.lg}` 24rpx corners.
- Avatar prominent (120-140rpx), name below in card-title, bio in body-sm.
- Optional gradient background overlay for decorative effect.

### Buttons

**`button-primary`** — Brand CTA (e.g., "Post", "Save", "Follow").
- Background `{colors.brand-primary}`, text `{colors.on-primary}`, type `{typography.button}`, padding 16rpx 32rpx, rounded `{rounded.md}`.
- Hover: `{colors.brand-primary-hover}`.

**`button-secondary`** — Outline/ghost CTA.
- Background `{colors.surface}`, text `{colors.ink}`, type `{typography.button}`, 1px `{colors.hairline}` border, rounded `{rounded.md}`, padding 16rpx 32rpx.

### FAB (Floating Action Button)
- Background `{colors.brand-gradient}`, text `{colors.on-primary}`, rounded `{rounded.pill}`.
- Enters with `scaleIn` animation.
- Active state: scale(0.9).

### Tab Bar
- Active tab uses `{colors.brand-primary}` color + 3px bottom border.
- Inactive tabs use `{colors.ink-muted}`.
- Tab items spaced evenly across the bar.

### Input Fields
- Background `{colors.surface-2}`, border 1px `{colors.hairline}`, `{rounded.xs}` 8rpx corners.
- Interior padding 20rpx vertical / 24rpx horizontal.
- Focus state: 1px `{colors.brand-primary}` border.

## Do's and Don'ts

### Do

- Reserve the brand gradient for CTAs, active states, and branding — not decorative fills.
- Use `{colors.surface}` for all card backgrounds against `{colors.canvas}` page.
- Apply 16rpx corner radius to cards as the default.
- Use hairline borders (`{colors.hairline}`) on cards and dividers.
- Stack cards with 16-24rpx vertical gaps.
- Use `{components.card-post-pressed.transform}` (scale 0.98) as the tactile card press.

### Don't

- Don't apply the gradient to passive elements (backgrounds, inactive tabs).
- Don't fill cards with the brand color.
- Don't use pure black `#000000` for text — always use `{colors.ink}` (#303133).
- Don't add shadows to nested elements inside cards.
- Don't animate layout properties — use only `transform` and `opacity`.
- Don't use `transition: all` — always specify properties.

## Responsive Behavior

### Breakpoints

| Name | Width | Key Changes |
|---|---|---|
| Mobile | 375-414px | Default layout, single column |
| Tablet | 768-1024px | Feed can show 2-column card grid |
| Desktop | 1280px+ | Side panels, chat sidebar visible |

### Touch Targets
- All interactive elements ≥ 44rpx tap height (44px at 1x scale).
- Buttons minimum 112rpx width.
- FAB at 112rpx × 112rpx minimum.

### Collapsing Strategy
- Feed always single-column on mobile; 2-column grid above 768px.
- Right-side panels (notifications, suggestions) collapse behind hamburger on mobile.
- User profile sections stack vertically on mobile; 2-column layout on wide screens.

## Iteration Guide

1. Reference any component by its `components:` token name.
2. When styling a new surface, decide first which level of the surface ladder it belongs to.
3. Default body text to `{typography.body}` at 28rpx weight 400.
4. Cards inherit `{components.card-post}` structure and vary by content type.
5. Apply brand gradient only to high-affordance interactive elements.
6. Use the stagger delay system (50ms per sibling) for list entrance animations.
7. Prefer `transform` and `opacity` for all animations.

## Known Gaps

- The project currently uses system fonts exclusively — no custom typefaces loaded.
- Dark mode is not yet implemented; all values reflect the light theme.
- Button active/focus states beyond pressed are not fully defined.
- Form validation styling (error borders, helper text) is minimal.
- The brand gradient is defined as SCSS, not as a CSS custom property — should be migrated.
- Component previews (storybook-style) are not implemented.

---

*Generated from `awesome-design-md` schema (9 sections). Inspired by Linear, Vercel, and Stripe design systems. See `uni.scss` for SCSS token implementation.*
