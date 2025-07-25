import React, { useState, useCallback } from 'react';
import { 
  TextField, 
  InputAdornment, 
  IconButton, 
  Box,
  useTheme,
  Tooltip
} from '@mui/material';
import { 
  Search as SearchIcon,
  Clear as ClearIcon
} from '@mui/icons-material';

interface SearchBarProps {
  placeholder?: string;
  value: string;
  onChange: (value: string) => void;
  onSearch?: (value: string) => void;
  disabled?: boolean;
  fullWidth?: boolean;
  size?: 'small' | 'medium';
  debounceMs?: number;
}

const SearchBar: React.FC<SearchBarProps> = ({
  placeholder = 'Search...',
  value,
  onChange,
  onSearch,
  disabled = false,
  fullWidth = true,
  size = 'medium',
  debounceMs = 300
}) => {
  const theme = useTheme();
  const [debounceTimer, setDebounceTimer] = useState<number | null>(null);

  const handleChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = event.target.value;
    onChange(newValue);

    // Debounce search if onSearch is provided
    if (onSearch) {
      if (debounceTimer) {
        clearTimeout(debounceTimer);
      }
      
      const timer = window.setTimeout(() => {
        onSearch(newValue);
      }, debounceMs);
      
      setDebounceTimer(timer);
    }
  }, [onChange, onSearch, debounceMs, debounceTimer]);

  const handleClear = useCallback(() => {
    onChange('');
    if (onSearch) {
      onSearch('');
    }
  }, [onChange, onSearch]);

  const handleKeyPress = useCallback((event: React.KeyboardEvent) => {
    if (event.key === 'Enter' && onSearch) {
      onSearch(value);
    }
  }, [onSearch, value]);

  return (
    <Box sx={{ width: fullWidth ? '100%' : 'auto' }}>
      <TextField
        placeholder={placeholder}
        value={value}
        onChange={handleChange}
        onKeyPress={handleKeyPress}
        disabled={disabled}
        size={size}
        fullWidth={fullWidth}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon 
                sx={{ 
                  color: theme.palette.text.secondary,
                  fontSize: size === 'small' ? '1.25rem' : '1.5rem',
                }} 
              />
            </InputAdornment>
          ),
          endAdornment: value && (
            <InputAdornment position="end">
              <Tooltip title="Clear search">
                <IconButton
                  size="small"
                  onClick={handleClear}
                  sx={{ 
                    color: theme.palette.text.secondary,
                    '&:hover': {
                      color: theme.palette.text.primary,
                    },
                  }}
                >
                  <ClearIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            </InputAdornment>
          ),
        }}
        sx={{
          '& .MuiOutlinedInput-root': {
            borderRadius: 2,
            backgroundColor: theme.palette.mode === 'dark' 
              ? 'rgba(255, 255, 255, 0.05)' 
              : 'rgba(0, 0, 0, 0.02)',
            '&:hover': {
              backgroundColor: theme.palette.mode === 'dark' 
                ? 'rgba(255, 255, 255, 0.08)' 
                : 'rgba(0, 0, 0, 0.04)',
            },
            '&.Mui-focused': {
              backgroundColor: theme.palette.mode === 'dark' 
                ? 'rgba(255, 255, 255, 0.1)' 
                : 'rgba(0, 0, 0, 0.06)',
            },
            '& fieldset': {
              borderColor: theme.palette.mode === 'dark' 
                ? 'rgba(255, 255, 255, 0.1)' 
                : 'rgba(0, 0, 0, 0.1)',
            },
            '&:hover fieldset': {
              borderColor: theme.palette.mode === 'dark' 
                ? 'rgba(255, 255, 255, 0.2)' 
                : 'rgba(0, 0, 0, 0.2)',
            },
            '&.Mui-focused fieldset': {
              borderColor: theme.palette.primary.main,
            },
          },
          '& .MuiInputBase-input': {
            fontSize: size === 'small' ? '0.875rem' : '1rem',
          },
        }}
      />
    </Box>
  );
};

export default SearchBar; 