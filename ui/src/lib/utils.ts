import { toast } from '@/components/ui/toast'

import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

import type { ModelError } from '@/client'

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

interface FetchError {
  name: string
  message: string
  response?: Response
}

export function handleError(title: string) {
  return function (error: FetchError) {
    error.response
      ?.json()
      .then((data: ModelError) => {
        toast({
          title: title,
          description: data.message,
          variant: 'destructive'
        })
      })
      .catch(() => {
        toast({
          title: title,
          description: error.message,
          variant: 'destructive'
        })
      })
  }
}
