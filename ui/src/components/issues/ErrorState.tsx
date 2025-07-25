import { memo } from 'react';
import { Alert } from '@mui/material';
import { useIssuePageStyles } from './IssuePageStyles';

interface ErrorStateProps {
  message: string;
}

const ErrorState = memo(({ message }: ErrorStateProps) => {
  const styles = useIssuePageStyles();

  return (
    <Alert severity="error" sx={styles.errorAlert}>
      {message}
    </Alert>
  );
});

ErrorState.displayName = 'ErrorState';

export default ErrorState;