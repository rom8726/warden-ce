import { keyframes } from '@emotion/react';
import { useTheme } from '@mui/material';

// Define the grow-from-bottom animation
const growFromBottom = keyframes`
  0% {
    transform: scaleY(0);
    transform-origin: bottom;
  }
  100% {
    transform: scaleY(1);
    transform-origin: bottom;
  }
`;

// Styles for the issue page
export const useIssuePageStyles = () => {
  const theme = useTheme();

  return {
    // Header styles
    headerContainer: {
      mb: 2
    },
    headerTitleContainer: {
      display: 'flex', 
      alignItems: 'center', 
      mb: 0.5
    },
    headerTitle: {
      ml: 0.5
    },
    headerMessage: {
      mb: 1.5
    },
    tagsContainer: {
      display: 'flex', 
      flexWrap: 'wrap', 
      gap: 0.5
    },

    // Stats styles
    statsContainer: {
      mb: 2
    },
    statsPaper: {
      p: 1.5, 
      height: '100%'
    },
    statsIconContainer: {
      display: 'flex', 
      alignItems: 'center', 
      mb: 0.5
    },
    statsIcon: {
      mr: 0.5
    },
    statsValue: {
      fontWeight: 'bold', 
      mb: 0.5
    },

    // Chart styles
    chartContainer: {
      mb: 2
    },
    chartHeader: {
      display: 'flex', 
      justifyContent: 'space-between', 
      alignItems: 'center', 
      mb: 1.5,
      p: 2
    },
    chartBox: {
      height: 250, 
      width: '100%',
      position: 'relative',
      borderRadius: 1,
      overflow: 'hidden',
      boxShadow: `0 4px 20px 0 ${theme.palette.mode === 'dark' ? 'rgba(0,0,0,0.3)' : 'rgba(0,0,0,0.1)'}`,
      animation: `${growFromBottom} 1200ms ease-in-out`,
      '&:hover': {
        '& .error-line': {
          strokeWidth: 3,
        },
        '& .warning-line': {
          strokeWidth: 3,
        },
        '& .info-line': {
          strokeWidth: 3,
        },
        '& .exception-line': {
          strokeWidth: 3,
        },
        '& .error-dot': {
          r: 4,
        },
        '& .warning-dot': {
          r: 4,
        },
        '& .info-dot': {
          r: 4,
        },
        '& .exception-dot': {
          r: 4,
        }
      }
    },

    // Loading and error styles
    loadingContainer: {
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center', 
      height: '40vh'
    },
    errorAlert: {
      mb: 2
    },

    // Tab content styles
    tabContent: {
      p: 1.5, 
      bgcolor: 'background.default', 
      borderRadius: 1
    },

    // Exception styles
    exceptionHeader: {
      p: 1,
      backgroundColor: theme.palette.mode === 'dark' 
        ? 'rgba(255, 0, 0, 0.1)' 
        : 'rgba(255, 0, 0, 0.05)',
      borderRadius: '4px 4px 0 0',
      fontWeight: 'bold',
      color: theme.palette.error.main,
      mb: 1.5
    },

    // Status change loading
    statusChangeLoadingOverlay: {
      position: 'absolute',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      backgroundColor: 'rgba(0, 0, 0, 0.3)',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      zIndex: 1000,
      borderRadius: 1
    }
  };
};