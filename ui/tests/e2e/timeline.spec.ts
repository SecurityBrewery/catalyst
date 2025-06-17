import { expect } from '@playwright/test'
import { randomUUID } from 'crypto'
import { login, test, createTicket } from './util'

const createTimeline = async (page, message: string) => {
  await page.getByRole('tab', { name: 'Timeline' }).click()
  await page.getByRole('button', { name: 'Add Timeline Item' }).click()
  await page.getByRole('tabpanel', { name: 'Timeline' }).getByRole('textbox').fill(message)
  await page.getByRole('button', { name: 'Save' }).click()
  await expect(page.getByText(message)).toBeVisible()
}

test('can create a timeline item', async ({ page }) => {
  await login(page)
  const ticketName = `playwright-${randomUUID()}`
  await createTicket(page, ticketName)
  const msg = `timeline-${randomUUID()}`
  await createTimeline(page, msg)
})

test.describe('update a timeline item', () => {
  test('can update message', async ({ page }) => {
    await login(page)
    const ticketName = `playwright-${randomUUID()}`
    await createTicket(page, ticketName)
    const msg = `timeline-${randomUUID()}`
    await createTimeline(page, msg)
    await page.getByRole('tab', { name: 'Timeline' }).click()
    await page.getByRole('button', { name: 'More' }).click()
    await page.getByRole('menuitem', { name: 'Edit' }).click()
    await page.locator('textarea').nth(1).fill('Updated Timeline')
    await page.getByRole('button', { name: 'Save' }).click()
    await expect(page.getByText('Updated Timeline')).toBeVisible()
  })
})

test('can delete a timeline item', async ({ page }) => {
  await login(page)
  const ticketName = `playwright-${randomUUID()}`
  await createTicket(page, ticketName)
  const msg = `timeline-${randomUUID()}`
  await createTimeline(page, msg)
  await page.getByRole('tab', { name: 'Timeline' }).click()
  await page.getByRole('button', { name: 'More' }).click()
  await page.getByRole('menuitem', { name: 'Delete' }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByText(msg)).toHaveCount(0)
})
