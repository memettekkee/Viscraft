import { Button, HStack, Text } from '@chakra-ui/react'
import { ERROR_MESSAGES } from '../../constants'

export interface ErrorBannerProps {
  errorCode: string
  onRetry?: () => void
  onDismiss?: () => void
}

export function ErrorBanner({ errorCode, onRetry, onDismiss }: ErrorBannerProps) {
  const message = ERROR_MESSAGES[errorCode] ?? 'An unexpected error occurred'

  return (
    <HStack
      w="100%"
      px="4"
      py="3"
      bg="parchment"
      borderWidth="1px"
      borderColor="amber"
      borderRadius="sm"
      align="center"
      justify="space-between"
      gap="3"
    >
      <HStack gap="2" flex="1" minW="0">
        <Text
          fontFamily="mono"
          fontSize="xs"
          color="warmgray"
          flexShrink={0}
        >
          [{errorCode}]
        </Text>
        <Text
          fontFamily="body"
          fontSize="sm"
          color="oxblood"
          lineClamp={2}
        >
          {message}
        </Text>
      </HStack>

      <HStack gap="2" flexShrink={0}>
        {onRetry && (
          <Button
            variant="outline"
            size="sm"
            onClick={onRetry}
            minW="44px"
            minH="44px"
          >
            Retry
          </Button>
        )}
        {onDismiss && (
          <Button
            variant="ghost"
            size="sm"
            onClick={onDismiss}
            minW="44px"
            minH="44px"
            aria-label="Dismiss error"
          >
            ✕
          </Button>
        )}
      </HStack>
    </HStack>
  )
}
