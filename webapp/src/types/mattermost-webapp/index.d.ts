export interface PluginRegistry {
  registerPostTypeComponent(typeName: string, component: React.ElementType);
  registerCallButtonAction(icon, action, dropdownText, tooltipText);
  registerChannelHeaderButtonAction(icon, action, dropdownText);
  // Add more if needed from https://developers.mattermost.com/extend/plugins/webapp/reference
}
