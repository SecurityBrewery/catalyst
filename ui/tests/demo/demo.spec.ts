import { expect } from '@playwright/test'
import { randomUUID } from 'crypto'
import { login, test, createTicket } from '../e2e/util'

// Test that file upload is disabled in demo mode

test('file upload is disabled', async ({ page }) => {
  await login(page)
  const ticketName = `playwright-${randomUUID()}`
  await createTicket(page, ticketName)
  await expect(page.getByText('Cannot upload files in demo mode')).toBeVisible()
})

// Test that reaction creation is disabled in demo mode

test('reaction creation is disabled', async ({ page }) => {
  await login(page)
  await page.goto('reactions')
  await page.getByRole('button', { name: 'New Reaction' }).click()
  await page.waitForURL('**/reactions/new')
  const saveBtn = page.getByRole('button', { name: 'Save' }).last()
  await expect(saveBtn).toBeDisabled()
  await expect(page.getByText('Reactions cannot be created or edited in demo mode')).toBeVisible()
})
