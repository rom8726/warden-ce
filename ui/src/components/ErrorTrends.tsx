import React, { useState, useEffect } from 'react';
import { 
  LineChart, 
  Line, 
  ResponsiveContainer, 
  Tooltip,
  XAxis,
  YAxis
} from 'recharts';
import { Box, useTheme, CircularProgress } from '@mui/material';
import { keyframes } from '@emotion/react';
import apiClient from '../api/apiClient';

interface ErrorTrendsProps {
  width?: number;
  height?: number;
  projectId: number;
}

// Define the grow-from-bottom animation
const growFromBottom = keyframes`
  0% {
    transform: scaleY(0);
    transform-origin: bottom;
  }
  100% {
    transform: scaleY(1);
    transform-origin: bottom;
  }
`;

const ErrorTrends: React.FC<ErrorTrendsProps> = ({ width = 100, height = 40, projectId }) => {
  const theme = useTheme();
  const [data, setData] = useState<{ time: string; errors: number }[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchErrorTrends = async () => {
      try {
        // Fetch data for the last 3 hours with 10-minute granularity
        const response = await apiClient.getEventsTimeseries('3h', '10m', projectId);

        // Filter to only include error level events
        const errorSeries = response.data.find(series => series.name !== 'info' && series.name != 'debug');

        if (errorSeries) {
          // Transform the data for the chart
          const chartData = errorSeries.occurrences.map((count, index) => {
            // Calculate the time for each data point
            const now = new Date();
            const minutesAgo = (errorSeries.occurrences.length - 1 - index) * 10;
            const time = new Date(now.getTime() - minutesAgo * 60000);

            return {
              time: time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
              errors: count
            };
          });

          setData(chartData);
        } else {
          // If no error data, set empty array
          setData([]);
        }
      } catch (err) {
        console.error('Error fetching error trends:', err);
        setData([]);
      } finally {
        setLoading(false);
      }
    };

    fetchErrorTrends();
  }, [projectId]);

  // Add a gradient for the line
  const gradientId = `errorLineGradient-${Math.random().toString(36).substring(2, 9)}`;

  if (loading) {
    return (
      <Box sx={{ 
        width, 
        height, 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center' 
      }}>
        <CircularProgress size={20} color="error" />
      </Box>
    );
  }

  // If no data or empty data, show a flat line
  if (!data || data.length === 0) {
    return (
      <Box sx={{ 
        width, 
        height,
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        color: theme.palette.text.secondary,
        fontSize: '0.7rem'
      }}>
        No data
      </Box>
    );
  }

  return (
    <Box sx={{ 
      width, 
      height,
      position: 'relative',
      borderRadius: 1,
      overflow: 'hidden',
      boxShadow: `0 2px 8px 0 ${theme.palette.mode === 'dark' ? 'rgba(0,0,0,0.2)' : 'rgba(0,0,0,0.05)'}`,
      animation: `${growFromBottom} 1200ms ease-in-out`,
      '&:hover': {
        '& .error-trend-line': {
          strokeWidth: 2.5,
        },
        '& .error-trend-dot': {
          r: 4,
        }
      }
    }}>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data} margin={{ top: 5, right: 5, bottom: 5, left: 5 }}>
          <defs>
            <linearGradient id={gradientId} x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#ff5252" stopOpacity={1}/>
              <stop offset="95%" stopColor={theme.palette.error.main} stopOpacity={0.6}/>
            </linearGradient>
          </defs>
          <Line 
            type="monotone" 
            dataKey="errors" 
            stroke={`url(#${gradientId})`}
            strokeWidth={2}
            dot={false}
            activeDot={{ r: 3.5, className: 'error-trend-dot', stroke: '#fff', strokeWidth: 1 }}
            className="error-trend-line"
            isAnimationActive={false}
          />
          <Tooltip 
            animationDuration={200}
            contentStyle={{ 
              backgroundColor: theme.palette.background.paper,
              border: 'none',
              borderRadius: 4,
              boxShadow: theme.shadows[4],
              fontSize: '0.75rem',
              padding: '6px 10px',
            }}
            labelStyle={{ 
              color: theme.palette.text.primary, 
              fontWeight: 'bold',
              marginBottom: '4px',
            }}
            itemStyle={{
              padding: '2px 0',
            }}
            formatter={(value) => [`${value} errors`, 'Errors']}
            labelFormatter={(label) => `Time: ${label}`}
            cursor={{ stroke: theme.palette.divider, strokeWidth: 1, strokeDasharray: '3 3' }}
          />
          <XAxis 
            dataKey="time" 
            hide={true}
          />
          <YAxis 
            hide={true}
          />
        </LineChart>
      </ResponsiveContainer>
    </Box>
  );
};

export default ErrorTrends;
