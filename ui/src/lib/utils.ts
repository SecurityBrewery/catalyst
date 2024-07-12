import { toast } from '@/components/ui/toast'

import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function human(bytes: number) {
  return bytes < 1024
    ? bytes + ' B'
    : bytes < 1048576
      ? (bytes / 1024).toFixed(1) + ' KB'
      : (bytes / 1048576).toFixed(1) + ' MB'
}

export function handleError(error: Error) {
  toast({
    title: error.name,
    description: error.message,
    variant: 'destructive'
  })
}
