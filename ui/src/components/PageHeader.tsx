import React from 'react';
import { 
  Box, 
  Typography, 
  Chip,
  useTheme,
  Skeleton
} from '@mui/material';

interface PageHeaderProps {
  title: string;
  subtitle?: string;
  icon?: React.ReactNode;
  badge?: {
    label: string;
    color: 'primary' | 'secondary' | 'error' | 'warning' | 'info' | 'success';
    variant?: 'filled' | 'outlined';
  };
  loading?: boolean;
  children?: React.ReactNode;
}

const PageHeader: React.FC<PageHeaderProps> = ({ 
  title, 
  subtitle, 
  icon, 
  badge, 
  loading = false,
  children 
}) => {
  const theme = useTheme();

  if (loading) {
    return (
      <Box sx={{ mb: 4, width: '100%' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1 }}>
          <Skeleton variant="circular" width={32} height={32} />
          <Skeleton variant="text" width={200} height={32} />
          {badge && <Skeleton variant="rectangular" width={80} height={24} sx={{ borderRadius: 1 }} />}
        </Box>
        {subtitle && <Skeleton variant="text" width={300} height={20} />}
      </Box>
    );
  }

  return (
    <Box sx={{ 
      mb: 4, 
      width: '100%',
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'flex-start',
      flexWrap: 'wrap',
      gap: 2
    }}>
      <Box sx={{ flex: 1, minWidth: 0 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1 }}>
          {icon && (
            <Box sx={{ 
              display: 'flex', 
              alignItems: 'center', 
              justifyContent: 'center',
              width: 40,
              height: 40,
              borderRadius: 2,
              backgroundColor: theme.palette.mode === 'dark' 
                ? 'rgba(130, 82, 255, 0.1)' 
                : 'rgba(130, 82, 255, 0.08)',
              color: theme.palette.primary.main,
            }}>
              {icon}
            </Box>
          )}
          <Typography 
            variant="h4" 
            component="h1" 
            sx={{ 
              fontWeight: 700,
              color: theme.palette.text.primary,
              lineHeight: 1.2,
            }}
          >
            {title}
          </Typography>
          {badge && (
            <Chip 
              label={badge.label} 
              color={badge.color} 
              variant={badge.variant || 'filled'}
              size="small"
              sx={{ 
                height: 24,
                fontSize: '0.75rem',
                fontWeight: 500,
              }}
            />
          )}
        </Box>
        {subtitle && (
          <Typography 
            variant="body1" 
            color="text.secondary"
            sx={{ 
              ml: icon ? 6 : 0,
              fontSize: '1rem',
              lineHeight: 1.5,
            }}
          >
            {subtitle}
          </Typography>
        )}
      </Box>
      {children && (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          {children}
        </Box>
      )}
    </Box>
  );
};

export default PageHeader; 