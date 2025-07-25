import React from 'react';
import { Box, useTheme } from '@mui/material';
import { useTheme as useAppTheme } from '../theme/ThemeContext';

interface LogoImgProps {
  size?: 'small' | 'medium' | 'large';
  variant?: 'full' | 'icon' | 'gray';
  showText?: boolean;
  sx?: any;
}

const sizeMap = {
  small: { height: 32 },
  medium: { height: 48 },
  large: { height: 58 }
};

const getLogoPath = (variant: string, mode: string) => {
  if (variant === 'icon') {
    return '/logo_icon.png';
  }
  if (mode == 'light') {
    return '/logo_full_gray.png';
  }
  // Full variant - use colored version
  return '/logo_full.png';
};

export const LogoImg: React.FC<LogoImgProps> = ({
  size = 'medium',
  variant = 'full',
  showText = true,
  sx
}) => {
  const theme = useTheme();
  const { mode } = useAppTheme();
  const {  height } = sizeMap[size];

  const logoPath = getLogoPath(variant, mode);

  if (variant === 'icon') {
    return (
      <Box
        component="img"
        src={logoPath}
        alt="Warden Logo"
        sx={{
          height: height,
          width: 'auto',
          objectFit: 'contain',
          transition: 'all 0.3s ease',
          '&:hover': {
            transform: 'scale(1.05)',
            filter: 'drop-shadow(0 2px 4px rgba(0,0,0,0.1))'
          },
          ...sx
        }}
      />
    );
  }

  return (
    <Box
      component="img"
      src={logoPath}
      alt="Warden Logo"
      sx={{
        height,
        width: 'auto',
        objectFit: 'contain',
        transition: 'all 0.3s ease',
        '&:hover': {
          transform: 'scale(1.02)',
          filter: 'drop-shadow(0 2px 4px rgba(0,0,0,0.1))'
        },
        ...sx
      }}
    />
  );
};

export default LogoImg;