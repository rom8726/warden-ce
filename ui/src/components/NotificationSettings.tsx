import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Button,
  CircularProgress,
  Divider,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Switch,
  FormControlLabel,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  Collapse,
  Card,
  CardContent,
  CardActions,
  Chip,
  Tabs,
  Tab,
  Tooltip
} from '@mui/material';
import {
  Add as AddIcon,
  Delete as DeleteIcon,
  Edit as EditIcon,
  Refresh as RefreshIcon,
  Close as CloseIcon,
  Send as SendIcon
} from '@mui/icons-material';
import apiClient from '../api/apiClient';
import type {
  NotificationSetting,
  NotificationRule,
  CreateNotificationSettingRequest,
  CreateNotificationRuleRequest
} from '../generated/api/client/api';

// Interface for our project data
interface Project {
  id: number;
  name: string;
  team_id?: number | null;
  team_name?: string | null;
  created_at: string;
}

interface NotificationSettingsProps {
  projectId: number;
}

// Extend the NotificationSetting interface to include rules
interface ExtendedNotificationSetting extends NotificationSetting {
  rules?: NotificationRule[];
}

interface EmailConfig {
  email_to: string;
}

interface MattermostConfig {
  webhook_url: string;
  channel_name: string;
}

interface WebhookConfig {
  webhook_url: string;
}

interface TelegramConfig {
  bot_token: string;
  chat_id: string;
}

interface SlackConfig {
  webhook_url: string;
  channel_name: string;
}

interface PachcaConfig {
  webhook_url: string;
}

