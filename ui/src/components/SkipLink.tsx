import React from 'react';
import { Box, Link } from '@mui/material';

interface SkipLinkProps {
  href: string;
  children: React.ReactNode;
}

const SkipLink: React.FC<SkipLinkProps> = ({ href, children }) => {
  return (
    <Box
      component={Link}
      href={href}
      className="skip-link"
      sx={{
        position: 'absolute',
        top: '-40px',
        left: '6px',
        background: 'primary.main',
        color: 'primary.contrastText',
        padding: 1,
        textDecoration: 'none',
        borderRadius: 0.5,
        zIndex: 10000,
        transition: 'top 0.3s',
        '&:focus': {
          top: '6px',
        },
        '&:hover': {
          textDecoration: 'none',
        },
      }}
    >
      {children}
    </Box>
  );
};

export default SkipLink; 