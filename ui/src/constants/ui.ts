// UI Constants for consistent styling and behavior

export const UI_CONSTANTS = {
  // Spacing
  SPACING: {
    XS: 0.5,
    SM: 1,
    MD: 2,
    LG: 3,
    XL: 4,
  },

  // Border radius
  BORDER_RADIUS: {
    SM: 4,
    MD: 8,
    LG: 12,
    XL: 16,
  },

  // Animation durations
  ANIMATION: {
    FAST: 150,
    NORMAL: 250,
    SLOW: 350,
  },

  // Z-index levels
  Z_INDEX: {
    DRAWER: 1200,
    APP_BAR: 1100,
    MODAL: 1300,
    TOOLTIP: 1500,
    SNACKBAR: 1400,
  },

  // Breakpoints
  BREAKPOINTS: {
    XS: 0,
    SM: 600,
    MD: 960,
    LG: 1280,
    XL: 1920,
  },

  // Page limits
  PAGINATION: {
    DEFAULT_PAGE_SIZE: 20,
    MAX_PAGE_SIZE: 100,
  },

  // Search
  SEARCH: {
    DEBOUNCE_MS: 300,
    MIN_LENGTH: 2,
  },

  // Notifications
  NOTIFICATIONS: {
    AUTO_HIDE_DURATION: 6000,
    MAX_DISPLAY: 5,
  },

  // Loading states
  LOADING: {
    SKELETON_DURATION: 1.5,
    SPINNER_SIZE: {
      SMALL: 24,
      MEDIUM: 32,
      LARGE: 48,
    },
  },

  // Colors (for reference)
  COLORS: {
    PRIMARY: 'rgb(130, 82, 255)',
    SECONDARY: '#FF5A8C',
    SUCCESS: '#4CAF50',
    WARNING: '#F9A825',
    ERROR: '#E53935',
    INFO: '#2196F3',
  },

  // Typography
  TYPOGRAPHY: {
    FONT_FAMILY: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    FONT_WEIGHTS: {
      LIGHT: 300,
      REGULAR: 400,
      MEDIUM: 500,
      SEMIBOLD: 600,
      BOLD: 700,
    },
  },

  // Layout
  LAYOUT: {
    DRAWER_WIDTH: 260,
    DRAWER_COLLAPSED_WIDTH: 72,
    APP_BAR_HEIGHT: 70,
    MAX_CONTENT_WIDTH: 1400,
  },

  // Icons
  ICONS: {
    SIZES: {
      SMALL: 16,
      MEDIUM: 24,
      LARGE: 32,
      XLARGE: 40,
    },
  },
} as const;

// Theme-specific constants
export const THEME_CONSTANTS = {
  LIGHT: {
    BACKGROUND: 'rgba(248, 249, 252, 0.95)',
    PAPER: 'rgba(255, 255, 255, 0.92)',
    TEXT_PRIMARY: 'rgb(40, 42, 54)',
    TEXT_SECONDARY: 'rgba(60, 65, 75, 0.75)',
  },
  DARK: {
    BACKGROUND: 'rgba(28, 30, 35, 0.95)',
    PAPER: 'rgba(38, 40, 48, 0.9)',
    TEXT_PRIMARY: 'rgb(245, 245, 250)',
    TEXT_SECONDARY: 'rgba(220, 220, 230, 0.7)',
  },
  BLUE: {
    BACKGROUND: 'rgba(16, 42, 66, 0.95)',
    PAPER: 'rgba(25, 55, 84, 0.9)',
    TEXT_PRIMARY: 'rgb(240, 245, 255)',
    TEXT_SECONDARY: 'rgba(220, 230, 255, 0.7)',
  },
  GREEN: {
    BACKGROUND: 'rgba(40, 85, 50, 0.9)',
    PAPER: 'rgba(50, 95, 65, 0.85)',
    TEXT_PRIMARY: 'rgb(240, 255, 245)',
    TEXT_SECONDARY: 'rgba(220, 255, 230, 0.7)',
  },
} as const;

// Issue level colors
export const ISSUE_LEVEL_COLORS = {
  fatal: '#D32F2F',
  error: '#E53935',
  exception: '#F44336',
  warning: '#F9A825',
  info: '#2196F3',
  debug: '#757575',
} as const;

// Status colors
export const STATUS_COLORS = {
  active: '#E53935',
  resolved: '#4CAF50',
  ignored: '#757575',
  muted: '#F9A825',
} as const;

// Common styles
export const COMMON_STYLES = {
  // Paper styles
  paper: {
    borderRadius: UI_CONSTANTS.BORDER_RADIUS.LG,
    backdropFilter: 'blur(12px)',
    boxShadow: '0 8px 24px 0 rgba(0, 0, 0, 0.08)',
  },

  // Button styles
  button: {
    borderRadius: UI_CONSTANTS.BORDER_RADIUS.MD,
    textTransform: 'none' as const,
    fontWeight: UI_CONSTANTS.TYPOGRAPHY.FONT_WEIGHTS.MEDIUM,
  },

  // Card styles
  card: {
    borderRadius: UI_CONSTANTS.BORDER_RADIUS.LG,
    transition: 'all 0.2s ease-in-out',
    '&:hover': {
      transform: 'translateY(-2px)',
      boxShadow: '0 12px 32px 0 rgba(0, 0, 0, 0.12)',
    },
  },

  // Input styles
  input: {
    borderRadius: UI_CONSTANTS.BORDER_RADIUS.MD,
    '& .MuiOutlinedInput-root': {
      borderRadius: UI_CONSTANTS.BORDER_RADIUS.MD,
    },
  },

  // Chip styles
  chip: {
    borderRadius: UI_CONSTANTS.BORDER_RADIUS.MD,
    fontWeight: UI_CONSTANTS.TYPOGRAPHY.FONT_WEIGHTS.MEDIUM,
  },
} as const; 