import React from 'react';
import {
  Paper,
  Typography,
  Box,
  CircularProgress,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  useTheme,
} from '@mui/material';
import {
  Analytics as AnalyticsChartIcon,
  Timeline as TimelineIcon,
} from '@mui/icons-material';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { IssueLevel } from '../../../../generated/api/client';
import { getLevelHexColor } from '../../../../utils/issues/issueUtils';
import { useTheme as useAppTheme } from '../../../../theme/ThemeContext';

interface TimeseriesData {
  period: {
    interval: string;
    granularity: string;
  };
  name: string;
  occurrences: number[];
}

interface ErrorsChartProps {
  selectedRelease: string | null;
  selectedLevel: IssueLevel | undefined;
  selectedTimeWindow: string;
  errorsTimeseries: TimeseriesData[];
  compareTimeseries: TimeseriesData[];
  compareMode: boolean;
  compareRelease: string | null;
  comparisonData: any;
  loading: boolean;
  onTimeWindowChange: (timeWindow: string) => void;
  onLevelChange: (level: IssueLevel | undefined) => void;
  onCompareModeToggle: () => void;
}

const timeWindowOptions = [
  { value: '1d', label: '1 day' },
  { value: '3d', label: '3 days' },
  { value: '5d', label: '5 days' },
  { value: '7d', label: '7 days' },
  { value: '10d', label: '10 days' },
  { value: '14d', label: '14 days' },
  { value: '20d', label: '20 days' },
  { value: '25d', label: '25 days' },
  { value: '30d', label: '30 days' },
  { value: '45d', label: '45 days' },
  { value: '60d', label: '60 days' },
  { value: '75d', label: '75 days' },
  { value: '90d', label: '90 days' },
];

const LEVELS: IssueLevel[] = ['fatal', 'error', 'exception', 'warning', 'info', 'debug'];

