import React, { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { useTheme } from '@mui/material';

interface AccessibilityContextType {
  highContrast: boolean;
  largeText: boolean;
  reducedMotion: boolean;
  toggleHighContrast: () => void;
  toggleLargeText: () => void;
  toggleReducedMotion: () => void;
  resetAccessibility: () => void;
}

const AccessibilityContext = createContext<AccessibilityContextType>({
  highContrast: false,
  largeText: false,
  reducedMotion: false,
  toggleHighContrast: () => {},
  toggleLargeText: () => {},
  toggleReducedMotion: () => {},
  resetAccessibility: () => {},
});

export const useAccessibility = () => useContext(AccessibilityContext);

interface AccessibilityProviderProps {
  children: ReactNode;
}

export const AccessibilityProvider: React.FC<AccessibilityProviderProps> = ({ children }) => {
  const [highContrast, setHighContrast] = useState(() => {
    return localStorage.getItem('highContrast') === 'true';
  });
  
  const [largeText, setLargeText] = useState(() => {
    return localStorage.getItem('largeText') === 'true';
  });
  
  const [reducedMotion, setReducedMotion] = useState(() => {
    return localStorage.getItem('reducedMotion') === 'true';
  });

  const theme = useTheme();

  // Apply accessibility settings to document and Material-UI theme
  useEffect(() => {
    const root = document.documentElement;
    
    // High Contrast Mode
    if (highContrast) {
      root.style.setProperty('--high-contrast', 'true');
      root.style.setProperty('--contrast-ratio', '4.5:1');
      root.style.setProperty('--text-primary', '#000000');
      root.style.setProperty('--text-secondary', '#333333');
      root.style.setProperty('--background-primary', '#ffffff');
      root.style.setProperty('--background-secondary', '#f0f0f0');
      root.style.setProperty('--border-color', '#000000');
      
      // Apply high contrast to Material-UI components
      const style = document.createElement('style');
      style.id = 'high-contrast-styles';
      style.textContent = `
        .MuiTypography-root {
          color: #000000 !important;
        }
        .MuiPaper-root {
          background-color: #ffffff !important;
          border: 2px solid #000000 !important;
        }
        .MuiButton-root {
          border: 2px solid #000000 !important;
          color: #000000 !important;
          background-color: #ffffff !important;
        }
        .MuiIconButton-root {
          border: 2px solid #000000 !important;
          color: #000000 !important;
        }
        .MuiSwitch-root .MuiSwitch-track {
          background-color: #000000 !important;
        }
        .MuiSwitch-root .MuiSwitch-thumb {
          background-color: #ffffff !important;
          border: 2px solid #000000 !important;
        }
        
        /* Exclude Snackbar and Alert from high contrast overrides */
        .MuiSnackbar-root,
        .MuiSnackbar-root *,
        .MuiAlert-root,
        .MuiAlert-root * {
          color: inherit !important;
          background-color: inherit !important;
          border: inherit !important;
        }
        
        /* Ensure Snackbar has high z-index */
        .MuiSnackbar-root {
          z-index: 1400 !important;
        }
        
        /* Ensure Alert text is visible */
        .MuiAlert-root .MuiAlert-message {
          color: inherit !important;
          font-weight: 500 !important;
        }
      `;
      document.head.appendChild(style);
    } else {
      root.style.removeProperty('--high-contrast');
      root.style.removeProperty('--contrast-ratio');
      root.style.removeProperty('--text-primary');
      root.style.removeProperty('--text-secondary');
      root.style.removeProperty('--background-primary');
      root.style.removeProperty('--background-secondary');
      root.style.removeProperty('--border-color');
      
      // Remove high contrast styles
      const existingStyle = document.getElementById('high-contrast-styles');
      if (existingStyle) {
        document.head.removeChild(existingStyle);
      }
    }

    // Large Text Mode
    if (largeText) {
      root.style.setProperty('--font-size-multiplier', '1.2');
      root.style.setProperty('--line-height-multiplier', '1.4');
      
      // Apply large text to Material-UI components
      const style = document.createElement('style');
      style.id = 'large-text-styles';
      style.textContent = `
        .MuiTypography-root {
          font-size: 1.2em !important;
          line-height: 1.6 !important;
        }
        .MuiButton-root {
          font-size: 1.1em !important;
          padding: 12px 20px !important;
        }
        .MuiIconButton-root {
          font-size: 1.2em !important;
        }
        .MuiMenuItem-root {
          font-size: 1.1em !important;
          padding: 12px 16px !important;
        }
        .MuiSwitch-root {
          transform: scale(1.2) !important;
        }
        
        /* Exclude Snackbar and Alert from large text overrides */
        .MuiSnackbar-root,
        .MuiSnackbar-root *,
        .MuiAlert-root,
        .MuiAlert-root * {
          font-size: inherit !important;
          line-height: inherit !important;
          padding: inherit !important;
          transform: inherit !important;
        }
      `;
      document.head.appendChild(style);
    } else {
      root.style.removeProperty('--font-size-multiplier');
      root.style.removeProperty('--line-height-multiplier');
      
      // Remove large text styles
      const existingStyle = document.getElementById('large-text-styles');
      if (existingStyle) {
        document.head.removeChild(existingStyle);
      }
    }

    // Reduced Motion Mode
    if (reducedMotion) {
      root.style.setProperty('--reduced-motion', 'true');
      root.style.setProperty('--animation-duration', '0.1s');
      
      // Apply reduced motion to Material-UI components
      const style = document.createElement('style');
      style.id = 'reduced-motion-styles';
      style.textContent = `
        * {
          animation-duration: 0.1s !important;
          transition-duration: 0.1s !important;
        }
        .MuiButton-root {
          transition: none !important;
        }
        .MuiIconButton-root {
          transition: none !important;
        }
        .MuiSwitch-root {
          transition: none !important;
        }
        .MuiPaper-root {
          transition: none !important;
        }
      `;
      document.head.appendChild(style);
    } else {
      root.style.removeProperty('--reduced-motion');
      root.style.removeProperty('--animation-duration');
      
      // Remove reduced motion styles
      const existingStyle = document.getElementById('reduced-motion-styles');
      if (existingStyle) {
        document.head.removeChild(existingStyle);
      }
    }
  }, [highContrast, largeText, reducedMotion]);

  // Save settings to localStorage
  useEffect(() => {
    localStorage.setItem('highContrast', highContrast.toString());
  }, [highContrast]);

  useEffect(() => {
    localStorage.setItem('largeText', largeText.toString());
  }, [largeText]);

  useEffect(() => {
    localStorage.setItem('reducedMotion', reducedMotion.toString());
  }, [reducedMotion]);

  const toggleHighContrast = () => {
    setHighContrast(prev => !prev);
  };

  const toggleLargeText = () => {
    setLargeText(prev => !prev);
  };

  const toggleReducedMotion = () => {
    setReducedMotion(prev => !prev);
  };

  const resetAccessibility = () => {
    setHighContrast(false);
    setLargeText(false);
    setReducedMotion(false);
  };

  // Add base accessibility styles
  useEffect(() => {
    const style = document.createElement('style');
    style.id = 'base-accessibility-styles';
    style.textContent = `
      :root {
        --focus-outline: 2px solid ${theme.palette.primary.main};
        --focus-outline-offset: 2px;
        --skip-link-bg: ${theme.palette.primary.main};
        --skip-link-color: ${theme.palette.primary.contrastText};
      }

      /* Skip link for keyboard navigation */
      .skip-link {
        position: absolute;
        top: -40px;
        left: 6px;
        background: var(--skip-link-bg);
        color: var(--skip-link-color);
        padding: 8px;
        text-decoration: none;
        border-radius: 4px;
        z-index: 10000;
        transition: top 0.3s;
      }

      .skip-link:focus {
        top: 6px;
      }

      /* Focus styles */
      *:focus {
        outline: var(--focus-outline);
        outline-offset: var(--focus-outline-offset);
      }

      /* Remove focus outline for mouse users */
      *:focus:not(:focus-visible) {
        outline: none;
      }

      /* Ensure focus is visible for keyboard users */
      *:focus-visible {
        outline: var(--focus-outline);
        outline-offset: var(--focus-outline-offset);
      }
    `;
    
    document.head.appendChild(style);
    
    return () => {
      const existingStyle = document.getElementById('base-accessibility-styles');
      if (existingStyle) {
        document.head.removeChild(existingStyle);
      }
    };
  }, [theme]);

  return (
    <AccessibilityContext.Provider
      value={{
        highContrast,
        largeText,
        reducedMotion,
        toggleHighContrast,
        toggleLargeText,
        toggleReducedMotion,
        resetAccessibility,
      }}
    >
      <div
        data-high-contrast={highContrast}
        data-large-text={largeText}
        data-reduced-motion={reducedMotion}
      >
        {children}
      </div>
    </AccessibilityContext.Provider>
  );
}; 