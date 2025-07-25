/**
 * Converts seconds to a human-readable time format
 * @param seconds - number of seconds (can be undefined)
 * @returns formatted time string or 'N/A' if value is undefined
 */
export const formatDuration = (seconds: number | undefined): string => {
  if (seconds === undefined || seconds === null || seconds <= 0) {
    return 'N/A';
  }
  
  if (seconds < 60) {
    return `${Math.round(seconds)}s`;
  }
  
  if (seconds < 3600) {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = Math.round(seconds % 60);
    if (remainingSeconds === 0) {
      return `${minutes}m`;
    }
    return `${minutes}m ${remainingSeconds}s`;
  }
  
  if (seconds < 86400) {
    const hours = Math.floor(seconds / 3600);
    const remainingMinutes = Math.floor((seconds % 3600) / 60);
    if (remainingMinutes === 0) {
      return `${hours}h`;
    }
    return `${hours}h ${remainingMinutes}m`;
  }
  
  const days = Math.floor(seconds / 86400);
  const remainingHours = Math.floor((seconds % 86400) / 3600);
  if (remainingHours === 0) {
    return `${days}d`;
  }
  return `${days}d ${remainingHours}h`;
};

/**
 * Converts seconds to a short format for UI display
 * @param seconds - number of seconds (can be undefined)
 * @returns short formatted time string or 'N/A' if value is undefined
 */
export const formatDurationShort = (seconds: number | undefined): string => {
  if (seconds === undefined || seconds === null || seconds <= 0) {
    return 'N/A';
  }
  
  if (seconds < 60) {
    return `${Math.round(seconds)}s`;
  }
  
  if (seconds < 3600) {
    const minutes = Math.floor(seconds / 60);
    return `${minutes}m`;
  }
  
  if (seconds < 86400) {
    const hours = Math.floor(seconds / 3600);
    return `${hours}h`;
  }
  
  const days = Math.floor(seconds / 86400);
  return `${days}d`;
}; 