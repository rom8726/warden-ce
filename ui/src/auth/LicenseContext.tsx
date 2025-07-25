import React, { createContext, useContext, type ReactNode } from 'react';

// Define LicenseStatusResponse type locally for community edition
interface LicenseStatusResponse {
  license: {
    is_valid: boolean;
    features: string[];
    type: string;
    expires_at: string;
  };
}

// Simplified LicenseContextType for community edition
interface LicenseContextType {
  licenseStatus: LicenseStatusResponse | null;
  isLicenseValid: boolean;
  isLoading: boolean;
  error: string | null;
  checkLicenseStatus: () => Promise<void>;
}

const LicenseContext = createContext<LicenseContextType | undefined>(undefined);

interface LicenseProviderProps {
  children: ReactNode;
}

// Simplified LicenseProvider for community edition
// Always returns a valid license status without making API calls
export const LicenseProvider: React.FC<LicenseProviderProps> = ({ children }) => {
  // Mock license status for community edition
  const mockLicenseStatus = {
    license: {
      is_valid: true,
      features: [],
      type: 'community',
      expires_at: '9999-12-31T23:59:59Z'
    }
  } as LicenseStatusResponse;

  // No-op function for compatibility
  const checkLicenseStatus = async () => {
    // No-op in community edition
    return Promise.resolve();
  };

  const value: LicenseContextType = {
    licenseStatus: mockLicenseStatus,
    isLicenseValid: true,
    isLoading: false,
    error: null,
    checkLicenseStatus,
  };

  return (
    <LicenseContext.Provider value={value}>
      {children}
    </LicenseContext.Provider>
  );
};

export const useLicense = (): LicenseContextType => {
  const context = useContext(LicenseContext);
  if (context === undefined) {
    throw new Error('useLicense must be used within a LicenseProvider');
  }
  return context;
}; 