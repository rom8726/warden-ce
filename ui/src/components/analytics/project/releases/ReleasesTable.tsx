import React from 'react';
import {
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  IconButton,
  Tooltip,
  Box,
  CircularProgress,
} from '@mui/material';
import {
  Assessment as AssessmentIcon,
  Visibility as ViewDetailsIcon,
  CompareArrows as CompareIcon,
} from '@mui/icons-material';

interface ReleaseSummary {
  version: string;
  known_issues_total: number;
  new_issues_total: number;
  regressions_total: number;
  resolved_in_version_total: number;
  users_affected: number;
  created_at: string;
}

interface ReleasesTableProps {
  releases: ReleaseSummary[];
  loading: boolean;
  error: string | null;
  selectedRelease: string | null;
  compareMode: boolean;
  comparisonData: any;
  onReleaseSelect: (version: string) => void;
  onCompareRelease: (releaseVersion: string) => void;
}

const ReleasesTable: React.FC<ReleasesTableProps> = ({
  releases,
  loading,
  error,
  selectedRelease,
  compareMode,
  comparisonData,
  onReleaseSelect,
  onCompareRelease,
}) => {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU');
  };

  if (loading) {
    return (
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
          <AssessmentIcon sx={{ mr: 1 }} />
          Release Summary
        </Typography>
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      </Paper>
    );
  }

  if (error) {
    return (
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
          <AssessmentIcon sx={{ mr: 1 }} />
          Release Summary
        </Typography>
        <Typography color="error">
          Error loading release data. Please try again.
        </Typography>
      </Paper>
    );
  }

  if (!releases || releases.length === 0) {
    return (
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
          <AssessmentIcon sx={{ mr: 1 }} />
          Release Summary
        </Typography>
        <Typography color="text.secondary">
          No releases found.
        </Typography>
      </Paper>
    );
  }

  return (
    <Paper sx={{ p: 3, mb: 3 }}>
      <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
        <AssessmentIcon sx={{ mr: 1 }} />
        Release Summary
      </Typography>
      
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Version</TableCell>
              <TableCell align="right">Known</TableCell>
              <TableCell align="right">New</TableCell>
              <TableCell align="right">Regress</TableCell>
              <TableCell align="right">Fixed</TableCell>
              <TableCell align="right">Users</TableCell>
              <TableCell>GeneratedAt</TableCell>
              {compareMode && comparisonData && (
                <TableCell align="center">Delta</TableCell>
              )}
              <TableCell align="center">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {releases.map((release) => {
              // Find delta for this release if in compare mode
              let deltaInfo = null;
              if (compareMode && comparisonData) {
                if (release.version === comparisonData.base.version) {
                  deltaInfo = { type: 'base', data: comparisonData.base };
                } else if (release.version === comparisonData.target.version) {
                  deltaInfo = { type: 'target', data: comparisonData.target };
                }
              }

              return (
                <TableRow 
                  key={release.version}
                  hover
                  selected={selectedRelease === release.version}
                  onClick={() => onReleaseSelect(release.version)}
                  sx={{ cursor: 'pointer' }}
                >
                  <TableCell component="th" scope="row">
                    <Typography variant="body2" fontWeight="medium">
                      {release.version}
                    </Typography>
                  </TableCell>
                  <TableCell align="right">{release.known_issues_total}</TableCell>
                  <TableCell align="right">{release.new_issues_total}</TableCell>
                  <TableCell align="right">{release.regressions_total}</TableCell>
                  <TableCell align="right">{release.resolved_in_version_total}</TableCell>
                  <TableCell align="right">{release.users_affected}</TableCell>
                  <TableCell>{formatDate(release.created_at)}</TableCell>
                  {compareMode && comparisonData && (
                    <TableCell align="center">
                      {deltaInfo ? (
                        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.5 }}>
                          {deltaInfo.type === 'base' && (
                            <Chip 
                              label="Base" 
                              size="small" 
                              color="primary" 
                              variant="outlined"
                            />
                          )}
                          {deltaInfo.type === 'target' && (
                            <Chip 
                              label="Target" 
                              size="small" 
                              color="secondary" 
                              variant="outlined"
                            />
                          )}
                        </Box>
                      ) : (
                        <Typography variant="caption" color="text.secondary">
                          -
                        </Typography>
                      )}
                    </TableCell>
                  )}
                  <TableCell align="center">
                    <Box sx={{ display: 'flex', gap: 1, justifyContent: 'center' }}>
                      <Tooltip title="View Details">
                        <IconButton 
                          size="small"
                          onClick={(e) => {
                            e.stopPropagation();
                            onReleaseSelect(release.version);
                          }}
                        >
                          <ViewDetailsIcon />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title={compareMode && selectedRelease === release.version ? "Exit Compare Mode" : "Compare to Previous"}>
                        <IconButton 
                          size="small"
                          color={compareMode && selectedRelease === release.version ? "primary" : "default"}
                          onClick={(e) => {
                            e.stopPropagation();
                            onCompareRelease(release.version);
                          }}
                        >
                          <CompareIcon />
                        </IconButton>
                      </Tooltip>
                    </Box>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>
    </Paper>
  );
};

export default ReleasesTable; 