import { Hono } from 'hono'

const app = new Hono()

app.get('/', (c) => {
  return c.text('Hello Hono! 3')
})

export default {
  port: 7893,
  fetch: app.fetch,
}