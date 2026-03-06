import { create } from 'zustand'

export interface WizardState {
  step: number
  projectName: string
  projectDescription: string
  frontend: string
  backend: string
  database: string
  claudeMdRules: string[]
  hooks: string[]
  cicd: string
  testing: string
  setStep: (step: number) => void
  nextStep: () => void
  prevStep: () => void
  setField: <K extends keyof WizardState>(key: K, value: WizardState[K]) => void
  reset: () => void
}

const initialState = {
  step: 0,
  projectName: '',
  projectDescription: '',
  frontend: 'react',
  backend: 'node',
  database: 'postgres',
  claudeMdRules: ['never-commit-without-build', 'never-expose-secrets', 'always-push-after-commit'],
  hooks: ['pre-commit-lint', 'pre-commit-typecheck'],
  cicd: 'github-actions',
  testing: 'vitest',
}

export const useWizardStore = create<WizardState>((set) => ({
  ...initialState,
  setStep: (step) => set({ step }),
  nextStep: () => set((s) => ({ step: Math.min(s.step + 1, 4) })),
  prevStep: () => set((s) => ({ step: Math.max(s.step - 1, 0) })),
  setField: (key, value) => set({ [key]: value }),
  reset: () => set(initialState),
}))
