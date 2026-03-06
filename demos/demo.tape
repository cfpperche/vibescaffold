# VibeScaffold — Full Demo
# Run: make demo
# Fluxo: onboarding → home → init wizard → scaffold → chat

Output demos/demo.gif

Set FontFamily "JetBrains Mono"
Set FontSize 18
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set TypingSpeed 60ms
Set Padding 20

# Setup: dir limpo, reseta onboarding
Hide
Type "rm -f ~/.vibescaffold/onboarding_seen"
Enter
Sleep 500ms
Type "cd /tmp && rm -rf vs-demo && mkdir vs-demo && cd vs-demo"
Enter
Sleep 500ms
Show

# Inicia o VS
Type "vs"
Enter
Sleep 3s
Screenshot demos/screenshots/onboarding.png

# Navega pelas categorias do onboarding
Down
Sleep 400ms
Down
Sleep 400ms
Down
Sleep 400ms
Down
Sleep 400ms
Down
Sleep 400ms

# Colapsa uma categoria
Enter
Sleep 1s

# Volta e expande
Enter
Sleep 1s

# Sai do onboarding → home
Type "q"
Sleep 3s
Screenshot demos/screenshots/home.png

# Doctor
Type "2"
Sleep 3s
Screenshot demos/screenshots/doctor.png

Escape
Sleep 2s

# Agent selector
Type "4"
Sleep 3s
Screenshot demos/screenshots/agent.png

Escape
Sleep 2s

# Init wizard
Type "1"
Sleep 2s
Screenshot demos/screenshots/init.png

# Project name
Type "meu-saas"
Enter
Sleep 500ms

# Description
Type "Plataforma SaaS para gestao de projetos"
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

# Project type - SaaS (first option)
Enter
Sleep 1s

# Stack - Go + Bubble Tea (first option)
Enter
Sleep 1s

# Principles - accept defaults
Enter
Sleep 1s

# Features - accept defaults
Enter
Sleep 3s
Screenshot demos/screenshots/done.png

# Enter para abrir chat
Enter
Sleep 3s
Screenshot demos/screenshots/chat.png

# Digita uma mensagem no chat
Type "/help"
Enter
Sleep 2s
Screenshot demos/screenshots/chat-help.png

# Sai
Type "/exit"
Enter
Sleep 1s
