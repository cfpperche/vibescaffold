import { createFileRoute } from '@tanstack/react-router'
import { Wizard } from '../components/wizard/Wizard'

export const Route = createFileRoute('/wizard')({
  component: Wizard,
})
