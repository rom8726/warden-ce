import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import EventsTab from '../EventsTab';
import * as issueUtils from '../../../utils/issues/issueUtils';

// Mock the utility functions
vi.mock('../../../utils/issues/issueUtils', () => ({
  getLevelColor: vi.fn().mockReturnValue('error'),
  formatDate: vi.fn().mockReturnValue('2023-01-01 12:00:00')
}));

// Mock the EventTags component
// @ts-ignore
vi.mock('../EventTags', () => ({
  default: ({ tags }) => <div data-testid="event-tags">{JSON.stringify(tags)}</div>
}));

describe('EventsTab', () => {
  const mockIssue = { id: '123', title: 'Test Issue' };
  
  test('renders alert when no events are available', () => {
    // @ts-ignore
    render(<EventsTab issueData={null} issue={mockIssue} />);
    expect(screen.getByText('No events available for this issue.')).toBeInTheDocument();
  });

  test('renders alert when events array is empty', () => {
    const issueData = { events: [] };
    // @ts-ignore
    render(<EventsTab issueData={issueData} issue={mockIssue} />);
    expect(screen.getByText('No events available for this issue.')).toBeInTheDocument();
  });

  test('renders table with events when events are available', () => {
    const issueData = {
      events: [
        {
          event_id: 'event1',
          timestamp: '2023-01-01T12:00:00Z',
          level: 'error',
          platform: 'javascript',
          environment: 'production',
          server_name: 'server1',
          tags: { tag1: 'value1' }
        },
        {
          event_id: 'event2',
          timestamp: '2023-01-02T12:00:00Z',
          level: 'warning',
          platform: 'python',
          environment: 'staging',
          server_name: 'server2',
          tags: { tag2: 'value2' }
        }
      ]
    };
    
    // @ts-ignore
    render(<EventsTab issueData={issueData} issue={mockIssue} />);
    
    // Check table headers
    expect(screen.getByText('Timestamp')).toBeInTheDocument();
    expect(screen.getByText('Event ID')).toBeInTheDocument();
    expect(screen.getByText('Level')).toBeInTheDocument();
    expect(screen.getByText('Platform')).toBeInTheDocument();
    expect(screen.getByText('Environment')).toBeInTheDocument();
    expect(screen.getByText('Server Name')).toBeInTheDocument();
    expect(screen.getByText('Tags')).toBeInTheDocument();
    
    // Check event data
    expect(screen.getByText('event1')).toBeInTheDocument();
    expect(screen.getByText('event2')).toBeInTheDocument();
    expect(screen.getByText('javascript')).toBeInTheDocument();
    expect(screen.getByText('python')).toBeInTheDocument();
    expect(screen.getByText('production')).toBeInTheDocument();
    expect(screen.getByText('staging')).toBeInTheDocument();
    expect(screen.getByText('server1')).toBeInTheDocument();
    expect(screen.getByText('server2')).toBeInTheDocument();
    
    // Verify utility functions were called
    expect(issueUtils.formatDate).toHaveBeenCalledWith('2023-01-01T12:00:00Z');
    expect(issueUtils.formatDate).toHaveBeenCalledWith('2023-01-02T12:00:00Z');
    expect(issueUtils.getLevelColor).toHaveBeenCalledWith('error');
    expect(issueUtils.getLevelColor).toHaveBeenCalledWith('warning');
    
    // Check EventTags component was rendered
    expect(screen.getAllByTestId('event-tags')).toHaveLength(2);
  });

  test('handles events without event_id', () => {
    const issueData = {
      events: [
        {
          timestamp: '2023-01-01T12:00:00Z',
          level: 'error',
          platform: 'javascript',
          environment: 'production',
          server_name: 'server1',
          tags: { tag1: 'value1' }
        }
      ]
    };
    
    render(<EventsTab issueData={issueData} issue={mockIssue} />);
    
    // The component should still render without errors
    expect(screen.getByText('javascript')).toBeInTheDocument();
    expect(screen.getByText('production')).toBeInTheDocument();
    expect(screen.getByText('server1')).toBeInTheDocument();
  });
});