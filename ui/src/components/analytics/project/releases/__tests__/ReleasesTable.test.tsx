import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import ReleasesTable from '../ReleasesTable';

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
    Typography: ({ children, variant, color, gutterBottom, sx }: any) => (
      <div data-testid="typography" data-variant={variant} data-color={color} style={sx}>
        {children}
      </div>
    ),
    Box: ({ children, sx }: any) => (
      <div data-testid="box" style={sx}>
        {children}
      </div>
    ),
    CircularProgress: () => <div data-testid="circular-progress">Loading...</div>,
    TableContainer: ({ children }: any) => (
      <div data-testid="table-container">{children}</div>
    ),
    Table: ({ children }: any) => (
      <table data-testid="table">{children}</table>
    ),
    TableHead: ({ children }: any) => (
      <thead data-testid="table-head">{children}</thead>
    ),
    TableBody: ({ children }: any) => (
      <tbody data-testid="table-body">{children}</tbody>
    ),
    TableRow: ({ children, hover, selected, onClick, sx }: any) => (
      <tr 
        data-testid="table-row" 
        data-hover={hover}
        data-selected={selected}
        onClick={onClick}
        style={sx}
      >
        {children}
      </tr>
    ),
    TableCell: ({ children, align, component, scope }: any) => (
      <td data-testid="table-cell" data-align={align} data-component={component} data-scope={scope}>
        {children}
      </td>
    ),
    Chip: ({ label, size, color, variant }: any) => (
      <div data-testid="chip" data-label={label} data-size={size} data-color={color} data-variant={variant}>
        {label}
      </div>
    ),
    IconButton: ({ children, size, color, onClick }: any) => (
      <button 
        data-testid="icon-button" 
        data-size={size} 
        data-color={color}
        onClick={onClick}
      >
        {children}
      </button>
    ),
    Tooltip: ({ children, title }: any) => (
      <div data-testid="tooltip" data-title={title}>
        {children}
      </div>
    ),
  };
});

// Mock Material-UI icons
vi.mock('@mui/icons-material', () => ({
  Assessment: () => <div data-testid="assessment-icon">AssessmentIcon</div>,
  Visibility: () => <div data-testid="visibility-icon">ViewDetailsIcon</div>,
  CompareArrows: () => <div data-testid="compare-icon">CompareIcon</div>,
}));