const NotificationSettings: React.FC<NotificationSettingsProps> = ({ projectId }) => {
  const [settings, setSettings] = useState<ExtendedNotificationSetting[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [expandedSettings, setExpandedSettings] = useState<number[]>([]);
  const [validationError, setValidationError] = useState<string | null>(null);
  const [openValidationDialog, setOpenValidationDialog] = useState(false);

  // State for field validation
  const [emailFieldError, setEmailFieldError] = useState(false);
  const [mattermostWebhookError, setMattermostWebhookError] = useState(false);
  const [mattermostChannelError, setMattermostChannelError] = useState(false);
  const [webhookUrlError, setWebhookUrlError] = useState(false);
  const [telegramBotTokenError, setTelegramBotTokenError] = useState(false);
  const [telegramChatIdError, setTelegramChatIdError] = useState(false);
  const [slackWebhookError, setSlackWebhookError] = useState(false);
  const [slackChannelError, setSlackChannelError] = useState(false);
  const [pachcaWebhookUrlError, setPachcaWebhookUrlError] = useState(false);
  const [project, setProject] = useState<Project>({
    id: projectId,
    name: `Project ${projectId}`,
    team_id: null,
    team_name: null,
    created_at: new Date().toISOString()
  });

  // State for add/edit setting dialog
  const [openSettingDialog, setOpenSettingDialog] = useState(false);
  const [settingDialogMode, setSettingDialogMode] = useState<'add' | 'edit'>('add');
  const [currentSetting, setCurrentSetting] = useState<ExtendedNotificationSetting | null>(null);
  const [settingType, setSettingType] = useState<string>('email');
  const [settingEnabled, setSettingEnabled] = useState(true);
  const [emailConfig, setEmailConfig] = useState<EmailConfig>({ email_to: '' });
  const [mattermostConfig, setMattermostConfig] = useState<MattermostConfig>({ webhook_url: '', channel_name: '' });
  const [webhookConfig, setWebhookConfig] = useState<WebhookConfig>({ webhook_url: '' });
  const [telegramConfig, setTelegramConfig] = useState<TelegramConfig>({ bot_token: '', chat_id: '' });
  const [slackConfig, setSlackConfig] = useState<SlackConfig>({ webhook_url: '', channel_name: '' });
  const [pachcaConfig, setPachcaConfig] = useState<PachcaConfig>({ webhook_url: '' });

  // State for tab selection
  const [selectedTab, setSelectedTab] = useState<string>('all');

  // State for add/edit rule dialog
  const [openRuleDialog, setOpenRuleDialog] = useState(false);
  const [ruleDialogMode, setRuleDialogMode] = useState<'add' | 'edit'>('add');
  const [currentRule, setCurrentRule] = useState<NotificationRule | null>(null);
  const [currentSettingId, setCurrentSettingId] = useState<number | null>(null);
  const [ruleEventLevel, setRuleEventLevel] = useState<string | null>(null);
  const [ruleFingerprint, setRuleFingerprint] = useState<string | null>(null);
  const [ruleIsNewError, setRuleIsNewError] = useState<boolean | null>(true);
  const [ruleIsRegression, setRuleIsRegression] = useState<boolean | null>(true);

  // Fetch notification rules for a specific setting (used for debugging)
  // This function is intentionally unused in the component but kept for debugging purposes
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const fetchRulesForSetting = async (settingId: number) => {
    try {
      const response = await apiClient.listNotificationRules(projectId, settingId);
      // Update the settings array with the fetched rules
      setSettings(prevSettings =>
        prevSettings.map(setting =>
          setting.id === settingId
            ? { ...setting, rules: response.data.notification_rules }
            : setting
        )
      );
    } catch (err) {
      console.error(`Error fetching rules for setting ${settingId}:`, err);
    }
  };

  // Fetch notification settings
  const fetchSettingsAndRules = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.listNotificationSettings(projectId);
      const fetchedSettings = response.data.notification_settings;
      setSettings(fetchedSettings);

      // After settings are loaded, fetch rules for each setting
      if (fetchedSettings.length > 0) {
        for (const setting of fetchedSettings) {
          try {
            const rulesResponse = await apiClient.listNotificationRules(projectId, setting.id);
            // Update the settings array with the fetched rules
            setSettings(prevSettings =>
              prevSettings.map(s =>
                s.id === setting.id
                  ? { ...s, rules: rulesResponse.data.notification_rules }
                  : s
              )
            );
          } catch (err) {
            console.error(`Error fetching rules for setting ${setting.id}:`, err);
          }
        }
      }
    } catch (err) {
      console.error('Error fetching notification settings:', err);
      setError('Failed to load notification settings');
    } finally {
      setLoading(false);
    }
  };

  // Fetch project details
  const fetchProjectDetails = async () => {
    try {
      const response = await apiClient.getProject(projectId);
      setProject(response.data.project);
    } catch (err) {
      console.error('Error fetching project details:', err);
    }
  };

  // Load data on component mount
  useEffect(() => {
    fetchProjectDetails();
    fetchSettingsAndRules();
  }, [projectId]);

  const hasEmailChannel = () => {
    return settings.some(setting => setting.type === 'email');
  };

  const handleAddSetting = (channelType?: string) => {
    setSettingDialogMode('add');
    setCurrentSetting(null);
    setSettingType(channelType || 'email');
    setSettingEnabled(true);
    setEmailConfig({ email_to: '' });
    setMattermostConfig({ webhook_url: '', channel_name: '' });
    setWebhookConfig({ webhook_url: '' });
    setTelegramConfig({ bot_token: '', chat_id: '' });
    setSlackConfig({ webhook_url: '', channel_name: '' });
    setPachcaConfig({ webhook_url: '' });
    setOpenSettingDialog(true);
  };

  // Helper functions for each channel type
  const handleAddEmailSetting = () => handleAddSetting('email');
  const handleAddTelegramSetting = () => handleAddSetting('telegram');
  const handleAddMattermostSetting = () => handleAddSetting('mattermost');
  const handleAddSlackSetting = () => handleAddSetting('slack');
  const handleAddWebhookSetting = () => handleAddSetting('webhook');
  const handleAddPachcaSetting = () => handleAddSetting('pachca');

  const handleEditSetting = (setting: ExtendedNotificationSetting) => {
    setSettingDialogMode('edit');
    setCurrentSetting(setting);
    setSettingType(setting.type);
    setSettingEnabled(setting.enabled);

    // Parse the config based on the setting type
    try {
      const config = JSON.parse(setting.config);
      switch (setting.type) {
        case 'email':
          setEmailConfig(config);
          break;
        case 'mattermost':
          setMattermostConfig(config);
          break;
        case 'webhook':
          setWebhookConfig(config);
          break;
        case 'telegram':
          setTelegramConfig(config);
          break;
        case 'slack':
          setSlackConfig(config);
          break;
        case 'pachca':
          setPachcaConfig(config);
          break;
      }
    } catch (e) {
      console.error('Error parsing setting config:', e);
    }

    setOpenSettingDialog(true);
  };

  const handleDeleteSetting = async (settingId: number) => {
    if (!confirm('Are you sure you want to delete this notification setting?')) {
      return;
    }

    try {
      await apiClient.deleteNotificationSetting(projectId, settingId);
      setSuccess('Notification setting deleted successfully');
      fetchSettingsAndRules(); // Refresh the list
    } catch (err) {
      console.error('Error deleting notification setting:', err);
      setError('Failed to delete notification setting');
    }
  };

  const handleSendTestNotification = async (settingId: number) => {
    try {
      await apiClient.sendTestNotification(projectId, settingId);
      setSuccess('Test notification sent successfully');
    } catch (err) {
      console.error('Error sending test notification:', err);
      setError('Failed to send test notification');
    }
  };

  const handleSaveSetting = async () => {
    // Validate the configuration based on the setting type
    let config = {};
    let isValid = true;

    switch (settingType) {
      case 'email':
        if (!emailConfig.email_to || !isValidEmail(emailConfig.email_to)) {
          setEmailFieldError(true);
          isValid = false;
        } else {
          setEmailFieldError(false);
          config = emailConfig;
        }
        break;
      case 'mattermost':
        if (!mattermostConfig.webhook_url || !isValidUrl(mattermostConfig.webhook_url)) {
          setMattermostWebhookError(true);
          isValid = false;
        } else {
          setMattermostWebhookError(false);
        }
        if (!mattermostConfig.channel_name) {
          setMattermostChannelError(true);
          isValid = false;
        } else {
          setMattermostChannelError(false);
        }
        if (isValid) {
          config = mattermostConfig;
        }
        break;
      case 'webhook':
        if (!webhookConfig.webhook_url || !isValidUrl(webhookConfig.webhook_url)) {
          setWebhookUrlError(true);
          isValid = false;
        } else {
          setWebhookUrlError(false);
          config = webhookConfig;
        }
        break;
      case 'telegram':
        if (!telegramConfig.bot_token) {
          setTelegramBotTokenError(true);
          isValid = false;
        } else {
          setTelegramBotTokenError(false);
        }
        if (!telegramConfig.chat_id) {
          setTelegramChatIdError(true);
          isValid = false;
        } else {
          setTelegramChatIdError(false);
        }
        if (isValid) {
          config = telegramConfig;
        }
        break;
      case 'slack':
        if (!slackConfig.webhook_url || !isValidUrl(slackConfig.webhook_url)) {
          setSlackWebhookError(true);
          isValid = false;
        } else {
          setSlackWebhookError(false);
        }
        if (!slackConfig.channel_name) {
          setSlackChannelError(true);
          isValid = false;
        } else {
          setSlackChannelError(false);
        }
        if (isValid) {
          config = slackConfig;
        }
        break;
      case 'pachca':
        if (!pachcaConfig.webhook_url || !isValidUrl(pachcaConfig.webhook_url)) {
          setPachcaWebhookUrlError(true);
          isValid = false;
        } else {
          setPachcaWebhookUrlError(false);
          config = pachcaConfig;
        }
        break;
    }

    if (!isValid) {
      setValidationError('Please fix the validation errors above');
      setOpenValidationDialog(true);
      return;
    }

    try {
      const request: CreateNotificationSettingRequest = {
        type: settingType as any,
        config: JSON.stringify(config),
        enabled: settingEnabled
      };

      if (settingDialogMode === 'add') {
        await apiClient.createNotificationSetting(projectId, request);
        setSuccess('Notification setting created successfully');
      } else {
        await apiClient.updateNotificationSetting(projectId, currentSetting!.id, request);
        setSuccess('Notification setting updated successfully');
      }

      setOpenSettingDialog(false);
      fetchSettingsAndRules(); // Refresh the list
    } catch (err) {
      console.error('Error saving notification setting:', err);
      setError('Failed to save notification setting');
    }
  };

  const handleAddRule = (settingId: number) => {
    setRuleDialogMode('add');
    setCurrentRule(null);
    setCurrentSettingId(settingId);
    setRuleEventLevel(null);
    setRuleFingerprint(null);
    setRuleIsNewError(true);
    setRuleIsRegression(true);
    setOpenRuleDialog(true);
  };

  const handleEditRule = (settingId: number, rule: NotificationRule) => {
    setRuleDialogMode('edit');
    setCurrentRule(rule);
    setCurrentSettingId(settingId);
    setRuleEventLevel(rule.event_level === '' ? 'any' : (rule.event_level || null));
    setRuleFingerprint(rule.fingerprint || null);
    setRuleIsNewError(rule.is_new_error || null);
    setRuleIsRegression(rule.is_regression || null);
    setOpenRuleDialog(true);
  };

  const handleDeleteRule = async (settingId: number, ruleId: number) => {
    if (!confirm('Are you sure you want to delete this notification rule?')) {
      return;
    }

    try {
      await apiClient.deleteNotificationRule(projectId, settingId, ruleId);
      setSuccess('Notification rule deleted successfully');
      fetchSettingsAndRules(); // Refresh the list
    } catch (err) {
      console.error('Error deleting notification rule:', err);
      setError('Failed to delete notification rule');
    }
  };

  const handleSaveRule = async () => {
    try {
      // Validation: at least one of is_new_error or is_regression must be enabled
      if (!ruleIsNewError && !ruleIsRegression) {
        setError('At least one of "For new errors" or "For regressions" must be enabled');
        return;
      }

      const request: CreateNotificationRuleRequest = {
        event_level: ruleEventLevel === null ? '' : (ruleEventLevel || undefined),
        fingerprint: ruleFingerprint || undefined,
        is_new_error: ruleIsNewError || undefined,
        is_regression: ruleIsRegression || undefined
      };

      if (ruleDialogMode === 'add') {
        await apiClient.createNotificationRule(projectId, currentSettingId!, request);
        setSuccess('Notification rule created successfully');
      } else {
        await apiClient.updateNotificationRule(projectId, currentSettingId!, currentRule!.id, request);
        setSuccess('Notification rule updated successfully');
      }

      setOpenRuleDialog(false);
      fetchSettingsAndRules(); // Refresh the list
    } catch (err) {
      console.error('Error saving notification rule:', err);
      setError('Failed to save notification rule');
    }
  };

  const renderSettingConfigForm = () => {
    switch (settingType) {
      case 'email':
        return (
          <TextField
            fullWidth
            margin="normal"
            label="Email Address"
            value={emailConfig.email_to}
            onChange={(e) => setEmailConfig({ ...emailConfig, email_to: e.target.value })}
            error={emailFieldError}
            helperText={emailFieldError ? 'Please enter a valid email address' : ''}
          />
        );
      case 'mattermost':
        return (
          <>
            <TextField
              fullWidth
              margin="normal"
              label="Webhook URL"
              value={mattermostConfig.webhook_url}
              onChange={(e) => setMattermostConfig({ ...mattermostConfig, webhook_url: e.target.value })}
              error={mattermostWebhookError}
              helperText={mattermostWebhookError ? 'Please enter a valid webhook URL' : ''}
            />
            <TextField
              fullWidth
              margin="normal"
              label="Channel Name"
              value={mattermostConfig.channel_name}
              onChange={(e) => setMattermostConfig({ ...mattermostConfig, channel_name: e.target.value })}
              error={mattermostChannelError}
              helperText={mattermostChannelError ? 'Please enter a channel name' : ''}
            />
          </>
        );
      case 'webhook':
        return (
          <TextField
            fullWidth
            margin="normal"
            label="Webhook URL"
            value={webhookConfig.webhook_url}
            onChange={(e) => setWebhookConfig({ ...webhookConfig, webhook_url: e.target.value })}
            error={webhookUrlError}
            helperText={webhookUrlError ? 'Please enter a valid webhook URL' : ''}
          />
        );
      case 'telegram':
        return (
          <>
            <TextField
              fullWidth
              margin="normal"
              label="Bot Token"
              value={telegramConfig.bot_token}
              onChange={(e) => setTelegramConfig({ ...telegramConfig, bot_token: e.target.value })}
              error={telegramBotTokenError}
              helperText={telegramBotTokenError ? 'Please enter a bot token' : ''}
            />
            <TextField
              fullWidth
              margin="normal"
              label="Chat ID"
              value={telegramConfig.chat_id}
              onChange={(e) => setTelegramConfig({ ...telegramConfig, chat_id: e.target.value })}
              error={telegramChatIdError}
              helperText={telegramChatIdError ? 'Please enter a chat ID' : ''}
            />
          </>
        );
      case 'slack':
        return (
          <>
            <TextField
              fullWidth
              margin="normal"
              label="Webhook URL"
              value={slackConfig.webhook_url}
              onChange={(e) => setSlackConfig({ ...slackConfig, webhook_url: e.target.value })}
              error={slackWebhookError}
              helperText={slackWebhookError ? 'Please enter a valid webhook URL' : ''}
            />
            <TextField
              fullWidth
              margin="normal"
              label="Channel Name"
              value={slackConfig.channel_name}
              onChange={(e) => setSlackConfig({ ...slackConfig, channel_name: e.target.value })}
              error={slackChannelError}
              helperText={slackChannelError ? 'Please enter a channel name' : ''}
            />
          </>
        );
      case 'pachca':
        return (
          <TextField
            fullWidth
            margin="normal"
            label="Webhook URL"
            value={pachcaConfig.webhook_url}
            onChange={(e) => setPachcaConfig({ ...pachcaConfig, webhook_url: e.target.value })}
            error={pachcaWebhookUrlError}
            helperText={pachcaWebhookUrlError ? 'Please enter a valid webhook URL' : ''}
          />
        );
      default:
        return null;
    }
  };

  const isValidEmail = (email: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const isValidUrl = (url: string): boolean => {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  };

  const toggleExpandSetting = (settingId: number) => {
    setExpandedSettings(prev =>
      prev.includes(settingId)
        ? prev.filter(id => id !== settingId)
        : [...prev, settingId]
    );
  };

  const handleTabChange = (_event: React.SyntheticEvent, newValue: string) => {
    setSelectedTab(newValue);
  };

  const renderRules = (setting: ExtendedNotificationSetting) => {
    if (!setting.rules || setting.rules.length === 0) {
      return (
        <Box sx={{ textAlign: 'center', py: 2 }}>
          <Typography variant="body2" color="text.secondary">
            No notification rules configured for this setting.
          </Typography>
        </Box>
      );
    }

    return (
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
        {setting.rules.map((rule) => (
          <Card key={rule.id} variant="outlined" sx={{ mb: 1 }}>
            <CardContent sx={{ pb: 1 }}>
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 1 }}>
                {rule.event_level !== undefined && rule.event_level !== null && (
                  <Chip
                    label={`Level: ${rule.event_level === '' ? 'Any Level' : rule.event_level}`}
                    size="small"
                    color="primary"
                    variant="outlined"
                  />
                )}
                {rule.fingerprint && (
                  <Chip
                    label={`Fingerprint: ${rule.fingerprint.substring(0, 8)}...`}
                    size="small"
                    color="secondary"
                    variant="outlined"
                  />
                )}
                {rule.is_new_error && (
                  <Chip
                    label="New Errors"
                    size="small"
                    color="info"
                    variant="outlined"
                  />
                )}
                {rule.is_regression && (
                  <Chip
                    label="Regressions"
                    size="small"
                    color="warning"
                    variant="outlined"
                  />
                )}
              </Box>
            </CardContent>
            <CardActions>
              <Button
                size="small"
                startIcon={<EditIcon />}
                onClick={() => handleEditRule(setting.id, rule)}
              >
                Edit
              </Button>
              <Button
                size="small"
                color="error"
                startIcon={<DeleteIcon />}
                onClick={() => handleDeleteRule(setting.id, rule.id)}
              >
                Delete
              </Button>
            </CardActions>
          </Card>
        ))}
      </Box>
    );
  };

  return (
    <Box>
      {/* Success and error messages */}
      <Collapse in={!!success || !!error}>
        <Box sx={{ mb: 2 }}>
          {success && (
            <Alert
              severity="success"
              action={
                <IconButton
                  aria-label="close"
                  color="inherit"
                  size="small"
                  onClick={() => setSuccess(null)}
                >
                  <CloseIcon fontSize="inherit" />
                </IconButton>
              }
            >
              {success}
            </Alert>
          )}
          {error && (
            <Alert
              severity="error"
              action={
                <IconButton
                  aria-label="close"
                  color="inherit"
                  size="small"
                  onClick={() => setError(null)}
                >
                  <CloseIcon fontSize="inherit" />
                </IconButton>
              }
            >
              {error}
            </Alert>
          )}
        </Box>
      </Collapse>

      {/* Header with refresh and add buttons */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="subtitle1">
          Configure notification channels for this project
        </Typography>
        <Box>
          <Button
            startIcon={<RefreshIcon />}
            onClick={fetchSettingsAndRules}
            disabled={loading}
            sx={{ mr: 1 }}
          >
            Refresh
          </Button>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => handleAddSetting()}
            disabled={loading}
          >
            Add Channel
          </Button>
        </Box>
      </Box>

      {/* Tabs for different notification channel types */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
        <Tabs
          value={selectedTab}
          onChange={handleTabChange}
          aria-label="notification channel tabs"
          variant="scrollable"
          scrollButtons="auto"
        >
          <Tab label="All" value="all" />
          <Tab label="Email" value="email" />
          <Tab label="Telegram" value="telegram" />
          <Tab label="Mattermost" value="mattermost" />
          <Tab label="Slack" value="slack" />
          <Tab label="Webhook" value="webhook" />
          <Tab label="Pachca" value="pachca" />
        </Tabs>
      </Box>

      {/* Loading indicator */}
      {loading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
          <CircularProgress />
        </Box>
      )}

      {/* No settings message */}
      {!loading && (settings.length === 0 || settings.filter(setting => selectedTab === 'all' || setting.type === selectedTab).length === 0) && (
        <Paper sx={{ p: 3, textAlign: 'center' }}>
          <Typography variant="body1" sx={{ mb: 2 }}>
            {settings.length === 0
              ? 'No notification settings configured for this project.'
              : `No ${selectedTab !== 'all' ? selectedTab : ''} notification settings configured for this project.`}
          </Typography>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => {
              switch (selectedTab) {
                case 'email':
                  handleAddEmailSetting();
                  break;
                case 'telegram':
                  handleAddTelegramSetting();
                  break;
                case 'mattermost':
                  handleAddMattermostSetting();
                  break;
                case 'slack':
                  handleAddSlackSetting();
                  break;
                case 'webhook':
                  handleAddWebhookSetting();
                  break;
                case 'pachca':
                  handleAddPachcaSetting();
                  break;
                default:
                  handleAddEmailSetting();
                  break;
              }
            }}
          >
            Add {selectedTab !== 'all' ? selectedTab.charAt(0).toUpperCase() + selectedTab.slice(1) : 'Email'} Channel
          </Button>
        </Paper>
      )}

      {/* Settings list */}
      {!loading && settings.length > 0 && (
        <Box>
          {settings
            .filter(setting => selectedTab === 'all' || setting.type === selectedTab)
            .map(setting => (
              <Paper key={setting.id} sx={{ p: 3, mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', flex: 1 }}>
                    <Box>
                      <Typography variant="h6">
                        {setting.type.charAt(0).toUpperCase() + setting.type.slice(1)} Notifications #{setting.id}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {setting.enabled ? 'Enabled' : 'Disabled'}
                      </Typography>
                    </Box>
                    <Box sx={{ ml: 4, flex: 1 }}>
                      {setting.type === 'email' && (
                        project.team_id ? (
                          <Typography variant="body2">
                            Notifications are sent to all project team members.
                          </Typography>
                        ) : (
                          <Typography variant="body2">
                            Email: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as EmailConfig;
                                return config.email_to;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                        )
                      )}
                      {setting.type === 'mattermost' && (
                        <Typography variant="body2">
                          Mattermost: {
                          (() => {
                            try {
                              const config = JSON.parse(setting.config) as MattermostConfig;
                              return config.channel_name;
                            } catch (e) {
                              return 'Invalid configuration';
                            }
                          })()
                        }
                        </Typography>
                      )}
                      {setting.type === 'webhook' && (
                        <Typography variant="body2">
                          Webhook configured
                        </Typography>
                      )}
                      {setting.type === 'telegram' && (
                        <Typography variant="body2">
                          Telegram bot configured
                        </Typography>
                      )}
                      {setting.type === 'slack' && (
                        <Typography variant="body2">
                          Slack: {
                          (() => {
                            try {
                              const config = JSON.parse(setting.config) as SlackConfig;
                              return config.channel_name;
                            } catch (e) {
                              return 'Invalid configuration';
                            }
                          })()
                        }
                        </Typography>
                      )}
                      {setting.type === 'pachca' && (
                        <Typography variant="body2">
                          Pachca bot configured
                        </Typography>
                      )}
                    </Box>
                  </Box>
                  <Box>
                    <Button
                      variant="outlined"
                      onClick={() => toggleExpandSetting(setting.id)}
                      sx={{ mr: 1 }}
                    >
                      {expandedSettings.includes(setting.id) ? 'Hide Rules' : 'Show Rules'}
                    </Button>
                    <IconButton
                      aria-label="edit"
                      onClick={() => handleEditSetting(setting)}
                      size="small"
                      sx={{ mr: 1 }}
                    >
                      <EditIcon />
                    </IconButton>
                    <Tooltip title="Send test notification">
                      <IconButton
                        aria-label="send test notification"
                        onClick={() => handleSendTestNotification(setting.id)}
                        size="small"
                        sx={{ mr: 1 }}
                        color="primary"
                      >
                        <SendIcon />
                      </IconButton>
                    </Tooltip>
                    <IconButton
                      aria-label="delete"
                      onClick={() => handleDeleteSetting(setting.id)}
                      size="small"
                      color="error"
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Box>
                </Box>

                <Collapse in={expandedSettings.includes(setting.id)} timeout="auto" unmountOnExit>
                  <Divider sx={{ my: 2 }} />

                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                    <Typography variant="subtitle2">
                      Notification Rules:
                    </Typography>
                    <Button
                      size="small"
                      startIcon={<AddIcon />}
                      onClick={() => handleAddRule(setting.id)}
                    >
                      Add Rule
                    </Button>
                  </Box>

                  {renderRules(setting)}

                  <Box sx={{ mt: 2 }}>
                    <Typography variant="subtitle2" gutterBottom>
                      Configuration Details:
                    </Typography>
                    <Box sx={{ pl: 2 }}>
                      {setting.type === 'email' && (
                        project.team_id ? (
                          <Typography variant="body2">
                            Notifications are sent to all project team members.
                          </Typography>
                        ) : (
                          <Typography variant="body2">
                            Email: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as EmailConfig;
                                return config.email_to;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                        )
                      )}
                      {setting.type === 'mattermost' && (
                        <>
                          <Typography variant="body2">
                            Webhook URL: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as MattermostConfig;
                                return config.webhook_url;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                          <Typography variant="body2">
                            Channel: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as MattermostConfig;
                                return config.channel_name;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                        </>
                      )}
                      {setting.type === 'webhook' && (
                        <Typography variant="body2">
                          Webhook URL: {
                          (() => {
                            try {
                              const config = JSON.parse(setting.config) as WebhookConfig;
                              return config.webhook_url;
                            } catch (e) {
                              return 'Invalid configuration';
                            }
                          })()
                        }
                        </Typography>
                      )}
                      {setting.type === 'telegram' && (
                        <>
                          <Typography variant="body2">
                            Bot Token: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as TelegramConfig;
                                return config.bot_token;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                          <Typography variant="body2">
                            Chat ID: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as TelegramConfig;
                                return config.chat_id;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                        </>
                      )}
                      {setting.type === 'slack' && (
                        <>
                          <Typography variant="body2">
                            Webhook URL: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as SlackConfig;
                                return config.webhook_url;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                          <Typography variant="body2">
                            Channel: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as SlackConfig;
                                return config.channel_name;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                        </>
                      )}
                      {setting.type === 'pachca' && (
                        <Typography variant="body2">
                          Webhook URL: {
                          (() => {
                            try {
                              const config = JSON.parse(setting.config) as PachcaConfig;
                              return config.webhook_url;
                            } catch (e) {
                              return 'Invalid configuration';
                            }
                          })()
                        }
                        </Typography>
                      )}
                    </Box>
                  </Box>
                </Collapse>
              </Paper>
            ))}
        </Box>
      )}

      {/* Add/Edit Setting Dialog */}
      <Dialog open={openSettingDialog} onClose={() => setOpenSettingDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          {settingDialogMode === 'add' ? 'Add Notification Channel' : 'Edit Notification Channel'}
        </DialogTitle>
        <DialogContent>
          <FormControl fullWidth margin="normal">
            <InputLabel id="setting-type-label">Notification Type</InputLabel>
            <Select
              labelId="setting-type-label"
              value={settingType}
              label="Notification Type"
              onChange={(e) => setSettingType(e.target.value)}
              disabled={settingDialogMode === 'edit'}
            >
              {/* Only show email option if we're editing an existing email channel or if no email channel exists */}
              {(settingDialogMode === 'edit' && currentSetting?.type === 'email') || (settingDialogMode === 'add' && !hasEmailChannel()) ? (
                <MenuItem value="email">Email</MenuItem>
              ) : null}
              <MenuItem value="mattermost">Mattermost</MenuItem>
              <MenuItem value="webhook">Webhook</MenuItem>
              <MenuItem value="telegram">Telegram</MenuItem>
              <MenuItem value="slack">Slack</MenuItem>
              <MenuItem value="pachca">Pachca</MenuItem>
            </Select>
          </FormControl>

          {renderSettingConfigForm()}

          <FormControlLabel
            control={
              <Switch
                checked={settingEnabled}
                onChange={(e) => setSettingEnabled(e.target.checked)}
              />
            }
            label="Enabled"
            sx={{ mt: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenSettingDialog(false)}>Cancel</Button>
          <Button
            onClick={handleSaveSetting}
            variant="contained"
            disabled={loading}
          >
            {loading ? <CircularProgress size={24} /> : 'Save'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Add/Edit Rule Dialog */}
      <Dialog open={openRuleDialog} onClose={() => setOpenRuleDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          {ruleDialogMode === 'add' ? 'Add Notification Rule' : 'Edit Notification Rule'}
        </DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Configure when notifications should be sent. Leave fields empty to match all values.
          </Typography>

          <FormControl fullWidth margin="normal">
            <InputLabel id="event-level-label">Event Level</InputLabel>
            <Select
              labelId="event-level-label"
              value={ruleEventLevel === null ? 'any' : (ruleEventLevel || '')}
              label="Event Level"
              onChange={(e) => setRuleEventLevel(e.target.value === 'any' ? null : e.target.value)}
            >
              <MenuItem value="any">Any Level</MenuItem>
              <MenuItem value="fatal">Fatal</MenuItem>
              <MenuItem value="error">Error</MenuItem>
              <MenuItem value="exception">Exception</MenuItem>
              <MenuItem value="warning">Warning</MenuItem>
              <MenuItem value="info">Info</MenuItem>
              <MenuItem value="debug">Debug</MenuItem>
            </Select>
          </FormControl>

          <FormControlLabel
            control={
              <Switch
                checked={!!ruleIsNewError}
                onChange={(e) => setRuleIsNewError(e.target.checked as boolean ? true : null)}
              />
            }
            label="For new errors"
            sx={{ mt: 2, display: 'block' }}
          />

          <FormControlLabel
            control={
              <Switch
                checked={!!ruleIsRegression}
                onChange={(e) => setRuleIsRegression(e.target.checked as boolean ? true : null)}
              />
            }
            label="For regressions (resolved → unresolved)"
            sx={{ mt: 1, display: 'block' }}
          />

          {!ruleIsNewError && !ruleIsRegression && (
            <Typography variant="body2" color="error" sx={{ mt: 1 }}>
              At least one of "For new errors" or "For regressions" must be enabled
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenRuleDialog(false)}>Cancel</Button>
          <Button
            onClick={handleSaveRule}
            variant="contained"
            disabled={loading}
          >
            {loading ? <CircularProgress size={24} /> : 'Save'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Validation Error Dialog */}
      <Dialog open={openValidationDialog} onClose={() => setOpenValidationDialog(false)}>
        <DialogTitle>Validation Error</DialogTitle>
        <DialogContent>
          <Typography>{validationError}</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenValidationDialog(false)} color="primary">
            OK
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default NotificationSettings;
