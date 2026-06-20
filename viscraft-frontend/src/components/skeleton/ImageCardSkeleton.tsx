import { Box, Text } from "@chakra-ui/react";
import { keyframes } from "@emotion/react";

const shimmer = keyframes`
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
`;

interface ImageCardSkeletonProps {
  label?: string;
}

export function ImageCardSkeleton({ label = "Mapping..." }: ImageCardSkeletonProps) {
  return (
    <Box
      position="relative"
      bg="parchment"
      borderWidth="1px"
      borderColor="amber"
      borderRadius="md"
      overflow="hidden"
      aspectRatio="4/3"
      display="flex"
      alignItems="center"
      justifyContent="center"
    >
      <Box
        position="absolute"
        inset="0"
        css={{
          background:
            "linear-gradient(90deg, transparent 25%, rgba(201, 118, 44, 0.08) 50%, transparent 75%)",
          backgroundSize: "200% 100%",
          animation: `${shimmer} 1.8s ease-in-out infinite`,
        }}
      />
      <Text
        fontFamily="mono"
        fontSize="sm"
        color="warmgray"
        position="relative"
        zIndex={1}
        userSelect="none"
      >
        {label}
      </Text>
    </Box>
  );
}
