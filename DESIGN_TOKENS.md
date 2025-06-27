# Design Tokens & Style Guide

This document defines the design system for MinIO Lite Admin, ensuring consistent visual design across all components and preventing theme breaking changes.

## Design Philosophy

**Modern/Minimal Dashboard Theme**
- Clean, professional appearance suitable for enterprise environments
- Subtle visual hierarchy with minimal distractions
- Focus on content and functionality over decorative elements
- Accessible design with proper contrast ratios

## Color Palette

### Neutral Colors (Primary)
```css
/* Gray Scale - Used for most UI elements */
bg-gray-50      /* #f9fafb - Light backgrounds, card backgrounds */
bg-gray-100     /* #f3f4f6 - Subtle backgrounds */
bg-gray-200     /* #e5e7eb - Border subtle */
bg-gray-300     /* #d1d5db - Border default */
bg-gray-400     /* #9ca3af - Text muted */
bg-gray-500     /* #6b7280 - Text secondary */
bg-gray-600     /* #4b5563 - Text default */
bg-gray-700     /* #374151 - Dark mode secondary backgrounds */
bg-gray-800     /* #1f2937 - Dark mode primary backgrounds */
bg-gray-900     /* #111827 - Dark mode page backgrounds */

/* Text Colors */
text-gray-900   /* #111827 - Primary text (light mode) */
text-gray-500   /* #6b7280 - Secondary text */
text-gray-400   /* #9ca3af - Muted text, placeholders */
text-white      /* #ffffff - Primary text (dark mode) */
```

### Accent Colors
```css
/* Blue - Primary actions, links, focus states */
border-blue-500  /* #3b82f6 - Primary action borders */
border-blue-600  /* #2563eb - Primary buttons, active states */
border-blue-400  /* #60a5fa - Hover states (dark mode) */

/* Red - Errors, destructive actions */
text-red-600     /* #dc2626 - Error text (light mode) */
text-red-400     /* #f87171 - Error text (dark mode) */
```

### Semantic Colors
```css
/* Status Indicators */
bg-green-100 text-green-800    /* Success states (light mode) */
bg-green-800 text-green-100    /* Success states (dark mode) */
bg-yellow-100 text-yellow-800  /* Warning states (light mode) */
bg-yellow-800 text-yellow-100  /* Warning states (dark mode) */
bg-red-100 text-red-800        /* Error states (light mode) */
bg-red-800 text-red-100        /* Error states (dark mode) */
```

## Typography

### Font Family
```css
/* Default system font stack (inherited from TailwindCSS defaults) */
font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
```

### Text Scales
```css
/* Headings */
text-xl font-semibold    /* Page titles (h1) - 20px */
text-lg font-medium      /* Section headers (h2) - 18px */
text-base font-medium    /* Subsection headers (h3) - 16px */

/* Body Text */
text-sm                  /* Primary body text - 14px */
text-xs                  /* Secondary text, captions - 12px */

/* Special Cases */
text-sm font-mono        /* Code, IDs, technical values */
```

### Font Weights
```css
font-semibold   /* 600 - Page titles, important headings */
font-medium     /* 500 - Section headers, button text */
font-normal     /* 400 - Body text, default weight */
```

## Spacing & Layout

### Container Widths
```css
max-w-7xl      /* 1280px - Main content container */
max-w-xl       /* 576px - Narrow content (forms, modals) */
```

### Padding & Margins
```css
/* Page Layout */
py-6 px-4 sm:px-6 lg:px-8    /* Page container padding */
p-6                          /* Card/section padding */
p-4                          /* Button/smaller element padding */

/* Component Spacing */
mb-6          /* Large vertical spacing between sections */
mb-4          /* Medium vertical spacing */
mb-2          /* Small vertical spacing */
space-x-4     /* Horizontal spacing between inline elements */
gap-4         /* Grid/flex gap spacing */
```

## Component Patterns

### Cards & Containers
```css
/* Primary Cards */
.card-primary {
  @apply bg-white dark:bg-gray-800 
         rounded-lg 
         shadow-sm 
         border border-gray-200 dark:border-gray-700 
         p-6;
}

/* Secondary Cards (nested within primary) */
.card-secondary {
  @apply bg-gray-50 dark:bg-gray-700 
         rounded-lg 
         p-4;
}
```

### Buttons
```css
/* Primary Action Buttons */
.btn-primary {
  @apply px-4 py-2 
         bg-blue-600 hover:bg-blue-700 
         text-white 
         rounded-lg 
         font-medium 
         transition-colors;
}

/* Secondary/Outline Buttons */
.btn-secondary {
  @apply px-4 py-2 
         border border-gray-200 dark:border-gray-600 
         hover:border-blue-500 dark:hover:border-blue-400 
         rounded-lg 
         font-medium 
         transition-colors;
}
```

### Interactive States
```css
/* Hover Effects */
hover:border-blue-500 dark:hover:border-blue-400    /* Interactive borders */
hover:bg-gray-50 dark:hover:bg-gray-700             /* Background hovers */

/* Focus States */
focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2

/* Disabled States */
disabled:opacity-50 disabled:cursor-not-allowed
```