describe('ReleasesTable', () => {
  const mockReleases = [
    {
      version: '1.0.0',
      known_issues_total: 10,
      new_issues_total: 5,
      regressions_total: 2,
      resolved_in_version_total: 8,
      users_affected: 150,
      created_at: '2024-01-01T00:00:00Z',
    },
    {
      version: '1.1.0',
      known_issues_total: 12,
      new_issues_total: 3,
      regressions_total: 1,
      resolved_in_version_total: 10,
      users_affected: 120,
      created_at: '2024-01-15T00:00:00Z',
    },
  ];

  const defaultProps = {
    releases: mockReleases,
    loading: false,
    error: null,
    selectedRelease: null,
    compareMode: false,
    comparisonData: null,
    onReleaseSelect: vi.fn(),
    onCompareRelease: vi.fn(),
  };

  test('renders loading state when loading is true', () => {
    render(<ReleasesTable {...defaultProps} loading={true} />);
    
    expect(screen.getByTestId('paper')).toBeInTheDocument();
    expect(screen.getByText('Release Summary')).toBeInTheDocument();
    expect(screen.getByTestId('circular-progress')).toBeInTheDocument();
  });

  test('renders error state when error is present', () => {
    const errorMessage = 'Failed to load releases';
    render(<ReleasesTable {...defaultProps} error={errorMessage} />);
    
    expect(screen.getByTestId('paper')).toBeInTheDocument();
    expect(screen.getByText('Release Summary')).toBeInTheDocument();
    expect(screen.getByText('Error loading release data. Please try again.')).toBeInTheDocument();
  });

  test('renders empty state when no releases', () => {
    render(<ReleasesTable {...defaultProps} releases={[]} />);
    
    expect(screen.getByTestId('paper')).toBeInTheDocument();
    expect(screen.getByText('Release Summary')).toBeInTheDocument();
    expect(screen.getByText('No releases found.')).toBeInTheDocument();
  });

  test('renders table with releases data', () => {
    render(<ReleasesTable {...defaultProps} />);
    
    expect(screen.getByTestId('table')).toBeInTheDocument();
    expect(screen.getByText('1.0.0')).toBeInTheDocument();
    expect(screen.getByText('1.1.0')).toBeInTheDocument();
    
    // Check for specific values in the first row
    const cells = screen.getAllByTestId('table-cell');
    const firstRowCells = cells.slice(8, 15); // Skip header row, get first data row
    expect(firstRowCells[1]).toHaveTextContent('10'); // known_issues_total for 1.0.0
    expect(firstRowCells[2]).toHaveTextContent('5'); // new_issues_total for 1.0.0
  });

  test('calls onReleaseSelect when row is clicked', () => {
    const onReleaseSelect = vi.fn();
    render(<ReleasesTable {...defaultProps} onReleaseSelect={onReleaseSelect} />);
    
    const rows = screen.getAllByTestId('table-row');
    // Skip header row, click first data row
    fireEvent.click(rows[1]);
    
    expect(onReleaseSelect).toHaveBeenCalledWith('1.0.0');
  });

  test('calls onReleaseSelect when view details button is clicked', () => {
    const onReleaseSelect = vi.fn();
    render(<ReleasesTable {...defaultProps} onReleaseSelect={onReleaseSelect} />);
    
    const viewButtons = screen.getAllByTestId('icon-button');
    fireEvent.click(viewButtons[0]); // First view button
    
    expect(onReleaseSelect).toHaveBeenCalledWith('1.0.0');
  });

  test('calls onCompareRelease when compare button is clicked', () => {
    const onCompareRelease = vi.fn();
    render(<ReleasesTable {...defaultProps} onCompareRelease={onCompareRelease} />);
    
    const compareButtons = screen.getAllByTestId('icon-button');
    fireEvent.click(compareButtons[1]); // First compare button
    
    expect(onCompareRelease).toHaveBeenCalledWith('1.0.0');
  });

  test('shows selected release with selected styling', () => {
    render(<ReleasesTable {...defaultProps} selectedRelease="1.0.0" />);
    
    const rows = screen.getAllByTestId('table-row');
    expect(rows[1]).toHaveAttribute('data-selected', 'true');
  });

  test('shows delta column when in compare mode with comparison data', () => {
    const comparisonData = {
      base: { version: '1.0.0' },
      target: { version: '1.1.0' },
    };
    
    render(<ReleasesTable {...defaultProps} compareMode={true} comparisonData={comparisonData} />);
    
    // Check that delta column headers are present
    expect(screen.getByText('Delta')).toBeInTheDocument();
  });

  test('shows base chip for base version in compare mode', () => {
    const comparisonData = {
      base: { version: '1.0.0' },
      target: { version: '1.1.0' },
    };
    
    render(<ReleasesTable {...defaultProps} compareMode={true} comparisonData={comparisonData} />);
    
    const chips = screen.getAllByTestId('chip');
    const baseChip = chips.find(chip => chip.getAttribute('data-label') === 'Base');
    expect(baseChip).toBeInTheDocument();
  });

  test('shows target chip for target version in compare mode', () => {
    const comparisonData = {
      base: { version: '1.0.0' },
      target: { version: '1.1.0' },
    };
    
    render(<ReleasesTable {...defaultProps} compareMode={true} comparisonData={comparisonData} />);
    
    const chips = screen.getAllByTestId('chip');
    const targetChip = chips.find(chip => chip.getAttribute('data-label') === 'Target');
    expect(targetChip).toBeInTheDocument();
  });

  test('formats date correctly', () => {
    render(<ReleasesTable {...defaultProps} />);
    
    // Check that dates are formatted (should contain slashes for ru-RU locale)
    const dateCells = screen.getAllByTestId('table-cell');
    const dateCell = dateCells.find(cell => 
      cell.textContent?.includes('01.01.2024') || cell.textContent?.includes('15.01.2024')
    );
    expect(dateCell).toBeInTheDocument();
  });

  test('prevents event propagation on button clicks', () => {
    const onReleaseSelect = vi.fn();
    const onCompareRelease = vi.fn();
    
    render(
      <ReleasesTable 
        {...defaultProps} 
        onReleaseSelect={onReleaseSelect}
        onCompareRelease={onCompareRelease}
      />
    );
    
    const viewButtons = screen.getAllByTestId('icon-button');
    fireEvent.click(viewButtons[0]);
    
    // Should only call onReleaseSelect, not trigger row click
    expect(onReleaseSelect).toHaveBeenCalledTimes(1);
  });
}); 