import { render, screen } from '@testing-library/react';
import TabPanel from '../TabPanel';

describe('TabPanel', () => {
  test('renders children when value matches index', () => {
    render(
      <TabPanel value={0} index={0}>
        <div data-testid="test-content">Test Content</div>
      </TabPanel>
    );

    expect(screen.getByTestId('test-content')).toBeInTheDocument();
    expect(screen.getByText('Test Content')).toBeInTheDocument();

    // Check accessibility attributes
    const tabpanel = screen.getByRole('tabpanel');
    expect(tabpanel).toHaveAttribute('id', 'issue-tabpanel-0');
    expect(tabpanel).toHaveAttribute('aria-labelledby', 'issue-tab-0');
    expect(tabpanel).not.toHaveAttribute('hidden');
  });

  test('does not render children when value does not match index', () => {
    render(
      <TabPanel value={1} index={0}>
        <div data-testid="test-content">Test Content</div>
      </TabPanel>
    );

    expect(screen.queryByTestId('test-content')).not.toBeInTheDocument();
    expect(screen.queryByText('Test Content')).not.toBeInTheDocument();

    // Check accessibility attributes
    const tabpanel = screen.getByRole('tabpanel', { hidden: true });
    expect(tabpanel).toHaveAttribute('id', 'issue-tabpanel-0');
    expect(tabpanel).toHaveAttribute('aria-labelledby', 'issue-tab-0');
    expect(tabpanel).toHaveAttribute('hidden');
  });

  test('passes additional props to the div', () => {
    render(
      <TabPanel value={0} index={0} data-custom="test-value">
        <div>Test Content</div>
      </TabPanel>
    );

    const tabpanel = screen.getByRole('tabpanel');
    expect(tabpanel).toHaveAttribute('data-custom', 'test-value');
  });
});
