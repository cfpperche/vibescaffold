# VibeScaffold — Full Demo
# Run: vhs demos/demo.tape

Output demos/demo.gif

Set FontFamily "JetBrains Mono"
Set FontSize 18
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set TypingSpeed 60ms
Set Padding 20

# Start the TUI
Hide
Type "cd /tmp && mkdir -p vs-demo && cd vs-demo"
Enter
Sleep 500ms
Show

Type "vs"
Enter
Sleep 2s
Screenshot demos/screenshots/home.png

# Navigate to doctor
Type "2"
Sleep 2s
Screenshot demos/screenshots/doctor.png

# Back to home
Type "q"
Sleep 1s

# Navigate to status
Type "3"
Sleep 2s
Screenshot demos/screenshots/status.png

# Back to home
Type "q"
Sleep 1s

# Navigate to init wizard
Type "1"
Sleep 2s
Screenshot demos/screenshots/init.png

# Fill in project name
Type "meu-saas"
Enter
Sleep 500ms

# Description
Type "Plataforma SaaS para gestao"
Enter
Sleep 500ms

# Author
Type "dev"
Enter
Sleep 500ms

# GitHub repo? No
Type "n"
Enter
Sleep 1s

# Select project type - SaaS
Enter
Sleep 1s

# Select stack
Enter
Sleep 1s

# Principles - accept defaults
Enter
Sleep 1s

# Features - accept defaults
Enter
Sleep 2s

Screenshot demos/screenshots/done.png
Sleep 3s

# Quit
Escape
Sleep 500ms
Type "q"
Sleep 1s
