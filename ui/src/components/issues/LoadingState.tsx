import { memo } from 'react';
import { Box, CircularProgress } from '@mui/material';
import { useIssuePageStyles } from './IssuePageStyles';

const LoadingState = memo(() => {
  const styles = useIssuePageStyles();

  return (
    <Box sx={styles.loadingContainer}>
      <CircularProgress />
    </Box>
  );
});

LoadingState.displayName = 'LoadingState';

export default LoadingState;