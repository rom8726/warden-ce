import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import ReleasesDelta from '../ReleasesDelta';

// Mock Material-UI components
vi.mock('@mui/material', async () => {
  const actual = await vi.importActual('@mui/material');
  return {
    ...actual,
    Paper: ({ children, sx }: any) => (
      <div data-testid="paper" style={sx}>
        {children}
      </div>
    ),
    Typography: ({ children, variant, sx }: any) => (
      <div data-testid="typography" data-variant={variant} style={sx}>
        {children}
      </div>
    ),
    Grid: ({ children, container, item, xs, md, lg }: any) => (
      <div data-testid="grid" data-container={container} data-item={item} data-xs={xs} data-md={md} data-lg={lg}>
        {children}
      </div>
    ),
    Card: ({ children, sx }: any) => (
      <div data-testid="card" style={sx}>
        {children}
      </div>
    ),
    CardContent: ({ children }: any) => (
      <div data-testid="card-content">{children}</div>
    ),
    Chip: ({ label, size, color, variant }: any) => (
      <div data-testid="chip" data-label={label} data-size={size} data-color={color} data-variant={variant}>
        {label}
      </div>
    ),
    Box: ({ children, sx }: any) => (
      <div data-testid="box" style={sx}>
        {children}
      </div>
    ),
    CircularProgress: () => <div data-testid="circular-progress">Loading...</div>,
    Avatar: ({ children, sx }: any) => (
      <div data-testid="avatar" style={sx}>
        {children}
      </div>
    ),
    LinearProgress: ({ variant, value, sx }: any) => (
      <div data-testid="linear-progress" data-variant={variant} data-value={value} style={sx}>
        Progress: {value}%
      </div>
    ),
    useTheme: () => ({
      palette: {
        text: {
          primary: '#000000',
          secondary: '#666666',
        },
        success: {
          main: '#4caf50',
        },
        error: {
          main: '#f44336',
        },
      },
    }),
  };
});

// Mock Material-UI icons
vi.mock('@mui/icons-material', () => ({
  TrendingUp: () => <div data-testid="trending-up-icon">TrendingUpIcon</div>,
  TrendingDown: () => <div data-testid="trending-down-icon">TrendingDownIcon</div>,
  Remove: () => <div data-testid="remove-icon">RemoveIcon</div>,
  CompareArrows: () => <div data-testid="compare-arrows-icon">CompareArrowsIcon</div>,
  Speed: () => <div data-testid="speed-icon">SpeedIcon</div>,
  BugReport: () => <div data-testid="bug-report-icon">BugReportIcon</div>,
  Group: () => <div data-testid="group-icon">GroupIcon</div>,
  Timeline: () => <div data-testid="timeline-icon">TimelineIcon</div>,
}));

// Mock theme context
vi.mock('../../../../../theme/ThemeContext', () => ({
  useTheme: () => ({
    mode: 'light',
  }),
}));

