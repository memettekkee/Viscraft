import { useState } from 'react'
import { Box, Button, Flex, Heading, Text } from '@chakra-ui/react'
import { AuthModal } from './components/AuthModal'

/**
 * Landing page with hero section and call-to-action to open the AuthModal.
 *
 * Validates: Requirements 1.1, 2.1
 */
export function LandingPage() {
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false)

  return (
    <Box bg="ink" minH="100vh">
      <Flex
        direction="column"
        align="center"
        justify="center"
        minH="100vh"
        px={{ base: '6', md: '8' }}
        textAlign="center"
      >
        <Heading
          as="h1"
          fontFamily="display"
          fontSize={{ base: '4xl', md: '6xl', lg: '7xl' }}
          color="parchment"
          mb="4"
        >
          Viscraft
        </Heading>

        <Text
          fontFamily="body"
          fontSize={{ base: 'lg', md: 'xl' }}
          color="warmgray"
          maxW="600px"
          mb="10"
        >
          AI-powered concept art generation for game developers, illustrators,
          and world-builders.
        </Text>

        <Button
          bg="amber"
          color="white"
          fontFamily="body"
          fontWeight="medium"
          fontSize={{ base: 'md', md: 'lg' }}
          px="8"
          py="6"
          borderRadius="sm"
          _hover={{ opacity: 0.9 }}
          _focusVisible={{
            outline: '2px solid',
            outlineColor: 'amber',
            outlineOffset: '2px',
          }}
          onClick={() => setIsAuthModalOpen(true)}
        >
          Get Started
        </Button>
      </Flex>

      <AuthModal
        isOpen={isAuthModalOpen}
        onClose={() => setIsAuthModalOpen(false)}
      />
    </Box>
  )
}
