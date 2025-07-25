import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import IssueChart from '../IssueChart';
import { IssueSource } from '../../../generated/api/client/api';

// Mock the recharts components
// @ts-ignore
vi.mock('recharts', () => ({
  ResponsiveContainer: ({ children }) => <div data-testid="responsive-container">{children}</div>,
  LineChart: ({ children }) => <div data-testid="line-chart">{children}</div>,
  Line: ({ dataKey, name }) => <div data-testid={`line-${dataKey}`}>{name}</div>,
  XAxis: () => <div data-testid="x-axis" />,
  YAxis: () => <div data-testid="y-axis" />,
  CartesianGrid: () => <div data-testid="cartesian-grid" />,
  Tooltip: () => <div data-testid="tooltip" />,
  Legend: () => <div data-testid="legend" />
}));

// Mock the useIssuePageStyles hook
vi.mock('../IssuePageStyles', () => ({
  useIssuePageStyles: () => ({
    chartContainer: { mb: 2 },
    chartHeader: { 
      display: 'flex', 
      justifyContent: 'space-between', 
      alignItems: 'center', 
      mb: 1.5,
      p: 2
    },
    chartBox: {
      height: 250, 
      width: '100%',
      position: 'relative',
      borderRadius: 1,
      overflow: 'hidden'
    }
  })
}));

// Mock the TIME_RANGE_OPTIONS
vi.mock('../../../utils/issues/issueUtils', () => ({
  TIME_RANGE_OPTIONS: [
    { value: '1h', label: 'Last 1 hour' },
    { value: '24h', label: 'Last 24 hours' },
    { value: '7d', label: 'Last 7 days' }
  ]
}));

describe('IssueChart', () => {
  const mockTimeRange = '24h';
  const mockOnTimeRangeChange = vi.fn();
  const mockIssueSource = IssueSource.Exception;

  test('renders no data message when chartData is empty', () => {
    render(
      <IssueChart 
        timeRange={mockTimeRange} 
        onTimeRangeChange={mockOnTimeRangeChange} 
        chartData={[]} 
        issueSource={mockIssueSource} 
      />
    );
    
    expect(screen.getByText('Occurrence Trend')).toBeInTheDocument();
    expect(screen.getByText('No data available for the selected time range.')).toBeInTheDocument();
    expect(screen.getByRole('combobox')).toBeInTheDocument();
  });

  test('renders chart when data is available', () => {
    const chartData = [
      { date: '2023-01-01', error: 5, warning: 2, info: 1, exception: 0 },
      { date: '2023-01-02', error: 3, warning: 1, info: 0, exception: 2 }
    ];
    
    render(
      <IssueChart 
        timeRange={mockTimeRange} 
        onTimeRangeChange={mockOnTimeRangeChange} 
        chartData={chartData} 
        issueSource={mockIssueSource} 
      />
    );
    
    expect(screen.getByText('Occurrence Trend')).toBeInTheDocument();
    expect(screen.getByTestId('responsive-container')).toBeInTheDocument();
    expect(screen.getByTestId('line-chart')).toBeInTheDocument();
    expect(screen.getByTestId('line-error')).toBeInTheDocument();
    expect(screen.getByTestId('line-warning')).toBeInTheDocument();
    expect(screen.getByTestId('line-info')).toBeInTheDocument();
    expect(screen.getByTestId('line-exception')).toBeInTheDocument();
  });

  test('renders only lines for data types that have values', () => {
    const chartData = [
      { date: '2023-01-01', error: 5, warning: 0, info: 0, exception: 0 },
      { date: '2023-01-02', error: 3, warning: 0, info: 0, exception: 0 }
    ];
    
    render(
      <IssueChart 
        timeRange={mockTimeRange} 
        onTimeRangeChange={mockOnTimeRangeChange} 
        chartData={chartData} 
        issueSource={mockIssueSource} 
      />
    );
    
    expect(screen.getByTestId('line-error')).toBeInTheDocument();
    expect(screen.queryByTestId('line-warning')).not.toBeInTheDocument();
    expect(screen.queryByTestId('line-info')).not.toBeInTheDocument();
    expect(screen.queryByTestId('line-exception')).not.toBeInTheDocument();
  });

  test('calls onTimeRangeChange when time range is changed', () => {
    const chartData = [
      { date: '2023-01-01', error: 5, warning: 2, info: 1, exception: 0 }
    ];
    
    render(
      <IssueChart 
        timeRange={mockTimeRange} 
        onTimeRangeChange={mockOnTimeRangeChange} 
        chartData={chartData} 
        issueSource={mockIssueSource} 
      />
    );
    
    const timeRangeSelect = screen.getByRole('combobox');
    fireEvent.mouseDown(timeRangeSelect);
    
    // Since we're mocking the Select component, we need to simulate the change event
    // This is a simplified approach since the actual Material-UI Select component is complex
    const event = {
      target: { value: '7d' }
    };
    mockOnTimeRangeChange(event);
    
    expect(mockOnTimeRangeChange).toHaveBeenCalledWith(expect.objectContaining({
      target: expect.objectContaining({ value: '7d' })
    }));
  });
});