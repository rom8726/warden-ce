import React, { useState, useMemo, memo, useCallback } from 'react';
import { 
  Box, 
  Typography, 
  useTheme,
  Paper,
  Collapse,
} from '@mui/material';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import * as prismStyles from 'react-syntax-highlighter/dist/cjs/styles/prism';

// Define the structure of a stack trace frame based on the issue description
interface StackTraceFrame {
  abs_path: string;
  context_line: string;
  function: string;
  in_app: boolean;
  lineno: number;
  module: string;
  post_context: string[];
  pre_context: string[];
}

interface StackTrace {
  frames: StackTraceFrame[];
}

interface StackTraceViewerProps {
  stacktrace: string;
  platform?: string;
}

// Extracted reusable styles
const useStackTraceStyles = () => {
  const theme = useTheme();

  return {
    container: {
      fontFamily: 'monospace',
      fontSize: '0.875rem',
      lineHeight: 1.5,
      overflow: 'auto',
      maxHeight: '500px',
      borderRadius: 1,
      border: 1,
      borderColor: 'divider',
      boxShadow: theme.palette.mode === 'dark' 
        ? '0 4px 20px rgba(0, 0, 0, 0.2)' 
        : '0 4px 20px rgba(0, 0, 0, 0.05)',
    },
    lineNumber: {
      minWidth: '2.5rem', 
      textAlign: 'right', 
      color: theme.palette.text.secondary,
      opacity: 0.7,
      pr: 1,
      borderRight: '1px solid',
      borderColor: 'divider',
      mr: 1.5,
      userSelect: 'none'
    },
    expandedLineNumber: {
      minWidth: '3rem', 
      textAlign: 'right', 
      color: theme.palette.text.secondary,
      pr: 1,
      borderRight: '1px solid',
      borderColor: 'divider',
      mr: 1.5,
      userSelect: 'none'
    },
    errorLineNumber: {
      minWidth: '3rem', 
      textAlign: 'right', 
      color: theme.palette.error.main,
      pr: 1,
      borderRight: '1px solid',
      borderColor: 'divider',
      mr: 1.5,
      userSelect: 'none'
    },
    frameLine: {
      p: 0.75,
      pl: 0,
      display: 'flex',
      backgroundColor: 'background.default',
      borderBottom: '1px solid',
      borderColor: 'divider',
    },
    clickableFrame: {
      cursor: 'pointer',
      '&:hover': {
        backgroundColor: theme.palette.mode === 'dark' 
          ? 'rgba(255, 255, 255, 0.05)' 
          : 'rgba(0, 0, 0, 0.04)',
      },
    },
    contextLine: {
      p: 0.5,
      pl: 0,
      display: 'flex',
      backgroundColor: theme.palette.mode === 'dark' 
        ? 'rgba(255, 0, 0, 0.1)' 
        : 'rgba(255, 0, 0, 0.05)',
      borderBottom: '1px solid',
      borderColor: 'divider',
      fontWeight: 'bold'
    },
    contextContainer: {
      m: 1, 
      backgroundColor: theme.palette.mode === 'dark' 
        ? 'rgba(0, 0, 0, 0.2)' 
        : 'rgba(0, 0, 0, 0.03)',
      borderRadius: 1,
      overflow: 'hidden'
    },
    primaryText: {
      color: theme.palette.primary.main,
      fontWeight: 500
    },
    secondaryText: {
      color: theme.palette.text.secondary,
      ml: 1
    },
    errorText: {
      color: theme.palette.error.main
    }
  };
};

// Line number component
const LineNumber = memo(({ number, isError = false, expanded = false }: { 
  number: number; 
  isError?: boolean;
  expanded?: boolean;
}) => {
  const styles = useStackTraceStyles();
  const style = isError 
    ? styles.errorLineNumber 
    : expanded 
      ? styles.expandedLineNumber 
      : styles.lineNumber;

  return (
    <Box sx={style}>
      {number}
    </Box>
  );
});

LineNumber.displayName = 'LineNumber';

