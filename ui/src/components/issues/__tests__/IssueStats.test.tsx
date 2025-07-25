import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import IssueStats from '../IssueStats';
import { IssueStatus } from '../../../generated/api/client/api';

// Mock the formatDate utility function
vi.mock('../../../utils/issues/issueUtils', () => ({
  formatDate: (dateString: any) => `Formatted: ${dateString}`
}));

// Mock the BaseInfoPanel component
// @ts-ignore
vi.mock('../BaseInfoPanel', () => ({
  default: ({ issueData }) => (
    <div data-testid="base-info-panel">
      {issueData ? 'Issue data provided' : 'No issue data'}
    </div>
  )
}));

// Mock the useIssuePageStyles hook
vi.mock('../IssuePageStyles', () => ({
  useIssuePageStyles: () => ({
    statsContainer: { mb: 2 },
    statsPaper: { p: 1.5, height: '100%' },
    statsIconContainer: { display: 'flex', alignItems: 'center', mb: 0.5 },
    statsIcon: { mr: 0.5 },
    statsValue: { fontWeight: 'bold', mb: 0.5 }
  })
}));

describe('IssueStats', () => {
  const mockIssue = {
    id: '123',
    title: 'Test Issue',
    count: 42,
    first_seen: '2023-01-01T00:00:00Z',
    last_seen: '2023-01-02T00:00:00Z',
    status: IssueStatus.Unresolved
  };

  test('renders occurrence count and dates', () => {
    render(<IssueStats issue={mockIssue} issueData={null} />);
    
    expect(screen.getByText('Occurrences')).toBeInTheDocument();
    expect(screen.getByText('42')).toBeInTheDocument();
    expect(screen.getByText('First seen: Formatted: 2023-01-01T00:00:00Z')).toBeInTheDocument();
    expect(screen.getByText('Last seen: Formatted: 2023-01-02T00:00:00Z')).toBeInTheDocument();
  });

  test('renders BaseInfoPanel with issueData', () => {
    const issueData = { events: [{ platform: 'javascript' }] };
    render(<IssueStats issue={mockIssue} issueData={issueData} />);
    
    const baseInfoPanel = screen.getByTestId('base-info-panel');
    expect(baseInfoPanel).toBeInTheDocument();
    expect(baseInfoPanel).toHaveTextContent('Issue data provided');
  });

  test('renders BaseInfoPanel without issueData', () => {
    render(<IssueStats issue={mockIssue} issueData={null} />);
    
    const baseInfoPanel = screen.getByTestId('base-info-panel');
    expect(baseInfoPanel).toBeInTheDocument();
    expect(baseInfoPanel).toHaveTextContent('No issue data');
  });

  test('renders resolved information when issue is resolved', () => {
    const resolvedIssue = {
      ...mockIssue,
      status: IssueStatus.Resolved,
      resolved_at: '2023-01-03T00:00:00Z',
      resolved_by: 'User123'
    };
    
    render(<IssueStats issue={resolvedIssue} issueData={null} />);
    
    expect(screen.getByText('Resolved at: Formatted: 2023-01-03T00:00:00Z by User123')).toBeInTheDocument();
  });

  test('renders resolved information without resolver when not available', () => {
    const resolvedIssue = {
      ...mockIssue,
      status: IssueStatus.Resolved,
      resolved_at: '2023-01-03T00:00:00Z',
      resolved_by: undefined
    };
    
    render(<IssueStats issue={resolvedIssue} issueData={null} />);
    
    expect(screen.getByText('Resolved at: Formatted: 2023-01-03T00:00:00Z')).toBeInTheDocument();
    expect(screen.queryByText(/by/)).not.toBeInTheDocument();
  });

  test('does not render resolved information when issue is not resolved', () => {
    render(<IssueStats issue={mockIssue} issueData={null} />);
    
    expect(screen.queryByText(/Resolved at/)).not.toBeInTheDocument();
  });
});