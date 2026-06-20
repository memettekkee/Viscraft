import { useEffect, useRef } from 'react'
import useSWR from 'swr'
import { useSWRConfig } from 'swr'
import { createProject } from '../../../service/project'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import { postFetcher } from '../../../helper/fetcher'
import type { ApiResponse, Project } from '../../../types'

const SEED_KEY = 'viscraft-sample-seeded'

/**
 * Auto-creates a sample campaign for new users + auto-selects first campaign.
 */
export function useSeedSampleCampaign() {
  const { mutate } = useSWRConfig()
  const activeProjectId = useWorkspaceStore((s) => s.activeProjectId)
  const setActiveProject = useWorkspaceStore((s) => s.setActiveProject)
  const seeding = useRef(false)

  const { data, isLoading } = useSWR<ApiResponse<Project[]>>(
    ['/projects/list'],
    postFetcher
  )

  const projects = data?.data ?? []

  // Auto-select first project if none is active
  useEffect(() => {
    if (!isLoading && projects.length > 0 && !activeProjectId) {
      setActiveProject(projects[0].id)
    }
  }, [isLoading, projects, activeProjectId, setActiveProject])

  // Auto-seed sample campaign for brand new users
  useEffect(() => {
    if (isLoading) return
    if (seeding.current) return
    if (localStorage.getItem(SEED_KEY)) return
    if (projects.length > 0) {
      localStorage.setItem(SEED_KEY, 'true')
      return
    }

    seeding.current = true
    localStorage.setItem(SEED_KEY, 'true')

    createProject({
      name: 'My First Campaign',
      description: 'Sample campaign to get you started',
      productCategory: 'beverage',
      visualStyle: 'Clean minimal',
    }).then((res) => {
      if (res.success && res.data) {
        mutate(['/projects/list'])
        setActiveProject(res.data.id)
      }
    }).catch(() => {
      // Silent fail
    })
  }, [isLoading, projects.length, mutate, setActiveProject])
}
