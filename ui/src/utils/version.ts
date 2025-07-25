import { versionInfo } from '../version';

/**
 * Get current frontend version
 */
export const getFrontendVersion = (): string => {
  return versionInfo.version;
};

/**
 * Get current frontend build time
 */
export const getFrontendBuildTime = (): string => {
  return versionInfo.buildTime;
};

/**
 * Get complete version info object
 */
export const getVersionInfo = () => {
  return {
    ...versionInfo,
    // Add additional metadata if needed
    userAgent: navigator.userAgent,
    timestamp: new Date().toISOString(),
  };
};

/**
 * Format version for display
 */
export const formatVersion = (version: string): string => {
  return `v${version}`;
};

/**
 * Check if current version is development
 */
export const isDevelopmentVersion = (): boolean => {
  return versionInfo.version === 'dev';
};

/**
 * Format frontend version for login page display
 * If it's a semantic version (tag), display as is (e.g., v1.2.3)
 * If it's a commit hash, display as v0.0.0
 */
export const formatFrontendVersionForLogin = (): string => {
  const version = versionInfo.version;
  
  // Check if it's a semantic version (e.g., 1.2.3, v1.2.3, 1.2.3-beta.1)
  const semanticVersionPattern = /^v?\d+\.\d+\.\d+(-[a-zA-Z0-9.-]+)?$/;
  
  if (semanticVersionPattern.test(version)) {
    // If it's already prefixed with 'v', return as is, otherwise add 'v'
    return version.startsWith('v') ? version : `v${version}`;
  }
  
  // If it's a commit hash or any other format, return v0.0.0
  return 'v0.0.0';
}; 