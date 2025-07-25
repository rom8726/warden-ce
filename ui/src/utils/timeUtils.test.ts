import { formatDuration, formatDurationShort } from './timeUtils';

describe('timeUtils', () => {
  describe('formatDuration', () => {
    it('should format seconds correctly', () => {
      expect(formatDuration(30)).toBe('30s');
      expect(formatDuration(45)).toBe('45s');
      expect(formatDuration(59)).toBe('59s');
    });

    it('should format minutes correctly', () => {
      expect(formatDuration(60)).toBe('1m');
      expect(formatDuration(90)).toBe('1m 30s');
      expect(formatDuration(120)).toBe('2m');
      expect(formatDuration(3599)).toBe('59m 59s');
    });

    it('should format hours correctly', () => {
      expect(formatDuration(3600)).toBe('1h');
      expect(formatDuration(3660)).toBe('1h 1m');
      expect(formatDuration(7200)).toBe('2h');
      expect(formatDuration(86399)).toBe('23h 59m');
    });

    it('should format days correctly', () => {
      expect(formatDuration(86400)).toBe('1d');
      expect(formatDuration(90000)).toBe('1d 1h');
      expect(formatDuration(172800)).toBe('2d');
      expect(formatDuration(176400)).toBe('2d 1h');
    });

    it('should handle edge cases', () => {
      expect(formatDuration(0)).toBe('N/A');
      expect(formatDuration(-1)).toBe('N/A');
      expect(formatDuration(undefined)).toBe('N/A');
      expect(formatDuration(null as any)).toBe('N/A');
    });
  });

  describe('formatDurationShort', () => {
    it('should format seconds correctly', () => {
      expect(formatDurationShort(30)).toBe('30s');
      expect(formatDurationShort(59)).toBe('59s');
    });

    it('should format minutes correctly', () => {
      expect(formatDurationShort(60)).toBe('1m');
      expect(formatDurationShort(120)).toBe('2m');
      expect(formatDurationShort(3599)).toBe('59m');
    });

    it('should format hours correctly', () => {
      expect(formatDurationShort(3600)).toBe('1h');
      expect(formatDurationShort(7200)).toBe('2h');
      expect(formatDurationShort(86399)).toBe('23h');
    });

    it('should format days correctly', () => {
      expect(formatDurationShort(86400)).toBe('1d');
      expect(formatDurationShort(172800)).toBe('2d');
    });

    it('should handle edge cases', () => {
      expect(formatDurationShort(0)).toBe('N/A');
      expect(formatDurationShort(-1)).toBe('N/A');
      expect(formatDurationShort(undefined)).toBe('N/A');
      expect(formatDurationShort(null as any)).toBe('N/A');
    });
  });
}); 