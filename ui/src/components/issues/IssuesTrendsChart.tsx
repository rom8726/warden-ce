import React from 'react';
import { Box, Typography, IconButton, FormControl, InputLabel, Select, MenuItem, useTheme, Collapse } from '@mui/material';
import { ResponsiveContainer, LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip as TooltipChart, Legend } from 'recharts';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { getLevelHexColor } from '../../utils/issues/issueUtils';

interface Level {
  key: string;
  label: string;
}

const LEVELS: Level[] = [
  { key: 'fatal', label: 'Fatal' },
  { key: 'error', label: 'Errors' },
  { key: 'exception', label: 'Exceptions' },
  { key: 'warning', label: 'Warnings' },
  { key: 'info', label: 'Info' },
  { key: 'debug', label: 'Debug' },
];

interface IssuesTrendsChartProps {
  data: any[];
  timeRange: string;
  onTimeRangeChange?: (event: any) => void;
  isChartExpanded?: boolean;
  onToggleChartExpanded?: () => void;
  chartTitle?: string;
  showTimeRangeSelector?: boolean;
  height?: number;
}

export const ChartTooltip = ({ active, payload, label }: any) => {
  if (!active || !payload || !payload.length) return null;
  return (
    <Box sx={{ p: 1, minWidth: 160 }}>
      <Typography variant="subtitle2" sx={{ mb: 1 }}>{label}</Typography>
      {LEVELS.map(({ key, label: levelLabel }) => {
        const item = payload.find((p: any) => p.dataKey === key);
        if (!item) return null;
        return (
          <Box key={key} sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 0.5 }}>
            <span style={{ color: getLevelHexColor(key), fontWeight: 500 }}>{levelLabel}</span>
            <span style={{ color: getLevelHexColor(key), fontWeight: 600 }}>{item.value}</span>
          </Box>
        );
      })}
    </Box>
  );
};

