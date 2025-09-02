import { expect } from '@playwright/test'
import { randomUUID } from 'crypto'
import { login, test, createTicket, createLink } from './util'

test('can create a link', async ({ page }) => {
  await login(page)
  const ticketName = `playwright-${randomUUID()}`
  await createTicket(page, ticketName)
  const linkName = `link-${randomUUID()}`
  await createLink(page, linkName, 'https://example.com')
})

test('can delete a link', async ({ page }) => {
  await login(page)
  const ticketName = `playwright-${randomUUID()}`
  await createTicket(page, ticketName)
  const linkName = `link-${randomUUID()}`
  await createLink(page, linkName, 'https://example.com')
  await page.locator('button', { hasText: 'Delete Link' }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByText(linkName)).toHaveCount(0)
})