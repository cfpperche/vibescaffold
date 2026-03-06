# VibeScaffold — Screenshots (no ffmpeg needed)
# Run: make screenshots
# Uses arrow+enter navigation (more reliable in VHS than number keys)

Set FontFamily "JetBrains Mono"
Set FontSize 18
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set Padding 20

Hide
Type "mkdir -p ~/.vibescaffold && echo 1 > ~/.vibescaffold/onboarding_seen"
Enter
Sleep 500ms
Type "cd /tmp && rm -rf vs-shots && mkdir vs-shots && cd vs-shots"
Enter
Sleep 1s
Show

Type "vs"
Enter
Sleep 8s
Screenshot demos/screenshots/home.png

# Doctor: Down moves to doctor (index 1), Enter navigates
Sleep 3s
Down
Sleep 3s
Enter
Sleep 8s
Screenshot demos/screenshots/doctor.png

Escape
Sleep 5s

# Status
Down
Sleep 500ms
Up
Sleep 500ms
Type "3"
Sleep 5s
Screenshot demos/screenshots/status.png

Escape
Sleep 5s

# Agent
Down
Sleep 500ms
Up
Sleep 500ms
Type "4"
Sleep 5s
Screenshot demos/screenshots/agent.png

Escape
Sleep 5s

# Onboarding
Down
Sleep 500ms
Up
Sleep 500ms
Type "5"
Sleep 5s
Screenshot demos/screenshots/onboarding-menu.png

# Quit
Escape
Sleep 2s
Type "q"
Sleep 1s
