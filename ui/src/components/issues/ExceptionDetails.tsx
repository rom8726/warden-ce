import { memo } from 'react';
import { Box } from '@mui/material';
import StackTraceViewer from '../StackTraceViewer';
import { useIssuePageStyles } from './IssuePageStyles';

interface ExceptionDetailsProps {
  exception: {
    exception_type: string;
    exception_value: string;
    stacktrace: string;
  };
  platform?: string;
}

const ExceptionDetails = memo(({ exception, platform }: ExceptionDetailsProps) => {
  const styles = useIssuePageStyles();

  if (!exception.stacktrace) {
    return null;
  }

  return (
    <Box>
      <Box sx={styles.exceptionHeader}>
        {exception.exception_type}: {exception.exception_value}
      </Box>
      <StackTraceViewer 
        stacktrace={exception.stacktrace} 
        platform={platform} 
      />
    </Box>
  );
});

ExceptionDetails.displayName = 'ExceptionDetails';

export default ExceptionDetails;