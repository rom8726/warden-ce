import React, { memo, useState, useEffect } from 'react';
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableContainer, 
  TableHead, 
  TableRow, 
  Chip, 
  Alert,
  Dialog,
  DialogTitle,
  DialogContent,
  IconButton,
  Typography,
  Grid,
  Box
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import { type Issue, type IssueResponse, type IssueEvent, IssueLevel } from '../../generated/api/client';
import { getLevelColor, formatDate } from '../../utils/issues/issueUtils';
import EventTags from './EventTags';

interface EventsTabProps {
  issueData: IssueResponse | null;
  issue: Issue;
}

// Component to display event details in a modal
const EventDetailsModal = ({ event, open, onClose }: { event: IssueEvent | null, open: boolean, onClose: () => void }) => {
  // State to track expanded fields
  const [expandedFields, setExpandedFields] = useState<Record<string, boolean>>({});
  // State to track if content is loaded
  const [contentLoaded, setContentLoaded] = useState(false);

  // Load content when modal opens
  useEffect(() => {
    if (open && !contentLoaded) {
      // Delay content loading slightly to ensure modal animation starts first
      const timer = setTimeout(() => {
        setContentLoaded(true);
      }, 100);
      return () => clearTimeout(timer);
    } else if (!open) {
      // Reset content loaded state when modal closes
      setContentLoaded(false);
      // Also reset expanded fields when modal closes
      setExpandedFields({});
    }
  }, [open, contentLoaded]);

  if (!event) return null;

  // Toggle expanded state for a field
  const toggleExpanded = (fieldName: string) => {
    setExpandedFields(prev => ({
      ...prev,
      [fieldName]: !prev[fieldName]
    }));
  };

  // Helper function to safely stringify objects with circular references
  const safeStringify = (obj: any, maxLength: number = 5000) => {
    try {
      const seen = new WeakSet();
      const result = JSON.stringify(obj, (_key, value) => {
        // Handle DOM elements and other non-serializable objects
        if (typeof value === 'object' && value !== null) {
          // Check for circular references
          if (seen.has(value)) {
            return '[Circular Reference]';
          }
          seen.add(value);

          // Handle DOM elements specifically
          if (value instanceof Element || value instanceof Node) {
            return '[DOM Element]';
          }
        }
        return value;
      }, 2);

      // Truncate if too long
      if (result.length > maxLength) {
        return result.substring(0, maxLength) + '... [truncated]';
      }

      return result;
    } catch (error) {
      return '[Error: Unable to stringify object]';
    }
  };

  // Helper function to render a field with label
  const renderField = (label: string, value: any) => {
    // Skip rendering if value is null or undefined
    if (value === null || value === undefined) return null;

    // Generate a field key for the expandedFields state
    const fieldKey = `field-${label}`;
    const isExpanded = expandedFields[fieldKey] || false;

    // Handle different types of values
    let displayValue: React.ReactNode = value;

    // Handle object types (like tags, request_headers)
    if (typeof value === 'object' && value !== null) {
      try {
        // Check if it's a large object
        const stringified = safeStringify(value);
        const isLargeObject = stringified.length > 1000;

        if (isLargeObject) {
          displayValue = (
            <Box>
              <Typography 
                variant="body2" 
                color="primary" 
                sx={{ cursor: 'pointer', textDecoration: 'underline' }}
                onClick={() => toggleExpanded(fieldKey)}
              >
                {isExpanded ? 'Collapse' : 'Expand'} {label} ({stringified.length} chars)
              </Typography>
              {isExpanded && (
                <Box sx={{ maxHeight: '200px', overflow: 'auto', mt: 1 }}>
                  <pre>{stringified}</pre>
                </Box>
              )}
            </Box>
          );
        } else {
          displayValue = (
            <Box sx={{ maxHeight: '200px', overflow: 'auto' }}>
              <pre>{stringified}</pre>
            </Box>
          );
        }
      } catch (error) {
        // Fallback if stringify still fails
        displayValue = '[Complex Object: Unable to display]';
      }
    }

    return (
      <Grid item xs={12} sm={6} md={4}>
        <Box sx={{ mb: 2 }}>
          <Typography variant="subtitle2" color="text.secondary">
            {label}
          </Typography>
          <Typography variant="body2">{displayValue}</Typography>
        </Box>
      </Grid>
    );
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="lg" fullWidth>
      <DialogTitle className="gradient-text-blue">
        Event Details
        <IconButton
          aria-label="close"
          onClick={onClose}
          sx={{ position: 'absolute', right: 8, top: 8 }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>
      <DialogContent dividers>
        {!contentLoaded ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <Typography>Loading event details...</Typography>
          </Box>
        ) : (
          <Grid container spacing={2}>
            {/* Required fields */}
            {renderField('Event ID', event.event_id)}
            {renderField('Timestamp', formatDate(event.timestamp))}
            {renderField('Project ID', event.project_id)}
            {renderField('Level', 
              // Convert level to string instead of passing a Chip component
              Object.entries(IssueLevel).find(([_, val]) => val === event.level)?.[0] || 
              (typeof event.level === 'string' ? event.level : String(event.level))
            )}
            {renderField('Source', event.source)}
            {renderField('Platform', event.platform)}
            {renderField('Message', event.message)}

            {/* Optional fields */}
            {renderField('Group Hash', event.group_hash)}
            {renderField('Server Name', event.server_name)}
            {renderField('Environment', event.environment)}
            {renderField('Release', event.release)}
            {renderField('Tags', event.tags)}

            {/* Optional and nullable fields */}
            {renderField('Exception Type', event.exception_type)}
            {renderField('Exception Value', event.exception_value)}
            {renderField('Request URL', event.request_url)}
            {renderField('Request Method', event.request_method)}
            {renderField('Request Query', event.request_query)}
            {renderField('Request Headers', event.request_headers)}
            {renderField('Request Data', event.request_data)}
            {renderField('Request Cookies', event.request_cookies)}
            {renderField('Request IP', event.request_ip)}
            {renderField('User Agent', event.user_agent)}
            {renderField('User ID', event.user_id)}
            {renderField('User Email', event.user_email)}
            {renderField('Runtime Name', event.runtime_name)}
            {renderField('Runtime Version', event.runtime_version)}
            {renderField('OS Name', event.os_name)}
            {renderField('OS Version', event.os_version)}
            {renderField('Browser Name', event.browser_name)}
            {renderField('Browser Version', event.browser_version)}
            {renderField('Device Architecture', event.device_arch)}
          </Grid>
        )}
      </DialogContent>
    </Dialog>
  );
};

