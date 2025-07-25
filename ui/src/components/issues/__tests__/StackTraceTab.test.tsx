import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import StackTraceTab from '../StackTraceTab';
import { IssueSource } from '../../../generated/api/client/api';

// Mock the ExceptionDetails component
vi.mock('../ExceptionDetails', () => ({
  default: ({ exception, platform }) => (
    <div data-testid="exception-details">
      <div>Type: {exception.exception_type}</div>
      <div>Value: {exception.exception_value}</div>
      <div>Platform: {platform || 'not specified'}</div>
    </div>
  )
}));

// Mock the StackTraceViewer component
vi.mock('../../StackTraceViewer', () => ({
  default: ({ stacktrace, platform }) => (
    <div data-testid="stack-trace-viewer">
      <div>Stacktrace: {stacktrace.substring(0, 20)}...</div>
      <div>Platform: {platform || 'not specified'}</div>
    </div>
  )
}));

// Mock the useIssuePageStyles hook
vi.mock('../IssuePageStyles', () => ({
  useIssuePageStyles: () => ({
    tabContent: {
      p: 1.5,
      bgcolor: 'background.default',
      borderRadius: 1
    }
  })
}));

describe('StackTraceTab', () => {
  const mockIssue = {
    id: '123',
    title: 'Test Issue',
    platform: 'javascript',
    source: IssueSource.Exception
  };

  test('renders ExceptionDetails for exception source with exception details', () => {
    const issueData = {
      events: [
        {
          exception_type: 'TypeError',
          exception_value: 'Cannot read property of undefined',
          exception_stacktrace: 'Error: Cannot read property of undefined\n    at Object.method (/path/to/file.js:10:15)'
        }
      ]
    };

    render(<StackTraceTab issueData={issueData} issue={mockIssue} />);

    expect(screen.getByTestId('exception-details')).toBeInTheDocument();
    expect(screen.getByText('Type: TypeError')).toBeInTheDocument();
    expect(screen.getByText('Value: Cannot read property of undefined')).toBeInTheDocument();
    expect(screen.getByText('Platform: javascript')).toBeInTheDocument();
  });

  test('renders StackTraceViewer when stackTrace prop is provided', () => {
    const stackTrace = 'Error: Something went wrong\n    at Object.method (/path/to/file.js:10:15)';

    render(<StackTraceTab issueData={null} issue={mockIssue} stackTrace={stackTrace} />);

    const stackTraceViewer = screen.getByTestId('stack-trace-viewer');
    expect(stackTraceViewer).toBeInTheDocument();

    // Check that the stacktrace text is contained within the viewer
    const stacktraceDiv = stackTraceViewer.querySelector('div:first-child');
    expect(stacktraceDiv).toHaveTextContent('Stacktrace: Error: Something we');

    expect(screen.getByText('Platform: javascript')).toBeInTheDocument();
  });

  test('renders event details for event source with events', () => {
    const eventIssue = {
      ...mockIssue,
      source: IssueSource.Event
    };

    const issueData = {
      events: [
        {
          message: 'Error occurred in application',
          exception_stacktrace: 'Error: Something went wrong\n    at Object.method (/path/to/file.js:10:15)'
        }
      ]
    };

    render(<StackTraceTab issueData={issueData} issue={eventIssue} />);

    expect(screen.getByText('Event Details')).toBeInTheDocument();
    expect(screen.getByText('Error occurred in application')).toBeInTheDocument();
    expect(screen.getByTestId('stack-trace-viewer')).toBeInTheDocument();
  });

  test('renders event data as JSON when no stack trace is available for event source', () => {
    const eventIssue = {
      ...mockIssue,
      source: IssueSource.Event
    };

    const issueData = {
      events: [
        {
          message: 'Error occurred in application',
          level: 'error',
          timestamp: '2023-01-01T12:00:00Z'
        }
      ]
    };

    render(<StackTraceTab issueData={issueData} issue={eventIssue} />);

    expect(screen.getByText('Event Details')).toBeInTheDocument();
    expect(screen.getByText('Error occurred in application')).toBeInTheDocument();
    expect(screen.getByText(/timestamp/)).toBeInTheDocument(); // Part of the JSON output
  });

  test('renders alert when no stack trace is available', () => {
    render(<StackTraceTab issueData={null} issue={mockIssue} />);

    expect(screen.getByText('No stack trace available for this issue.')).toBeInTheDocument();
  });
});
