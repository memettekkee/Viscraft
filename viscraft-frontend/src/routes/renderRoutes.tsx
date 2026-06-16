import { RouterProvider } from 'react-router-dom'
import { router } from './routesConfig'

export function AppRouter() {
  return <RouterProvider router={router} />
}
