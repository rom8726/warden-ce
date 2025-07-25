import { type ChipProps } from '@mui/material';
import { IssueStatus, IssueSource, type TimeseriesData } from '../../generated/api/client/api';
import { alpha } from '@mui/material/styles';

// Types
export interface ChartDataPoint {
  date: string;
  error: number;
  warning: number;
  info: number;
  exception: number;
  fatal: number;
  debug: number;
}

export interface TimeRangeOption {
  value: string;
  label: string;
}

// Constants
export const TIME_RANGE_OPTIONS: TimeRangeOption[] = [
  { value: '5m', label: 'Last 5 minutes' },
  { value: '10m', label: 'Last 10 minutes' },
  { value: '30m', label: 'Last 30 minutes' },
  { value: '60m', label: 'Last 60 minutes' },
  { value: '3h', label: 'Last 3 hours' },
  { value: '6h', label: 'Last 6 hours' },
  { value: '12h', label: 'Last 12 hours' },
  { value: '24h', label: 'Last 24 hours' },
  { value: '7d', label: 'Last 7 days' },
  { value: '14d', label: 'Last 14 days' },
  { value: '30d', label: 'Last 30 days' },
];

// Returns Material-UI chip color based on issue level
export const getLevelColor = (level: string): ChipProps['color'] => {
  switch (level) {
    case 'fatal':
      return 'error';
    case 'error':
      return 'error';
    case 'exception':
      return 'warning';
    case 'warning':
      return 'warning';
    case 'info':
      return 'info';
    case 'debug':
      return 'default';
    default:
      return 'default';
  }
};

// Returns hex color code based on issue level for custom styling
export const getLevelHexColor = (level: string): string => {
  switch (level) {
    case 'fatal':
      return '#d32f2f'; // Dark red
    case 'error':
      return '#f25454'; // Red (lighter than fatal)
    case 'exception':
      return '#ff9800'; // Orange
    case 'warning':
      return '#ffc107'; // Dark yellow
    case 'info':
      return '#03a9f4'; // Light blue
    case 'debug':
      return '#9e9e9e'; // Gray
    default:
      return '#9e9e9e'; // Gray
  }
};


// Returns styling for level badges
export const getLevelBadgeStyles = (level: string, size: 'small' | 'medium' = 'medium') => {
  const primaryColor = getLevelHexColor(level);

  return {
    // Base styles
    height: size === 'small' ? 24 : 32, // Match the height of a standard MUI Chip component
    fontWeight: 500,
    color: primaryColor,
    width: size === 'small' ? 70 : 90, // Fixed width based on "exception" (9 characters)

    // Clean background
    backgroundColor: alpha(primaryColor, 0.2),

    // Border styling similar to the outlined variant
    border: `1px solid ${primaryColor}`,

    // Rounded corners
    borderRadius: '8px', // Match the border radius of a standard MUI Chip component

    // Text styling
    '& .MuiChip-label': {
      px: size === 'small' ? 1 : 1.2,
      fontSize: size === 'small' ? '0.7rem' : '0.75rem',
      letterSpacing: '0.02em',
      textAlign: 'center',
      display: 'block', // Ensure the label takes full width
    },

    // Simple hover effect
    transition: 'all 0.2s ease-in-out',
    '&:hover': {
      backgroundColor: alpha(primaryColor, 0.30),
    },
  };
};

export const getStatusColor = (status: string): ChipProps['color'] => {
  switch (status) {
    case IssueStatus.Resolved:
      return 'success';
    case IssueStatus.Unresolved:
      return 'error';
    case IssueStatus.Ignored:
      return 'default';
    default:
      return 'default';
  }
};

// Date formatting utility
export const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleString();
};