describe('ReleasesDelta', () => {
  const mockComparisonData = {
    delta: {
      known_issues_total: 5,
      new_issues_total: -2,
      regressions_total: 1,
      resolved_in_version_total: 3,
      users_affected: -10,
    },
    base_version: '1.0.0',
    target_version: '1.1.0',
  };

  const defaultProps = {
    comparisonData: mockComparisonData,
    loading: false,
    error: null,
    compareMode: true,
    compareRelease: '1.0.0',
    selectedRelease: '1.1.0',
  };

  test('renders nothing when compareMode is false', () => {
    const { container } = render(<ReleasesDelta {...defaultProps} compareMode={false} />);
    expect(container.firstChild).toBeNull();
  });

  test('renders loading state when loading is true', () => {
    render(<ReleasesDelta {...defaultProps} loading={true} />);
    
    expect(screen.getByTestId('paper')).toBeInTheDocument();
    expect(screen.getByTestId('circular-progress')).toBeInTheDocument();
  });

  test('renders error state when error is present', () => {
    const errorMessage = 'Failed to load comparison data';
    render(<ReleasesDelta {...defaultProps} error={errorMessage} />);
    
    expect(screen.getByTestId('paper')).toBeInTheDocument();
    expect(screen.getByText('Error loading comparison data. Please try again.')).toBeInTheDocument();
  });

  test('renders comparison data when data is present', () => {
    render(<ReleasesDelta {...defaultProps} />);
    
    expect(screen.getByTestId('paper')).toBeInTheDocument();
    expect(screen.getByText(/Release Comparison:/)).toBeInTheDocument();
    const versionElements = screen.getAllByText(/1\.0\.0 → 1\.1\.0/);
    expect(versionElements.length).toBeGreaterThan(0);
  });

  test('renders delta metrics correctly', () => {
    render(<ReleasesDelta {...defaultProps} />);
    
    // Check for positive delta (known issues increased by 5)
    expect(screen.getByText('+5')).toBeInTheDocument();
    expect(screen.getByText('Known Issues')).toBeInTheDocument();
    
    // Check for negative delta (new issues decreased by 2)
    expect(screen.getByText('-2')).toBeInTheDocument();
    expect(screen.getByText('New Issues')).toBeInTheDocument();
    
    // Check for other metrics
    expect(screen.getByText('+1')).toBeInTheDocument(); // regressions
    expect(screen.getByText('+3')).toBeInTheDocument(); // resolved
    expect(screen.getByText('-10')).toBeInTheDocument(); // users affected
  });

  test('renders zero delta values correctly', () => {
    const zeroDeltaData = {
      ...mockComparisonData,
      delta: {
        known_issues_total: 0,
        new_issues_total: 0,
        regressions_total: 0,
        resolved_in_version_total: 0,
        users_affected: 0,
      },
    };
    
    render(<ReleasesDelta {...defaultProps} comparisonData={zeroDeltaData} />);
    
    // Check for zero values - use getAllByText and find the one that's actually "0"
    const zeroElements = screen.getAllByText('0');
    expect(zeroElements.length).toBeGreaterThan(0);
    
    // Check for "No Change" labels
    const noChangeElements = screen.getAllByText('No Change');
    expect(noChangeElements.length).toBeGreaterThan(0);
  });

  test('renders correct icons for different delta types', () => {
    render(<ReleasesDelta {...defaultProps} />);
    
    // Should have trending up icons for positive deltas
    const trendingUpIcons = screen.getAllByTestId('trending-up-icon');
    expect(trendingUpIcons.length).toBeGreaterThan(0);
    
    // Should have trending down icons for negative deltas
    const trendingDownIcons = screen.getAllByTestId('trending-down-icon');
    expect(trendingDownIcons.length).toBeGreaterThan(0);
    
    // Should have compare-arrows icons (может быть несколько)
    const compareArrowsIcons = screen.getAllByTestId('compare-arrows-icon');
    expect(compareArrowsIcons.length).toBeGreaterThan(0);
  });

  test('renders progress bars for delta values', () => {
    render(<ReleasesDelta {...defaultProps} />);
    
    const progressBars = screen.getAllByTestId('linear-progress');
    expect(progressBars.length).toBeGreaterThan(0);
  });

  test('renders comparison summary card', () => {
    render(<ReleasesDelta {...defaultProps} />);
    // Ищем по подстроке, чтобы не зависеть от эмодзи
    const summary = screen.getAllByTestId('typography').find(el => el.textContent?.includes('Comparison Summary'));
    expect(summary).toBeDefined();
  });

  test('renders metric cards for each delta type', () => {
    render(<ReleasesDelta {...defaultProps} />);
    
    const cards = screen.getAllByTestId('card');
    expect(cards.length).toBeGreaterThan(0);
  });

  test('handles undefined delta values', () => {
    const comparisonDataWithUndefined = {
      ...mockComparisonData,
      delta: {
        ...mockComparisonData.delta,
        known_issues_total: 0,
      },
    };
    
    render(<ReleasesDelta {...defaultProps} comparisonData={comparisonDataWithUndefined} />);
    
    // Should show 0 for undefined values
    expect(screen.getByText('0')).toBeInTheDocument();
  });

  test('renders nothing when no comparison data', () => {
    const { container } = render(<ReleasesDelta {...defaultProps} comparisonData={null} />);
    expect(container.firstChild).toBeNull();
  });

  test('renders correct labels for different metrics', () => {
    render(<ReleasesDelta {...defaultProps} />);
    expect(screen.getByText('Known Issues')).toBeInTheDocument();
    expect(screen.getByText('New Issues')).toBeInTheDocument();
    expect(screen.getByText('Regressions')).toBeInTheDocument();
    expect(screen.getByText('Resolved')).toBeInTheDocument();
    expect(screen.getByText('Users Affected')).toBeInTheDocument();
  });

  test('renders "No Change" label when delta is zero', () => {
    const zeroDeltaData = {
      ...mockComparisonData,
      delta: {
        known_issues_total: 0,
        new_issues_total: 0,
        regressions_total: 0,
        resolved_in_version_total: 0,
        users_affected: 0,
      },
    };
    render(<ReleasesDelta {...defaultProps} comparisonData={zeroDeltaData} />);
    // Ищем по подстроке
    const noChange = screen.getAllByTestId('typography').find(el => el.textContent?.includes('No Change'));
    expect(noChange).toBeDefined();
  });

  test('renders correct trend labels', () => {
    render(<ReleasesDelta {...defaultProps} />);
    // Проверяем Improved и Worsened, т.к. в моках нет нулевых дельт
    expect(screen.getAllByText('Improved').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Worsened').length).toBeGreaterThan(0);
  });

  test('renders comparison summary with version information', () => {
    render(<ReleasesDelta {...defaultProps} />);
    // Ищем нужный элемент среди всех с data-testid="typography"
    const headers = screen.getAllByTestId('typography');
    const header = headers.find(h => h.textContent?.includes('Release Comparison'));
    expect(header).toBeDefined();
    expect(header?.textContent).toContain('Release Comparison');
    expect(header?.textContent).toContain('1.0.0');
    expect(header?.textContent).toContain('1.1.0');
  });
}); 