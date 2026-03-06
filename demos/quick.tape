# VibeScaffold — Quick Demo (home + navegação)
# Run: make demo-quick

Output demos/quick.gif

Set FontFamily "JetBrains Mono"
Set FontSize 20
Set Width 900
Set Height 500
Set Theme "Catppuccin Mocha"
Set TypingSpeed 80ms
Set Padding 20

# Dir limpo sem CLAUDE.md para ir pra home
Hide
Type "cd /tmp && rm -rf vs-quick && mkdir vs-quick && cd vs-quick"
Enter
Sleep 500ms
# Marca onboarding como visto
Type "mkdir -p ~/.vibescaffold && echo 1 > ~/.vibescaffold/onboarding_seen"
Enter
Sleep 500ms
Show

Type "vs"
Enter
Sleep 2s

# Navigate with arrows
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Up
Sleep 500ms
Up
Sleep 500ms

# Go to doctor
Type "2"
Sleep 2s

# Back
Escape
Sleep 1s

# Onboarding
Type "5"
Sleep 2s

# Navigate files
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms

# Back
Type "q"
Sleep 1s

# Quit
Type "q"
Sleep 1s
