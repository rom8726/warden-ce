import React, { createContext, useContext, type ReactNode } from 'react';

// Using Record<string, never> as a type-safe empty object
type ConfigContextType = Record<string, never>;

// Create the context with a default value
const ConfigContext = createContext<ConfigContextType>({});

// Custom hook to use the config context
export const useConfig = (): ConfigContextType => {
  return useContext(ConfigContext);
};

interface ConfigProviderProps {
  children: ReactNode;
}

export const ConfigProvider: React.FC<ConfigProviderProps> = ({ children }) => {
  return (
    <ConfigContext.Provider value={{}}>
      {children}
    </ConfigContext.Provider>
  );
};