import { expect } from '@playwright/test'
import { randomUUID } from 'crypto'
import path from 'path'
import { fileURLToPath } from 'url'
import { login, test, createTicket, uploadFile } from './util'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const filePath = path.join(__dirname, '../assets/file.txt')

test('can upload a file', async ({ page }) => {
  await login(page)
  const ticketName = `playwright-${randomUUID()}`
  await createTicket(page, ticketName)
  await uploadFile(page, filePath)
})

test('can delete a file', async ({ page }) => {
  await login(page)
  const ticketName = `playwright-${randomUUID()}`
  await createTicket(page, ticketName)
  await uploadFile(page, filePath)
  await page.locator('button', { hasText: 'Delete File' }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByText(path.basename(filePath))).toHaveCount(0)
})

