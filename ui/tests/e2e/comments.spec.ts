import { expect } from '@playwright/test'
import { randomUUID } from 'crypto'
import { login, test, createTicket, createComment } from './util'

test('can create a comment', async ({ page }) => {
  await login(page)
  const ticketName = `playwright-${randomUUID()}`
  await createTicket(page, ticketName)
  const message = `comment-${randomUUID()}`
  await createComment(page, message)
})

test.describe('update a comment', () => {
  const updates = [
    {
      field: 'message',
      update: async (page) => {
        await page.getByRole('button', { name: 'More' }).click()
        await page.getByRole('menuitem', { name: 'Edit' }).click()
        await page.locator('.CodeMirror textarea').first().fill('Updated Comment')
        await page.getByRole('button', { name: 'Save' }).click()
      },
      assert: async (page) => {
        await expect(page.getByText('Updated Comment')).toBeVisible()
      }
    }
  ]

  for (const { field, update, assert } of updates) {
    test(`can update ${field}`, async ({ page }) => {
      await login(page)
      const ticketName = `playwright-${randomUUID()}`
      await createTicket(page, ticketName)
      const message = `comment-${randomUUID()}`
      await createComment(page, message)
      await update(page, message)
      await assert(page)
    })
  }
})

test('can delete a comment', async ({ page }) => {
  await login(page)
  const ticketName = `playwright-${randomUUID()}`
  await createTicket(page, ticketName)
  const message = `comment-${randomUUID()}`
  await createComment(page, message)
  await page.getByRole('button', { name: 'More' }).click()
  await page.getByRole('menuitem', { name: 'Delete' }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByText(message)).toHaveCount(0)
})