const ErrorsChart: React.FC<ErrorsChartProps> = ({
  selectedRelease,
  selectedLevel,
  selectedTimeWindow,
  errorsTimeseries,
  compareTimeseries,
  compareMode,
  compareRelease,
  comparisonData,
  loading,
  onTimeWindowChange,
  onLevelChange,
  onCompareModeToggle,
}) => {
  const theme = useTheme();
  const { mode } = useAppTheme();

  const getErrorsChartBackground = () => {
    switch (mode) {
      case 'dark':
        return 'linear-gradient(135deg, rgba(45, 48, 56, 0.9) 0%, rgba(35, 38, 46, 0.9) 100%)';
      case 'blue':
        return 'linear-gradient(135deg, rgba(25, 55, 84, 0.9) 0%, rgba(16, 42, 66, 0.9) 100%)';
      case 'green':
        return 'linear-gradient(135deg, rgba(50, 95, 65, 0.9) 0%, rgba(40, 85, 50, 0.9) 70%, rgba(60, 105, 75, 0.9) 100%)';
      case 'light':
      default:
        return 'linear-gradient(135deg, rgba(155, 89, 182, 0.1) 0%, rgba(142, 68, 173, 0.1) 100%)';
    }
  };

  const getErrorsChartTextColor = () => {
    switch (mode) {
      case 'light':
        return theme.palette.text.primary;
      default:
        return 'white';
    }
  };

  const getErrorsChartBorderColor = () => {
    switch (mode) {
      case 'light':
        return 'rgba(155, 89, 182, 0.2)';
      default:
        return 'rgba(255, 255, 255, 0.2)';
    }
  };

  const renderDelta = (value: number | undefined, label: string) => {
    if (value === undefined || value === 0) {
      return (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
          <Typography variant="body2" color="text.secondary">
            {label}: 0
          </Typography>
        </Box>
      );
    }
    
    // Определяем, является ли изменение "хорошим" или "плохим"
    const isGoodChange = (val: number, metricLabel: string): boolean => {
      switch (metricLabel) {
        case 'Known Issues':
        case 'New Issues':
        case 'Regressions':
        case 'Users':
          // Для этих метрик уменьшение - это хорошо
          return val < 0;
        case 'Resolved':
          // Для resolved issues увеличение - это хорошо
          return val > 0;
        default:
          // По умолчанию считаем, что уменьшение - это хорошо
          return val < 0;
      }
    };
    
    const isGood = isGoodChange(value, label);
    const color = isGood ? 'success.main' : 'error.main';
    
    return (
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
        <Typography variant="body2" color={color}>
          {label}: {value > 0 ? '+' : ''}{value}
        </Typography>
      </Box>
    );
  };

  const ChartTooltip = ({ active, payload, label }: any) => {
    if (!active || !payload || !payload.length) return null;
    return (
      <Box sx={{ p: 1, minWidth: 160, bgcolor: 'background.paper', border: 1, borderColor: 'divider', borderRadius: 1 }}>
        <Typography variant="subtitle2" sx={{ mb: 1 }}>{label}</Typography>
        {payload.map((item: any) => (
          <Box key={item.dataKey} sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 0.5 }}>
            <span style={{ color: item.color, fontWeight: 500 }}>{item.name}</span>
            <span style={{ color: item.color, fontWeight: 600 }}>{item.value}</span>
          </Box>
        ))}
      </Box>
    );
  };

  const CustomLegend = ({ payload }: any) => {
    if (!payload || payload.length === 0) return null;
    
    return (
      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, justifyContent: 'center', mt: 1 }}>
        {payload.map((entry: any) => {
          const isComparison = entry.dataKey.includes('_compare');
          
          return (
            <Box
              key={entry.dataKey}
              sx={{
                display: 'flex',
                alignItems: 'center',
                gap: 0.5,
                cursor: 'pointer',
                p: 0.5,
                borderRadius: 0.5,
                '&:hover': {
                  bgcolor: 'action.hover'
                }
              }}
              onClick={() => onLevelChange(entry.dataKey.replace('_compare', '') as IssueLevel)}
            >
              <svg width="20" height="4" style={{ overflow: 'visible' }}>
                {isComparison ? (
                  // Dashed line for comparison
                  <line
                    x1="0"
                    y1="2"
                    x2="20"
                    y2="2"
                    stroke={entry.color}
                    strokeWidth="2"
                    strokeDasharray="3,2"
                  />
                ) : (
                  // Solid line for main release
                  <line
                    x1="0"
                    y1="2"
                    x2="20"
                    y2="2"
                    stroke={entry.color}
                    strokeWidth="2"
                  />
                )}
              </svg>
              <Typography variant="caption" sx={{ color: entry.color, fontWeight: 500 }}>
                {entry.value}
              </Typography>
            </Box>
          );
        })}
      </Box>
    );
  };

  const prepareChartData = (timeseries: TimeseriesData[], isComparison: boolean = false) => {
    if (!timeseries || timeseries.length === 0) return [];
    
    console.log('Raw timeseries data:', timeseries, 'isComparison:', isComparison);
    
    // Group by level
    const levelGroups: Record<string, TimeseriesData> = {};
    
    timeseries.forEach(series => {
      const levelKey = isComparison ? `${series.name}_compare` : series.name;
      levelGroups[levelKey] = series;
    });
    
    // Create chart data
    const chartData = [];
    const occurrencesLength = timeseries.length > 0 ? timeseries[0].occurrences.length : 0;
    
    // Generate time labels based on the selected time window
    const generateTimeLabel = (index: number) => {
      const now = new Date();
      let timeOffset = 0;
      
      // Determine time offset based on a selected window and granularity
      if (selectedTimeWindow === '1d' || selectedTimeWindow === '3d' || selectedTimeWindow === '5d' || selectedTimeWindow === '7d' || selectedTimeWindow === '10d') {
        // 1 hour granularity
        timeOffset = (occurrencesLength - 1 - index) * 60 * 60 * 1000;
      } else if (selectedTimeWindow === '14d' || selectedTimeWindow === '20d') {
        // 4 hours granularity
        timeOffset = (occurrencesLength - 1 - index) * 4 * 60 * 60 * 1000;
      } else if (selectedTimeWindow === '25d' || selectedTimeWindow === '30d') {
        // 6 hours granularity
        timeOffset = (occurrencesLength - 1 - index) * 6 * 60 * 60 * 1000;
      } else if (selectedTimeWindow === '45d' || selectedTimeWindow === '60d') {
        // 12 hours granularity
        timeOffset = (occurrencesLength - 1 - index) * 12 * 60 * 60 * 1000;
      } else {
        // 1 day granularity
        timeOffset = (occurrencesLength - 1 - index) * 24 * 60 * 60 * 1000;
      }
      
      const date = new Date(now.getTime() - timeOffset);
      
      // Format based on granularity
      if (selectedTimeWindow === '1d' || selectedTimeWindow === '3d' || selectedTimeWindow === '5d' || selectedTimeWindow === '7d' || selectedTimeWindow === '10d') {
        return date.toLocaleString('ru-RU', { 
          month: 'short', 
          day: 'numeric',
          hour: '2-digit',
          minute: '2-digit'
        });
      } else if (selectedTimeWindow === '14d' || selectedTimeWindow === '20d') {
        return date.toLocaleString('ru-RU', { 
          month: 'short', 
          day: 'numeric',
          hour: '2-digit'
        });
      } else if (selectedTimeWindow === '25d' || selectedTimeWindow === '30d') {
        return date.toLocaleString('ru-RU', { 
          month: 'short', 
          day: 'numeric',
          hour: '2-digit'
        });
      } else if (selectedTimeWindow === '45d' || selectedTimeWindow === '60d') {
        return date.toLocaleString('ru-RU', { 
          month: 'short', 
          day: 'numeric',
          hour: '2-digit'
        });
      } else {
        return date.toLocaleDateString('ru-RU', { 
          month: 'short', 
          day: 'numeric'
        });
      }
    };
    
    for (let i = 0; i < occurrencesLength; i++) {
      const dataPoint: Record<string, any> = { 
        index: i,
        timeLabel: generateTimeLabel(i)
      };
      
      // Add counts for each level
      Object.keys(levelGroups).forEach(level => {
        const baseLevel = level.replace('_compare', '');
        // If a specific level is selected, only add that level's data
        if (!selectedLevel || selectedLevel === baseLevel) {
          dataPoint[level] = levelGroups[level].occurrences[i];
        }
      });
      
      chartData.push(dataPoint);
    }
    
    console.log('Processed chart data:', chartData);
    return chartData;
  };

  const chartData = (() => {
    const mainData = prepareChartData(errorsTimeseries || [], false);
    const compareData = compareMode ? prepareChartData(compareTimeseries || [], true) : [];
    
    if (!compareMode || compareData.length === 0) {
      return mainData;
    }
    
    // Merge data for comparison
    const mergedData = [];
    const maxLength = Math.max(mainData.length, compareData.length);
    
    for (let i = 0; i < maxLength; i++) {
      const dataPoint: Record<string, any> = { 
        index: i,
        timeLabel: mainData[i]?.timeLabel || compareData[i]?.timeLabel || `Day ${i}`
      };
      
      // Add main release data
      if (mainData[i]) {
        Object.keys(mainData[i]).forEach(key => {
          if (key !== 'index' && key !== 'timeLabel') {
            dataPoint[key] = mainData[i][key];
          }
        });
      }
      
      // Add comparison release data
      if (compareData[i]) {
        Object.keys(compareData[i]).forEach(key => {
          if (key !== 'index' && key !== 'timeLabel') {
            dataPoint[key] = compareData[i][key];
          }
        });
      }
      
      mergedData.push(dataPoint);
    }
    
    return mergedData;
  })();

  if (!selectedRelease) {
    return null;
  }

  return (
    <Paper sx={{ 
      p: 3,
      background: getErrorsChartBackground(),
      color: getErrorsChartTextColor(),
      position: 'relative',
      overflow: 'hidden'
    }}>
      <Box sx={{
        position: 'absolute',
        top: -20,
        right: -20,
        width: 100,
        height: 100,
        borderRadius: '50%',
        background: mode === 'light' ? 'rgba(155, 89, 182, 0.1)' : 'rgba(255,255,255,0.1)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center'
      }}>
        <TimelineIcon sx={{ 
          fontSize: 40, 
          opacity: mode === 'light' ? 0.4 : 0.3,
          color: mode === 'light' ? '#9b59b6' : 'white'
        }} />
      </Box>
      
      <Box sx={{ position: 'relative', zIndex: 1 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 2, flexWrap: 'wrap', gap: 2 }}>
          <Typography variant="h6" sx={{ 
            display: 'flex', 
            alignItems: 'center',
            color: getErrorsChartTextColor(),
            fontWeight: 600
          }}>
            <AnalyticsChartIcon sx={{ mr: 1 }} />
            Errors Over Time (per level)
            {selectedLevel && (
              <Chip 
                label={`Filtered: ${selectedLevel}`}
                size="small"
                color="primary"
                onDelete={() => onLevelChange(undefined)}
                sx={{ ml: 1 }}
              />
            )}
            {compareMode && compareRelease && comparisonData && (
              <Chip 
                label={`${compareRelease} → ${selectedRelease}`}
                size="small"
                color="secondary"
                onDelete={onCompareModeToggle}
                sx={{ ml: 1 }}
              />
            )}
          </Typography>
          
          {/* Time Window Selector */}
          <FormControl sx={{ minWidth: 150 }}>
            <InputLabel sx={{ color: getErrorsChartTextColor() }}>Time Window</InputLabel>
            <Select
              value={selectedTimeWindow}
              label="Time Window"
              onChange={(e) => onTimeWindowChange(e.target.value)}
              size="small"
              sx={{
                color: getErrorsChartTextColor(),
                '& .MuiOutlinedInput-notchedOutline': {
                  borderColor: getErrorsChartBorderColor()
                },
                '&:hover .MuiOutlinedInput-notchedOutline': {
                  borderColor: mode === 'light' ? 'rgba(155, 89, 182, 0.4)' : 'rgba(255, 255, 255, 0.4)'
                }
              }}
            >
              {timeWindowOptions.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Box>
        
        {/* Comparison Summary for Chart */}
        {compareMode && comparisonData && (
          <Box sx={{ 
            mb: 2, 
            p: 2, 
            borderRadius: 2,
            bgcolor: mode === 'light' ? 'rgba(155, 89, 182, 0.05)' : 'rgba(255, 255, 255, 0.1)',
            border: `1px solid ${getErrorsChartBorderColor()}`
          }}>
            <Typography variant="subtitle2" gutterBottom sx={{ 
              color: getErrorsChartTextColor(),
              fontWeight: 600
            }}>
              📊 Comparison Summary:
            </Typography>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2 }}>
              {renderDelta(comparisonData.delta.known_issues_total, 'Known Issues')}
              {renderDelta(comparisonData.delta.new_issues_total, 'New Issues')}
              {renderDelta(comparisonData.delta.regressions_total, 'Regressions')}
              {renderDelta(comparisonData.delta.resolved_in_version_total, 'Resolved')}
              {renderDelta(comparisonData.delta.users_affected, 'Users')}
            </Box>
          </Box>
        )}
        
        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress sx={{ color: getErrorsChartTextColor() }} />
          </Box>
        ) : chartData.length > 0 ? (
          <Box sx={{ height: 400 }}>
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" stroke={mode === 'light' ? 'rgba(0,0,0,0.1)' : 'rgba(255,255,255,0.1)'} />
                <XAxis 
                  dataKey="timeLabel" 
                  tick={{ fill: getErrorsChartTextColor() }}
                  axisLine={{ stroke: getErrorsChartBorderColor() }}
                />
                <YAxis 
                  tick={{ fill: getErrorsChartTextColor() }}
                  axisLine={{ stroke: getErrorsChartBorderColor() }}
                />
                <Tooltip content={ChartTooltip} />
                <Legend content={CustomLegend} />
                {LEVELS.map((level) => {
                  // Check if this level has any data
                  const hasData = chartData.some(item => (item as any)[level] !== undefined && (item as any)[level] > 0);
                  const hasCompareData = compareMode && chartData.some(item => (item as any)[`${level}_compare`] !== undefined && (item as any)[`${level}_compare`] > 0);
                  
                  if (!hasData && !hasCompareData) return null;
                  
                  // If a specific level is selected, show only that level
                  if (selectedLevel && selectedLevel !== level) return null;
                  
                  const lines = [];
                  
                  // Main release line
                  if (hasData) {
                    lines.push(
                      <Line
                        key={level}
                        type="monotone"
                        dataKey={level}
                        stroke={getLevelHexColor(level)}
                        strokeWidth={2}
                        dot={{ fill: getLevelHexColor(level) }}
                        name={`${level.charAt(0).toUpperCase() + level.slice(1)} (${selectedRelease})`}
                      />
                    );
                  }
                  
                  // Comparison release line
                  if (hasCompareData) {
                    lines.push(
                      <Line
                        key={`${level}_compare`}
                        type="monotone"
                        dataKey={`${level}_compare`}
                        stroke={getLevelHexColor(level)}
                        strokeWidth={2}
                        strokeDasharray="5 5"
                        dot={{ fill: getLevelHexColor(level) }}
                        name={`${level.charAt(0).toUpperCase() + level.slice(1)} (${compareRelease})`}
                        legendType="line"
                      />
                    );
                  }
                  
                  return lines;
                })}
              </LineChart>
            </ResponsiveContainer>
          </Box>
        ) : (
          <Typography sx={{ color: getErrorsChartTextColor() }}>
            No data.
          </Typography>
        )}
      </Box>
    </Paper>
  );
};

export default ErrorsChart; 