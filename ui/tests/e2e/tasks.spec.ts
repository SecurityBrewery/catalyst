import { expect } from '@playwright/test'
import { randomUUID } from 'crypto'
import { login, test, createTicket, createTask } from './util'

test('can create a task', async ({ page }) => {
  await login(page)
  const ticketName = `playwright-${randomUUID()}`
  await createTicket(page, ticketName)
  const taskName = `task-${randomUUID()}`
  await createTask(page, taskName, false)
})

test.describe('update a task', () => {
  const updates = [
    {
      field: 'name',
      update: async (page, taskName: string) => {
        await page.getByText("Toggle Sidebar").click()

        await page.getByRole('tab', { name: 'Tasks' }).click()
        await page.getByText(taskName).click()
        await page.getByRole('tabpanel', { name: 'Tasks' }).getByRole('textbox').fill('Updated Task')
        await page.keyboard.press('Enter')
      },
      assert: async (page) => {
        await expect(page.getByText('Updated Task')).toBeVisible()
      }
    },
    {
      field: 'status',
      update: async (page) => {
        await page.getByRole('tab', { name: 'Tasks' }).click()
        const cb = page.getByRole('checkbox').first()
        await cb.click()
      },
      assert: async (page) => {
        const cb = page.getByRole('checkbox').first()
        await expect(cb).toHaveAttribute('data-state', 'checked')
      }
    }
  ]

  for (const { field, update, assert } of updates) {
    test(`can update ${field}`, async ({ page }) => {
      await login(page)
      const ticketName = `playwright-${randomUUID()}`
      await createTicket(page, ticketName)
      const taskName = `task-${randomUUID()}`
      await createTask(page, taskName, false)
      await update(page, taskName)
      await assert(page)
    })
  }
})

test('can delete a task', async ({ page }) => {
  await login(page)
  const ticketName = `playwright-${randomUUID()}`
  await createTicket(page, ticketName)
  const taskName = `task-${randomUUID()}`
  await createTask(page, taskName, false)
  await page.getByRole('tab', { name: 'Tasks' }).click()
  await page.locator('button', { hasText: 'Delete Task' }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByText(taskName)).toHaveCount(0)
})
