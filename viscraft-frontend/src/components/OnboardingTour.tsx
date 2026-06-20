import { useState, useEffect, useCallback } from 'react'
import { Box, Button, Text, HStack, Flex } from '@chakra-ui/react'
import { createPortal } from 'react-dom'
import { useAuthStore } from '../store/authStore'
import { completeTour } from '../service/auth'

interface TourStep {
  target: string
  title: string
  content: string
  placement: 'right' | 'bottom' | 'left' | 'top'
}

const STEPS: TourStep[] = [
  {
    target: '[data-tour="sidebar"]',
    title: '👋 Welcome to Viscraft!',
    content: 'This is your campaign list. Each campaign groups ad shots for one product. Click "+ New Campaign" to create your first campaign.',
    placement: 'right',
  },
  {
    target: '[data-tour="generate-card"]',
    title: '✨ Generate Ad Shots',
    content: 'Click here to create a new ad shot. You\'ll describe your product, pick background, lighting, mood, and angle — then the AI builds a professional prompt and generates your ad photo.',
    placement: 'bottom',
  },
  {
    target: '[data-tour="sidebar"]',
    title: '🔄 Regenerate Anytime',
    content: 'After generating, your ad shots appear in the grid. Click any image to view it full-size, or hit "Regenerate" to tweak the prompt and create a new version. Old shots are never deleted.',
    placement: 'right',
  },
]

function getTargetRect(selector: string): DOMRect | null {
  const el = document.querySelector(selector)
  return el ? el.getBoundingClientRect() : null
}

function TooltipBox({
  step,
  index,
  total,
  onNext,
  onSkip,
}: {
  step: TourStep
  index: number
  total: number
  onNext: () => void
  onSkip: () => void
}) {
  const [rect, setRect] = useState<DOMRect | null>(null)

  useEffect(() => {
    const update = () => setRect(getTargetRect(step.target))
    update()
    window.addEventListener('resize', update)
    return () => window.removeEventListener('resize', update)
  }, [step.target])

  if (!rect) return null

  const PADDING = 8
  const TOOLTIP_W = 300
  const TOOLTIP_H = 200
  const scrollY = window.scrollY

  let top = 0
  let left = 0

  switch (step.placement) {
    case 'right':
      top = rect.top + scrollY + rect.height / 2 - TOOLTIP_H / 2
      left = rect.right + PADDING
      break
    case 'bottom':
      top = rect.bottom + scrollY + PADDING
      left = rect.left + rect.width / 2 - TOOLTIP_W / 2
      break
    case 'left':
      top = rect.top + scrollY + rect.height / 2 - TOOLTIP_H / 2
      left = rect.left - TOOLTIP_W - PADDING
      break
    case 'top':
      top = rect.top + scrollY - TOOLTIP_H - PADDING
      left = rect.left + rect.width / 2 - TOOLTIP_W / 2
      break
  }

  left = Math.max(PADDING, Math.min(left, window.innerWidth - TOOLTIP_W - PADDING))
  top = Math.max(PADDING, top)

  const isLast = index === total - 1

  return createPortal(
    <>
      {/* Overlay */}
      <Box
        position="fixed"
        inset="0"
        bg="rgba(22,20,15,0.65)"
        zIndex={9998}
        onClick={onSkip}
        pointerEvents="auto"
      />
      {/* Spotlight */}
      <Box
        position="fixed"
        zIndex={9999}
        style={{
          top: rect.top - 4,
          left: rect.left - 4,
          width: rect.width + 8,
          height: rect.height + 8,
          borderRadius: 6,
          boxShadow: '0 0 0 9999px rgba(22,20,15,0.65)',
          pointerEvents: 'none',
        }}
      />
      {/* Tooltip */}
      <Box
        position="absolute"
        zIndex={10000}
        style={{ top, left, width: TOOLTIP_W }}
        bg="#FAF6EC"
        borderWidth="1px"
        borderColor="#C9762C"
        borderRadius="md"
        p="4"
        boxShadow="lg"
        pointerEvents="auto"
      >
        <Text fontFamily="display" fontSize="sm" fontWeight="bold" color="#16140F" mb="1">
          {step.title}
        </Text>
        <Text fontFamily="body" fontSize="xs" color="#16140F" lineHeight="tall" mb="3" whiteSpace="pre-line">
          {step.content}
        </Text>
        <Flex justify="space-between" align="center">
          <Text fontFamily="mono" fontSize="2xs" color="#6B6555">
            {index + 1} / {total}
          </Text>
          <HStack gap="2">
            <Button
              size="xs"
              variant="ghost"
              color="#6B6555"
              onClick={onSkip}
              minH="28px"
            >
              Skip
            </Button>
            <Button
              size="xs"
              bg="#C9762C"
              color="white"
              onClick={onNext}
              minH="28px"
              _hover={{ opacity: 0.9 }}
            >
              {isLast ? 'Got it!' : 'Next →'}
            </Button>
          </HStack>
        </Flex>
      </Box>
    </>,
    document.body
  )
}

export function OnboardingTour() {
  const user = useAuthStore((s) => s.user)
  const updateUser = useAuthStore((s) => s.updateUser)
  const [step, setStep] = useState<number | null>(null)

  useEffect(() => {
    // Tour hanya muncul kalau user belum complete tour (dari DB)
    if (user && !user.tourCompleted) {
      const timer = setTimeout(() => setStep(0), 2000)
      return () => clearTimeout(timer)
    }
  }, [user])

  const markDone = useCallback(async () => {
    setStep(null)
    // Update DB
    try {
      await completeTour()
    } catch {
      // silent fail, tidak kritis
    }
    // Update local state biar tidak muncul lagi tanpa reload
    if (user) {
      updateUser({ ...user, tourCompleted: true })
    }
  }, [user, updateUser])

  const handleNext = useCallback(() => {
    if (step === null) return
    if (step >= STEPS.length - 1) {
      markDone()
    } else {
      setStep(step + 1)
    }
  }, [step, markDone])

  const handleSkip = useCallback(() => {
    markDone()
  }, [markDone])

  if (step === null || step >= STEPS.length) return null

  return (
    <TooltipBox
      step={STEPS[step]}
      index={step}
      total={STEPS.length}
      onNext={handleNext}
      onSkip={handleSkip}
    />
  )
}
