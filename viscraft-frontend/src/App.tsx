import { ChakraProvider } from '@chakra-ui/react'
import { SWRConfig } from 'swr'
import { system } from './components/styles/theme'
import { AppRouter } from './routes'
import { ViscraftToaster } from './components/CustomToast'

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
        <ViscraftToaster />
      </SWRConfig>
    </ChakraProvider>
  )
}

export default App
