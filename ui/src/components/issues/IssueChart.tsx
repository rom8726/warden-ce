import { memo, useMemo } from 'react';
import { 
  Box, 
  Typography, 
  Paper,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  useTheme
} from '@mui/material';
import type { SelectChangeEvent } from '@mui/material/Select';
import { 
  ResponsiveContainer, 
  LineChart, 
  Line, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  Legend 
} from 'recharts';
import { IssueSource } from '../../generated/api/client';
import { useIssuePageStyles } from './IssuePageStyles';
import {type ChartDataPoint} from '../../utils/issues/issueUtils';

interface IssueChartProps {
  timeRange: string;
  onTimeRangeChange: (event: SelectChangeEvent) => void;
  chartData: ChartDataPoint[];
  issueSource: IssueSource;
}

const IssueChart = memo(({ timeRange, onTimeRangeChange, chartData}: IssueChartProps) => {
  const styles = useIssuePageStyles();
  const theme = useTheme();

  // Check if there's data for each level
  const hasErrorData = useMemo(() => 
    chartData.some(point => point.error && point.error > 0), 
    [chartData]
  );
  
  const hasWarningData = useMemo(() => 
    chartData.some(point => point.warning && point.warning > 0), 
    [chartData]
  );
  
  const hasInfoData = useMemo(() => 
    chartData.some(point => point.info && point.info > 0), 
    [chartData]
  );
  
  const hasExceptionData = useMemo(() => 
    chartData.some(point => point.exception && point.exception > 0), 
    [chartData]
  );

  const hasFatalData = useMemo(() => 
    chartData.some(point => point.fatal && point.fatal > 0), 
    [chartData]
  );
  const hasDebugData = useMemo(() => 
    chartData.some(point => point.debug && point.debug > 0), 
    [chartData]
  );

  // If no data, show a message
  if (chartData.length === 0) {
    return (
      <Paper sx={styles.chartContainer}>
        <Box sx={styles.chartHeader}>
          <Typography variant="h6" className="gradient-subtitle-blue">
            Occurrence Trend
          </Typography>
          <FormControl size="small" sx={{ minWidth: 120 }}>
            <InputLabel htmlFor="time-range-select">Time Range</InputLabel>
            <Select
              id="time-range-select"
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
        </Box>
        <Box sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="body1" color="text.secondary">
            No data available for the selected time range.
          </Typography>
        </Box>
      </Paper>
    );
  }

  return (
    <Paper sx={styles.chartContainer}>
      <Box sx={styles.chartHeader}>
        <Typography variant="h6">
          Occurrence Trend
        </Typography>
        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel htmlFor="time-range-select">Time Range</InputLabel>
          <Select
            id="time-range-select"
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
      </Box>
      <Box sx={styles.chartBox}>
        <ResponsiveContainer width="100%" height="100%">
          <LineChart
            data={chartData}
            margin={{
              top: 20,
              right: 30,
              left: 20,
              bottom: 10,
            }}
          >
            <defs>
              <linearGradient id="colorError" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#ff5252" stopOpacity={1}/>
                <stop offset="95%" stopColor="#f44336" stopOpacity={0.8}/>
              </linearGradient>
              <linearGradient id="colorWarning" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#ffb74d" stopOpacity={1}/>
                <stop offset="95%" stopColor="#ff9800" stopOpacity={0.8}/>
              </linearGradient>
              <linearGradient id="colorInfo" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#64b5f6" stopOpacity={1}/>
                <stop offset="95%" stopColor="#2196f3" stopOpacity={0.8}/>
              </linearGradient>
              <linearGradient id="colorException" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#ff5252" stopOpacity={1}/>
                <stop offset="95%" stopColor="#b71c1c" stopOpacity={0.8}/>
              </linearGradient>
              <linearGradient id="colorFatal" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#d32f2f" stopOpacity={1}/>
                <stop offset="95%" stopColor="#d32f2f" stopOpacity={0.8}/>
              </linearGradient>
              <linearGradient id="colorDebug" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#9e9e9e" stopOpacity={1}/>
                <stop offset="95%" stopColor="#9e9e9e" stopOpacity={0.8}/>
              </linearGradient>
            </defs>
            <CartesianGrid 
              strokeDasharray="3 3" 
              stroke={theme.palette.mode === 'dark' ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.1)'} 
            />
            <XAxis 
              dataKey="date" 
              tick={{ fill: theme.palette.text.secondary, fontSize: 12 }}
              tickLine={{ stroke: theme.palette.divider }}
              axisLine={{ stroke: theme.palette.divider }}
            />
            <YAxis 
              tick={{ fill: theme.palette.text.secondary, fontSize: 12 }}
              tickLine={{ stroke: theme.palette.divider }}
              axisLine={{ stroke: theme.palette.divider }}
            />
            <Tooltip 
              animationDuration={300}
              contentStyle={{ 
                backgroundColor: theme.palette.background.paper,
                border: 'none',
                borderRadius: 8,
                boxShadow: theme.shadows[6],
                padding: '12px 16px',
                fontSize: '13px',
              }}
              itemStyle={{
                padding: '4px 0',
              }}
              labelStyle={{
                fontWeight: 'bold',
                marginBottom: '8px',
              }}
              cursor={{ stroke: theme.palette.divider, strokeWidth: 1, strokeDasharray: '5 5' }}
            />
            <Legend 
              verticalAlign="top"
              height={36}
              iconSize={12}
              iconType="circle"
              wrapperStyle={{
                paddingTop: '10px',
                fontSize: '13px',
              }}
            />
            {/* Display all available data types regardless of source */}
            {hasErrorData && (
              <Line
                type="monotone"
                dataKey="error"
                stroke="url(#colorError)"
                strokeWidth={2}
                dot={{ r: 0 }}
                activeDot={{ r: 4, className: 'error-dot', stroke: '#fff', strokeWidth: 1 }}
                className="error-line"
                name="Error"
                isAnimationActive={false}
              />
            )}
            {hasWarningData && (
              <Line
                type="monotone"
                dataKey="warning"
                stroke="url(#colorWarning)"
                strokeWidth={2}
                dot={{ r: 0 }}
                activeDot={{ r: 4, className: 'warning-dot', stroke: '#fff', strokeWidth: 1 }}
                className="warning-line"
                name="Warning"
                isAnimationActive={false}
              />
            )}
            {hasInfoData && (
              <Line
                type="monotone"
                dataKey="info"
                stroke="url(#colorInfo)"
                strokeWidth={2}
                dot={{ r: 0 }}
                activeDot={{ r: 4, className: 'info-dot', stroke: '#fff', strokeWidth: 1 }}
                className="info-line"
                name="Info"
                isAnimationActive={false}
              />
            )}
            {hasExceptionData && (
              <Line
                type="monotone"
                dataKey="exception"
                stroke="url(#colorException)"
                strokeWidth={2}
                dot={{ r: 0 }}
                activeDot={{ r: 4, className: 'exception-dot', stroke: '#fff', strokeWidth: 1 }}
                className="exception-line"
                name="Exception"
                isAnimationActive={false}
              />
            )}
            {hasFatalData && (
              <Line
                type="monotone"
                dataKey="fatal"
                stroke="url(#colorFatal)"
                strokeWidth={2}
                dot={{ r: 0 }}
                activeDot={{ r: 4, className: 'fatal-dot', stroke: '#fff', strokeWidth: 1 }}
                className="fatal-line"
                name="Fatal"
                isAnimationActive={false}
              />
            )}
            {hasDebugData && (
              <Line
                type="monotone"
                dataKey="debug"
                stroke="url(#colorDebug)"
                strokeWidth={2}
                dot={{ r: 0 }}
                activeDot={{ r: 4, className: 'debug-dot', stroke: '#fff', strokeWidth: 1 }}
                className="debug-line"
                name="Debug"
                isAnimationActive={false}
              />
            )}
          </LineChart>
        </ResponsiveContainer>
      </Box>
    </Paper>
  );
});

IssueChart.displayName = 'IssueChart';

export default IssueChart;