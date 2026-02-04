---
name: minimal-ui-design
description: Use when designing or refining UIs that must be visually minimal, low-noise, and icon-forward while staying understandable for new users, especially when reducing text, consolidating controls, or streamlining dialogs, toolbars, search panels, or list results.
---

# Minimal UI Design

## Overview

Create UIs that feel calm and minimal without sacrificing clarity. Prioritize reduced copy, icon-led actions, and subtle guidance that keeps new users oriented.

## When to Use

- The user asks for minimal, clean, or low-noise UI.
- You are removing redundant text or simplifying dense layouts.
- You are switching text buttons to icon-only controls.
- You need to keep a UI beginner-friendly after simplification.
- You are redesigning search dialogs, toolbars, or action bars.

Do NOT use when the goal is a marketing/brand-heavy or storytelling UI.

## Core Principles

1. **Reduce, don’t hide**
   - Remove duplicates and verbose labels first.
   - Keep one minimal cue per area (placeholder, tooltip, or keycap row).

2. **Icon-first, never icon-only without labels**
   - Use icons for primary actions.
   - Always include `aria-label`/`title` for meaning and onboarding.
   - Prefer icon-only + tooltip; use visible text only when confusion remains.

3. **Guidance must be low-noise**
   - Use subtle keycap rows or tiny tooltips over large help panels.
   - Keep hint text to 1–4 words.

4. **Density with breathing**
   - Compact spacing, but avoid reducing font size to “fake minimal”.
   - Use consistent spacing scale; keep touch targets ≥ 44×44.

5. **Keep the editing flow keyboard-first**
   - Shortcut hints belong in a single, quiet strip.
   - Show only shortcuts that are usable in the current surface.
   - Avoid showing the same shortcut in multiple places.

## Workflow (Minimal UI Pass)

1. **Delete redundancy**
   - Remove repeated labels, “empty state” descriptions, and duplicate headers.
2. **Collapse controls**
   - Merge similar actions into a single toggle or icon button.
3. **Replace text with icons**
   - Convert button text to icons, keep labels via `aria-label`/`title`.
4. **Re-add minimal guidance**
   - Choose ONE: placeholder, tooltip, or keycap row.
   - If you add a keycap row, remove other shortcut hints nearby.
5. **Check clarity**
   - Can a new user pass a 10-second scan test?

## Quick Reference

**Do**
- Use a single shortcut bar (keycaps) when helpful.
- Keep labels in accessibility attributes.
- Use small, consistent icon sizes.
- Use subtle active states (theme-tinted, low opacity).

**Don’t**
- Add a settings toggle just to show labels.
- Show multi-line help blocks inside minimal dialogs.
- Shrink font size to “look clean”.
- Duplicate hints in multiple locations.
- Reintroduce visible labels everywhere by default.

## Common Mistakes

- **Over-hiding**: icon-only without labels or tooltips.
- **Help bloat**: adding a help panel defeats minimalism.
- **False minimal**: tiny fonts and tight line-height reduce readability.
- **Conflicting cues**: shortcut hints shown in both header and footer.
- **Global hints**: showing shortcuts that don’t apply in the current UI.

## Example (Condensed)

User: “搜索框要极简、图标化，但新手也要能用。”  
Response: “移除冗余文案，保留 1 个引导（占位符或 keycap 条），按钮改图标并补 `aria-label`，减少信息密度但保留清晰路径。”
