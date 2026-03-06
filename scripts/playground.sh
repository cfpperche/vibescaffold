#!/usr/bin/env bash
# Cria um diretório temporário limpo e roda o vs lá dentro.
# Uso: make playground
#
# O diretório é criado em .playground/ (gitignored).
# Cada execução limpa e recria do zero para testar o fluxo completo.

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
PLAYGROUND="$ROOT_DIR/.playground"
VS_BIN="$ROOT_DIR/dist/vs"

# Garante que o binário existe
if [ ! -f "$VS_BIN" ]; then
    echo "Binário não encontrado. Rodando make build..."
    make -C "$ROOT_DIR" build
fi

# Limpa playground anterior
rm -rf "$PLAYGROUND"
mkdir -p "$PLAYGROUND"

# Remove flag de onboarding para testar first-run
rm -f ~/.vibescaffold/onboarding_seen

echo "╔══════════════════════════════════════════╗"
echo "║  VibeScaffold Playground                 ║"
echo "║                                          ║"
echo "║  Diretório: .playground/                 ║"
echo "║  Onboarding: resetado                    ║"
echo "║  Binário: dist/vs                        ║"
echo "║                                          ║"
echo "║  Teste o fluxo completo:                 ║"
echo "║  1. Onboarding (first-run)               ║"
echo "║  2. Home → Init → Wizard                 ║"
echo "║  3. Scaffold cria projeto                ║"
echo "║  4. Chat abre automaticamente            ║"
echo "║                                          ║"
echo "╚══════════════════════════════════════════╝"
echo ""

cd "$PLAYGROUND"
exec "$VS_BIN"
