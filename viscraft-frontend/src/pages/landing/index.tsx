import { useState } from 'react'
import { Box, Button, Flex, Heading, Text, VStack, HStack, Image } from '@chakra-ui/react'
import { motion } from 'framer-motion'
import { AuthModal } from './components/AuthModal'
import heroImage from '../../assets/hero.png'

const MotionBox = motion.create(Box)
const MotionFlex = motion.create(Flex)

export function LandingPage() {
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false)

  return (
    <Box bg="ink" minH="100vh" overflow="hidden">
      {/* Hero Section */}
      <Flex
        direction={{ base: 'column', lg: 'row' }}
        align="center"
        justify="center"
        minH="100vh"
        px={{ base: '6', md: '12', lg: '20' }}
        gap={{ base: '8', lg: '16' }}
        py={{ base: '16', lg: '0' }}
      >
        {/* Left: Text content */}
        <MotionFlex
          direction="column"
          align={{ base: 'center', lg: 'flex-start' }}
          textAlign={{ base: 'center', lg: 'left' }}
          flex="1"
          maxW="600px"
          initial={{ opacity: 0, x: -30 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
        >
          {/* App name — big and prominent */}
          <Heading
            as="h1"
            fontFamily="display"
            fontSize={{ base: '5xl', md: '6xl', lg: '7xl' }}
            color="parchment"
            mb="3"
            lineHeight="1"
            letterSpacing="tight"
          >
            Viscraft
          </Heading>

          <Heading
            as="h2"
            fontFamily="display"
            fontSize={{ base: '2xl', md: '3xl', lg: '4xl' }}
            color="parchment"
            mb="5"
            lineHeight="1.2"
          >
            Product Photos{' '}
            <Text as="span" color="amber">
              Made Easy
            </Text>
          </Heading>

          <Text
            fontFamily="body"
            fontSize={{ base: 'md', md: 'lg' }}
            color="warmgray"
            mb="8"
            lineHeight="tall"
          >
            Generate stunning product photography with AI.
            Pick your background, lighting, and mood — get studio-quality
            ad shots in seconds without a photographer.
          </Text>

          <HStack gap="4" flexWrap="wrap" justify={{ base: 'center', lg: 'flex-start' }}>
            <Button
              bg="amber"
              color="white"
              fontFamily="body"
              fontWeight="medium"
              fontSize="md"
              px="8"
              py="6"
              borderRadius="sm"
              _hover={{ opacity: 0.9, transform: 'translateY(-1px)' }}
              _focusVisible={{
                outline: '2px solid',
                outlineColor: 'amber',
                outlineOffset: '2px',
              }}
              onClick={() => setIsAuthModalOpen(true)}
            >
              Sign In
            </Button>
          </HStack>
        </MotionFlex>

        {/* Right: Hero image / preview */}
        <MotionBox
          flex="1"
          maxW={{ base: '400px', lg: '500px' }}
          initial={{ opacity: 0, x: 30 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6, delay: 0.2 }}
        >
          <Box
            position="relative"
            borderRadius="lg"
            overflow="hidden"
            boxShadow="0 25px 50px -12px rgba(201, 118, 44, 0.25)"
            border="1px solid"
            borderColor="rgba(201, 118, 44, 0.3)"
          >
            <Image
              src={heroImage}
              alt="Viscraft storyboard preview"
              w="100%"
              h="auto"
              objectFit="cover"
            />
            {/* Gradient overlay at bottom */}
            <Box
              position="absolute"
              bottom="0"
              left="0"
              right="0"
              h="40%"
              bgGradient="to-t"
              gradientFrom="ink"
              gradientTo="transparent"
            />
          </Box>
        </MotionBox>
      </Flex>

      {/* Feature highlights section */}
      <Box py="20" px={{ base: '6', md: '12' }}>
        <VStack gap="12" maxW="900px" mx="auto">
          <Heading
            as="h2"
            fontFamily="display"
            fontSize={{ base: '2xl', md: '3xl' }}
            color="parchment"
            textAlign="center"
          >
            How It Works
          </Heading>

          <Flex
            direction={{ base: 'column', md: 'row' }}
            gap="8"
            w="100%"
          >
            <FeatureCard
              step="1"
              title="Describe"
              description="Tell us about your product — what it is, its key features, and the vibe you want."
            />
            <FeatureCard
              step="2"
              title="Style"
              description="Pick background, lighting, mood, and angle. We build the perfect prompt for the AI."
            />
            <FeatureCard
              step="3"
              title="Generate"
              description="Get studio-quality product photos in seconds. Tweak and regenerate until it's perfect."
            />
          </Flex>
        </VStack>
      </Box>

      <AuthModal
        isOpen={isAuthModalOpen}
        onClose={() => setIsAuthModalOpen(false)}
      />
    </Box>
  )
}

function FeatureCard({ step, title, description }: { step: string; title: string; description: string }) {
  return (
    <Box
      flex="1"
      bg="rgba(201, 118, 44, 0.06)"
      border="1px solid"
      borderColor="rgba(201, 118, 44, 0.2)"
      borderRadius="md"
      p="6"
      textAlign="center"
    >
      <Box
        w="40px"
        h="40px"
        borderRadius="full"
        bg="amber"
        display="flex"
        alignItems="center"
        justifyContent="center"
        mx="auto"
        mb="4"
      >
        <Text fontFamily="mono" fontSize="sm" fontWeight="bold" color="white">
          {step}
        </Text>
      </Box>
      <Text fontFamily="display" fontSize="lg" color="parchment" mb="2">
        {title}
      </Text>
      <Text fontFamily="body" fontSize="sm" color="warmgray" lineHeight="tall">
        {description}
      </Text>
    </Box>
  )
}
