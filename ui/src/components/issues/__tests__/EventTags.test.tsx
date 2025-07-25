import { render, screen } from '@testing-library/react';
import EventTags from '../EventTags';

describe('EventTags', () => {
  test('renders nothing when tags is undefined', () => {
    const { container } = render(<EventTags />);
    expect(container.firstChild).toBeNull();
  });

  test('renders nothing when tags is an empty object', () => {
    const { container } = render(<EventTags tags={{}} />);
    expect(container.firstChild).toBeNull();
  });

  test('renders chips for each tag', () => {
    const tags = {
      'tag1': 'value1',
      'tag2': 'value2',
      'tag3': 'value3'
    };
    
    render(<EventTags tags={tags} />);
    
    expect(screen.getByText('tag1: value1')).toBeInTheDocument();
    expect(screen.getByText('tag2: value2')).toBeInTheDocument();
    expect(screen.getByText('tag3: value3')).toBeInTheDocument();
  });

  test('renders a single tag correctly', () => {
    const tags = {
      'singleTag': 'singleValue'
    };
    
    render(<EventTags tags={tags} />);
    
    expect(screen.getByText('singleTag: singleValue')).toBeInTheDocument();
  });
});