const IssuesTrendsChart: React.FC<IssuesTrendsChartProps> = ({
  data,
  timeRange,
  onTimeRangeChange,
  isChartExpanded = true,
  onToggleChartExpanded,
  chartTitle = 'Issues Over Time',
  showTimeRangeSelector = true,
  height = 300,
}) => {
  const theme = useTheme();

  return (
    <Box sx={{ p: 2, borderBottom: `1px solid ${theme.palette.divider}` }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2, flexWrap: 'wrap', gap: 1 }}>
        <Typography variant="h6" sx={{ fontWeight: 600 }} className="gradient-subtitle-blue">{chartTitle}</Typography>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, flexWrap: 'wrap' }}>
          {showTimeRangeSelector && (
            <FormControl size="small" sx={{ minWidth: 120 }}>
              <InputLabel>Time Range</InputLabel>
              <Select
                value={timeRange}
                label="Time Range"
                onChange={onTimeRangeChange}
                sx={{ fontSize: '0.875rem' }}
              >
                <MenuItem value="10m">Last 10 minutes</MenuItem>
                <MenuItem value="30m">Last 30 minutes</MenuItem>
                <MenuItem value="1h">Last hour</MenuItem>
                <MenuItem value="3h">Last 3 hours</MenuItem>
                <MenuItem value="6h">Last 6 hours</MenuItem>
                <MenuItem value="12h">Last 12 hours</MenuItem>
                <MenuItem value="24h">Last 24 hours</MenuItem>
                <MenuItem value="3d">Last 3 days</MenuItem>
                <MenuItem value="7d">Last 7 days</MenuItem>
                <MenuItem value="14d">Last 14 days</MenuItem>
                <MenuItem value="30d">Last 30 days</MenuItem>
              </Select>
            </FormControl>
          )}
          {onToggleChartExpanded && (
            <IconButton 
              size="small" 
              onClick={onToggleChartExpanded}
              sx={{ 
                bgcolor: theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.05)' : 'rgba(0, 0, 0, 0.04)',
                '&:hover': {
                  bgcolor: theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.08)',
                }
              }}
            >
              {isChartExpanded ? <ExpandLessIcon fontSize="small" /> : <ExpandMoreIcon fontSize="small" />}
            </IconButton>
          )}
        </Box>
      </Box>
      <Collapse in={isChartExpanded}>
        <Box sx={{ height, width: '100%' }}>
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={data}>
              <defs>
                <linearGradient id="colorFatal" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor={getLevelHexColor('fatal')} stopOpacity={0.8}/>
                  <stop offset="95%" stopColor={getLevelHexColor('fatal')} stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorError" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor={getLevelHexColor('error')} stopOpacity={0.8}/>
                  <stop offset="95%" stopColor={getLevelHexColor('error')} stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorException" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor={getLevelHexColor('exception')} stopOpacity={0.8}/>
                  <stop offset="95%" stopColor={getLevelHexColor('exception')} stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorWarning" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor={getLevelHexColor('warning')} stopOpacity={0.8}/>
                  <stop offset="95%" stopColor={getLevelHexColor('warning')} stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorInfo" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor={getLevelHexColor('info')} stopOpacity={0.8}/>
                  <stop offset="95%" stopColor={getLevelHexColor('info')} stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorDebug" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor={getLevelHexColor('debug')} stopOpacity={0.8}/>
                  <stop offset="95%" stopColor={getLevelHexColor('debug')} stopOpacity={0}/>
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke={theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)'} />
              <XAxis 
                dataKey="timeLabel" 
                stroke={theme.palette.text.secondary}
                fontSize={12}
                tick={{ fill: theme.palette.text.secondary }}
              />
              <YAxis 
                stroke={theme.palette.text.secondary}
                fontSize={12}
                tick={{ fill: theme.palette.text.secondary }}
              />
              <TooltipChart 
                content={ChartTooltip}
                contentStyle={{
                  backgroundColor: theme.palette.mode === 'dark' ? 'rgba(45, 48, 56, 0.95)' : 'rgba(255, 255, 255, 0.95)',
                  border: `1px solid ${theme.palette.divider}`,
                  borderRadius: 8,
                }}
                labelStyle={{ color: theme.palette.text.primary }}
              />
              <Legend 
                verticalAlign="top"
                height={30}
                iconSize={10}
                iconType="circle"
              />
              <Line
                type="monotone"
                dataKey="fatal"
                stroke={getLevelHexColor('fatal')}
                strokeWidth={2}
                dot={{ r: 3 }}
                activeDot={{ r: 5 }}
                name="Fatal"
              />
              <Line
                type="monotone"
                dataKey="error"
                stroke={getLevelHexColor('error')}
                strokeWidth={2}
                dot={{ r: 3 }}
                activeDot={{ r: 5 }}
                name="Errors"
              />
              <Line
                type="monotone"
                dataKey="exception"
                stroke={getLevelHexColor('exception')}
                strokeWidth={2}
                dot={{ r: 3 }}
                activeDot={{ r: 5 }}
                name="Exceptions"
              />
              <Line
                type="monotone"
                dataKey="warning"
                stroke={getLevelHexColor('warning')}
                strokeWidth={2}
                dot={{ r: 3 }}
                activeDot={{ r: 5 }}
                name="Warnings"
              />
              <Line
                type="monotone"
                dataKey="info"
                stroke={getLevelHexColor('info')}
                strokeWidth={2}
                dot={{ r: 3 }}
                activeDot={{ r: 5 }}
                name="Info"
              />
              <Line
                type="monotone"
                dataKey="debug"
                stroke={getLevelHexColor('debug')}
                strokeWidth={2}
                dot={{ r: 3 }}
                activeDot={{ r: 5 }}
                name="Debug"
              />
            </LineChart>
          </ResponsiveContainer>
        </Box>
      </Collapse>
    </Box>
  );
};

export default IssuesTrendsChart;
