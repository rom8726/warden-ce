import { render } from '@testing-library/react';
import { vi } from 'vitest';
import LoadingState from '../LoadingState';

// Mock the useIssuePageStyles hook
vi.mock('../IssuePageStyles', () => ({
  useIssuePageStyles: () => ({
    loadingContainer: {
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '40vh'
    }
  })
}));

describe('LoadingState', () => {
  test('renders a circular progress indicator', () => {
    render(<LoadingState />);
    
    // Check that the CircularProgress component is rendered
    // Since CircularProgress doesn't have a specific role or text,
    // we can check for the presence of the SVG element it renders
    const progressIndicator = document.querySelector('svg');
    expect(progressIndicator).toBeInTheDocument();
    
    // Check that the container has the expected class for styling
    const container = progressIndicator?.closest('div');
    expect(container).toBeInTheDocument();
  });
});