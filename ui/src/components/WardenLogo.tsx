import React from 'react';
import { Typography, Box } from '@mui/material';
import { Link } from 'react-router-dom';
import './WardenLogo.css';

interface WardenLogoProps {
  variant?: 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6';
  logoSize?: number;
  showLink?: boolean;
  className?: string;
  sx?: any;
}

const WardenLogo: React.FC<WardenLogoProps> = ({
                                                 variant = 'h4',
                                                 logoSize = 32,
                                                 showLink = true,
                                                 className = '',
                                                 sx = {}
                                               }) => {

  const logoContent = (
    <Typography
      variant={variant}
      className={`gradient-text ${className}`}
      sx={{
        textDecoration: 'none',
        fontWeight: 700,
        display: 'flex',
        alignItems: 'center',
        gap: 1,
        ...sx
      }}
    >
      <img
        src="/logo_icon.png"
        style={{
          width: logoSize,
          height: Math.round(logoSize * 1.22),
          display: 'block'
        }}
        alt="Warden Logo"
      />
      Warden
    </Typography>
  );

  if (showLink) {
    return (
      <Box component={Link} to="/" sx={{ textDecoration: 'none' }}>
        {logoContent}
      </Box>
    );
  }

  return logoContent;
};

export default WardenLogo;