## Dark Mode Implementation

### Strategy
- Use TailwindCSS `dark:` utility classes
- Assume `class` strategy (manual toggle vs system preference)
- Provide equivalent contrast in both modes

### Dark Mode Patterns
```css
/* Background Progression (light → dark) */
bg-white → bg-gray-800        /* Primary backgrounds */
bg-gray-50 → bg-gray-700      /* Secondary backgrounds */
bg-gray-100 → bg-gray-600     /* Tertiary backgrounds */

/* Text Progression (light → dark) */
text-gray-900 → text-white    /* Primary text */
text-gray-600 → text-gray-300 /* Secondary text */
text-gray-400 → text-gray-400 /* Muted text (often same) */

/* Border Progression (light → dark) */
border-gray-200 → border-gray-700    /* Subtle borders */
border-gray-300 → border-gray-600    /* Prominent borders */
```

## Responsive Design

### Breakpoints (TailwindCSS defaults)
```css
sm:   640px   /* Small devices (landscape phones) */
md:   768px   /* Medium devices (tablets) */
lg:   1024px  /* Large devices (desktops) */
xl:   1280px  /* Extra large devices */
2xl:  1536px  /* 2X Extra large devices */
```

### Responsive Patterns
```css
/* Grid Layouts */
grid-cols-1 md:grid-cols-3        /* Stack on mobile, 3 cols on desktop */
grid-cols-1 sm:grid-cols-2 lg:grid-cols-4  /* Responsive grid progression */

/* Spacing */
px-4 sm:px-6 lg:px-8             /* Progressive horizontal padding */
py-4 md:py-6                     /* Progressive vertical padding */
```

## Animation & Transitions

### Loading States
```css
/* Spinner */
.spinner {
  @apply animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600;
}
```

### Transition Effects
```css
/* Standard Transitions */
transition-colors    /* Color changes (hover, focus) */
transition-all       /* Multiple property changes */
duration-200         /* Fast transitions (200ms) */
duration-300         /* Standard transitions (300ms) */
```

## Accessibility Guidelines

### Color Contrast
- Ensure minimum 4.5:1 contrast ratio for normal text
- Ensure minimum 3:1 contrast ratio for large text
- Provide visual alternatives to color-only information

### Focus Management
```css
/* Focus Rings */
focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
```

### Screen Reader Support
- Use semantic HTML elements
- Provide `aria-label` for icon-only buttons
- Use proper heading hierarchy (`h1` → `h2` → `h3`)

## Component-Specific Guidelines

### Dashboard Layout
```css
/* Page Structure */
.dashboard-layout {
  @apply min-h-screen bg-gray-50 dark:bg-gray-900;
}

/* Header */
.dashboard-header {
  @apply bg-white dark:bg-gray-800 
         shadow-sm 
         border-b border-gray-200 dark:border-gray-700;
}

/* Main Content */
.dashboard-main {
  @apply max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8;
}
```

### Status Indicators
```css
/* Server Status Cards */
.status-card {
  @apply bg-gray-50 dark:bg-gray-700 rounded-lg p-4;
}

.status-label {
  @apply text-sm font-medium text-gray-500 dark:text-gray-400;
}

.status-value {
  @apply mt-1 text-lg font-semibold text-gray-900 dark:text-white;
}
```

## Usage Examples

### Creating New Components
When creating new components, follow these patterns:

```vue
<template>
  <!-- Use semantic HTML -->
  <section class="card-primary">
    <!-- Header with proper typography -->
    <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
      Section Title
    </h2>
    
    <!-- Content with proper spacing -->
    <div class="space-y-4">
      <!-- Interactive elements with hover states -->
      <button class="btn-secondary">
        Action Button
      </button>
    </div>
  </section>
</template>
```

### Custom Utility Classes
For repeated patterns, consider creating custom utility classes:

```css
/* In a separate CSS file or style block */
@layer components {
  .status-grid {
    @apply grid grid-cols-1 md:grid-cols-3 gap-4;
  }
  
  .action-grid {
    @apply grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4;
  }
}
```

## Maintenance Notes

### Adding New Colors
- Ensure both light and dark mode variants
- Test contrast ratios with WebAIM Contrast Checker
- Update this document with new color definitions

### Breaking Changes
- Major color palette changes require version bump
- Component pattern changes should be backward compatible
- Always test in both light and dark modes

### Testing Checklist
- [ ] Test in both light and dark modes
- [ ] Verify responsive behavior on mobile/tablet/desktop
- [ ] Check color contrast ratios
- [ ] Validate with screen reader
- [ ] Test keyboard navigation

## Tools & Resources

### Development Tools
- **Browser DevTools**: Chrome/Firefox dark mode simulation
- **Tailwind Docs**: https://tailwindcss.com/docs
- **Contrast Checker**: https://webaim.org/resources/contrastchecker/

### VS Code Extensions
- **Tailwind CSS IntelliSense**: Auto-completion and class validation
- **Headwind**: Automatic class sorting

This design system ensures consistency and maintainability across the MinIO Lite Admin interface while providing clear guidelines for future development.