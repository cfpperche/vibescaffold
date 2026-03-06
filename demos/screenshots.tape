# VibeScaffold — Screenshots Only (no ffmpeg needed)
# Run: make screenshots

Set FontFamily "JetBrains Mono"
Set FontSize 18
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set Padding 20

Type "vs"
Enter
Sleep 3s
Screenshot demos/screenshots/home.png

# Doctor
Type "2"
Sleep 3s
Screenshot demos/screenshots/doctor.png

# Back to home
Escape
Sleep 2s

# Status
Type "3"
Sleep 3s
Screenshot demos/screenshots/status.png

# Back to home
Escape
Sleep 2s

# Agent
Type "4"
Sleep 3s
Screenshot demos/screenshots/agent.png

# Quit
Escape
Sleep 1s
Type "q"
Sleep 500ms