// Context line component
const ContextLine = memo(({ 
  lineNumber, 
  content, 
  isError = false,
  language = 'markup',
  codeStyle = prismStyles.vs
}: {
  lineNumber: number;
  content: string;
  isError?: boolean;
  language?: string;
  codeStyle?: any;
}) => {
  const styles = useStackTraceStyles();
  const theme = useTheme();

  return (
    <Box 
      sx={{
        ...styles.frameLine,
        p: 0,
        pl: 0,
        opacity: isError ? 1 : 0.7,
        ...(isError && styles.contextLine)
      }}
    >
      <LineNumber number={lineNumber} isError={isError} expanded />
      <Box sx={{ flex: 1, overflow: 'hidden' }}>
        <SyntaxHighlighter
          language={language}
          style={codeStyle}
          customStyle={{
            margin: 0,
            padding: '0.5rem',
            background: 'transparent',
            border: 'none',
            ...(isError && {
              color: theme.palette.error.main,
              fontWeight: 'bold'
            })
          }}
        >
          {content}
        </SyntaxHighlighter>
      </Box>
    </Box>
  );
});

ContextLine.displayName = 'ContextLine';

// Helper function to determine the appropriate language based on platform
const getPrismLanguage = (platform?: string): string => {
  if (!platform) return 'markup'; // Default fallback

  // Map platform to Prism language identifier
  switch (platform.toLowerCase()) {
    case 'javascript':
    case 'js':
      return 'javascript';
    case 'typescript':
    case 'ts':
      return 'typescript';
    case 'python':
    case 'py':
      return 'python';
    case 'java':
      return 'java';
    case 'csharp':
    case 'c#':
      return 'csharp';
    case 'php':
      return 'php';
    case 'ruby':
      return 'ruby';
    case 'go':
      return 'go';
    case 'rust':
      return 'rust';
    case 'kotlin':
      return 'kotlin';
    case 'swift':
      return 'swift';
    default:
      return 'markup'; // Default fallback
  }
};

