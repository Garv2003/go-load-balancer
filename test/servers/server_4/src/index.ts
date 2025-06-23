import { Hono } from 'hono'

const app = new Hono()

app.get('/', (c) => {
  return c.text('Hello Hono! 4')
})

export default {
  port: 7894,
  fetch: app.fetch,
}