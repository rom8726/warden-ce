import React from 'react';
import {
  Alert,
  Box,
  Typography,
  CircularProgress,
  Button,
} from '@mui/material';
import {
  SwapHoriz as SwapIcon,
} from '@mui/icons-material';

interface ComparisonAlertProps {
  compareMode: boolean;
  loading: boolean;
  error: string | null;
  comparisonData: any;
  onSwitchComparison?: () => void;
}

const ComparisonAlert: React.FC<ComparisonAlertProps> = ({
  compareMode,
  loading,
  error,
  comparisonData,
  onSwitchComparison,
}) => {
  if (!compareMode) {
    return null;
  }

  if (loading) {
    return (
      <Alert severity="info" sx={{ mb: 2 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <CircularProgress size={16} />
          <Typography variant="body2">
            Loading comparison data...
          </Typography>
        </Box>
      </Alert>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        <Typography variant="body2">
          Error loading comparison data. Please try again.
        </Typography>
      </Alert>
    );
  }

  if (comparisonData) {
    return (
      <Alert severity="info" sx={{ mb: 2 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: 2 }}>
          <Typography variant="body2">
            Comparing releases: <strong>{comparisonData.base_version}</strong> → <strong>{comparisonData.target_version}</strong>
          </Typography>
          {onSwitchComparison && (
            <Button
              variant="outlined"
              size="small"
              startIcon={<SwapIcon />}
              onClick={onSwitchComparison}
              sx={{ minWidth: 'auto', whiteSpace: 'nowrap' }}
            >
              Switch Direction
            </Button>
          )}
        </Box>
      </Alert>
    );
  }

  return null;
};

export default ComparisonAlert; 