// Get runtime configuration
const getConfig = () => {
  // Try to get from window.WARDEN_CONFIG first (runtime)
  if (typeof window !== 'undefined' && window.WARDEN_CONFIG) {
    return window.WARDEN_CONFIG;
  }
  
  // Fallback to build-time environment variables
  return {
    VERSION: import.meta.env.VITE_VERSION || 'dev',
    BUILD_TIME: import.meta.env.VITE_BUILD_TIME || new Date().toISOString(),
  };
};

const config = getConfig();

// Version information from runtime or build time
export const VERSION = config.VERSION;

// Parse build time with fallback
const parseBuildTime = (buildTimeStr: string | undefined): string => {
  if (!buildTimeStr) {
    return new Date().toISOString();
  }
  
  // Try to parse the build time string
  const parsed = new Date(buildTimeStr);
  if (isNaN(parsed.getTime())) {
    // If parsing fails, return current time
    return new Date().toISOString();
  }
  
  return parsed.toISOString();
};

export const BUILD_TIME = parseBuildTime(config.BUILD_TIME);

// Version info object
export const versionInfo = {
  version: VERSION,
  buildTime: BUILD_TIME,
};

export default versionInfo; 