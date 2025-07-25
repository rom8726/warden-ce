import React from 'react';
import { Box, useTheme } from '@mui/material';
import { useTheme as useAppTheme } from '../theme/ThemeContext';

interface LogoProps {
  size?: 'small' | 'medium' | 'large';
  variant?: 'full' | 'icon' | 'text';
  color?: 'primary' | 'inherit' | 'white' | 'black';
  showText?: boolean;
}

const sizeMap = {
  small: { width: 32, height: 32, fontSize: 16 },
  medium: { width: 48, height: 48, fontSize: 20 },
  large: { width: 58, height: 58, fontSize: 24 }
};

export const Logo: React.FC<LogoProps> = ({
  size = 'medium',
  variant = 'full',
  color = 'primary',
  showText = true
}) => {
  const theme = useTheme();
  useAppTheme();
  const { width, height, fontSize } = sizeMap[size];

  const getColor = () => {
    switch (color) {
      case 'primary':
        return theme.palette.primary.main;
      case 'white':
        return '#ffffff';
      case 'black':
        return '#000000';
      case 'inherit':
      default:
        return 'inherit';
    }
  };

  const logoColor = getColor();

  if (variant === 'icon') {
    return (
      <Box
        sx={{
          width,
          height,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          borderRadius: '50%',
          backgroundColor: logoColor,
          color: color === 'primary' ? '#ffffff' : theme.palette.background.paper,
          fontSize: fontSize * 0.6,
          fontWeight: 'bold',
          boxShadow: theme.shadows[2],
          transition: 'all 0.3s ease',
          '&:hover': {
            transform: 'scale(1.05)',
            boxShadow: theme.shadows[4]
          }
        }}
      >
        W
      </Box>
    );
  }

  if (variant === 'text') {
    return (
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          gap: 1,
          color: logoColor,
          fontSize,
          fontWeight: 'bold',
          fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
          letterSpacing: '0.5px',
          transition: 'all 0.3s ease',
          '&:hover': {
            transform: 'scale(1.02)'
          }
        }}
      >
        <Box
          component="span"
          sx={{
            width: height * 0.4,
            height: height * 0.4,
            borderRadius: '50%',
            backgroundColor: logoColor,
            color: color === 'primary' ? '#ffffff' : theme.palette.background.paper,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: fontSize * 0.4,
            fontWeight: 'bold'
          }}
        >
          W
        </Box>
        {showText && (
          <Box component="span" sx={{ fontSize: fontSize * 0.8 }}>
            Warden
          </Box>
        )}
      </Box>
    );
  }

  // Full variant (default)
  return (
    <Box
      sx={{
        display: 'flex',
        alignItems: 'center',
        gap: 1.5,
        color: logoColor,
        fontSize,
        fontWeight: 'bold',
        fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
        letterSpacing: '0.5px',
        transition: 'all 0.3s ease',
        '&:hover': {
          transform: 'scale(1.02)'
        }
      }}
    >
      <Box
        component="span"
        sx={{
          width: height * 0.5,
          height: height * 0.5,
          borderRadius: '50%',
          backgroundColor: logoColor,
          color: color === 'primary' ? '#ffffff' : theme.palette.background.paper,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: fontSize * 0.5,
          fontWeight: 'bold',
          boxShadow: theme.shadows[1]
        }}
      >
        W
      </Box>
      {showText && (
        <Box component="span">
          Warden
        </Box>
      )}
    </Box>
  );
};

export default Logo;
