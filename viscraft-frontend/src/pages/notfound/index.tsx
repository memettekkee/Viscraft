import { Box, Button, Flex, Text } from '@chakra-ui/react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../../store/authStore'


export function NotFoundPage() {
  const navigate = useNavigate()
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)

  return (
    <Box bg="ink" minH="100vh" display="flex" alignItems="center" justifyContent="center">
      <Flex direction="column" align="center" textAlign="center" px="6" gap="4">
        {/* Decorative number */}
        <Text
          fontFamily="display"
          fontSize={{ base: '7xl', md: '9xl' }}
          color="rgba(201,118,44,0.15)"
          lineHeight="1"
          fontWeight="bold"
          userSelect="none"
        >
          404
        </Text>

        {/* Stamp badge */}
        <Box
          px="3"
          py="1"
          border="1px solid"
          borderColor="amber"
          borderRadius="sm"
          transform="rotate(-2deg)"
          mt="-6"
        >
          <Text fontFamily="mono" fontSize="xs" color="amber" textTransform="uppercase" letterSpacing="wider">
            Page Not Found
          </Text>
        </Box>

        <Text
          fontFamily="display"
          fontSize={{ base: 'xl', md: '2xl' }}
          color="parchment"
          mt="2"
        >
          This page doesn't exist on the map.
        </Text>

        <Text fontFamily="body" fontSize="sm" color="warmgray" maxW="360px" lineHeight="tall">
          The page you're looking for may have been moved, deleted, or never existed.
        </Text>

        <Button
          mt="2"
          bg="amber"
          color="white"
          fontFamily="body"
          px="6"
          py="5"
          minH="44px"
          borderRadius="sm"
          _hover={{ opacity: 0.9 }}
          onClick={() => navigate(isAuthenticated ? '/workspace' : '/')}
        >
          {isAuthenticated ? 'Back to Workspace' : 'Back to Home'}
        </Button>
      </Flex>
    </Box>
  )
}
