import { expect } from '@playwright/test'
import { randomUUID } from 'crypto'
import { login, test, createTicket } from './util'

const createLink = async (page, name: string, url: string) => {
  await page.getByRole('button', { name: 'Add item' }).first().click()
  await page.locator('input[name="name"]').fill(name)
  await page.locator('#url').fill(url)
  await page.getByRole('button', { name: 'Save' }).click()
  await expect(page.getByText(name)).toBeVisible()
}

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