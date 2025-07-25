import React, { Component, type ReactNode } from 'react';
import { 
  Box, 
  Typography, 
  Button, 
  Paper,
  useTheme
} from '@mui/material';
import { 
  Error as ErrorIcon,
  Refresh as RefreshIcon
} from '@mui/icons-material';

// Add Node.js types for process.env
declare global {
  namespace NodeJS {
    interface ProcessEnv {
      NODE_ENV: 'development' | 'production' | 'test';
    }
  }
}

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
  errorInfo?: React.ErrorInfo;
}

class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('ErrorBoundary caught an error:', error, errorInfo);
    this.setState({ error, errorInfo });
  }

  handleRetry = () => {
    this.setState({ hasError: false, error: undefined, errorInfo: undefined });
  };

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback;
      }

      return <ErrorFallback error={this.state.error} onRetry={this.handleRetry} />;
    }

    return this.props.children;
  }
}

interface ErrorFallbackProps {
  error?: Error;
  onRetry: () => void;
}

const ErrorFallback: React.FC<ErrorFallbackProps> = ({ error, onRetry }) => {
  const theme = useTheme();

  return (
    <Box sx={{ 
      display: 'flex', 
      alignItems: 'center', 
      justifyContent: 'center',
      minHeight: '400px',
      p: 3,
    }}>
      <Paper sx={{ 
        p: 4, 
        maxWidth: 500,
        textAlign: 'center',
        background: theme.palette.mode === 'dark' 
          ? 'linear-gradient(135deg, rgba(45, 48, 56, 0.8), rgba(35, 38, 46, 0.9))'
          : 'linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(248, 249, 252, 0.9))',
        backdropFilter: 'blur(10px)',
        boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
      }}>
        <Box sx={{ 
          display: 'flex', 
          alignItems: 'center', 
          justifyContent: 'center',
          width: 80,
          height: 80,
          borderRadius: '50%',
          backgroundColor: theme.palette.mode === 'dark' 
            ? 'rgba(244, 67, 54, 0.1)' 
            : 'rgba(244, 67, 54, 0.08)',
          color: theme.palette.error.main,
          mb: 3,
          mx: 'auto',
        }}>
          <ErrorIcon sx={{ fontSize: 40 }} />
        </Box>
        
        <Typography 
          variant="h5" 
          component="h2" 
          gutterBottom
          sx={{ 
            fontWeight: 600,
            color: theme.palette.text.primary,
            mb: 2,
          }}
        >
          Something went wrong
        </Typography>
        
        <Typography 
          variant="body1" 
          color="text.secondary"
          sx={{ 
            mb: 3,
            lineHeight: 1.6,
          }}
        >
          We encountered an unexpected error. Please try refreshing the page or contact support if the problem persists.
        </Typography>

        {error && process.env.NODE_ENV === 'development' && (
          <Box sx={{ 
            mb: 3, 
            p: 2, 
            backgroundColor: theme.palette.mode === 'dark' 
              ? 'rgba(0, 0, 0, 0.3)' 
              : 'rgba(0, 0, 0, 0.05)',
            borderRadius: 1,
            textAlign: 'left',
          }}>
            <Typography variant="caption" color="text.secondary" sx={{ fontWeight: 600 }}>
              Error Details (Development):
            </Typography>
            <Typography variant="caption" component="pre" sx={{ 
              display: 'block', 
              mt: 1,
              color: theme.palette.error.main,
              fontSize: '0.75rem',
              whiteSpace: 'pre-wrap',
              wordBreak: 'break-word',
            }}>
              {error.message}
            </Typography>
          </Box>
        )}
        
        <Button
          variant="contained"
          onClick={onRetry}
          startIcon={<RefreshIcon />}
          sx={{
            px: 3,
            py: 1.5,
            borderRadius: 2,
            textTransform: 'none',
            fontWeight: 500,
          }}
        >
          Try Again
        </Button>
      </Paper>
    </Box>
  );
};

export default ErrorBoundary; 