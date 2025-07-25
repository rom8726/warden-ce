import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import ComparisonAlert from '../ComparisonAlert';

// Mock Material-UI components
vi.mock('@mui/material', async () => {
  const actual = await vi.importActual('@mui/material');
  return {
    ...actual,
    Alert: ({ children, severity, sx }: any) => (
      <div data-testid="alert" data-severity={severity} style={sx}>
        {children}
      </div>
    ),
    Box: ({ children, sx }: any) => (
      <div data-testid="box" style={sx}>
        {children}
      </div>
    ),
    Typography: ({ children, variant }: any) => (
      <div data-testid="typography" data-variant={variant}>
        {children}
      </div>
    ),
    CircularProgress: ({ size }: any) => (
      <div data-testid="circular-progress" data-size={size}>
        Loading...
      </div>
    ),
    Button: ({ children, onClick, variant, size, startIcon }: any) => (
      <button 
        data-testid="button" 
        data-variant={variant} 
        data-size={size}
        onClick={onClick}
      >
        {startIcon}
        {children}
      </button>
    ),
  };
});

// Mock Material-UI icons
vi.mock('@mui/icons-material', () => ({
  SwapHoriz: () => <div data-testid="swap-icon">SwapIcon</div>,
}));

describe('ComparisonAlert', () => {
  const defaultProps = {
    compareMode: false,
    loading: false,
    error: null,
    comparisonData: null,
    onSwitchComparison: vi.fn(),
  };

  test('renders nothing when compareMode is false', () => {
    const { container } = render(<ComparisonAlert {...defaultProps} />);
    expect(container.firstChild).toBeNull();
  });

  test('renders loading state when compareMode is true and loading is true', () => {
    render(<ComparisonAlert {...defaultProps} compareMode={true} loading={true} />);
    
    expect(screen.getByTestId('alert')).toBeInTheDocument();
    expect(screen.getByTestId('alert')).toHaveAttribute('data-severity', 'info');
    expect(screen.getByTestId('circular-progress')).toBeInTheDocument();
    expect(screen.getByText('Loading comparison data...')).toBeInTheDocument();
  });

  test('renders error state when compareMode is true and error is present', () => {
    const errorMessage = 'Failed to load comparison data';
    render(<ComparisonAlert {...defaultProps} compareMode={true} error={errorMessage} />);
    
    expect(screen.getByTestId('alert')).toBeInTheDocument();
    expect(screen.getByTestId('alert')).toHaveAttribute('data-severity', 'error');
    expect(screen.getByText('Error loading comparison data. Please try again.')).toBeInTheDocument();
  });

  test('renders comparison data when compareMode is true and comparisonData is present', () => {
    const comparisonData = {
      base_version: '1.0.0',
      target_version: '1.1.0',
    };
    
    render(<ComparisonAlert {...defaultProps} compareMode={true} comparisonData={comparisonData} />);
    
    expect(screen.getByTestId('alert')).toBeInTheDocument();
    expect(screen.getByTestId('alert')).toHaveAttribute('data-severity', 'info');
    
    // Check for the text content using a more flexible approach
    const typography = screen.getByTestId('typography');
    expect(typography).toHaveTextContent('Comparing releases:');
    expect(typography).toHaveTextContent('1.0.0');
    expect(typography).toHaveTextContent('1.1.0');
    
    expect(screen.getByTestId('button')).toBeInTheDocument();
    expect(screen.getByText('Switch Direction')).toBeInTheDocument();
  });

  test('calls onSwitchComparison when switch button is clicked', () => {
    const onSwitchComparison = vi.fn();
    const comparisonData = {
      base_version: '1.0.0',
      target_version: '1.1.0',
    };
    
    render(
      <ComparisonAlert 
        {...defaultProps} 
        compareMode={true} 
        comparisonData={comparisonData}
        onSwitchComparison={onSwitchComparison}
      />
    );
    
    const switchButton = screen.getByTestId('button');
    fireEvent.click(switchButton);
    
    expect(onSwitchComparison).toHaveBeenCalledTimes(1);
  });

  test('does not render switch button when onSwitchComparison is not provided', () => {
    const comparisonData = {
      base_version: '1.0.0',
      target_version: '1.1.0',
    };
    
    render(
      <ComparisonAlert 
        {...defaultProps} 
        compareMode={true} 
        comparisonData={comparisonData}
        onSwitchComparison={undefined}
      />
    );
    
    expect(screen.queryByTestId('button')).not.toBeInTheDocument();
    
    // Check for the text content using a more flexible approach
    const typography = screen.getByTestId('typography');
    expect(typography).toHaveTextContent('Comparing releases:');
    expect(typography).toHaveTextContent('1.0.0');
    expect(typography).toHaveTextContent('1.1.0');
  });

  test('renders nothing when compareMode is true but no data, loading, or error', () => {
    const { container } = render(
      <ComparisonAlert 
        {...defaultProps} 
        compareMode={true} 
        loading={false}
        error={null}
        comparisonData={null}
      />
    );
    expect(container.firstChild).toBeNull();
  });
}); 