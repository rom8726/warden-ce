import React, { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { ThemeProvider as MuiThemeProvider } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';
import { getTheme, type ThemeMode } from '../theme';

// Available themes in the application
export const AVAILABLE_THEMES: ThemeMode[] = ['light', 'dark'];

type ThemeContextType = {
  mode: ThemeMode;
  toggleTheme: () => void;
  setTheme: (theme: ThemeMode) => void;
  availableThemes: ThemeMode[];
};

const ThemeContext = createContext<ThemeContextType>({
  mode: 'dark',
  toggleTheme: () => {},
  setTheme: () => {},
  availableThemes: AVAILABLE_THEMES,
});

export const useTheme = () => useContext(ThemeContext);

type ThemeProviderProps = {
  children: ReactNode;
};

export const ThemeProvider: React.FC<ThemeProviderProps> = ({ children }) => {
  // Get the theme mode from localStorage or use 'dark' as default
  const [mode, setMode] = useState<ThemeMode>(() => {
    const savedMode = localStorage.getItem('themeMode');
    // Validate that the saved mode is one of the available themes
    return AVAILABLE_THEMES.includes(savedMode as ThemeMode) 
      ? (savedMode as ThemeMode) 
      : 'dark';
  });

  // Set a specific theme
  const setTheme = (theme: ThemeMode) => {
    if (AVAILABLE_THEMES.includes(theme)) {
      setMode(theme);
    }
  };

  // Toggle through all available themes
  const toggleTheme = () => {
    setMode((prevMode) => {
      const currentIndex = AVAILABLE_THEMES.indexOf(prevMode);
      const nextIndex = (currentIndex + 1) % AVAILABLE_THEMES.length;
      return AVAILABLE_THEMES[nextIndex];
    });
  };

  // Save the theme mode to localStorage whenever it changes
  useEffect(() => {
    localStorage.setItem('themeMode', mode);
  }, [mode]);

  // Get the theme object based on the current mode
  const theme = getTheme(mode);

  return (
    <ThemeContext.Provider value={{ mode, toggleTheme, setTheme, availableThemes: AVAILABLE_THEMES }}>
      <MuiThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </MuiThemeProvider>
    </ThemeContext.Provider>
  );
};
