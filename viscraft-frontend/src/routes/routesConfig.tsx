import { createBrowserRouter } from 'react-router-dom'
import { LandingPage } from '../pages/landing'
import { WorkspacePage } from '../pages/workspace'
import { NotFoundPage } from '../pages/notfound'
import { ProtectedRoute } from '../components/ProtectedRoute'
import { Layout } from '../components/layout/Layout'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <LandingPage />,
  },
  {
    path: '/workspace',
    element: (
      <ProtectedRoute>
        <Layout>
          <WorkspacePage />
        </Layout>
      </ProtectedRoute>
    ),
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])
