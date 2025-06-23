import { Hono } from 'hono'

const app = new Hono()

app.get('/', (c) => {
  return c.text('Hello Hono! 2')
})

export default {
  port: 7892,
  fetch: app.fetch,
}