# VibeForge — Onboarding screenshots
# Run: make screenshots-onboarding
# Separate tape because onboarding is first-run only.

Set FontFamily "JetBrains Mono"
Set FontSize 18
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set Padding 20

Hide
Type "rm -f ~/.vibeforge/onboarding_seen"
Enter
Sleep 500ms
Type "rm -rf /tmp/vs-onb && mkdir /tmp/vs-onb && cd /tmp/vs-onb"
Enter
Sleep 1s
Show

Type "vs"
Enter
Sleep 4s
Screenshot demos/screenshots/onboarding.png

# Navigate to a file to show detail panel
Down
Sleep 400ms
Down
Sleep 400ms
Down
Sleep 400ms
Down
Sleep 400ms
Down
Sleep 2s
Screenshot demos/screenshots/onboarding-detail.png

# Collapse a category
Enter
Sleep 2s
Screenshot demos/screenshots/onboarding-collapsed.png

Type "q"
Sleep 1s
