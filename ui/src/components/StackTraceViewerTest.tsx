import React from 'react';
import { Box, Typography, Button } from '@mui/material';
import StackTraceViewer from './StackTraceViewer';

// Sample stack trace data based on the issue description
const sampleStackTrace = JSON.stringify({
  frames: [
    {
      abs_path: "/Users/roman/projects/rom8726/sentry-sender/main.go",
      context_line: "\tsentry.CaptureException(fmt.Errorf(\"This is test error2\"))",
      function: "main",
      in_app: true,
      lineno: 27,
      module: "main",
      post_context: [
        "\tsentry.CaptureMessage(\"This is test message\")",
        "\t//",
        "\terr = someFunction()",
        "\tif err != nil {",
        "\t\tsentry.CaptureException(err)"
      ],
      pre_context: [
        "",
        "\tsentry.CaptureEvent(&sentry.Event{",
        "\t\tMessage: \"This is test event\",",
        "\t})",
        "\tsentry.CaptureException(fmt.Errorf(\"This is test error\"))"
      ]
    }
  ]
});

// Sample stack trace as a string (for fallback testing)
const sampleStackTraceString = `TypeError: Cannot read property 'data' of undefined
    at processResponse (/app/src/utils/api.js:25:10)
    at async Function.handleResponse (/app/src/utils/api.js:40:21)
    at async fetchData (/app/src/services/dataService.js:15:16)
    at async loadUserData (/app/src/components/UserProfile.js:32:18)
    at async UserProfile (/app/src/components/UserProfile.js:58:5)`;

const StackTraceViewerTest: React.FC = () => {
  const [useJsonFormat, setUseJsonFormat] = React.useState(true);

  const toggleFormat = () => {
    setUseJsonFormat(!useJsonFormat);
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Stack Trace Viewer Test
      </Typography>
      
      <Button 
        variant="contained" 
        onClick={toggleFormat} 
        sx={{ mb: 2 }}
      >
        Toggle Format: {useJsonFormat ? 'JSON' : 'String'}
      </Button>
      
      <Typography variant="h6" gutterBottom>
        {useJsonFormat ? 'JSON Format (Structured)' : 'String Format (Fallback)'}
      </Typography>
      
      <StackTraceViewer 
        stacktrace={useJsonFormat ? sampleStackTrace : sampleStackTraceString} 
        platform={useJsonFormat ? 'go' : 'javascript'} 
      />
    </Box>
  );
};

export default StackTraceViewerTest;