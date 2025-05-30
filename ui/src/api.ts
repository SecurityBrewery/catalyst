import { Configuration, DefaultApi } from '@/client'

const config = new Configuration({
  basePath: 'http://localhost:8090/api'
})

export const api = new DefaultApi(config)
