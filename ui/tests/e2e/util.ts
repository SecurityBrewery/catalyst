import { expect } from '@playwright/test'
import { test as baseTest } from 'playwright/test'
import path from 'path'

export const test = baseTest.extend({
  page: async ({ page }, use) => {
    page.on('console', (msg) => console.log(msg.text()))
    page.on('response', async (response) => {
      if (response.status() >= 400) {
        console.error('Error response:', response.status(), response.url())
      }
      expect(response.status()).toBeLessThan(400)
    })

    await use(page)
  }
})

export const login = async (page, admin: boolean = true) => {
  await page.goto('login')
  if (admin) {
    await page.getByPlaceholder('Email').fill('admin@catalyst-soar.com')
  } else {
    await page.getByPlaceholder('Email').fill('user@catalyst-soar.com')
  }
  await page.getByPlaceholder('Password').fill('1234567890')
  await page.getByRole('button', { name: 'Login' }).click()
  await page.waitForURL('**/dashboard')
}

export const createTicket = async (page, name: string) => {
  await page.goto('tickets/incident')
  await page.getByRole('button', { name: 'New Ticket' }).click()
  await page.locator('#name').fill(name)
  await page.locator('#description').fill('Suspicious behavior detected by user in HR department.')
  await page.locator('#severity').selectOption('Low')
  await page.getByRole('button', { name: 'Save' }).click()
  await page.waitForURL('**/tickets/incident/t*')
}

export const createTimeline = async (page, message: string) => {
  await page.getByRole('tab', { name: 'Timeline' }).click()
  await page.getByRole('button', { name: 'Add Timeline Item' }).click()
  await page.getByRole('tabpanel', { name: 'Timeline' }).getByRole('textbox').fill(message)
  await page.getByRole('button', { name: 'Save' }).click()
  await expect(page.getByText(message)).toBeVisible()
}

export const createComment = async (page, message: string) => {
  await page.getByRole('tab', { name: 'Comments' }).click()
  await page.getByRole('button', { name: 'Add Comment' }).click()
  await page.getByRole('tabpanel', { name: 'Comments' }).getByRole('textbox').fill(message)
  await page.getByRole('button', { name: 'Save' }).click()
  await expect(page.getByText(message)).toBeVisible()
}

export const createTask = async (page, name: string, done: boolean) => {
  await page.getByRole('tab', { name: 'Tasks' }).click()
  await page.getByRole('button', { name: 'Add Task' }).click()
  await page.getByPlaceholder('Add a task...').fill(name)
  await page.getByRole('button', { name: 'Save' }).click()
  if (done) {
    await page.getByRole('checkbox').last().click()
  }
  await expect(page.getByText(name)).toBeVisible()
}

export const createLink = async (page, name: string, url: string) => {
  await page.getByRole('button', { name: 'Add item' }).first().click()
  await page.locator('input[name="name"]').fill(name)
  await page.locator('#url').fill(url)
  await page.getByRole('button', { name: 'Save' }).click()
  await expect(page.getByText(name)).toBeVisible()
}

export const uploadFile = async (page, filePath: string) => {
  await page.getByRole('button', { name: 'Add item' }).last().click()
  await page.setInputFiles('input[type="file"]', filePath)
  await page.getByRole('button', { name: 'Upload 1 file' }).first().click()
  await page.getByRole('button', { name: 'Close' }).first().click()
  const name = path.basename(filePath)
  await expect(page.getByText(name, { exact: true })).toBeVisible()
}
