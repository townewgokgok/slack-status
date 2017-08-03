package settings

var SettingsExample = `slack:
  # Your token can be created at https://api.slack.com/custom-integrations/legacy-tokens
  token: xoxp-...

templates:
  lunch: ':fork_and_knife: Having lunch'
  home: ':house: Working remotely'

# The settings after here are necessary to obtain music information.
# You can leave them unchanged unless you use this feature.
#
# Format placeholders:
#   %A : Artist
#   %a : Album
#   %t : Title

itunes:
  template_id: itunes
  watch_interval_sec: 3
  format: ':musical_note: %A - %t (from "%a")'

lastfm:
  template_id: lastfm
  watch_interval_sec: 15
  format: ':musical_note: %A - %t (from "%a")'
  user_name: ...
  api_key: ...
`
