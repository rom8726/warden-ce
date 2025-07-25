import React from 'react';

// Simplified LicenseGuard for community edition
// Always allows access without checking license status
interface LicenseGuardProps {
  children: React.ReactNode;
}

const LicenseGuard: React.FC<LicenseGuardProps> = ({ children }) => {
  // In community edition, always render children
  return <>{children}</>;
};

export default LicenseGuard; 