const EventsTab = memo(({ issueData }: EventsTabProps) => {
  const [selectedEvent, setSelectedEvent] = useState<IssueEvent | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  const handleEventClick = (event: IssueEvent) => {
    // Use setTimeout to defer state updates and prevent UI blocking
    setTimeout(() => {
      setSelectedEvent(event);
      setModalOpen(true);
    }, 0);
  };

  const handleCloseModal = () => {
    setModalOpen(false);
    // Clear the selected event after modal is closed to free up memory
    setTimeout(() => {
      setSelectedEvent(null);
    }, 300); // Wait for modal close animation to complete
  };

  // If we have events from the API response
  if (issueData?.events && issueData.events.length > 0) {
    return (
      <>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Timestamp</TableCell>
                <TableCell>Event ID</TableCell>
                <TableCell>Level</TableCell>
                <TableCell>Platform</TableCell>
                <TableCell>Environment</TableCell>
                <TableCell>Server Name</TableCell>
                <TableCell>Tags</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {issueData.events.map((event) => (
                <TableRow 
                  key={event.event_id || `${event.timestamp}-${Math.random()}`}
                  onClick={() => handleEventClick(event)}
                  hover
                  sx={{ cursor: 'pointer' }}
                >
                  <TableCell>{formatDate(event.timestamp)}</TableCell>
                  <TableCell>{event.event_id}</TableCell>
                  <TableCell>
                    <Chip 
                      label={
                        // Find the enum key by value
                        Object.entries(IssueLevel).find(([_, val]) => val === event.level)?.[0] || 
                        (typeof event.level === 'string' ? event.level : JSON.stringify(event.level).substring(0, 20) + '...')
                      } 
                      color={getLevelColor(event.level)} 
                      size="small"
                    />
                  </TableCell>
                  <TableCell>{event.platform}</TableCell>
                  <TableCell>{event.environment}</TableCell>
                  <TableCell>{event.server_name}</TableCell>
                  <TableCell>
                    <EventTags tags={event.tags} />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        {/* Modal for displaying event details */}
        <EventDetailsModal 
          event={selectedEvent} 
          open={modalOpen} 
          onClose={handleCloseModal} 
        />
      </>
    );
  }

  // No events available
  return (
    <Alert severity="info" sx={{ m: 2 }}>
      No events available for this issue.
    </Alert>
  );
});

EventsTab.displayName = 'EventsTab';

export default EventsTab;
