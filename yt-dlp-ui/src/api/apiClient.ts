import axios from 'axios'
import { config } from '@/config/config'

const axiosInstance = axios.create({
  baseURL: config.apiBaseUrl.replace(/\/$/, ''),
  timeout: 10_000,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
})

export default axiosInstance
