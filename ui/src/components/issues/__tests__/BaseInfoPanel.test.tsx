import { render, screen } from '@testing-library/react';
import BaseInfoPanel from '../BaseInfoPanel';

describe('BaseInfoPanel', () => {
  test('renders nothing when issueData is null', () => {
    const { container } = render(<BaseInfoPanel issueData={null} />);
    expect(container.firstChild).toBeNull();
  });

  test('renders nothing when events array is empty', () => {
    const issueData = { events: [] };
    const { container } = render(<BaseInfoPanel issueData={issueData} />);
    expect(container.firstChild).toBeNull();
  });

  test('renders basic info from the first event', () => {
    const issueData = {
      events: [
        {
          platform: 'javascript',
          server_name: 'web-server-1',
          environment: 'production'
        }
      ]
    };
    
    render(<BaseInfoPanel issueData={issueData} />);
    
    // Check that the title is rendered
    expect(screen.getByText('Base Info')).toBeInTheDocument();
    
    // Check that the basic info is rendered
    expect(screen.getByText('Platform')).toBeInTheDocument();
    expect(screen.getByText('javascript')).toBeInTheDocument();
    expect(screen.getByText('Server Name')).toBeInTheDocument();
    expect(screen.getByText('web-server-1')).toBeInTheDocument();
    expect(screen.getByText('Environment')).toBeInTheDocument();
    expect(screen.getByText('production')).toBeInTheDocument();
  });

  test('renders N/A for missing values', () => {
    const issueData = {
      events: [
        {
          platform: 'javascript',
          server_name: null,
          environment: undefined
        }
      ]
    };
    
    render(<BaseInfoPanel issueData={issueData} />);
    
    expect(screen.getByText('Platform')).toBeInTheDocument();
    expect(screen.getByText('javascript')).toBeInTheDocument();
    expect(screen.getByText('Server Name')).toBeInTheDocument();
    expect(screen.getAllByText('N/A')).toHaveLength(2); // For both server_name and environment
  });

  test('renders exception information when available', () => {
    const issueData = {
      events: [
        {
          platform: 'javascript',
          server_name: 'web-server-1',
          environment: 'production',
          exception_type: 'TypeError',
          exception_value: 'Cannot read property of undefined'
        }
      ]
    };
    
    render(<BaseInfoPanel issueData={issueData} />);
    
    // Check that the exception info is rendered
    expect(screen.getByText('Exception Type')).toBeInTheDocument();
    expect(screen.getByText('TypeError')).toBeInTheDocument();
    expect(screen.getByText('Exception Value')).toBeInTheDocument();
    expect(screen.getByText('Cannot read property of undefined')).toBeInTheDocument();
  });
});