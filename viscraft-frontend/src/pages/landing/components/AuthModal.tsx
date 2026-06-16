import { Box, Button, Flex, Input, Text } from '@chakra-ui/react'
import { ReusableModal } from '../../../components/ReusableModal'
import { useAuthForm } from '../hooks/useAuthForm'

/**
 * AuthModal with login/register tabs.
 * UI-only component — logic lives in useAuthForm hook.
 *
 * Validates: Requirements 1.1, 1.2, 1.3, 1.4, 1.5, 2.1, 2.2, 2.3, 2.4, 2.5
 */

interface AuthModalProps {
  isOpen: boolean
  onClose: () => void
  defaultTab?: 'login' | 'register'
}

export function AuthModal({ isOpen, onClose }: AuthModalProps) {
  const {
    activeTab,
    email,
    setEmail,
    password,
    setPassword,
    name,
    setName,
    emailError,
    setEmailError,
    passwordError,
    setPasswordError,
    apiError,
    isSubmitting,
    hasValidationErrors,
    handleTabSwitch,
    handleSubmit,
  } = useAuthForm(onClose)

  return (
    <ReusableModal
      isOpen={isOpen}
      onClose={onClose}
      title={activeTab === 'login' ? 'Welcome Back' : 'Create Account'}
    >
      {/* Tab switcher */}
      <Flex gap="2" mb="6">
        <Button
          flex="1"
          minH="44px"
          variant={activeTab === 'login' ? 'solid' : 'outline'}
          onClick={() => handleTabSwitch('login')}
        >
          Login
        </Button>
        <Button
          flex="1"
          minH="44px"
          variant={activeTab === 'register' ? 'solid' : 'outline'}
          onClick={() => handleTabSwitch('register')}
        >
          Register
        </Button>
      </Flex>

      {/* API error */}
      {apiError && (
        <Box
          bg="red.50"
          border="1px solid"
          borderColor="oxblood"
          borderRadius="sm"
          p="3"
          mb="4"
        >
          <Text color="oxblood" fontSize="sm" fontFamily="mono">
            {apiError}
          </Text>
        </Box>
      )}

      {/* Email field */}
      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Email
        </Text>
        <Input
          type="email"
          placeholder="you@example.com"
          value={email}
          onChange={(e) => {
            setEmail(e.target.value)
            setEmailError(null)
          }}
          borderColor={emailError ? 'oxblood' : undefined}
        />
        {emailError && (
          <Text color="oxblood" fontSize="xs" mt="1">
            {emailError}
          </Text>
        )}
      </Box>

      {/* Password field */}
      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Password
        </Text>
        <Input
          type="password"
          placeholder="Min 8 characters"
          value={password}
          onChange={(e) => {
            setPassword(e.target.value)
            setPasswordError(null)
          }}
          borderColor={passwordError ? 'oxblood' : undefined}
        />
        {passwordError && (
          <Text color="oxblood" fontSize="xs" mt="1">
            {passwordError}
          </Text>
        )}
      </Box>

      {/* Name field (register only) */}
      {activeTab === 'register' && (
        <Box mb="4">
          <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
            Name (optional)
          </Text>
          <Input
            type="text"
            placeholder="Your name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </Box>
      )}

      {/* Submit button */}
      <Button
        width="full"
        variant="solid"
        mt="2"
        minH="44px"
        disabled={hasValidationErrors || isSubmitting}
        onClick={handleSubmit}
      >
        {isSubmitting
          ? 'Please wait...'
          : activeTab === 'login'
            ? 'Log In'
            : 'Create Account'}
      </Button>
    </ReusableModal>
  )
}
