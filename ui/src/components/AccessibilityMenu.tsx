import React, { useState } from 'react';
import { 
  IconButton, 
  Menu, 
  MenuItem, 
  ListItemIcon, 
  ListItemText, 
  Divider,
  Typography,
  Switch,
  Tooltip,
  useTheme
} from '@mui/material';
import { 
  Accessibility as AccessibilityIcon,
  Visibility as VisibilityIcon,
  TextFields as TextFieldsIcon,
  Animation as AnimationIcon,
  Settings as SettingsIcon
} from '@mui/icons-material';
import { useAccessibility } from './AccessibilityProvider';

const AccessibilityMenu: React.FC = () => {
  const theme = useTheme();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);
  
  const {
    highContrast,
    largeText,
    reducedMotion,
    toggleHighContrast,
    toggleLargeText,
    toggleReducedMotion,
    resetAccessibility,
  } = useAccessibility();

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleReset = () => {
    resetAccessibility();
    handleClose();
  };

  return (
    <>
      <Tooltip title="Accessibility settings">
        <IconButton
          onClick={handleClick}
          sx={{
            color: theme.palette.mode === 'dark' ? 'inherit' : 'primary.main',
            '&:hover': {
              backgroundColor: theme.palette.mode === 'dark' 
                ? 'rgba(255, 255, 255, 0.08)' 
                : 'rgba(130, 82, 255, 0.08)',
            },
          }}
        >
          <AccessibilityIcon />
        </IconButton>
      </Tooltip>
      
      <Menu
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        PaperProps={{
          elevation: 3,
          sx: {
            width: 280,
            borderRadius: 2,
            mt: 1.5,
            '& .MuiMenuItem-root': {
              px: 2,
              py: 1.5,
            },
          },
        }}
        transformOrigin={{ horizontal: 'right', vertical: 'top' }}
        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
      >
        <MenuItem>
          <ListItemIcon>
            <SettingsIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>
            <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
              Accessibility Settings
            </Typography>
          </ListItemText>
        </MenuItem>
        
        <Divider />
        
        <MenuItem onClick={toggleHighContrast}>
          <ListItemIcon>
            <VisibilityIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>
            <Typography variant="body2">High Contrast</Typography>
            <Typography variant="caption" color="text.secondary">
              Improve text visibility
            </Typography>
          </ListItemText>
          <Switch 
            checked={highContrast} 
            size="small"
            sx={{ ml: 1 }}
          />
        </MenuItem>
        
        <MenuItem onClick={toggleLargeText}>
          <ListItemIcon>
            <TextFieldsIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>
            <Typography variant="body2">Large Text</Typography>
            <Typography variant="caption" color="text.secondary">
              Increase font size
            </Typography>
          </ListItemText>
          <Switch 
            checked={largeText} 
            size="small"
            sx={{ ml: 1 }}
          />
        </MenuItem>
        
        <MenuItem onClick={toggleReducedMotion}>
          <ListItemIcon>
            <AnimationIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>
            <Typography variant="body2">Reduced Motion</Typography>
            <Typography variant="caption" color="text.secondary">
              Minimize animations
            </Typography>
          </ListItemText>
          <Switch 
            checked={reducedMotion} 
            size="small"
            sx={{ ml: 1 }}
          />
        </MenuItem>
        
        <Divider />
        
        <MenuItem onClick={handleReset}>
          <ListItemIcon>
            <SettingsIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>
            <Typography variant="body2">Reset to Default</Typography>
          </ListItemText>
        </MenuItem>
      </Menu>
    </>
  );
};

export default AccessibilityMenu; 