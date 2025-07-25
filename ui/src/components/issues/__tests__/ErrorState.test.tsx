import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import ErrorState from '../ErrorState';

// Mock the useIssuePageStyles hook
vi.mock('../IssuePageStyles', () => ({
  useIssuePageStyles: () => ({
    errorAlert: {
      mb: 2
    }
  })
}));

describe('ErrorState', () => {
  test('renders an error alert with the provided message', () => {
    const errorMessage = 'An error occurred while loading the data';
    render(<ErrorState message={errorMessage} />);
    
    // Check that the Alert component is rendered with the correct message
    expect(screen.getByText(errorMessage)).toBeInTheDocument();
    
    // Check that the Alert has the correct severity
    const alert = screen.getByRole('alert');
    expect(alert).toHaveClass('MuiAlert-standardError');
  });
  
  test('renders with a different error message', () => {
    const errorMessage = 'Network error: Failed to fetch';
    render(<ErrorState message={errorMessage} />);
    
    expect(screen.getByText(errorMessage)).toBeInTheDocument();
  });
});