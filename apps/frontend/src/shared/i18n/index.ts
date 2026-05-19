import { messages as pt } from './messages.pt'
import { messages as en } from './messages.en'
import type { MessageKey } from './messages.pt'

const dictionaries: Record<string, Record<string, string>> = {
  pt,
  en,
}

function getBrowserLanguage(): string {
  if (typeof window === 'undefined') return 'en'
  const lang = navigator.language || 'en'
  if (lang.startsWith('pt')) return 'pt'
  return 'en'
}

export function getMessage(key: MessageKey | string): string {
  const lang = getBrowserLanguage()
  const dict = dictionaries[lang] || dictionaries.en
  return (dict[key as MessageKey] || key) as string
}
