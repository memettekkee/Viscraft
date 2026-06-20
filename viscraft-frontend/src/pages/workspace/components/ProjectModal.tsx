import { Box, Button, Input, Text, Textarea, NativeSelectRoot, NativeSelectField } from '@chakra-ui/react'
import { ReusableModal } from '../../../components/ReusableModal'
import { useProjectForm } from '../hooks/useProjectForm'

const PRODUCT_CATEGORIES = [
  { value: 'general', label: 'General' },
  { value: 'food', label: 'Food & Snacks' },
  { value: 'beverage', label: 'Beverages' },
  { value: 'cosmetics', label: 'Cosmetics & Skincare' },
  { value: 'fashion', label: 'Fashion & Accessories' },
  { value: 'electronics', label: 'Electronics' },
  { value: 'home', label: 'Home & Living' },
]

export interface ProjectModalProps {
  isOpen: boolean
  onClose: () => void
}

export function ProjectModal({ isOpen, onClose }: ProjectModalProps) {
  const {
    name,
    setName,
    description,
    setDescription,
    productCategory,
    setProductCategory,
    visualStyle,
    setVisualStyle,
    nameError,
    apiError,
    isSubmitting,
    isValid,
    handleSubmit,
    resetForm,
  } = useProjectForm(onClose)

  function handleClose() {
    resetForm()
    onClose()
  }

  return (
    <ReusableModal isOpen={isOpen} onClose={handleClose} title="New Campaign">
      {apiError && (
        <Box bg="red.50" border="1px solid" borderColor="oxblood" borderRadius="sm" p="3" mb="4">
          <Text color="oxblood" fontSize="sm" fontFamily="mono">{apiError}</Text>
        </Box>
      )}

      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Campaign Name *
        </Text>
        <Input
          type="text"
          placeholder="e.g. Summer Coffee Launch"
          value={name}
          minHeight="44px"
          onChange={(e) => setName(e.target.value)}
          borderColor={nameError ? 'oxblood' : undefined}
        />
        {nameError && <Text color="oxblood" fontSize="xs" mt="1">{nameError}</Text>}
      </Box>

      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Description (optional)
        </Text>
        <Textarea
          placeholder="Brief description of your campaign"
          value={description}
          minHeight="70px"
          onChange={(e) => setDescription(e.target.value)}
        />
      </Box>

      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Product Category
        </Text>
        <NativeSelectRoot disabled={isSubmitting} width="100%" height="44px">
          <NativeSelectField
            value={productCategory}
            onChange={(e) => setProductCategory(e.target.value)}
            px="3"
            fontFamily="body"
            fontSize="sm"
            bg="parchment"
            color="ink"
            borderWidth="1px"
            borderColor="amber"
            borderRadius="sm"
            height="44px"
          >
            {PRODUCT_CATEGORIES.map((opt) => (
              <option key={opt.value} value={opt.value}>{opt.label}</option>
            ))}
          </NativeSelectField>
        </NativeSelectRoot>
      </Box>

      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Visual Style (optional)
        </Text>
        <Input
          type="text"
          placeholder="e.g. Clean minimal, Luxury dark, Playful colorful"
          value={visualStyle}
          minHeight="44px"
          onChange={(e) => setVisualStyle(e.target.value)}
          disabled={isSubmitting}
        />
        <Text color="warmgray" fontSize="xs" mt="1" fontStyle="italic">
          Influences the generated prompt style
        </Text>
      </Box>

      <Button
        width="full"
        variant="solid"
        mt="2"
        minHeight="44px"
        disabled={!isValid || isSubmitting}
        onClick={handleSubmit}
      >
        {isSubmitting ? 'Creating...' : 'Create Campaign'}
      </Button>
    </ReusableModal>
  )
}
