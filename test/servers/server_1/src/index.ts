import { Hono } from 'hono'

const app = new Hono()

app.get('/', (c) => {
  return c.text('Hello Hono! 1')
})

export default {
  port: 7891,
  fetch: app.fetch,
}
