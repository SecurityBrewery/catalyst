import { test, expect } from '@playwright/test'
import { login, createTicket, createLink, createComment, createTimeline, createTask } from '../tests/e2e/util'

const screenshot = (name: string) => `../docs/screenshots/${name}.png`

const createFullTicket = async (page, name: string) => {
  await createTicket(page, name)
  await createLink(page, 'https://hr-portal.example.com', 'HR Portal Access Log')
  await createTask(page, 'Review access logs for HR user', true)
  await createTask(page, 'Interview HR staff involved', false)
  await createComment(page, 'Suspicious login detected from unrecognized device in HR department')
  await createComment(page, 'User reported unusual activity on their HR account')
  await createTimeline(page, 'Security team notified and initial investigation started')
  await createTimeline(page, 'Access to sensitive HR data temporarily restricted')
  await createTimeline(page, 'Awaiting further analysis from IT forensics')
}

test('dashboard screenshot', async ({ page }) => {
  await login(page)
  await page.waitForTimeout(7000)
  await page.screenshot({ path: screenshot('dashboard') })
})

test('ticket screenshot', async ({ page }) => {
  await login(page)
  const name = 'INCIDENT-1234'
  await createFullTicket(page, name)
  await page.getByText("Toggle Sidebar").click()
  await page.waitForTimeout(7000)
  await page.screenshot({ path: screenshot('ticket') })
})

test('tasks screenshot', async ({ page }) => {
  await login(page)
  const ticketName = 'INCIDENT-1234'
  await createFullTicket(page, ticketName)
  await createTask(page, 'Follow up with HR department', false)
  await page.getByText("Toggle Sidebar").click()
  await page.waitForTimeout(7000)
  await page.screenshot({ path: screenshot('tasks') })
})

test('reactions screenshot', async ({ page }) => {
  await login(page)
  await page.goto('reactions')
  await page.getByText('Assign new Tickets').click()
  await expect(page.getByRole('heading', { name: 'Reactions' })).toBeVisible()
  await page.getByText("Toggle Sidebar").click()
  await page.waitForTimeout(7000)
  await page.screenshot({ path: screenshot('reactions') })
})
