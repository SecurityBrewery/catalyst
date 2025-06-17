import { cn, human } from '@/lib/utils'

describe('cn utility function', () => {
  it('returns merged class names', () => {
    expect(cn('class1', 'class2')).toBe('class1 class2')
  })

  it('handles empty inputs gracefully', () => {
    expect(cn()).toBe('')
  })
})

describe('human utility function', () => {
  it('formats bytes less than 1024 as B', () => {
    expect(human(512)).toBe('512 B')
  })

  it('formats bytes between 1024 and 1048576 as KB', () => {
    expect(human(2048)).toBe('2.0 KB')
  })

  it('formats bytes greater than or equal to 1048576 as MB', () => {
    expect(human(2097152)).toBe('2.0 MB')
  })

  it('handles zero bytes correctly', () => {
    expect(human(0)).toBe('0 B')
  })

  it('handles negative byte values gracefully', () => {
    expect(human(-1024)).toBe('-1024 B')
  })
})
