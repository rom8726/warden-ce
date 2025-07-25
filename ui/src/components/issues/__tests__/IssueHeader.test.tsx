import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import IssueHeader from '../IssueHeader';
import { IssueStatus, IssueSource } from '../../../generated/api/client/api';

// Mock the utility functions
vi.mock('../../../utils/issues/issueUtils', () => ({
  getLevelColor: (level: any) => {
    switch (level) {
      case 'error': return 'error';
      case 'warning': return 'warning';
      case 'info': return 'info';
      default: return 'error';
    }
  },
  getLevelHexColor: (level: any) => {
    switch (level) {
      case 'error': return '#d32f2f';
      case 'warning': return '#ed6c02';
      case 'info': return '#0288d1';
      case 'fatal': return '#9c27b0';
      case 'exception': return '#f57c00';
      default: return '#d32f2f';
    }
  },
  getLevelBadgeStyles: (level: any) => {
    return {
      backgroundColor: level === 'error' ? '#d32f2f' : '#ed6c02',
      color: 'white'
    };
  },
  getStatusColor: (status: any) => {
    switch (status) {
      case IssueStatus.Resolved: return 'success';
      case IssueStatus.Unresolved: return 'error';
      case IssueStatus.Ignored: return 'default';
      default: return 'default';
    }
  }
}));

// Mock the useIssuePageStyles hook
vi.mock('../IssuePageStyles', () => ({
  useIssuePageStyles: () => ({
    headerContainer: { mb: 2 },
    headerTitleContainer: { display: 'flex', alignItems: 'center', mb: 0.5 },
    headerTitle: { ml: 0.5 },
    headerMessage: { mb: 1.5 },
    tagsContainer: { display: 'flex', flexWrap: 'wrap', gap: 0.5 }
  })
}));

describe('IssueHeader', () => {
  const mockIssue = {
    id: '123',
    title: 'Test Issue',
    message: 'This is a test issue message',
    level: 'error',
    status: IssueStatus.Unresolved,
    source: IssueSource.Exception,
    platform: 'javascript'
  };

  test('renders issue title and message', () => {
    // @ts-ignore
    render(<IssueHeader issue={mockIssue} />);
    
    expect(screen.getByText('Test Issue')).toBeInTheDocument();
    expect(screen.getByText('This is a test issue message')).toBeInTheDocument();
  });

  test('renders level, status, and source chips', () => {
    // @ts-ignore
    render(<IssueHeader issue={mockIssue} />);
    
    expect(screen.getByText('error')).toBeInTheDocument();
    expect(screen.getByText(IssueStatus.Unresolved)).toBeInTheDocument();
    expect(screen.getByText(IssueSource.Exception)).toBeInTheDocument();
  });

  test('renders tags when provided', () => {
    const tags = {
      'browser': 'Chrome',
      'os': 'Windows',
      'version': '1.0.0'
    };
    
    // @ts-ignore
    render(<IssueHeader issue={mockIssue} tags={tags} />);
    
    expect(screen.getByText('browser: Chrome')).toBeInTheDocument();
    expect(screen.getByText('os: Windows')).toBeInTheDocument();
    expect(screen.getByText('version: 1.0.0')).toBeInTheDocument();
  });

  test('renders status chip as clickable when onStatusChange is provided', () => {
    const handleStatusChange = vi.fn();
    // @ts-ignore
    render(<IssueHeader issue={mockIssue} onStatusChange={handleStatusChange} />);
    
    const statusChip = screen.getByText(IssueStatus.Unresolved);
    expect(statusChip).toBeInTheDocument();
    
    // Click the status chip to open the menu
    fireEvent.click(statusChip);
    
    // Check that the menu options are rendered
    expect(screen.getByText(`Change to ${IssueStatus.Resolved}`)).toBeInTheDocument();
    expect(screen.getByText(`Change to ${IssueStatus.Ignored}`)).toBeInTheDocument();
  });

  test('calls onStatusChange when a new status is selected', () => {
    const handleStatusChange = vi.fn();
    // @ts-ignore
    render(<IssueHeader issue={mockIssue} onStatusChange={handleStatusChange} />);
    
    // Click the status chip to open the menu
    fireEvent.click(screen.getByText(IssueStatus.Unresolved));
    
    // Click the "Change to Resolved" option
    fireEvent.click(screen.getByText(`Change to ${IssueStatus.Resolved}`));
    
    // Check that onStatusChange was called with the new status
    expect(handleStatusChange).toHaveBeenCalledWith(IssueStatus.Resolved);
  });

  test('disables status chip when statusChangeLoading is true', () => {
    const handleStatusChange = vi.fn();
    // @ts-ignore
    render(
      <IssueHeader 
        issue={mockIssue} 
        onStatusChange={handleStatusChange} 
        statusChangeLoading={true} 
      />
    );
    
    // Try to click the status chip
    fireEvent.click(screen.getByText(IssueStatus.Unresolved));
    
    // Check that the menu doesn't open (no menu items are rendered)
    expect(screen.queryByText(`Change to ${IssueStatus.Resolved}`)).not.toBeInTheDocument();
    expect(screen.queryByText(`Change to ${IssueStatus.Ignored}`)).not.toBeInTheDocument();
  });

  test('renders different available status options based on current status', () => {
    const handleStatusChange = vi.fn();
    
    // Test with Resolved status
    const resolvedIssue = { ...mockIssue, status: IssueStatus.Resolved };
    const { rerender } = render(
      <IssueHeader issue={resolvedIssue} onStatusChange={handleStatusChange} />
    );
    
    // Click the status chip to open the menu
    fireEvent.click(screen.getByText(IssueStatus.Resolved));
    
    // Check that only the "Change to Unresolved" option is available
    expect(screen.getByText(`Change to ${IssueStatus.Unresolved}`)).toBeInTheDocument();
    expect(screen.queryByText(`Change to ${IssueStatus.Ignored}`)).not.toBeInTheDocument();
    
    // Close the menu
    fireEvent.click(document.body);
    
    // Test with Ignored status
    const ignoredIssue = { ...mockIssue, status: IssueStatus.Ignored };
    rerender(<IssueHeader issue={ignoredIssue} onStatusChange={handleStatusChange} />);
    
    // Click the status chip to open the menu
    fireEvent.click(screen.getByText(IssueStatus.Ignored));
    
    // Check that both "Change to Resolved" and "Change to Unresolved" options are available
    expect(screen.getByText(`Change to ${IssueStatus.Resolved}`)).toBeInTheDocument();
    expect(screen.getByText(`Change to ${IssueStatus.Unresolved}`)).toBeInTheDocument();
  });
});