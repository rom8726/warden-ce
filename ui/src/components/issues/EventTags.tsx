import { memo } from 'react';
import { Chip } from '@mui/material';

interface EventTagsProps {
  tags?: Record<string, string>;
}

const EventTags = memo(({ tags }: EventTagsProps) => {
  if (!tags || Object.keys(tags).length === 0) return null;

  return (
    <>
      {Object.entries(tags).map(([key, value]) => (
        <Chip 
          key={`${key}-${value}`}
          label={`${key}: ${value}`} 
          size="small"
          variant="outlined"
          sx={{ mr: 0.5, mb: 0.5 }}
        />
      ))}
    </>
  );
});

EventTags.displayName = 'EventTags';

export default EventTags;