const StackTraceViewer: React.FC<StackTraceViewerProps> = ({ stacktrace, platform }) => {
  // const theme = useTheme();
  const styles = useStackTraceStyles();
  const [expandedFrame, setExpandedFrame] = useState<number | null>(null);

  // Determine the language for syntax highlighting
  const language = useMemo(() => getPrismLanguage(platform), [platform]);

  // Select appropriate style based on theme
  const theme = useTheme();
  const codeStyle = useMemo(() => 
    theme.palette.mode === 'dark' ? prismStyles.vscDarkPlus : prismStyles.vs, 
  [theme.palette.mode]);

  // Try to parse the stacktrace as JSON
  const parsedStackTrace = useMemo(() => {
    let result: StackTrace | null = null;
    try {
      // Check if the stacktrace is already a JSON object or a string that needs parsing
      if (typeof stacktrace === 'string') {
        result = JSON.parse(stacktrace) as StackTrace;
      } else {
        console.error('Stacktrace is not a string:', stacktrace);
      }
    } catch (e) {
      // If parsing fails, we'll fall back to the old display method
      console.error('Failed to parse stacktrace as JSON:', e);
    }

    // Validate that the parsed stacktrace has the expected structure
    if (result && (!result.frames || !Array.isArray(result.frames) || result.frames.length === 0)) {
      console.error('Parsed stacktrace does not have valid frames:', result);
      return null;
    }

    return result;
  }, [stacktrace]);

  const toggleFrame = useCallback((index: number) => {
    setExpandedFrame(prevIndex => prevIndex === index ? null : index);
  }, []);

  // Simple stack trace line component
  const SimpleStackTraceLine = memo(({ line, index, totalLines }: {
    line: string;
    index: number;
    totalLines: number;
  }) => {
    const styles = useStackTraceStyles();
    const theme = useTheme();
    const isErrorLine = index === 0;

    // Parse file paths in the line
    const filePathMatch = line.match(/(\s+at\s+.+\s+\()([^)]+)(\))/);

    return (
      <Box 
        key={index} 
        sx={{ 
          ...styles.frameLine,
          borderBottom: index < totalLines - 1 ? '1px solid' : 'none',
          whiteSpace: 'pre-wrap',
          wordBreak: 'break-all',
          ...(isErrorLine && {
            fontWeight: 'bold',
          }),
          padding: 0
        }}
      >
        <LineNumber number={index + 1} />
        <Box sx={{ flex: 1, overflow: 'hidden' }}>
          {filePathMatch ? (
            <Box sx={{ display: 'flex', flexWrap: 'wrap' }}>
              <Typography component="span">
                {filePathMatch[1]}
              </Typography>
              <Typography component="span" sx={styles.primaryText}>
                {filePathMatch[2]}
              </Typography>
              <Typography component="span">
                {filePathMatch[3]}
              </Typography>
            </Box>
          ) : (
            <SyntaxHighlighter
              language={language}
              style={codeStyle}
              customStyle={{
                margin: 0,
                padding: '0.75rem',
                background: 'transparent',
                border: 'none',
                ...(isErrorLine && {
                  color: theme.palette.error.main,
                })
              }}
            >
              {line}
            </SyntaxHighlighter>
          )}
        </Box>
      </Box>
    );
  });

  SimpleStackTraceLine.displayName = 'SimpleStackTraceLine';

  // Stack frame component
  const StackFrame = memo(({ frame, index, isExpanded, onToggle }: {
    frame: StackTraceFrame;
    index: number;
    isExpanded: boolean;
    onToggle: () => void;
  }) => {
    const styles = useStackTraceStyles();
    const startLineNumber = frame.lineno - (frame.pre_context?.length || 0);

    return (
      <Box>
        {/* Frame header - always visible */}
        <Box 
          onClick={onToggle}
          sx={{ 
            ...styles.frameLine,
            ...styles.clickableFrame
          }}
        >
          <LineNumber number={index + 1} />
          <Box sx={{ flex: 1 }}>
            <Typography component="span" sx={styles.primaryText}>
              {frame.abs_path}
            </Typography>
            <Typography component="span" sx={styles.secondaryText}>
              :{frame.lineno}
            </Typography>
            <Typography component="span" sx={styles.secondaryText}>
              in {frame.function}
            </Typography>
            {frame.module && (
              <Typography component="span" sx={styles.secondaryText}>
                ({frame.module})
              </Typography>
            )}
          </Box>
        </Box>

        {/* Expanded code view */}
        <Collapse in={isExpanded}>
          <Paper elevation={0} sx={styles.contextContainer}>
            {/* Pre-context lines */}
            {frame.pre_context?.map((line, lineIndex) => (
              <ContextLine 
                key={`pre-${lineIndex}`}
                lineNumber={startLineNumber + lineIndex}
                content={line}
                language={language}
                codeStyle={codeStyle}
              />
            ))}

            {/* Context line (the error line) */}
            <ContextLine 
              lineNumber={frame.lineno}
              content={frame.context_line}
              isError
              language={language}
              codeStyle={codeStyle}
            />

            {/* Post-context lines */}
            {frame.post_context?.map((line, lineIndex) => (
              <ContextLine 
                key={`post-${lineIndex}`}
                lineNumber={frame.lineno + lineIndex + 1}
                content={line}
                language={language}
                codeStyle={codeStyle}
              />
            ))}
          </Paper>
        </Collapse>
      </Box>
    );
  });

  StackFrame.displayName = 'StackFrame';

  // If we couldn't parse the stacktrace, fall back to the old display method
  if (!parsedStackTrace || !parsedStackTrace.frames) {
    const lines = stacktrace.split('\n');

    return (
      <Box sx={styles.container}>
        {lines.map((line, index) => (
          <SimpleStackTraceLine 
            key={index}
            line={line}
            index={index}
            totalLines={lines.length}
          />
        ))}
      </Box>
    );
  }

  // Render the structured stack trace
  return (
    <Box sx={styles.container}>
      {parsedStackTrace.frames.map((frame, frameIndex) => (
        <StackFrame
          key={frameIndex}
          frame={frame}
          index={frameIndex}
          isExpanded={expandedFrame === frameIndex}
          onToggle={() => toggleFrame(frameIndex)}
        />
      ))}
    </Box>
  );
};

export default StackTraceViewer;
