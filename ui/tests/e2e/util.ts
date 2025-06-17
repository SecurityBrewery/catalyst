import { test as baseTest } from 'playwright/test'
import { expect } from '@playwright/test'

export const test = baseTest.extend({
  page: async ({ page }, use) => {
    page.on('console', msg => console.log(msg.text()));
    // page.on('request', request => console.log('>>', request.method(), request.url(), request.postData()))
    page.on('response', async response => {
      let body = ''

      if (response.url().includes('/api/')) {
        try {
          body = await response.text()
        } catch (error) {
          // Ignore errors in reading the response body
        }
      }

      // console.log('<<', response.status(), response.url(), body.trim())

      if (response.status() >= 400) {
        console.error('Error response:', response.status(), response.url())
      }
      expect(response.status()).toBeLessThan(400)
    })

    await use(page)
  }
})

export const login = async (page) => {
  await page.goto('login')
  await page.getByPlaceholder('Username').fill('user@catalyst-soar.com')
  await page.getByPlaceholder('Password').fill('1234567890')
  await page.getByRole('button', { name: 'Login' }).click()
  await page.waitForURL('**/dashboard')
}

export const createTicket = async (page, name: string) => {
  await page.goto('tickets/incident')
  await page.getByRole('button', { name: 'New Ticket' }).click()
  await page.locator('#name').fill(name)
  await page.locator('#description').fill('Test description')
  await page.locator('#severity').selectOption('Low')
  await page.getByRole('button', { name: 'Save' }).click()
  await page.waitForURL('**/tickets/incident/incident*')
}