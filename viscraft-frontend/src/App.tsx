import { ChakraProvider } from '@chakra-ui/react'
import { SWRConfig } from 'swr'
import { system } from './components/styles/theme'
import { AppRouter } from './routes'

function App() {
  return (
    <ChakraProvider value={system}>
      <SWRConfig
        value={{
          revalidateOnFocus: false,
          shouldRetryOnError: false,
          dedupingInterval: 2000,
        }}
      >
        <AppRouter />
      </SWRConfig>
    </ChakraProvider>
  )
}

export default App
