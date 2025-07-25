import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { ThemeProvider } from '@mui/material/styles';
import { theme } from '../../../theme';
import ProjectEditForm from '../ProjectEditForm';

// Mock the API client
jest.mock('../../../api/apiClient', () => ({
  updateProject: jest.fn()
}));

const mockApiClient = require('../../../api/apiClient');

const defaultProps = {
  projectId: 1,
  initialName: 'Test Project',
  initialDescription: 'Test Description',
  onSave: jest.fn(),
  onCancel: jest.fn()
};

const renderWithTheme = (component: React.ReactElement) => {
  return render(
    <ThemeProvider theme={theme}>
      {component}
    </ThemeProvider>
  );
};

describe('ProjectEditForm', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders with initial values', () => {
    renderWithTheme(<ProjectEditForm {...defaultProps} />);
    
    expect(screen.getByDisplayValue('Test Project')).toBeInTheDocument();
    expect(screen.getByDisplayValue('Test Description')).toBeInTheDocument();
    expect(screen.getByText('Edit Project Details')).toBeInTheDocument();
  });

  it('validates required fields', async () => {
    renderWithTheme(<ProjectEditForm {...defaultProps} />);
    
    const nameInput = screen.getByDisplayValue('Test Project');
    fireEvent.change(nameInput, { target: { value: '' } });
    
    const saveButton = screen.getByText('Save Changes');
    fireEvent.click(saveButton);
    
    await waitFor(() => {
      expect(screen.getByText('Project name is required')).toBeInTheDocument();
    });
  });

  it('validates minimum name length', async () => {
    renderWithTheme(<ProjectEditForm {...defaultProps} />);
    
    const nameInput = screen.getByDisplayValue('Test Project');
    fireEvent.change(nameInput, { target: { value: 'A' } });
    
    const saveButton = screen.getByText('Save Changes');
    fireEvent.click(saveButton);
    
    await waitFor(() => {
      expect(screen.getByText('Project name must be at least 2 characters long')).toBeInTheDocument();
    });
  });

  it('calls onCancel when cancel button is clicked', () => {
    renderWithTheme(<ProjectEditForm {...defaultProps} />);
    
    const cancelButton = screen.getByText('Cancel');
    fireEvent.click(cancelButton);
    
    expect(defaultProps.onCancel).toHaveBeenCalled();
  });

  it('calls API and onSave when form is valid', async () => {
    const mockResponse = {
      data: {
        project: {
          name: 'Updated Project',
          description: 'Updated Description'
        }
      }
    };
    
    mockApiClient.updateProject.mockResolvedValue(mockResponse);
    
    renderWithTheme(<ProjectEditForm {...defaultProps} />);
    
    const saveButton = screen.getByText('Save Changes');
    fireEvent.click(saveButton);
    
    await waitFor(() => {
      expect(mockApiClient.updateProject).toHaveBeenCalledWith(1, {
        name: 'Test Project',
        description: 'Test Description'
      });
      expect(defaultProps.onSave).toHaveBeenCalledWith('Updated Project', 'Updated Description');
    });
  });

  it('shows error message when API call fails', async () => {
    const errorMessage = 'API Error';
    mockApiClient.updateProject.mockRejectedValue(new Error(errorMessage));
    
    renderWithTheme(<ProjectEditForm {...defaultProps} />);
    
    const saveButton = screen.getByText('Save Changes');
    fireEvent.click(saveButton);
    
    await waitFor(() => {
      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });
  });

  it('trims whitespace from inputs', async () => {
    const mockResponse = {
      data: {
        project: {
          name: '  Updated Project  ',
          description: '  Updated Description  '
        }
      }
    };
    
    mockApiClient.updateProject.mockResolvedValue(mockResponse);
    
    renderWithTheme(<ProjectEditForm {...defaultProps} />);
    
    const nameInput = screen.getByDisplayValue('Test Project');
    const descriptionInput = screen.getByDisplayValue('Test Description');
    
    fireEvent.change(nameInput, { target: { value: '  Updated Project  ' } });
    fireEvent.change(descriptionInput, { target: { value: '  Updated Description  ' } });
    
    const saveButton = screen.getByText('Save Changes');
    fireEvent.click(saveButton);
    
    await waitFor(() => {
      expect(mockApiClient.updateProject).toHaveBeenCalledWith(1, {
        name: 'Updated Project',
        description: 'Updated Description'
      });
    });
  });
}); 