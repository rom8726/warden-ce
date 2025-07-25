import React from 'react';
import { Box, Typography } from '@mui/material';

interface DividerProps {
  text?: string;
}

const Divider: React.FC<DividerProps> = ({ text = 'or' }) => {
  return (
    <Box sx={{ display: 'flex', alignItems: 'center', my: 2 }}>
      <Box sx={{ flex: 1, height: 1, backgroundColor: 'divider' }} />
      <Typography
        variant="body2"
        sx={{
          px: 2,
          color: 'text.secondary',
          textTransform: 'uppercase',
          fontSize: '0.75rem',
          fontWeight: 500,
        }}
      >
        {text}
      </Typography>
      <Box sx={{ flex: 1, height: 1, backgroundColor: 'divider' }} />
    </Box>
  );
};

export default Divider; 