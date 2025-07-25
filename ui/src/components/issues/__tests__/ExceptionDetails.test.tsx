import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import ExceptionDetails from '../ExceptionDetails';

// Mock the StackTraceViewer component
// @ts-ignore
vi.mock('../../StackTraceViewer', () => ({
  default: ({ stacktrace, platform }) => (
    <div data-testid="stack-trace-viewer">
      <div>Stacktrace: {stacktrace}</div>
      <div>Platform: {platform || 'not specified'}</div>
    </div>
  )
}));

// Mock the useIssuePageStyles hook
vi.mock('../IssuePageStyles', () => ({
  useIssuePageStyles: () => ({
    exceptionHeader: {
      p: 1,
      backgroundColor: 'rgba(255, 0, 0, 0.1)',
      borderRadius: '4px 4px 0 0',
      fontWeight: 'bold',
      color: 'red',
      mb: 1.5
    }
  })
}));

describe('ExceptionDetails', () => {
  test('renders nothing when stacktrace is empty', () => {
    const exception = {
      exception_type: 'TypeError',
      exception_value: 'Cannot read property of undefined',
      stacktrace: ''
    };

    const { container } = render(<ExceptionDetails exception={exception} />);
    expect(container.firstChild).toBeNull();
  });

  test('renders exception details when stacktrace is provided', () => {
    const exception = {
      exception_type: 'TypeError',
      exception_value: 'Cannot read property of undefined',
      stacktrace: 'Error: Cannot read property of undefined\n    at Object.method (/path/to/file.js:10:15)'
    };

    render(<ExceptionDetails exception={exception} />);

    // Check that the exception header is rendered
    expect(screen.getByText('TypeError: Cannot read property of undefined')).toBeInTheDocument();

    // Check that the StackTraceViewer is rendered with the correct props
    const stackTraceViewer = screen.getByTestId('stack-trace-viewer');
    expect(stackTraceViewer).toBeInTheDocument();

    // Check that the stacktrace text is contained within the viewer
    const stacktraceDiv = stackTraceViewer.querySelector('div:first-child');
    expect(stacktraceDiv).toHaveTextContent('Stacktrace: Error: Cannot read property of undefined');

    expect(screen.getByText('Platform: not specified')).toBeInTheDocument();
  });

  test('passes platform to StackTraceViewer when provided', () => {
    const exception = {
      exception_type: 'TypeError',
      exception_value: 'Cannot read property of undefined',
      stacktrace: 'Error: Cannot read property of undefined\n    at Object.method (/path/to/file.js:10:15)'
    };

    render(<ExceptionDetails exception={exception} platform="javascript" />);

    expect(screen.getByText('Platform: javascript')).toBeInTheDocument();
  });
});