// Get granularity based on time range
export const getGranularity = (timeRange: string): string => {
  if (['5m', '10m', '30m'].includes(timeRange)) {
    return '1m';
  } else if (timeRange === '60m') {
    return '10m';
  } else if (['3h', '6h'].includes(timeRange)) {
    return '10m';
  } else if (['12h', '24h'].includes(timeRange)) {
    return '1h';
  } else if (['7d', '14d', '30d'].includes(timeRange)) {
    return '1d';
  }
  return '1h'; // Default
};

// Transform timeseries data for the chart
export const transformTimeseriesData = (
  timeseriesData: TimeseriesData[] | null,
  issueSource: IssueSource | undefined
): ChartDataPoint[] => {
  if (!timeseriesData || !issueSource) {
    return [];
  }

  // Handle array format
  const seriesList = timeseriesData.filter(series => Array.isArray(series.occurrences) && series.occurrences.length > 0);

  if (seriesList.length === 0) {
    return [];
  }

  // Use the first series to determine period and generate dates
  const firstSeries = seriesList[0];
  const { granularity } = firstSeries.period || { interval: '7d', granularity: '1d' };
  const occurrencesLength = firstSeries.occurrences.length;

  // Generate dates for the chart
  const now = new Date();
  const dates: string[] = [];

  // Calculate the start date based on the interval and granularity
  const startDate = new Date(now);
  if (granularity === '1d') {
    startDate.setDate(startDate.getDate() - occurrencesLength);
  } else if (granularity === '1h') {
    startDate.setHours(startDate.getHours() - occurrencesLength);
  } else if (granularity === '10m') {
    startDate.setMinutes(startDate.getMinutes() - occurrencesLength * 10);
  } else if (granularity === '1m') {
    startDate.setMinutes(startDate.getMinutes() - occurrencesLength);
  }

  // Generate dates for each data point
  for (let i = 0; i < occurrencesLength; i++) {
    const date = new Date(startDate);
    if (granularity === '1d') {
      date.setDate(date.getDate() + i);
      dates.push(date.toLocaleDateString());
    } else if (granularity === '1h') {
      date.setHours(date.getHours() + i);
      dates.push(date.toLocaleString('default', { month: 'numeric', day: 'numeric', hour: 'numeric' }));
    } else if (granularity === '10m') {
      date.setMinutes(date.getMinutes() + i * 10);
      dates.push(date.toLocaleString('default', { hour: 'numeric', minute: 'numeric' }));
    } else if (granularity === '1m') {
      date.setMinutes(date.getMinutes() + i);
      dates.push(date.toLocaleString('default', { hour: 'numeric', minute: 'numeric' }));
    } else {
      dates.push(date.toLocaleDateString());
    }
  }

  // Create a result object with all data types initialized to 0
  const result: ChartDataPoint[] = dates.map(date => ({
    date,
    error: 0,
    warning: 0,
    info: 0,
    exception: 0,
    fatal: 0,
    debug: 0
  }));

  // Directly use the data from each series
  seriesList.forEach(series => {
    const name = series.name.toLowerCase();

    // Map the series name to the corresponding data type
    if (name.includes('error')) {
      series.occurrences.forEach((value, index) => {
        if (index < result.length) {
          result[index].error = value;
        }
      });
    } else if (name.includes('warning')) {
      series.occurrences.forEach((value, index) => {
        if (index < result.length) {
          result[index].warning = value;
        }
      });
    } else if (name.includes('info')) {
      series.occurrences.forEach((value, index) => {
        if (index < result.length) {
          result[index].info = value;
        }
      });
    } else if (name.includes('exception')) {
      series.occurrences.forEach((value, index) => {
        if (index < result.length) {
          result[index].exception = value;
        }
      });
    } else if (name.includes('fatal')) {
      series.occurrences.forEach((value, index) => {
        if (index < result.length) {
          result[index].fatal = value;
        }
      });
    } else if (name.includes('debug')) {
      series.occurrences.forEach((value, index) => {
        if (index < result.length) {
          result[index].debug = value;
        }
      });
    }
  });

  return result;
};
