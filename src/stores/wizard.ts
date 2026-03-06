import { create } from 'zustand';

export type ProjectType = 'saas' | 'api' | 'game' | 'cli' | 'mobile' | 'lib';

export type Runtime = 'bun' | 'node';
export type Linter = 'biome' | 'eslint-prettier';

export interface WizardState {
  step: number;
  projectName: string;
  projectDescription: string;
  projectType: ProjectType;
  runtime: Runtime;
  linter: Linter;
  frontend: string;
  backend: string;
  database: string;
  principles: string[];
  claudeTools: string[];
  setStep: (step: number) => void;
  nextStep: () => void;
  prevStep: () => void;
  setField: <K extends keyof WizardState>(key: K, value: WizardState[K]) => void;
  reset: () => void;
}

const initialState = {
  step: 0,
  projectName: '',
  projectDescription: '',
  projectType: 'saas' as ProjectType,
  runtime: 'bun' as Runtime,
  linter: 'biome' as Linter,
  frontend: 'react',
  backend: 'node',
  database: 'postgres',
  principles: [
    'never-commit-without-build',
    'never-expose-secrets',
    'always-push-after-commit',
    'tests-before-merge',
  ],
  claudeTools: ['claude-md', 'context-docs', 'git-hooks', 'github-actions'],
};

export const useWizardStore = create<WizardState>((set) => ({
  ...initialState,
  setStep: (step) => set({ step }),
  nextStep: () => set((s) => ({ step: Math.min(s.step + 1, 4) })),
  prevStep: () => set((s) => ({ step: Math.max(s.step - 1, 0) })),
  setField: (key, value) => set({ [key]: value }),
  reset: () => set(initialState),
}));
